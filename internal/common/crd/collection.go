package crd

// + Implements github.com/zydee3/stockdb/internal/common/crd.CRD interface

type DataCollection struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   DataCollectionMetaData `yaml:"metadata"`
	Spec       DataCollectionSpec     `yaml:"spec"`
}

type DataCollectionMetaData struct {
	Name string `yaml:"name"`
}

type DataCollectionSpec struct {
	Source   DataCollectionSource   `yaml:"source"`
	Targets  DataCollectionTargets  `yaml:"targets"`
	Schedule DataCollectionSchedule `yaml:"schedule"`
	Options  DataCollectionOptions  `yaml:"options"`
}

type DataCollectionSource struct {
	Type       string            `yaml:"type"`
	Endpoint   string            `yaml:"endpoint"`
	Parameters map[string]string `yaml:"parameters,omitempty"`
}

type DataCollectionTargets struct {
	Securities []DataCollectionSecurity `yaml:"securities"`
}

type DataCollectionSecurity struct {
	Symbol string `yaml:"symbol"`
}

type DataCollectionSchedule struct {
	Type      string `yaml:"type"                json:"type"`
	Frequency string `yaml:"frequency,omitempty" json:"frequency,omitempty"`
	StartFrom string `yaml:"startFrom,omitempty" json:"startFrom,omitempty"`
	StartDate string `yaml:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate   string `yaml:"endDate,omitempty"   json:"endDate,omitempty"`
}

type DataCollectionOptions struct {
	Timeout  string `yaml:"timeout"  json:"timeout"`
	Retries  int    `yaml:"retries"  json:"retries"`
	Priority int    `yaml:"priority" json:"priority"`
}

func (dc *DataCollection) GetAPIVersion() string {
	return dc.APIVersion
}

func (dc *DataCollection) GetKind() string {
	return dc.Kind
}

func (dc *DataCollection) GetName() string {
	return dc.Metadata.Name
}

func (dc *DataCollection) GetSource() DataCollectionSource {
	return dc.Spec.Source
}

func (dc *DataCollection) GetSchedule() DataCollectionSchedule {
	return dc.Spec.Schedule
}

func (dc *DataCollection) GetSecurities() []DataCollectionSecurity {
	return dc.Spec.Targets.Securities
}

func (dc *DataCollection) GetOptions() DataCollectionOptions {
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
