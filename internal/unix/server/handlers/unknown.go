package handlers

import (
	"fmt"

	"github.com/zydee3/stockdb/internal/unix/messages"
)

func OnUnknownRequest(cmd messages.Command) messages.Response {
	return messages.Response{
		Type:    messages.ResponseTypeSuccess,
		Message: fmt.Sprintf("Received Unknown Command: %s", cmd.Type.String()),
	}
}
