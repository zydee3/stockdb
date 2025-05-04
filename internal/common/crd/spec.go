package crd

import "fmt"

type Spec struct {
	Source   Source   `db:"source" mapstructure:"source" yaml:"source"`
	Targets  Targets  `db:"targets" mapstructure:"targets" yaml:"targets"`
	Schedule Schedule `db:"schedule" mapstructure:"schedule" yaml:"schedule"`
	Options  Options  `db:"options" mapstructure:"options" yaml:"options"`
}

func (s *Spec) String() string {
	return fmt.Sprintf("Spec(Source: %s, Targets: %s, Schedule: %s, Options: %s)",
		s.Source.String(),
		s.Targets.String(),
		s.Schedule.String(),
		s.Options.String(),
	)
}
