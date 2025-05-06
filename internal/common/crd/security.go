package crd

import "fmt"

type Security struct {
	Symbol string `db:"symbol" mapstructure:"symbol" yaml:"symbol"`
}

func (s *Security) String() string {
	return fmt.Sprintf("Security(Symbol: %s)",
		s.Symbol,
	)
}
