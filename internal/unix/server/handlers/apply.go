package handlers

import (
	"github.com/mitchellh/mapstructure"

	"github.com/zydee3/stockdb/internal/common/crd"
	"github.com/zydee3/stockdb/internal/factory/manager"
	"github.com/zydee3/stockdb/internal/unix/messages"
)

func OnApplyRequest(cmd messages.Command, manager *manager.Manager) messages.Response {
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
			Message: err.Error(),
		}
	}

	// Save the job to the database
	saveID, saveError := manager.SaveCRD(job)
	if saveError != nil {
		return messages.Response{
			Type:    messages.ResponseTypeError,
			Message: saveError.Error(),
		}
	}

	return messages.Response{
		Type:    messages.ResponseTypeSuccess,
		Message: saveID,
	}
}
