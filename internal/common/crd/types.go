package crd

import (
	"strings"
)

type Kind string

const (
	KindDataCollection Kind = "datacollection"
	KindUnknown        Kind = "unknown"
)

func ParseDataCollectionKind(s string) Kind {
	s = strings.ToLower(s)
	t := Kind(s)

	switch t {
	case KindDataCollection,
		KindUnknown:
		return t
	default:
		return KindUnknown
	}
}

func (t Kind) String() string {
	return string(t)
}

type ScheduleType string

const (
	ScheduleTypeInterval  ScheduleType = "interval"
	ScheduleTypeRecurring ScheduleType = "recurring"
	ScheduleTypeUnknown   ScheduleType = "unknown"
)

func ParseDataCollectionScheduleType(s string) ScheduleType {
	s = strings.ToLower(s)
	t := ScheduleType(s)

	switch t {
	case ScheduleTypeInterval,
		ScheduleTypeRecurring,
		ScheduleTypeUnknown:
		return t
	default:
		return ScheduleTypeUnknown
	}
}

func (t ScheduleType) String() string {
	return string(t)
}
