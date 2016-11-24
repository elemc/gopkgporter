package jobs

import (
	"fmt"
	"gopkgporter/app/common"
	"gopkgporter/app/models"
	"os/exec"
	"strings"

	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"
)

type PushPackages struct {
	jobs.Job
}

func (c PushPackages) Run() {
	dbgorm, err := common.GetGORM()
	if err != nil {
		revel.ERROR.Printf("Error in JOB PushPackages get GORM: %s", err)
		return
	}
	defer dbgorm.Close()

	var pkgs []models.PackagesToPush

	ctx := dbgorm.Find(&pkgs, "done=?", false)
	if ctx.Error != nil {
		revel.ERROR.Printf("Error in JOB PushPackages fetch packages: %s", err)
		return
	}

	script := revel.Config.StringDefault("porter.push_script", "./koji-pp")
	revel.INFO.Printf("Porter use script from [%s]", script)

	for _, pkg := range pkgs {
		buildID := fmt.Sprintf("--id %d", pkg.BuildID)
		ver := fmt.Sprintf("--ver %s", pkg.Version)
		repo := fmt.Sprintf("--repo %s", pkg.Repository)
		branch := fmt.Sprintf("--branch %s", pkg.Branch)
		dist := fmt.Sprintf("--dist %s", pkg.Distributive)

		cmd := exec.Command(script, buildID, ver, repo, branch, dist)

		revel.INFO.Printf("Command: %s", strings.Join(cmd.Args, " "))
		data, err := cmd.CombinedOutput()
		revel.INFO.Printf("Command output: [%s]", data)
		if err != nil {
			revel.ERROR.Printf("Error start command script: %s", err)
			continue
		}

		pkg.Done = true
		ctx = dbgorm.Save(&pkg)
		if ctx.Error != nil {
			revel.ERROR.Printf("Error save pushed packages done flag: %s", err)
			continue
		}
	}
}

func init() {
	revel.OnAppStart(func() {
		jobs.Schedule("@every 1m", PushPackages{})
	})
}
