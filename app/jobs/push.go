package jobs

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/elemc/gopkgporter/app/common"
	"github.com/elemc/gopkgporter/app/models"

	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"
)

// PushPackages job
type PushPackages struct {
	*jobs.Job
}

// Run functions is a main function for PushPackages job
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
	if len(pkgs) == 0 {
		return
	}

	script := revel.Config.StringDefault("porter.push_script", "./koji-pp")
	revel.INFO.Printf("Porter use script from [%s]", script)

	dists := make(map[string]int)
	for _, pkg := range pkgs {
		val, ok := dists[pkg.Distributive]
		if !ok {
			val = 0
		}
		dists[pkg.Distributive] = val + 1
	}

	for k := range dists {
		err = c.createDotFile(".startpush", k)
		if err != nil {
			revel.ERROR.Printf("Error in creating start dotfile: %s", err)
			return
		}
	}

	for _, pkg := range pkgs {
		var argList []string
		argList = append(argList, fmt.Sprintf("--id %d", pkg.BuildID))
		argList = append(argList, fmt.Sprintf("--ver %s", pkg.Version))
		argList = append(argList, fmt.Sprintf("--repo %s", pkg.Repository))
		argList = append(argList, fmt.Sprintf("--branch %s", pkg.Branch))
		argList = append(argList, fmt.Sprintf("--dist %s", pkg.Distributive))

		cmd := exec.Command(script, strings.Join(argList, " "))

		revel.INFO.Printf("Command: %s", strings.Join(cmd.Args, " "))
		var data []byte
		data, err = cmd.CombinedOutput()
		if err != nil {
			revel.ERROR.Printf("Error start command script: %s", err)
			revel.INFO.Printf("Command output: [%s]", data)
			continue
		}

		pkg.Done = true
		ctx = dbgorm.Save(&pkg)
		if ctx.Error != nil {
			revel.ERROR.Printf("Error save pushed packages done flag: %s", err)
			continue
		}
	}

	for k := range dists {
		err = c.createDotFile(".endpush", k)
		if err != nil {
			revel.ERROR.Printf("Error in creating start dotfile: %s", err)
			return
		}
	}
	if err != nil {
		revel.ERROR.Printf("Error in creating end dotfile: %s", err)
		return
	}
}

func init() {
	revel.OnAppStart(func() {
		jobs.Schedule("@every 1m", PushPackages{})
	})
}

func (c *PushPackages) createDotFile(filename string, dist string) (err error) {
	tempDir, err := ioutil.TempDir("", "gopkgporter")
	if err != nil {
		return err
	}

	fn := filepath.Join(tempDir, filename)
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	cmd := exec.Command("rsync", "-av4", fn, fmt.Sprintf("koji.tigro.info::%s", dist))
	data, err := cmd.CombinedOutput()
	if err != nil {
		revel.WARN.Printf("Output command [%s]: %s", strings.Join(cmd.Args, " "), data)
		return err
	}

	return
}
