package daemon

import (
	"fmt"
	"os"

	"context"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/urfave/cli"
	"github.com/zydee3/stockdb/internal/api/server"
	"github.com/zydee3/stockdb/internal/api/socket"
	"github.com/zydee3/stockdb/internal/common/logger"
)

type Daemon struct {
	ctx           context.Context
	cancelFunc    context.CancelFunc
	serviceGroup  sync.WaitGroup
	errors        chan error
	shutdownTimer *time.Timer
}

func NewDaemon() *Daemon {
	ctx, cancel := context.WithCancel(context.Background())
	return &Daemon{
		ctx:        ctx,
		cancelFunc: cancel,
		errors:     make(chan error, 10), // Buffer for component errors
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

// runSocketServer launches the socket API server as a managed component
func (d *Daemon) runSocketServer() {
	defer d.serviceGroup.Done()

	err := server.StartServer(socket.SOCKET_PATH, d.ctx)

	// If the context is cancelled, it means the daemon is shutting down
	// and we don't want to report that as an error
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
		logger.Errorf("Component error: %s", err.Error())
		return d.Shutdown()
	}
}

func (d *Daemon) Shutdown() error {
	// Set shutdown deadline - don't wait forever
	d.shutdownTimer = time.NewTimer(30 * time.Second)

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
	app := cli.NewApp()
	app.Name = "stockd"
	app.Usage = "Daemon for StockDB"
	app.Version = "0.1.0"

	app.Action = func(c *cli.Context) error {
		d := NewDaemon()
		if err := d.Run(); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logger.Errorf("Error: %s", err.Error())
		os.Exit(1)
	}
}
