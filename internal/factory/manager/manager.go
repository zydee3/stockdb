package manager

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zydee3/stockdb/internal/common/crd"
	"github.com/zydee3/stockdb/internal/common/logger"
)

type Manager struct {
	ctx      context.Context
	cancel   context.CancelFunc
	database *Storage
	wg       sync.WaitGroup
	mu       sync.RWMutex
}

func NewManager(ctx context.Context) (*Manager, error) {
	mctx, cancel := context.WithCancel(ctx)

	storage, err := NewStorage()
	if err != nil {
		cancel()
		return nil, err
	}

	return &Manager{
		ctx:      mctx,
		cancel:   cancel,
		wg:       sync.WaitGroup{},
		mu:       sync.RWMutex{},
		database: storage,
	}, nil
}

func (m *Manager) Shutdown(timeout time.Duration) error {
	logger.Info("Shutting down manager...")

	// Signal all background goroutines to stop by canceling the context
	m.cancel()

	// Wait for operations to complete with timeout
	c := make(chan struct{})
	go func() {
		// Wait for all goroutines to finish
		m.wg.Wait()
		close(c)
	}()

	// Wait for either completion or timeout
	select {
	case <-c:
		logger.Info("Manager stopped successfully")
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("manager shutdown timed out after %s", timeout)
	}
}

func (m *Manager) SaveCRD(resource crd.CRD) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate the resource

	// Save the resource to the database
	logger.Infof("Saving CRD: %s", resource.GetName())

	return resource.GetName(), nil
}
