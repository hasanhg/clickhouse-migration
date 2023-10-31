package database

import (
	"log"
	"time"

	"github.com/google/uuid"
	null "gopkg.in/guregu/null.v3"
)

type Audit struct {
	ID            string    `db:"id" json:"id"`
	WorkspaceID   string    `db:"workspace_id" json:"workspace_id"`
	UserID        string    `db:"user_id" json:"user_id"`
	Category      string    `db:"category" json:"category"`
	Action        string    `db:"action" json:"action"`
	Description   string    `db:"description" json:"description"`
	Data          []byte    `db:"data" json:"data"`
	PreviousState []byte    `db:"previous_state" json:"previous_state"`
	NextState     []byte    `db:"next_state" json:"next_state"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

func (j *Audit) Create() error {
	now := time.Now()
	j.ID = uuid.NewString()
	j.CreatedAt = now
	return ch.Create(j).Error
}

type AuditPG struct {
	ID            string    `db:"id" json:"id"`
	UserID        string    `db:"user_id" json:"user_id"`
	Username      string    `db:"full_name" json:"username"`
	WorkspaceID   string    `db:"workspace_id" json:"workspace_id"`
	Category      string    `db:"category" json:"category"`
	Action        string    `db:"action" json:"action"`
	Description   string    `db:"description" json:"description"`
	Data          []byte    `db:"data" json:"data"`
	DoneAt        null.Time `db:"done_at" json:"done_at"`
	PreviousState []byte    `db:"previous_state" json:"previous_state"`
	NextState     []byte    `db:"next_state" json:"next_state"`
}

func MigrateAudits() {
	var (
		size = 1_000_000
		page = 1
	)

	for {
		var (
			audits   []*AuditPG
			chAudits = []*Audit{}
		)

		query := `select a.*, u.full_name from workspace.audit a
			inner join workspace.users u on a.user_id = u.id 
			offset $1 rows fetch next $2 rows only`

		if _, err := db.Select(&audits, query, (page-1)*size, size); err != nil {
			log.Println("Migrate audits failed:", err)
			return
		}

		log.Println("Selected from pg")

		for _, a := range audits {
			chAudits = append(chAudits, &Audit{
				ID:            a.ID,
				WorkspaceID:   a.WorkspaceID,
				UserID:        a.UserID,
				Category:      a.Category,
				Action:        a.Action,
				Description:   a.Description,
				Data:          a.Data,
				PreviousState: a.PreviousState,
				NextState:     a.NextState,
				CreatedAt:     a.DoneAt.ValueOrZero(),
			})
		}

		if len(chAudits) == 0 {
			log.Println("Migration completed")
			return
		}

		log.Println("Data convertion done")

		err := ch.Create(chAudits).Error
		if err != nil {
			log.Println("Migrate audits failed on create", err)
			return
		}

		log.Println("Migration page completed:", page)
		page++
	}
}
