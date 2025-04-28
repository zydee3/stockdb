package messages

type CommandType string

const (
	CommandTypeApply   CommandType = "apply"
	CommandTypeUnknown CommandType = "unknown"
)

func (t CommandType) String() string {
	return string(t)
}

func NewCommandType(s string) CommandType {
	switch s {
	case "apply":
		return CommandTypeApply
	default:
		return CommandTypeUnknown
	}
}

type Command struct {
	Type       CommandType       `json:"type"`
	Parameters map[string]string `json:"parameters"`
	Data       any               `json:"data,omitempty"`
}
