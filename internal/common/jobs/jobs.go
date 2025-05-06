package jobs

import (
	"time"

	"github.com/zydee3/stockdb/internal/common/crd"
)

// Note: Oscar
// We dont need a JobType here because we can use the CRD Kind as the identity.

// Job is a struct containing the job being handled by the manager.
type Job struct {
	CRD       crd.CRD   `json:"crd"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Status    string    `json:"status"`
}

// type Job struct {
// 	ID                int64               `db:"id"`
// 	JobID             string              `db:"job_id"`
// 	JobKind           crd              `db:"job_type"`
// 	CreatedAt         time.Time           `db:"created_at"`
// 	UpdatedAt         time.Time           `db:"updated_at"`
// 	ScheduleType      crd.CRDScheduleType `db:"schedule_type"`
// 	ScheduleStartAt   time.Time           `db:"schedule_start_at"`
// 	ScheduleEndAt     *time.Time          `db:"schedule_end_at"`
// 	ScheduleFrequency string              `db:"schedule_frequency"`
// 	SpecJSON          string              `db:"spec_json"`
// 	Attempts          int                 `db:"attempts"`
// 	MaxRetries        int                 `db:"max_retries"`
// }
