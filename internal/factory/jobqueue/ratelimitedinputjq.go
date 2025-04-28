package jobqueue

import (
	"context"
	"github.com/zydee3/stockdb/internal/common/logger"
	rate "golang.org/x/time/rate"
)

type rateLimitedInputJQ struct {
	inputJobQueue InputJobQueue
	rateLimiter   *rate.Limiter
}

func (r *rateLimitedInputJQ) Add(context context.Context, jobDefinition JobDefinition) error {
	err := r.rateLimiter.Wait(context)
	if err != nil {
		logger.Debugf("Failed to add job definition to rate limited job queue: %v", err)
		return err
	}

	return r.inputJobQueue.Add(context, jobDefinition)
}

func NewRateLimitedInputJobQueue(targetJQ InputJobQueue, limiter *rate.Limiter) InputJobQueue {
	return &rateLimitedInputJQ{
		inputJobQueue: targetJQ,
		rateLimiter:   limiter,
	}
}
