package jobs

import (
	"database/sql"
	"fmt"
	"gopkgporter/app/common"
	"gopkgporter/app/models"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"

	_ "github.com/lib/pq"
)

type KojiBuild struct {
	ID             uint
	PkgID          uint
	Version        string
	Release        string
	Epoch          sql.NullString
	CreateEvent    uint
	CompletionTime time.Time
	State          uint
	TaskID         sql.NullInt64
	OwnerID        uint
}

type GetFromKoji struct {
	jobs.Job
}

var (
	dbgorm *gorm.DB
	koji   *sql.DB
)

func (c GetFromKoji) Run() {
	err := getBuilds()
	if err != nil {
		revel.ERROR.Printf("Get builds from koji failed: %s", err)
	}
}

func getKojiBuilds() (builds []KojiBuild, err error) {
	queryStr := `SELECT id, pkg_id, version, release, epoch, create_event,
                        completion_time, state, task_id, owner
                     FROM build
                     WHERE state='1'`

	rows, err := koji.Query(queryStr)
	if err != nil {
		revel.ERROR.Printf("Error in query [%s]: %s", queryStr, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		build := KojiBuild{}
		err = rows.Scan(&build.ID, &build.PkgID, &build.Version, &build.Release,
			&build.Epoch, &build.CreateEvent, &build.CompletionTime,
			&build.State, &build.TaskID, &build.OwnerID)
		if err != nil {
			revel.ERROR.Printf("Error in fetching information: %s", err)
			return
		}
		builds = append(builds, build)
	}
	return
}

func getBuilds() (err error) {
	// GORM database
	dbgorm, err = common.GetGORM()
	if err != nil {
		return
	}
	defer dbgorm.Close()

	// KOJI database
	koji, err = sql.Open("postgres", "user=alex dbname=koji sslmode=disable")
	if err != nil {
		revel.ERROR.Printf("Error connetction to koji database: %s", err)
		return
	}
	defer koji.Close()

	builds, err := getKojiBuilds()
	if err != nil {
		return
	}

	for _, build := range builds {
		buildedPackage := models.BuildedPackage{}
		//revel.INFO.Printf("Get builded package with ID=%d", build.ID)
		d := dbgorm.First(&buildedPackage, build.ID)
		if err := d.Error; err != nil {
			buildedPackage.BuildID = build.ID
			buildedPackage.Owner = getOwner(build.OwnerID)
			buildedPackage.BuildPackage = getPackage(build.PkgID, buildedPackage.Owner)
			buildedPackage.Version = build.Version
			buildedPackage.Release = build.Release
			buildedPackage.Epoch = build.Epoch.String
			buildedPackage.CompletionTime = build.CompletionTime
			buildedPackage.TaskID = uint(build.TaskID.Int64)
			buildedPackage.TagName = getTagNameForBuild(build.ID)
			buildedPackage.Pushed = false
			dbgorm.LogMode(true)
			dbgorm.Create(&buildedPackage)
		}
	}
	return
}

func getOwner(id uint) (owner models.Owner) {
	d := dbgorm.First(&owner, id)
	if err := d.Error; err != nil {
		queryStr := fmt.Sprintf(`SELECT name
            FROM users
            WHERE id=%d`, id)

		row := koji.QueryRow(queryStr)
		err := row.Scan(&owner.OwnerName)
		if err != nil {
			revel.ERROR.Printf("Error scan from query [%s]: %s", queryStr, err)
			return
		}
		owner.ID = id
		dbgorm.Create(&owner)
	}
	return
}

func getPackageName(id uint) (name string) {
	queryStr := fmt.Sprintf(`SELECT name
                             FROM package
                             WHERE id=%d`, id)

	row := koji.QueryRow(queryStr)
	err := row.Scan(&name)
	if err != nil {
		revel.ERROR.Printf("Error scan from query [%s]: %s", queryStr, err)
		return
	}
	return
}

func getPackageOwnerAndTagID(id uint) (repo models.Repo) {
	queryStr := fmt.Sprintf(`SELECT owner, tag_id
                             FROM tag_packages
                             WHERE package_id=%d`, id)

	var (
		ownerID uint
		tagID   uint
	)

	row := koji.QueryRow(queryStr)
	err := row.Scan(&ownerID, &tagID)
	if err != nil {
		revel.ERROR.Printf("Error scan from query [%s]: %s", queryStr, err)
		return
	}
	return
}

func getPackage(pkgID uint, owner models.Owner) (pkg models.Package) {
	d := dbgorm.First(&pkg, pkgID)
	if err := d.Error; err != nil {
		pkg.ID = pkgID
		pkg.PkgName = getPackageName(pkgID)
		pkg.PkgOwner = owner
		dbgorm.Create(&pkg)
	}

	return
}

func getTagNameForBuild(id uint) (name string) {
	queryStr := fmt.Sprintf(`SELECT tag_id
                             FROM tag_listing
                             WHERE build_id=%d`, id)

	var tagID uint

	row := koji.QueryRow(queryStr)
	err := row.Scan(&tagID)
	if err != nil {
		revel.ERROR.Printf("Error scan from query [%s]: %s", queryStr, err)
		return
	}

	queryStr = fmt.Sprintf(`SELECT name FROM tag WHERE id=%d`, tagID)
	row = koji.QueryRow(queryStr)
	err = row.Scan(&name)
	if err != nil {
		revel.ERROR.Printf("Error scan from query [%s]: %s", queryStr, err)
		return
	}

	return
}

func init() {
	revel.OnAppStart(func() {
		jobs.Schedule("@every 10m", GetFromKoji{})
	})
}
