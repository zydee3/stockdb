package handlers

import (
	"fmt"

	"github.com/zydee3/stockdb/internal/factory/manager"
	"github.com/zydee3/stockdb/internal/unix/messages"
)

func OnUnknownRequest(cmd messages.Command, _ *manager.Manager) messages.Response {
	return messages.Response{
		Type:    messages.ResponseTypeSuccess,
		Message: fmt.Sprintf("Received Unknown Command: %s", cmd.Type.String()),
	}
}
