package io_multiplexing

type Operation uint32

const OpRead Operation = 0
const OpWrite Operation = 1

type Event struct {
	Fd int
	Op Operation
}

type IOMultiplexer interface {
	Monitor(event Event) error
	Wait() ([]Event, error)
	Close() error
}
