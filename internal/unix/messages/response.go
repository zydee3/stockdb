package messages

import (
	"fmt"
)

type ResponseType string

const (
	ResponseTypeSuccess ResponseType = "success"
	ResponseTypeError   ResponseType = "error"
	ResponseTypeUnknown ResponseType = "unknown"
)

type Response struct {
	Type    ResponseType `json:"type"`
	Message string       `json:"message,omitempty"`
	Data    any          `json:"data,omitempty"`
}

func NewResponseType(s string) ResponseType {
	switch s {
	case "success":
		return ResponseTypeSuccess
	case "error":
		return ResponseTypeError
	default:
		return ResponseTypeUnknown
	}
}

func (t ResponseType) String() string {
	return string(t)
}

func (r *Response) String() string {
	return fmt.Sprintf("Response (Type: %s, Message: %s, Data: %v)", r.Type, r.Message, r.Data)
}
