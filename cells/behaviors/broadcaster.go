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
// BROADCASTER BEHAVIOR
//--------------------

// broadcasterBehavior is a simple repeater.
type broadcasterBehavior struct {
	id      string
	emitter mesh.Emitter
}

// NewBroadcasterBehavior creates a broadcasting behavior that just emits every
// received event. It's intended to work as an entry point for events, which
// shall be immediately processed by several subscribers.
func NewBroadcasterBehavior(id string) mesh.Behavior {
	return &broadcasterBehavior{
		id: id,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *broadcasterBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *broadcasterBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *broadcasterBehavior) Terminate() error {
	return nil
}

// Process emits the event to all subscribers.
func (b *broadcasterBehavior) Process(evt *event.Event) error {
	return b.emitter.Broadcast(evt)
}

// Recover from an error.
func (b *broadcasterBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
