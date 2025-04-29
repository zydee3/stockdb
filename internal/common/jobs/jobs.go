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
