package jobqueue_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"golang.org/x/time/rate"

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

func TestInputJobQueueImplementations(t *testing.T) {
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
	var outputImplementations []outputFactory

	//nolint:govet // Ignore output implementations is empty until an output implementation is added.
	for _, impl := range outputImplementations {
		t.Run(impl.name+"/GetOutputChannel", func(t *testing.T) {
			testGetOutputChannelSucceeds(t, impl.newQ())
		})
	}
}

func TestFullJobQueueImplementations(t *testing.T) {
	var fullImplementations = []fullFactory{
		{
			name: "UnifiedJobQueue",
			newQ: func() jobqueue.FullJobQueue { return jobqueue.NewUnifiedJobQueue(10) },
		},
	}

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
		t.Run(impl.name+"/ConcurrentAdds", func(t *testing.T) {
			testFullJobQueueConcurrentAdds(t, impl.newQ())
		})
		t.Run(impl.name+"/ConcurrentReads", func(t *testing.T) {
			testFullJobQueueMultipleOutputChannel(t, impl.newQ())
		})
	}
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
	<-ctx.Done()

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

func testFullJobQueueConcurrentAdds(t *testing.T, q jobqueue.FullJobQueue) {
	outCh, err := q.GetOutputChannel()
	if err != nil {
		t.Fatalf("GetOutputChannel returned error: %v", err)
	}

	const numProducers = 10
	var wg sync.WaitGroup
	wg.Add(numProducers)

	for range numProducers {
		job := testJobDefinition()
		go func(j jobqueue.JobDefinition) {
			defer wg.Done()
			addErr := q.Add(context.Background(), j)
			if addErr != nil {
				t.Errorf("Add returned error for job %v: %v", j, addErr)
			}
		}(job)
	}

	// TODO: Get unique identifier of job to compare inputs end up matching outputs.
	received := 0
	for range numProducers {
		select {
		case <-outCh:
			received++
		case <-time.After(2 * time.Second):
			t.Fatalf("Timed out waiting to receive job")
		}
	}

	wg.Wait()

	if received != numProducers {
		t.Errorf("Expected %d jobs, but received %d", numProducers, received)
	}
}

func testFullJobQueueMultipleOutputChannel(t *testing.T, q jobqueue.FullJobQueue) {
	const numElements = 20

	outCh1, err1 := q.GetOutputChannel()
	outCh2, err2 := q.GetOutputChannel()
	if err1 != nil || err2 != nil {
		t.Fatalf("GetOutputChannel returned an error (err1=%v, err2=%v)", err1, err2)
	}
	if outCh1 == nil || outCh2 == nil {
		t.Fatal("Output channel is nil")
	}
	if outCh1 != outCh2 {
		t.Errorf("GetOutputChannel returned different channels on successive calls")
	}

	var wg sync.WaitGroup
	wg.Add(2)
	recvCount := 0
	var mu sync.Mutex
	consume := func(ch <-chan jobqueue.JobDefinition) {
		defer wg.Done()
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					return
				}
				mu.Lock()
				recvCount++
				mu.Unlock()
			case <-time.After(100 * time.Millisecond):
				// If no new job arrives for a while, assume no more jobs and stop.
				return
			}
		}
	}
	go consume(outCh1)
	go consume(outCh2)

	for range numElements {
		job := testJobDefinition()
		if err := q.Add(context.Background(), job); err != nil {
			t.Errorf("Unexpected error on Add(%v): %v", job, err)
		}
	}

	wg.Wait()

	if recvCount != numElements {
		t.Errorf("Expected %d jobs, but received %d", numElements, recvCount)
	}
}
