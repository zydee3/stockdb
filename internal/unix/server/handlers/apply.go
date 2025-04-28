package handlers

import (
	"fmt"

	"github.com/zydee3/stockdb/internal/unix/messages"
)

func OnApplyRequest(cmd messages.Command) messages.Response {
	return messages.Response{
		Type:    messages.ResponseTypeSuccess,
		Message: fmt.Sprintf("Received Apply Command: %s", cmd.Type.String()),
	}
}
