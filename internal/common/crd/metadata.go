package crd

import "fmt"

type MetaData struct {
	Name string `db:"name" mapstructure:"name" yaml:"name"`
}

func (m *MetaData) String() string {
	return fmt.Sprintf("MetaData(Name: %s)", m.Name)
}
