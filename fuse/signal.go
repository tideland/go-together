// Tideland Go Together - Fuse
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package fuse // import "tideland.dev/go/together/fuse"

//--------------------
// IMPORTS
//--------------------

import (
	"sync"
	"time"

	"tideland.dev/go/trace/failure"
)

//--------------------
// STATUS
//--------------------

// Status describes the status of a background groroutine.
type Status int

// Different statuses of a background goroutine.
const (
	Unknown Status = iota
	Starting
	Ready
	Working
	Stopping
	Stopped
)

// statusStr contains the string representation of a status.
var statusStr = map[Status]string{
	Unknown:  "unknown",
	Starting: "starting",
	Ready:    "ready",
	Working:  "working",
	Stopping: "stopping",
	Stopped:  "stopped",
}

// String implements the fmt.Stringer interface.
func (s Status) String() string {
	if str, ok := statusStr[s]; ok {
		return str
	}
	return "invalid"
}

//--------------------
// SIGNALER
//--------------------

// Signaler describes types with a status and able to
// notify others about its changes. It's the read-only
// interface to a Signal.
type Signaler interface {
	// Status returns the current goroutine status.
	Status() Status

	// Done waits until the given status is notified.
	Done(status Status) <-chan struct{}

	// Wait waits until the given status or duration, what comes first.
	Wait(status Status, timeout time.Duration) error
}

//--------------------
// SIGNAL
//--------------------

// Signal allows code to be notified about status changes.
type Signal struct {
	mu      sync.RWMutex
	status  Status
	signals [6]chan struct{}
}

// NewSignal creates a new Signal instance.
func NewSignal() *Signal {
	s := &Signal{
		status: Unknown,
	}
	// Create signal channels but close the unknown
	// one immediately for correct tests.
	for i := range s.signals {
		s.signals[i] = make(chan struct{})
	}
	close(s.signals[Unknown])
	return s
}

// Notify sets the new status and informs listeners.
func (s *Signal) Notify(status Status) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if status <= s.status || status > Stopped {
		return
	}
	for i := s.status + 1; i <= status; i++ {
		close(s.signals[i])
	}
	s.status = status
}

// Status implements Signaler.
func (s *Signal) Status() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

// Done implements Signaler.
func (s *Signal) Done(status Status) <-chan struct{} {
	if status > Stopped {
		status = Stopped
	}
	signalc := s.signals[status]
	return signalc
}

// Wait implements Signaler.
func (s *Signal) Wait(status Status, timeout time.Duration) error {
	if status < Unknown || status > Stopped {
		return failure.New("waiting signal: invalid status %v", status)
	}
	select {
	case <-s.signals[status]:
		return nil
	case <-time.After(timeout):
		return failure.New("waiting signal for %v: timeout", status)
	}
}

// EOF
