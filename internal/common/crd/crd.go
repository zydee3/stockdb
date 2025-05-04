package crd

type CRD interface {
	GetAPIVersion() string
	GetKind() Kind
	GetName() string
	GetSource() Source
	GetSchedule() Schedule
	GetSecurities() []Security
	GetOptions() Options
	GetJobCount() int
	Split(batchSize int) []CRD
	String() string
}
