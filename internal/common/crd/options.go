package crd

import "fmt"

type Options struct {
	Timeout string `db:"timeout" mapstructure:"timeout" yaml:"timeout"`
	Retries int    `db:"retries" mapstructure:"retries" yaml:"retries"`
}

func (o *Options) String() string {
	return fmt.Sprintf("Options(Timeout: %s, Retries: %d)",
		o.Timeout,
		o.Retries,
	)
}
