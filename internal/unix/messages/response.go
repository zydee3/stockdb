package messages

type ResponseType string

const (
	ResponseTypeSuccess ResponseType = "success"
	ResponseTypeError   ResponseType = "error"
	ResponseTypeUnknown ResponseType = "unknown"
)

func (t ResponseType) String() string {
	return string(t)
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

type Response struct {
	Type    ResponseType `json:"type"`
	Message string       `json:"message,omitempty"`
	Data    any          `json:"data,omitempty"`
}
