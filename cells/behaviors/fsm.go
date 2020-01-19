// Tideland Go Together - Cells - Behaviors
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// FSM BEHAVIOR
//--------------------

// FSMProcessor is the signature of a function or method which processes
// an event and returns the following status or an error.
type FSMProcessor func(emitter mesh.Emitter, evt *event.Event) FSMStatus

// FSMStatus describes the current status of a finite state machine.
// It also contains a reference to the current process function.
type FSMStatus struct {
	Info    string
	Process FSMProcessor
	Error   error
}

// Done returns true if the status contains no processor anymore.
func (s FSMStatus) Done() bool {
	return s.Process == nil || s.Error != nil
}

// FSMInfo contains information about the current status of the FSM.
type FSMInfo struct {
	Info  string
	Done  bool
	Error error
}

// fsmBehavior runs the finite state machine.
type fsmBehavior struct {
	id      string
	emitter mesh.Emitter
	status  FSMStatus
}

// NewFSMBehavior creates a finite state machine behavior based on the
// passed initial status. The process function is called with the event
// and has to return the next status, which can be the same or a different
// one.
func NewFSMBehavior(id string, status FSMStatus) mesh.Behavior {
	return &fsmBehavior{
		id:     id,
		status: status,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *fsmBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *fsmBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *fsmBehavior) Terminate() error {
	return nil
}

// Process executes the state function and stores
// the returned new state.
func (b *fsmBehavior) Process(evt *event.Event) error {
	// Check if done.
	if b.status.Done() {
		return nil
	}
	// Process event and determine next status.
	switch evt.Topic() {
	case TopicFSMStatus:
		// Emit information.
		return b.emitter.Broadcast(event.New(
			event.TopicStatus,
			"info", b.status.Info,
			"done", b.status.Done(),
			"error", b.status.Error,
		))
	default:
		// Process event.
		b.status = b.status.Process(b.emitter, evt)
		return b.status.Error
	}
}

// Recover from an error.
func (b *fsmBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
