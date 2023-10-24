package database

import (
	"log"
	"time"

	"github.com/google/uuid"
	null "gopkg.in/guregu/null.v3"
)

var (
	RobotTypes = map[string]int64{"development": 0, "ondemand": 1, "production": 2}
	STATUS     = map[int64]string{0: "%", 1: "Running", 2: "Success", 3: "Failed", 4: "Force Stopped"}
)

type Job struct {
	ID              string     `db:"id" json:"id"`
	RobotID         string     `db:"robot_id" json:"robot_id"`
	WorkspaceID     string     `db:"workspace_id" json:"workspace_id"`
	FlowID          *string    `db:"flow_id" json:"flow_id"`
	PublishedFlowID *string    `db:"published_flow_id" json:"published_flow_id"`
	RobotType       int64      `db:"robot_type" json:"robot_type"`
	RunAt           time.Time  `db:"run_at" json:"run_at"`
	StoppedAt       *time.Time `db:"stopped_at" json:"stopped_at"`
	RunningTime     int64      `db:"running_time" json:"running_time"`
	Status          string     `db:"status" json:"status"`
	Data            string     `db:"data" json:"data"`
	RobotName       string     `db:"robot_name" json:"robot_name"`
	FlowName        string     `db:"flow_name" json:"flow_name"`
	VersionName     *string    `db:"version_name" json:"version_name"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	IsDeleted       bool       `db:"is_deleted" json:"is_deleted"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at"`
}

func (j *Job) Insert() error {
	now := time.Now()
	j.ID = uuid.NewString()
	j.CreatedAt = now
	j.UpdatedAt = now
	return ch.Create(j).Error
}

func (j *Job) Update() error {
	j.UpdatedAt = time.Now()
	return ch.Create(j).Error
}

func (j *Job) Delete() (*Job, error) {
	now := time.Now()
	j.IsDeleted = true
	j.DeletedAt = &now
	return j, ch.Create(j).Error
}

type JobPG struct {
	ID              string    `db:"id" json:"id"`
	RobotID         string    `db:"robot_id" json:"robot_id"`
	WorkspaceID     string    `db:"workspace_id" json:"workspace_id"`
	FlowID          *string   `db:"flow_id" json:"flow_id"`
	PublishedFlowID *string   `db:"published_flow_id" json:"published_flow_id"`
	RobotType       int       `db:"robot_type" json:"robot_type"`
	RunAt           null.Time `db:"run_at" json:"run_at"`
	StoppedAt       null.Time `db:"stopped_at" json:"stopped_at"`
	RunningTime     int64     `db:"running_time" json:"running_time"`
	Status          string    `db:"status" json:"status"`
	Data            JSONB     `db:"data" json:"data"`
}

type JobEntry struct {
	JobPG
	RobotName   string  `db:"robot_name" json:"robot_name"`
	FlowName    string  `db:"flow_name" json:"flow_name"`
	VersionName *string `db:"version_name" json:"version_name"`
}

func MigrateJobs() {
	var (
		size = 100
		page = 1
	)

	for {
		var (
			jobs   []*JobEntry
			chJobs = []*Job{}
		)

		query := `select r."name" as robot_name, COALESCE(f."name", 'Untitled') as flow_name, j."id", j.robot_id, j.workspace_id,
		j.flow_id, j.published_flow_id, j.robot_type, j.run_at, j.running_time, j.status, j.stopped_at, j."data", 
		fv."name" as version_name from workspace.jobs j 
		inner join workspace.robots r on j.robot_id = r.id
		left join workspace.flows f on j.flow_id = f.id 
		left join workspace.published_flows pf on pf.id = j.published_flow_id 
		left join workspace.flows_versions fv on fv.id = version_id 
		order by run_at desc offset $1 rows fetch next $2 rows only`

		if _, err := db.Select(&jobs, query, (page-1)*size, size); err != nil {
			log.Println("Migrate jobs failed:", err)
			return
		}

		for _, j := range jobs {
			chJobs = append(chJobs, &Job{
				ID:              j.ID,
				RobotID:         j.RobotID,
				WorkspaceID:     j.WorkspaceID,
				FlowID:          j.FlowID,
				PublishedFlowID: j.PublishedFlowID,
				RobotType:       int64(j.RobotType),
				RunAt:           j.RunAt.ValueOrZero(),
				StoppedAt:       j.StoppedAt.Ptr(),
				RunningTime:     j.RunningTime,
				Status:          j.Status,
				Data:            string(j.Data),
				RobotName:       j.RobotName,
				FlowName:        j.FlowName,
				VersionName:     j.VersionName,
				CreatedAt:       j.RunAt.Time,
				UpdatedAt:       j.RunAt.Time,
				IsDeleted:       false,
				DeletedAt:       nil,
			})
		}

		if len(chJobs) == 0 {
			log.Println("Migration completed")
			return
		}

		err := ch.Create(chJobs).Error
		if err != nil {
			log.Println("Migrate jobs failed on create", err)
			return
		}

		log.Println("Migration page:", page)
		page++
	}
}
