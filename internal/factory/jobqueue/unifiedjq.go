package jobqueue

import (
	"context"
	"github.com/zydee3/stockdb/internal/common/logger"
)

type unifiedJobQueue struct {
	jobQueueChannel chan JobDefinition
}

func (u *unifiedJobQueue) Add(context context.Context, jobDefinition JobDefinition) error {
	select {
	case u.jobQueueChannel <- jobDefinition:
		return nil
	case <-context.Done():
		err := context.Err()
		logger.Debugf("Failed to add job definition to unified job queue: %v", err)
		return err
	}
}

func (u *unifiedJobQueue) GetOutputChannel() (<-chan JobDefinition, error) {
	return u.jobQueueChannel, nil
}

func NewUnifiedJobQueue() FullJobQueue {
	return &unifiedJobQueue{
		jobQueueChannel: make(chan JobDefinition),
	}
}
