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
)

//--------------------
// ONCE BEHAVIOR
//--------------------

// OneTimer describes the function called after the first event.
type OneTimer func(emitter mesh.Emitter, evt *event.Event) error

// onceBehavior implements the once behavior.
type onceBehavior struct {
	id       string
	emitter  mesh.Emitter
	oneTimer OneTimer
}

// NewOnceBehavior creates a behavior where the cell calls the one-timer
// function for the first received event. Afterwards it will never be called
// again.
func NewOnceBehavior(id string, oneTimer OneTimer) mesh.Behavior {
	return &onceBehavior{
		id:       id,
		oneTimer: oneTimer,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *onceBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *onceBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *onceBehavior) Terminate() error {
	return nil
}

// Process calls the one-timer, but only for the first received event.
func (b *onceBehavior) Process(evt *event.Event) {
	if b.oneTimer != nil {
		err := b.oneTimer(b.emitter, evt)
		fuse.Trigger(err)
		b.oneTimer = nil
	}
}

// Recover from an error.
func (b *onceBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
