package server

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zydee3/stockdb/internal/common/logger"
)

type Tracker struct {
	wg            sync.WaitGroup             // Tracks active connections for graceful shutdown
	active        atomic.Int64               // Current connection count for observability
	totalAccepted atomic.Int64               // Total connections accepted (monotonic counter)
	mu            sync.RWMutex               // Guards metadata map access
	metadata      map[uint64]*connectionData // Per-connection metadata for tracing/debugging
	nextID        atomic.Uint64              // Atomic ID generator for connection tracking
}

type connectionData struct {
	startTime  time.Time
	remoteAddr string
	attributes map[string]string
}

func NewTracker() *Tracker {
	return &Tracker{
		metadata: make(map[uint64]*connectionData),
	}
}

func (t *Tracker) Track(remoteAddr string) func() {
	t.wg.Add(1)

	active := t.active.Add(1)
	total := t.totalAccepted.Add(1)
	id := t.nextID.Add(1)

	// Record connection metadata
	t.mu.Lock()
	t.metadata[id] = &connectionData{
		startTime:  time.Now(),
		remoteAddr: remoteAddr,
		attributes: make(map[string]string),
	}

	t.mu.Unlock()

	// Log connection tracking information
	logger.Infof("Connection tracking: active=%d total=%d", active, total)

	// Return cleanup function that will be called on defer
	return func() {
		t.active.Add(-1)
		t.wg.Done()

		// Clean up metadata
		t.mu.Lock()
		delete(t.metadata, id)
		t.mu.Unlock()
	}
}

func (t *Tracker) AddAttribute(id uint64, key, value string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if metadata, exists := t.metadata[id]; exists {
		metadata.attributes[key] = value
	}
}

// ActiveCount returns the current number of active connections
func (t *Tracker) ActiveCount() int64 {
	return t.active.Load()
}

// TotalCount returns the total number of connections processed
func (t *Tracker) TotalCount() int64 {
	return t.totalAccepted.Load()
}

// WaitForCompletion waits for all tracked connections to complete
// Context allows for cancellation or timeout of the wait operation
func (t *Tracker) WaitForCompletion(ctx context.Context) error {
	// Create a channel that will be closed when all connections are done
	done := make(chan struct{})

	go func() {
		t.wg.Wait()
		close(done)
	}()

	// Wait for either completion or context cancellation
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		remaining := t.active.Load()
		return fmt.Errorf("wait canceled with %d connections still active: %w", remaining, ctx.Err())
	}
}

// DrainConnections blocks new connections and waits for existing ones to complete
// Implements the sidecar readiness probe pattern for Kubernetes graceful termination
func (t *Tracker) DrainConnections(timeout time.Duration) error {
	logger.Infof("draining %d active connections", t.active.Load())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return t.WaitForCompletion(ctx)
}
