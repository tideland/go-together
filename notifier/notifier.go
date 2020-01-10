// Tideland Go Together - Signalbox
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package notifier // import "tideland.dev/go/together/notifier"

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

// Signaler describes types with a status and able to
// notify others about its changes.
type Signaler interface {
	// Status returns the current goroutine status.
	Status() Status

	// Done waits until the given status is notified.
	Done(status Status) <-chan struct{}

	// Wait waits until the given status or duration, what comes first.
	Wait(status Status, timeout time.Duration) error
}

//--------------------
// SIGNALBOX
//--------------------

// Signalbox allows code to be notified about status changes.
type Signalbox struct {
	mu      sync.RWMutex
	status  Status
	signals [6]chan struct{}
}

// NewSignalbox creates a new Signalbox instance.
func NewSignalbox() *Signalbox {
	sb := &Signalbox{
		status: Unknown,
	}
	// Create signal channels but close the unknown
	// one immediately for correct tests.
	for s := range sb.signals {
		sb.signals[s] = make(chan struct{})
	}
	close(sb.signals[Unknown])
	return sb
}

// Notify sets the new status and informs listeners.
func (sb *Signalbox) Notify(status Status) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	if status <= sb.status || status > Stopped {
		return
	}
	for s := sb.status + 1; s <= status; s++ {
		close(sb.signals[s])
	}
	sb.status = status
}

// Status implements Signaler.
func (sb *Signalbox) Status() Status {
	sb.mu.RLock()
	defer sb.mu.RUnlock()
	return sb.status
}

// Done implements Signaler.
func (sb *Signalbox) Done(status Status) <-chan struct{} {
	if status > Stopped {
		status = Stopped
	}
	dC := sb.signals[status]
	return dC
}

// Wait implements Signaler.
func (sb *Signalbox) Wait(status Status, timeout time.Duration) error {
	if status < Unknown || status > Stopped {
		return failure.New("waiting signalbox: invalid status %v", status)
	}
	select {
	case <-sb.signals[status]:
		return nil
	case <-time.After(timeout):
		return failure.New("waiting signalbox for %v: timeout", status)
	}
}

//--------------------
// BUNDLE
//--------------------

// Bundle distributes status changes to a number of signalboxes.
type Bundle struct {
	mu          sync.RWMutex
	signalbox   *Signalbox
	signalboxes map[*Signalbox]struct{}
}

// NewBundle creates an empty bundle.
func NewBundle() *Bundle {
	b := &Bundle{
		signalbox:   NewSignalbox(),
		signalboxes: make(map[*Signalbox]struct{}),
	}
	b.signalboxes[b.signalbox] = struct{}{}
	return b
}

// Add appends one or more signalboxes to the bundle.
func (b *Bundle) Add(signalboxes ...*Signalbox) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, signalbox := range signalboxes {
		b.signalboxes[signalbox] = struct{}{}
	}
}

// Notify informs all notifiers.
func (b *Bundle) Notify(status Status) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for signalbox := range b.signalboxes {
		signalbox.Notify(status)
	}
}

// Status implements Signaler.
func (b *Bundle) Status() Status {
	return b.signalbox.Status()
}

// Done implements Signaler.
func (b *Bundle) Done(status Status) <-chan struct{} {
	return b.signalbox.Done(status)
}

// Wait implements Signaler.
func (b *Bundle) Wait(status Status, timeout time.Duration) error {
	return b.signalbox.Wait(status, timeout)
}

// EOF
