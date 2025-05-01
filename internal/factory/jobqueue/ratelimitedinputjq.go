package jobqueue

import (
	"context"

	"golang.org/x/time/rate"

	"github.com/zydee3/stockdb/internal/common/jobs"
	"github.com/zydee3/stockdb/internal/common/logger"
)

type rateLimitedInputJQ struct {
	inputJobQueue InputJobQueue
	rateLimiter   *rate.Limiter
}

func (r *rateLimitedInputJQ) Add(ctx context.Context, jobDefinition jobs.Job) error {
	err := r.rateLimiter.Wait(ctx)
	if err != nil {
		logger.Debugf("Failed to add job definition to rate limited job queue: %v", err)
		return err
	}

	return r.inputJobQueue.Add(ctx, jobDefinition)
}

func NewRateLimitedInputJobQueue(targetJQ InputJobQueue, limiter *rate.Limiter) InputJobQueue {
	return &rateLimitedInputJQ{
		inputJobQueue: targetJQ,
		rateLimiter:   limiter,
	}
}
