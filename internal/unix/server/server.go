package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/zydee3/stockdb/internal/common/logger"
	"github.com/zydee3/stockdb/internal/common/utility"
	"github.com/zydee3/stockdb/internal/factory/manager"
	"github.com/zydee3/stockdb/internal/unix/messages"
	"github.com/zydee3/stockdb/internal/unix/server/handlers"
)

// StartServer initializes and runs the Unix socket server ctx provides
// lifecycle control from the parent daemon.
func StartServer(ctx context.Context, socketPath string) error {
	if socketPath == "" {
		return errors.New("socket path is not set")
	}

	if err := createSocketDirectory(socketPath); err != nil {
		return err
	}

	listener, err := createSocketListener(socketPath)
	if err != nil {
		return err
	}

	defer listener.Close()

	manager, err := manager.NewManager(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer func() {
		const shutdownTimeout = 30 * time.Second
		if shutdownError := manager.Shutdown(shutdownTimeout); shutdownError != nil {
			logger.Error("%w", shutdownError)
		}
	}()

	logger.Infof("Socket server started on %s", socketPath)

	return runServer(ctx, manager, listener, socketPath)
}

func createSocketDirectory(socketPath string) error {
	// Delete the socket file if it exists
	if _, err := os.Stat(socketPath); err == nil {
		// File exists, try to remove it
		if removeError := os.Remove(socketPath); removeError != nil {
			return fmt.Errorf("%w", removeError)
		}
	} else if !os.IsNotExist(err) {
		// Some other error occurred that isnt a "file not found" error
		logger.Error("%w", err)
	}

	// Create the socket directory if it doesn't exist
	const (
		socketDirPerm = 0755
	)
	if err := utility.CreateParentDir(socketPath, socketDirPerm); err != nil {
		return err
	}

	return nil
}

func createSocketListener(socketPath string) (net.Listener, error) {
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, err
	}

	// Set the socket permissions
	//nolint:gosec // Both ends are reading and writing to the socket
	if chmodError := os.Chmod(socketPath, 0660); chmodError != nil {
		if unixCloseError := listener.Close(); unixCloseError != nil {
			logger.Warn("%w", unixCloseError)
			logger.Error("%w", unixCloseError)
		}

		// Remove the socket file
		return nil, chmodError
	}

	return listener, nil
}

func runServer(ctx context.Context, mgr *manager.Manager, listener net.Listener, socketPath string) error {
	const (
		drainTimeout = 30 * time.Second
	)

	// Create the connection tracker
	tracker := NewTracker()

	// Create a child context for go routines for graceful shutdown
	acceptCtx, cancelAccept := context.WithCancel(ctx)
	defer cancelAccept()

	// Start accepting connections in a goroutine
	acceptDone := make(chan error, 1)
	go func() {
		acceptDone <- acceptConnections(acceptCtx, mgr, listener, tracker)
	}()

	// Wait for either parent context cancellation or acceptor error
	var err error
	select {
	case <-ctx.Done():
		logger.Info("Daemon shutdown initiated")

		// Cancel the acceptor to stop new connections
		cancelAccept()

		// Allow 30 seconds for graceful drain of active connections
		drainCtx, cancelDrain := context.WithTimeout(context.Background(), drainTimeout)
		defer cancelDrain()

		// Wait for either drain completion or timeout
		drainErr := tracker.WaitForCompletion(drainCtx)
		if drainErr != nil {
			logger.Error("%w", drainErr)
		} else {
			logger.Info("Drain completed successfully")
		}

		// Wait for acceptor to exit and capture any error
		err = <-acceptDone

	case err = <-acceptDone:
		// Acceptor exited with error
		logger.Error("%w", err)
	}

	// clean up the socket file
	cleanupErr := cleanupSocket(socketPath)
	if err == nil {
		err = cleanupErr
	}

	return err
}

func acceptConnections(ctx context.Context, mgr *manager.Manager, listener net.Listener, tracker *Tracker) error {
	for {
		// Use acceptChan pattern to make listener.Accept() cancellable
		acceptChan := make(chan net.Conn, 1)
		acceptErrChan := make(chan error, 1)

		go func() {
			connection, err := listener.Accept()
			if err != nil {
				acceptErrChan <- err
				return
			}
			acceptChan <- connection
		}()

		// Wait for either new connection or shutdown signal
		select {
		case <-ctx.Done():
			// Shut down
			return ctx.Err()

		case connection := <-acceptChan:
			// New connection
			go handleConnection(connection, mgr, tracker)

		case err := <-acceptErrChan:
			// If we're shutting down, ignore accept errors
			if ctx.Err() != nil {
				return ctx.Err()
			}

			// Otherwise, log the errors and continue accepting connections
			logger.Error("%w", err)
		}
	}
}

func handleConnection(connection net.Conn, mgr *manager.Manager, tracker *Tracker) {
	var requestHandlers = map[messages.CommandType]func(messages.Command, *manager.Manager) messages.Response{
		messages.CommandTypeApply:   handlers.OnApplyRequest,
		messages.CommandTypeUnknown: handlers.OnUnknownRequest,
	}

	// Register connection with tracker and get completion function
	cleanupFn := tracker.Track(connection.RemoteAddr().String())

	// cleanupFn will decrement the WaitGroup when connection handling is done
	defer cleanupFn()

	defer connection.Close()

	cmd, err := parseCommand(connection)
	if err != nil {
		logger.Error("%w", err)
		sendErrorResponse(connection, "failed to decode command")
		return
	}

	response := requestHandlers[cmd.Type](*cmd, mgr)

	// Send response back to client
	if respError := sendResponse(connection, response); respError != nil {
		logger.Error("%w", respError)
	}
}

func parseCommand(connection net.Conn) (*messages.Command, error) {
	command := &messages.Command{}
	decoder := json.NewDecoder(connection)
	err := decoder.Decode(&command)
	return command, err
}

func sendResponse(connection net.Conn, response messages.Response) error {
	encoder := json.NewEncoder(connection)
	return encoder.Encode(response)
}

func sendErrorResponse(connection net.Conn, message string) {
	response := messages.Response{
		Type:    messages.ResponseTypeError,
		Message: message,
	}

	if err := sendResponse(connection, response); err != nil {
		logger.Error("%w", err)
	}
}

func cleanupSocket(socketPath string) error {
	if err := os.RemoveAll(socketPath); err != nil {
		logger.Error("%w", err)
		return err
	}

	return nil
}
