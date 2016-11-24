package models

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

type PackagesToPush struct {
	gorm.Model
	BuildID      uint   `gorm:"column:build_id"`
	Version      string `gorm:"column:ver"`
	Repository   string `gorm:"column:repo"`
	Branch       string `gorm:"column:branch"`
	Distributive string `gorm:"column:dist"`
	Done         bool   `gorm:"column:done"`
}

var (
	dists map[string]string
)

func init() {
	dists = make(map[string]string)
	dists["dist-rfr"] = "rf"
	dists["dist-el"] = "el"
}

func (push *PackagesToPush) Fill(build BuildedPackage, branch string) {
	for prefix, distValue := range dists {
		begin := strings.Index(build.TagName, prefix)
		if begin < 0 {
			continue
		}

		push.Distributive = distValue

		// version
		verBegin := begin + len(prefix)
		verPart := build.TagName[verBegin:]
		revel.INFO.Printf("Version part: %s", verPart)
		if strings.Contains(verPart, "rawhide") {
			push.Version = "rawhide"
		} else if strings.Contains(verPart, "devel") {
			push.Version = "pre"
		} else {
			push.Version = verPart
		}
	}

	push.Repository = build.BuildPackage.PkgRepo.RepoName
	push.Branch = branch
	push.Done = false
	push.BuildID = build.BuildID
}
