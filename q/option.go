package q

import (
	"fmt"
	"time"
)

type TaskOptionType int

const (
	UniqueKeyOpt TaskOptionType = iota
)

type PubSubOptionType int

const (
	OrderedOpt PubSubOptionType = iota
)

type CloudTasksOptionType int

const (
	ProcessAtOpt CloudTasksOptionType = iota
)

type TaskOption interface {
	// String returns a string representation of the option.
	String() string

	// Type describes the type of the option.
	Type() TaskOptionType

	// Value returns a value used to create this option.
	Value() interface{}
}

type PubSubOption interface {
	// String returns a string representation of the option.
	String() string

	// Type describes the type of the option.
	Type() PubSubOptionType

	// Value returns a value used to create this option.
	Value() interface{}
}

type CloudTasksOption interface {
	// String returns a string representation of the option.
	String() string

	// Type describes the type of the option.
	Type() CloudTasksOptionType

	// Value returns a value used to create this option.
	Value() interface{}
}

// Internal option representations.
type (
	uniqueKeyOption string
	orderedOption   string
	processAtOption time.Time
)

// UniqueKey returns an option to specify the unique key.
func UniqueKey(key string) TaskOption {
	return uniqueKeyOption(key)
}

func (key uniqueKeyOption) String() string       { return fmt.Sprintf("UniqueKey(%q)", string(key)) }
func (key uniqueKeyOption) Type() TaskOptionType { return UniqueKeyOpt }
func (key uniqueKeyOption) Value() interface{}   { return string(key) }

// Ordered returns an option to specify the ordered key.
func Ordered(key string) PubSubOption {
	return orderedOption(key)
}

func (key orderedOption) String() string         { return fmt.Sprintf("Ordered(%q)", string(key)) }
func (key orderedOption) Type() PubSubOptionType { return OrderedOpt }
func (key orderedOption) Value() interface{}     { return string(key) }

// ProcessAt returns an option to specify when to process the given task.
func ProcessAt(t time.Time) CloudTasksOption {
	return processAtOption(t)
}

func (t processAtOption) String() string {
	return fmt.Sprintf("ProcessAt(%v)", time.Time(t).Format(time.UnixDate))
}
func (t processAtOption) Type() CloudTasksOptionType { return ProcessAtOpt }
func (t processAtOption) Value() interface{}         { return time.Time(t) }
