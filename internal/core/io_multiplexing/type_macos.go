//go:build darwin

package io_multiplexing

import "syscall"

func (e Event) toNative(flags uint16) syscall.Kevent_t {
	filter := syscall.EVFILT_READ
	if e.Op == OpWrite {
		filter = syscall.EVFILT_WRITE
	}
	return syscall.Kevent_t{
		Ident:  uint64(e.Fd),
		Filter: int16(filter),
		Flags:  flags,
	}
}

func createEvent(kq syscall.Kevent_t) Event {
	op := OpRead
	if kq.Filter == syscall.EVFILT_WRITE {
		op = OpWrite
	}
	return Event{
		Fd: int(kq.Ident),
		Op: op,
	}
}
