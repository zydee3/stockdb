package jobqueue

import "context"

type JobDefinition struct {
}

type FullJobQueue interface {
	InputJobQueue
	OutputJobQueue
}

type InputJobQueue interface {
	Add(context context.Context, jobDefinition JobDefinition) error
}

type OutputJobQueue interface {
	GetOutputChannel() (<-chan JobDefinition, error)
}
