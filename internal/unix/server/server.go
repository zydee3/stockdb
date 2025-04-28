package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/zydee3/stockdb/internal/common/logger"
	"github.com/zydee3/stockdb/internal/common/utility"
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

	logger.Infof("Socket server started on %s", socketPath)

	return runServer(ctx, listener, socketPath)
}

func createSocketDirectory(socketPath string) error {
	// Delete the socket file if it exists
	if _, err := os.Stat(socketPath); err == nil {
		// File exists, try to remove it
		if removeError := os.Remove(socketPath); removeError != nil {
			return fmt.Errorf("failed to remove socket: %s", removeError.Error())
		}
	} else if !os.IsNotExist(err) {
		// Some other error occurred that isnt a "file not found" error
		logger.Error(err.Error())
	}

	// Create the socket directory if it doesn't exist
	const (
		socketDirPerm = 0755
	)
	if err := utility.CreateParentDir(socketPath, socketDirPerm); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func createSocketListener(socketPath string) (net.Listener, error) {
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("failed to create unix socket: %s", err.Error()), 1)
	}

	// Set the socket permissions
	//nolint:gosec // Both ends are reading and writing to the socket
	if chmodError := os.Chmod(socketPath, 0660); chmodError != nil {
		if unixCloseError := listener.Close(); unixCloseError != nil {
			logger.Errorf("failed to close socket: %v", unixCloseError)
		}

		// Remove the socket file
		return nil, cli.NewExitError(fmt.Sprintf("failed to set socket permissions: %s", chmodError.Error()), 1)
	}

	return listener, nil
}

func runServer(ctx context.Context, listener net.Listener, socketPath string) error {
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
		acceptDone <- acceptConnections(acceptCtx, listener, tracker)
	}()

	// Wait for either parent context cancellation or acceptor error
	var err error
	select {
	case <-ctx.Done():
		logger.Info("daemon shutdown initiated")

		// Cancel the acceptor to stop new connections
		cancelAccept()

		// Allow 30 seconds for graceful drain of active connections
		drainCtx, cancelDrain := context.WithTimeout(context.Background(), drainTimeout)
		defer cancelDrain()

		// Wait for either drain completion or timeout
		drainErr := tracker.WaitForCompletion(drainCtx)
		if drainErr != nil {
			logger.Errorf("drain failed to complete: %v", drainErr)
		} else {
			logger.Info("drain completed successfully")
		}

		// Wait for acceptor to exit and capture any error
		err = <-acceptDone

	case err = <-acceptDone:
		// Acceptor exited with error
		logger.Errorf("acceptor exited with error: %s", err.Error())
	}

	// clean up the socket file
	cleanupErr := cleanupSocket(socketPath)
	if err == nil {
		err = cleanupErr
	}

	return err
}

func acceptConnections(ctx context.Context, listener net.Listener, tracker *Tracker) error {
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
			go handleConnection(connection, tracker)

		case err := <-acceptErrChan:
			// If we're shutting down, ignore accept errors
			if ctx.Err() != nil {
				return ctx.Err()
			}

			// Otherwise, log the errors and continue accepting connections
			logger.Errorf("Error accepting connection: %s", err.Error())
		}
	}
}

func handleConnection(connection net.Conn, tracker *Tracker) {
	var requestHandlers = map[messages.CommandType]func(messages.Command) messages.Response{
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
		logger.Errorf("Error parsing command: %s", err.Error())
		sendErrorResponse(connection, "failed to decode command")
		return
	}

	logger.Infof("Received command: %+v", *cmd)

	response := requestHandlers[cmd.Type](*cmd)

	// Send response back to client
	if respError := sendResponse(connection, response); respError != nil {
		logger.Errorf("Error sending response: %s", respError.Error())
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
		logger.Errorf("Error sending error response: %s", err.Error())
	}
}

func cleanupSocket(socketPath string) error {
	if err := os.RemoveAll(socketPath); err != nil {
		logger.Errorf("Failed to remove socket directory: %v", err)
		return err
	}

	return nil
}
