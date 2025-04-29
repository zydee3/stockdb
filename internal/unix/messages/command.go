package messages

type CommandType string

const (
	CommandTypeApply   CommandType = "apply"
	CommandTypeUnknown CommandType = "unknown"
)

type Command struct {
	Type       CommandType       `json:"type"`
	Parameters map[string]string `json:"parameters"`
	Data       any               `json:"data,omitempty"`
}

func NewCommandType(s string) CommandType {
	switch s {
	case "apply":
		return CommandTypeApply
	default:
		return CommandTypeUnknown
	}
}

func (t CommandType) String() string {
	return string(t)
}
