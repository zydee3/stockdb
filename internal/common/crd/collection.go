package crd

import "fmt"

// + Implements github.com/zydee3/stockdb/internal/common/crd.CRD interface

type DataCollection struct {
	APIVersion string   `db:"api_version" mapstructure:"APIVersion" yaml:"apiVersion"`
	Kind       Kind     `db:"kind" mapstructure:"Kind" yaml:"kind"`
	Metadata   MetaData `db:"metadata" mapstructure:"Metadata" yaml:"metadata"`
	Spec       Spec     `db:"spec" mapstructure:"Spec" yaml:"spec"`
}

func (dc *DataCollection) GetAPIVersion() string {
	return dc.APIVersion
}

func (dc *DataCollection) GetKind() Kind {
	return dc.Kind
}

func (dc *DataCollection) GetName() string {
	return dc.Metadata.Name
}

func (dc *DataCollection) GetSource() Source {
	return dc.Spec.Source
}

func (dc *DataCollection) GetSchedule() Schedule {
	return dc.Spec.Schedule
}

func (dc *DataCollection) GetSecurities() []Security {
	return dc.Spec.Targets.Securities
}

func (dc *DataCollection) GetOptions() Options {
	return dc.Spec.Options
}

// TODO: Oscar
func (dc *DataCollection) GetJobCount() int {
	return 0
}

// TODO: Oscar
func (dc *DataCollection) Split(batchSize int) []CRD {
	jobCount := dc.GetJobCount()

	splitSize := jobCount / batchSize
	if jobCount%batchSize != 0 {
		splitSize++
	}

	splitCRDs := make([]CRD, splitSize)
	return splitCRDs
}

func (dc *DataCollection) String() string {
	return fmt.Sprintf("DataCollection(APIVersion: %s, Kind: %s, Metadata: %s, Spec: %s)",
		dc.APIVersion,
		dc.Kind.String(),
		dc.Metadata.String(),
		dc.Spec.String(),
	)
}
