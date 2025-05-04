package crd

import "fmt"

type Schedule struct {
	Type      ScheduleType `db:"type" mapstructure:"type" yaml:"type"`
	Frequency string       `db:"frequency" mapstructure:"frequency" yaml:"frequency,omitempty"`
	StartDate string       `db:"startDate" mapstructure:"startDate" yaml:"startDate,omitempty"`
	EndDate   string       `db:"endDate" mapstructure:"endDate" yaml:"endDate,omitempty"`
}

func (s *Schedule) String() string {
	return fmt.Sprintf("Schedule(Type: %s, Frequency: %s, StartDate: %s, EndDate: %s)",
		s.Type.String(),
		s.Frequency,
		s.StartDate,
		s.EndDate,
	)
}
