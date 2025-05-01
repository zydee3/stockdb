package jobqueue

import (
	"context"

	"github.com/zydee3/stockdb/internal/common/jobs"
	"github.com/zydee3/stockdb/internal/common/logger"
)

type unifiedJobQueue struct {
	jobQueueChannel chan jobs.Job
}

func (u *unifiedJobQueue) Add(ctx context.Context, jobDefinition jobs.Job) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	select {
	case u.jobQueueChannel <- jobDefinition:
		return nil
	case <-ctx.Done():
		err := ctx.Err()
		logger.Debugf("Failed to add job definition to unified job queue: %v", err)
		return err
	}
}

func (u *unifiedJobQueue) GetOutputChannel() (<-chan jobs.Job, error) {
	return u.jobQueueChannel, nil
}

func NewUnifiedJobQueue(size uint) FullJobQueue {
	return &unifiedJobQueue{
		jobQueueChannel: make(chan jobs.Job, size),
	}
}
