package crd

import "fmt"

type Targets struct {
	Securities []Security `db:"securities" mapstructure:"securities" yaml:"securities"`
}

func (t *Targets) String() string {
	return fmt.Sprintf("Targets(Securities: %v)",
		t.Securities,
	)
}
