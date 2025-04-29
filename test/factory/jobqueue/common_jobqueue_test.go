package jobqueue

import (
	"context"
	"golang.org/x/time/rate"
	"testing"
	"time"

	"github.com/zydee3/stockdb/internal/factory/jobqueue"
)

type inputFactory struct {
	name string
	newQ func() jobqueue.InputJobQueue
}

type outputFactory struct {
	name string
	newQ func() jobqueue.OutputJobQueue
}

type fullFactory struct {
	name string
	newQ func() jobqueue.FullJobQueue
}

var inputImplementations = []inputFactory{
	{
		name: "RateLimitJobQueue",
		newQ: func() jobqueue.InputJobQueue {
			return jobqueue.NewRateLimitedInputJobQueue(
				jobqueue.NewUnifiedJobQueue(10),
				rate.NewLimiter(rate.Every(1*time.Second), 10),
			)
		},
	},
}

var outputImplementations []outputFactory

var fullImplementations = []fullFactory{
	{
		name: "UnifiedJobQueue",
		newQ: func() jobqueue.FullJobQueue { return jobqueue.NewUnifiedJobQueue(10) },
	},
}

func testJobDefinition() jobqueue.JobDefinition {
	return jobqueue.JobDefinition{}
}

func testAddSucceeds(t *testing.T, q jobqueue.InputJobQueue) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := q.Add(ctx, testJobDefinition()); err != nil {
		t.Fatalf("Add() failed: %v", err)
	}
}

func testAddHonorsCancel(t *testing.T, q jobqueue.InputJobQueue) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := q.Add(ctx, testJobDefinition())
	if err == nil {
		t.Fatal("expected error due to canceled context, got nil")
	}
}

func testGetOutputChannelSucceeds(t *testing.T, q jobqueue.OutputJobQueue) {
	_, err := q.GetOutputChannel()
	if err != nil {
		t.Fatalf("GetOutputChannel() error: %v", err)
	}
}

func TestInputJobQueueImplementations(t *testing.T) {
	for _, impl := range inputImplementations {
		t.Run(impl.name+"/AddSucceeds", func(t *testing.T) {
			testAddSucceeds(t, impl.newQ())
		})
		t.Run(impl.name+"/AddHonorsCancel", func(t *testing.T) {
			testAddHonorsCancel(t, impl.newQ())
		})
	}
}

func TestOutputJobQueueImplementations(t *testing.T) {
	for _, impl := range outputImplementations {
		t.Run(impl.name+"/GetOutputChannel", func(t *testing.T) {
			testGetOutputChannelSucceeds(t, impl.newQ())
		})
	}
}

func TestFullJobQueueImplementations(t *testing.T) {
	for _, impl := range fullImplementations {
		t.Run(impl.name+"/AddSucceeds", func(t *testing.T) {
			testAddSucceeds(t, impl.newQ())
		})
		t.Run(impl.name+"/AddHonorsCancel", func(t *testing.T) {
			testAddHonorsCancel(t, impl.newQ())
		})
		t.Run(impl.name+"/GetOutputChannel", func(t *testing.T) {
			testGetOutputChannelSucceeds(t, impl.newQ())
		})
	}
}
