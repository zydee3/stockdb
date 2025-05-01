package daemon

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/zydee3/stockdb/internal/common/logger"
	"github.com/zydee3/stockdb/internal/common/version"
	daemonConfig "github.com/zydee3/stockdb/internal/config"
	"github.com/zydee3/stockdb/internal/unix/server"
	"github.com/zydee3/stockdb/internal/unix/socket"
)

type Daemon struct {
	ctx           context.Context
	cancelFunc    context.CancelFunc
	serviceGroup  sync.WaitGroup
	errors        chan error
	shutdownTimer *time.Timer
}

func NewDaemon(ctx context.Context) *Daemon {
	const (
		errorChannelSize = 10
	)
	ctx, cancel := context.WithCancel(ctx)
	return &Daemon{
		ctx:        ctx,
		cancelFunc: cancel,
		errors:     make(chan error, errorChannelSize), // Buffer for component errors
	}
}

func (d *Daemon) Start() error {
	pid := os.Getpid()
	logger.Infof("Starting Daemon (PID: %d)", pid)

	services := []func(){
		d.runSocketServer,
		// todo: add other services for manager and worker
	}

	// Initialize and start each service
	for _, service := range services {
		d.serviceGroup.Add(1)
		go service()
	}

	return nil
}

// runSocketServer launches the socket API server as a managed component.
func (d *Daemon) runSocketServer() {
	defer d.serviceGroup.Done()

	err := server.StartServer(d.ctx, socket.SocketPath)

	// If the context is cancelled, it means the daemon is shutting down
	// and we don't want to report that as an error.
	if err != nil && d.ctx.Err() == nil {
		d.errors <- fmt.Errorf("socket server error: %w", err)
	}
}

func (d *Daemon) Run() error {
	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start all services
	if err := d.Start(); err != nil {
		return err
	}

	// Block until either a signal or component error
	select {
	case sig := <-sigChan:
		// Perform graceful shutdown of all services
		logger.Infof("Received shutdown signal: %s", sig)
		return d.Shutdown()

	case err := <-d.errors:
		// A component reported an error
		logger.Error("%w", err)
		return d.Shutdown()
	}
}

func (d *Daemon) Shutdown() error {
	// Set shutdown deadline - don't wait forever
	d.shutdownTimer = time.NewTimer(daemonConfig.DaemonShutdownTimeout)

	// Cancel context to signal all services to stop
	d.cancelFunc()

	// Wait for services to exit or timeout
	shutdownComplete := make(chan struct{})
	go func() {
		d.serviceGroup.Wait()
		close(shutdownComplete)
	}()

	select {
	case <-shutdownComplete:
		logger.Info("shutdown complete")
		d.shutdownTimer.Stop()
	case <-d.shutdownTimer.C:
		logger.Error("shutdown timed out")
		// todo add handling here
	}

	return nil
}

func Init() {
	cmd := &cli.Command {
		Name: "stockd",
		Description: "Daemon for StockDB",
		Version: version.GetVersion(),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			d := NewDaemon(ctx)

			if err := d.Run(); err != nil {
				return cli.Exit(err, 1)
			}

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Error("%w", err)
		os.Exit(1)
	}
}
