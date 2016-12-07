package jobs

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/elemc/gopkgporter/app/common"
	"github.com/elemc/gopkgporter/app/models"
	"github.com/jinzhu/gorm"
	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"

	// import for PostgreSQL
	_ "github.com/lib/pq"
)

// KojiBuild struct for build table from koji db
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

// GetFromKoji job
type GetFromKoji struct {
	*jobs.Job
}

var (
	dbgorm *gorm.DB
	koji   *sql.DB
)

// Run functions is a main function for all jobs
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
	koji, err = sql.Open("postgres", "user=koji dbname=koji sslmode=disable")
	if err != nil {
		revel.ERROR.Printf("Error connetction to koji database: %s", err)
		return
	}
	koji.SetMaxIdleConns(-1)
	koji.SetMaxOpenConns(5)
	defer koji.Close()

	builds, err := getKojiBuilds()
	if err != nil {
		return
	}

	for _, build := range builds {
		if isLiveMedia(build) {
			continue
		}
		err = createNewBuildPackage(build)
		if err != nil {
			log.Printf("Error in create new build package: %s", err)
		}
	}
	return
}

func isLiveMedia(build KojiBuild) (result bool) {
	queryStr := fmt.Sprintf("SELECT method FROM task WHERE id=%d", build.TaskID.Int64)

	rows, err := koji.Query(queryStr)
	if err != nil {
		revel.ERROR.Printf("Error in query [%s]: %s", queryStr, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var method string
		err = rows.Scan(&method)
		if err != nil {
			revel.ERROR.Printf("Error in fetching information: %s", err)
			continue
		}
		if method == "livemedia" {
			result = true
		}
	}
	return
}

func createNewBuildPackage(build KojiBuild) (err error) {
	buildedPackage := models.BuildedPackage{}
	d := dbgorm.First(&buildedPackage, "build_id=?", build.ID)
	if err = d.Error; err != nil {
		tagName := getTagNameForBuild(build.ID)
		if tagName == "" {
			return fmt.Errorf("in build with id=%d tag name is empty, skip it", build.ID)
		}
		buildedPackage.BuildID = build.ID
		buildedPackage.Owner = getOwner(build.OwnerID)
		buildedPackage.BuildPackage = getPackage(build.PkgID, buildedPackage.Owner)
		buildedPackage.Version = build.Version
		buildedPackage.Release = build.Release
		buildedPackage.Epoch = build.Epoch.String
		buildedPackage.CompletionTime = build.CompletionTime
		buildedPackage.TaskID = uint(build.TaskID.Int64)
		buildedPackage.TagName = tagName
		buildedPackage.Pushed = false

		ctx := dbgorm.Create(&buildedPackage)
		if ctx.Error != nil {
			err = ctx.Error
			revel.ERROR.Printf("Error in create builded package: %s", err)
		} else {
			err = nil
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
	var unknownRepo models.Repo
	q := dbgorm.First(&unknownRepo, 1)
	if q.Error != nil {
		return
	}
	d := dbgorm.First(&pkg, pkgID)
	if err := d.Error; err != nil {
		pkg.ID = pkgID
		pkg.PkgName = getPackageName(pkgID)
		pkg.PkgOwner = owner
		pkg.PkgRepo = unknownRepo
		dbgorm.Create(&pkg)
	}

	return
}

func getTagNameForBuild(id uint) (name string) {
	queryStr := fmt.Sprintf(`SELECT tag_id
                             FROM tag_listing
                             WHERE build_id=%d`, id)

	rows, err := koji.Query(queryStr)
	if err != nil {
		revel.ERROR.Printf("Error in query [%s]: %s", queryStr, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tagID uint
		err = rows.Scan(&tagID)
		if err != nil {
			revel.ERROR.Printf("Error scan from query [%s]: %s", queryStr, err)
			continue
		}

		queryStr = fmt.Sprintf(`SELECT name FROM tag WHERE id=%d`, tagID)
		row := koji.QueryRow(queryStr)
		err = row.Scan(&name)
		if err != nil {
			revel.ERROR.Printf("Error scan from query [%s]: %s", queryStr, err)
			return
		}
	}

	return
}

func init() {
	revel.OnAppStart(func() {
		jobs.Schedule("@every 1m", GetFromKoji{})
	})
}
