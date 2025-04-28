package crd

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
	Type      string `yaml:"type" json:"type"`
	Frequency string `yaml:"frequency,omitempty" json:"frequency,omitempty"`
	StartFrom string `yaml:"startFrom,omitempty" json:"startFrom,omitempty"`
	StartDate string `yaml:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate   string `yaml:"endDate,omitempty" json:"endDate,omitempty"`
}

type DataCollectionOptions struct {
	Timeout  string `yaml:"timeout" json:"timeout"`
	Retries  int    `yaml:"retries" json:"retries"`
	Priority int    `yaml:"priority" json:"priority"`
}
