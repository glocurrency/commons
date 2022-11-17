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
	OrderedKeyOpt PubSubOptionType = iota
	OrderedByTaskNameOpt
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
	uniqueKeyOption         string
	orderedByTaskNameOption bool
	orderedKeyOption        string
	processAtOption         time.Time
)

// UniqueKey returns an option to specify the unique key.
func UniqueKey(key string) TaskOption {
	return uniqueKeyOption(key)
}

func (key uniqueKeyOption) String() string       { return fmt.Sprintf("UniqueKey(%q)", string(key)) }
func (key uniqueKeyOption) Type() TaskOptionType { return UniqueKeyOpt }
func (key uniqueKeyOption) Value() interface{}   { return string(key) }

// Ordered returns an option to specify the ordered key.
func OrderedByTaskName() PubSubOption {
	return orderedByTaskNameOption(true)
}

func (orderedByTaskNameOption) String() string         { return "OrderedByTaskName()" }
func (orderedByTaskNameOption) Type() PubSubOptionType { return OrderedByTaskNameOpt }
func (orderedByTaskNameOption) Value() interface{}     { return true }

// Ordered returns an option to specify the ordered key.
func OrderedKey(key string) PubSubOption {
	return orderedKeyOption(key)
}

func (key orderedKeyOption) String() string         { return fmt.Sprintf("OrderedKey(%q)", string(key)) }
func (key orderedKeyOption) Type() PubSubOptionType { return OrderedKeyOpt }
func (key orderedKeyOption) Value() interface{}     { return string(key) }

// ProcessAt returns an option to specify when to process the given task.
func ProcessAt(t time.Time) CloudTasksOption {
	return processAtOption(t)
}

func (t processAtOption) String() string {
	return fmt.Sprintf("ProcessAt(%v)", time.Time(t).Format(time.UnixDate))
}
func (t processAtOption) Type() CloudTasksOptionType { return ProcessAtOpt }
func (t processAtOption) Value() interface{}         { return time.Time(t) }
