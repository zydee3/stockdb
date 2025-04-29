package crd

type CRD interface {
	GetAPIVersion() string
	GetKind() string
	GetName() string
	GetSource() DataCollectionSource
	GetSchedule() DataCollectionSchedule
	GetSecurities() []DataCollectionSecurity
	GetOptions() DataCollectionOptions
	GetJobCount() int
	Split(batchSize int) []CRD
}
