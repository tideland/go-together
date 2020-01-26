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
	"tideland.dev/go/together/fuse"
	"tideland.dev/go/trace/logger"
)

//--------------------
// CALLBACK BEHAVIOR
//--------------------

// Callbacker is a function called by the behavior when it receives an event.
type Callbacker func(emitter mesh.Emitter, evt *event.Event) error

// callbackBehavior is an event processor calling all stored functions
// if it receives an event.
type callbackBehavior struct {
	id        string
	emitter   mesh.Emitter
	callbacks []Callbacker
}

// NewCallbackBehavior creates a behavior with a number of callback functions.
// Each time an event is received those functions are called in the same order
// they have been passed.
func NewCallbackBehavior(id string, callbacks ...Callbacker) mesh.Behavior {
	if len(callbacks) == 0 {
		logger.Errorf("callback behavior %q created without callback functions", id)
	}
	return &callbackBehavior{
		id:        id,
		callbacks: callbacks,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *callbackBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *callbackBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *callbackBehavior) Terminate() error {
	return nil
}

// ProcessEvent calls a callback functions with the event data.
func (b *callbackBehavior) Process(evt *event.Event) {
	for _, callback := range b.callbacks {
		fuse.Trigger(callback(b.emitter, evt))
	}
}

// Recover from an error.
func (b *callbackBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
