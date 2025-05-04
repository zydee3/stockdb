package handlers

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/zydee3/stockdb/internal/common/crd"
	"github.com/zydee3/stockdb/internal/common/logger"
	"github.com/zydee3/stockdb/internal/unix/messages"
)

func OnApplyRequest(cmd messages.Command) messages.Response {
	jobData := cmd.Data
	if jobData == nil {
		return messages.Response{
			Type:    messages.ResponseTypeError,
			Message: "No job data provided",
		}
	}

	// Process the job data
	var job crd.DataCollection
	if err := mapstructure.Decode(jobData, &job); err != nil {
		return messages.Response{
			Type:    messages.ResponseTypeError,
			Message: fmt.Sprintf("Failed to decode job data: %v", err),
		}
	}

	logger.Infof("Processing job: %s", job.String())

	return messages.Response{
		Type:    messages.ResponseTypeSuccess,
		Message: fmt.Sprintf("Received Apply Command: %s", cmd.Type.String()),
	}
}
