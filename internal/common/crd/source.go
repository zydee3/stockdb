package crd

import "fmt"

type Source struct {
	Type       string            `db:"type" mapstructure:"type" yaml:"type"`
	Endpoint   string            `db:"endpoint" mapstructure:"endpoint" yaml:"endpoint"`
	Parameters map[string]string `db:"parameters" mapstructure:"parameters" yaml:"parameters,omitempty"`
}

func (s *Source) String() string {
	return fmt.Sprintf("Source(Type: %s, Endpoint: %s, Parameters: %v)",
		s.Type,
		s.Endpoint,
		s.Parameters,
	)
}
