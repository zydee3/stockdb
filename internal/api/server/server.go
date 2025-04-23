package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/urfave/cli"
	"github.com/zydee3/stockdb/internal/api/messages"
	"github.com/zydee3/stockdb/internal/api/server/handlers"
	"github.com/zydee3/stockdb/internal/common/utility"
)

var requestHandlers = map[messages.CommandType]func(messages.Command) messages.Response{
	messages.CommandTypeApply:   handlers.OnApplyRequest,
	messages.CommandTypeUnknown: handlers.OnUnknownRequest,
}

// StartServer initializes and runs the Unix socket server
// ctx provides lifecycle control from the parent daemon
func StartServer(socketPath string, ctx context.Context) error {
	if socketPath == "" {
		return fmt.Errorf("socket path is not set")
	}

	if err := createSocketDirectory(socketPath); err != nil {
		return err
	}

	listener, err := createSocketListener(socketPath)
	if err != nil {
		return err
	}

	defer listener.Close()

	fmt.Printf("StockDB daemon started and listening on %s\n", socketPath)

	return runServer(listener, socketPath, ctx)
}

func createSocketDirectory(socketPath string) error {
	// Delete the socket file if it exists
	if err := os.RemoveAll(socketPath); err != nil {
		return fmt.Errorf("failed to remove socket: %s", err.Error())
	}

	// Create the socket directory if it doesn't exist
	if err := utility.CreateParentDir(socketPath); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func createSocketListener(socketPath string) (net.Listener, error) {
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, cli.NewExitError(fmt.Sprintf("failed to create unix socket: %v", err), 1)
	}

	// Set the socket permissions
	if err := os.Chmod(socketPath, 0660); err != nil {
		listener.Close()
		return nil, cli.NewExitError(fmt.Sprintf("failed to set socket permissions: %v", err), 1)
	}

	return listener, nil
}

func runServer(listener net.Listener, socketPath string, ctx context.Context) error {
	// Create the connection tracker
	tracker := NewTracker()

	// Create a child context for go routines for graceful shutdown
	acceptCtx, cancelAccept := context.WithCancel(ctx)
	defer cancelAccept()

	// Start accepting connections in a goroutine
	acceptDone := make(chan error, 1)
	go func() {
		acceptDone <- acceptConnections(listener, acceptCtx, tracker)
	}()

	// Wait for either parent context cancellation or acceptor error
	var err error
	select {
	case <-ctx.Done():
		fmt.Println("server shutdown initiated")

		// Cancel the acceptor to stop new connections
		cancelAccept()

		// Allow 30 seconds for graceful drain of active connections
		drainCtx, cancelDrain := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelDrain()

		// Wait for either drain completion or timeout
		drainErr := tracker.WaitForCompletion(drainCtx)
		if drainErr != nil {
			fmt.Printf("drain failed to complete: %v\n", drainErr)
		} else {
			fmt.Println("drain completed successfully")
		}

		// Wait for acceptor to exit and capture any error
		err = <-acceptDone

	case err = <-acceptDone:
		// Acceptor exited with error
		fmt.Printf("Acceptor exited: %v\n", err)
	}

	// clean up the socket file
	cleanupErr := cleanupSocket(socketPath)
	if err == nil {
		err = cleanupErr
	}

	return err
}

func acceptConnections(listener net.Listener, ctx context.Context, tracker *Tracker) error {
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
			log.Printf("Error accepting connection: %v", err)
		}
	}
}

func handleConnection(connection net.Conn, tracker *Tracker) {
	// Register connection with tracker and get completion function
	cleanupFn := tracker.Track(connection.RemoteAddr().String())

	// cleanupFn will decrement the WaitGroup when connection handling is done
	defer cleanupFn()

	defer connection.Close()

	cmd, err := parseCommand(connection)
	if err != nil {
		log.Printf("Error decoding command: %v", err)
		sendErrorResponse(connection, "failed to decode command")
		return
	}

	log.Printf("Received command: %+v", *cmd)

	response := requestHandlers[cmd.Type](*cmd)

	// Send response back to client
	if err := sendResponse(connection, response); err != nil {
		log.Printf("Error encoding response: %v", err)
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
		log.Printf("Error sending error response: %s", err.Error())
	}
}

func cleanupSocket(socketPath string) error {
	if err := os.RemoveAll(socketPath); err != nil {
		log.Printf("Failed to remove socket file: %v", err)
		return err
	}

	return nil
}
