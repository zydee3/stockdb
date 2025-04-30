package jobqueue

import (
	"context"

	"github.com/zydee3/stockdb/internal/common/jobs"
	"github.com/zydee3/stockdb/internal/common/logger"
)

type unifiedJobQueue struct {
	jobQueueChannel chan jobs.Job
}

func (u *unifiedJobQueue) Add(context context.Context, jobDefinition jobs.Job) error {
	if context.Err() != nil {
		return context.Err()
	}

	select {
	case u.jobQueueChannel <- jobDefinition:
		return nil
	case <-context.Done():
		err := context.Err()
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
