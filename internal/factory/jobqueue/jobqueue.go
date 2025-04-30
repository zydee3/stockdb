package jobqueue

import (
	"context"

	"github.com/zydee3/stockdb/internal/common/jobs"
)

type FullJobQueue interface {
	InputJobQueue
	OutputJobQueue
}

type InputJobQueue interface {
	Add(context context.Context, jobDefinition jobs.Job) error
}

type OutputJobQueue interface {
	GetOutputChannel() (<-chan jobs.Job, error)
}
