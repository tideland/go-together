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
// MAPPER BEHAVIOR
//--------------------

// Mapper is a function type mapping an event to another one.
type Mapper func(evt *event.Event) (*event.Event, error)

// mapperBehavior maps the received event to a new event.
type mapperBehavior struct {
	id      string
	emitter mesh.Emitter
	mapper  Mapper
}

// NewMapperBehavior creates a map behavior based on the passed function.
// It emits the mapped events.
func NewMapperBehavior(id string, mapper Mapper) mesh.Behavior {
	return &mapperBehavior{
		id:     id,
		mapper: mapper,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *mapperBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *mapperBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *mapperBehavior) Terminate() error {
	return nil
}

// Process maps the received event to a new one and emits it.
func (b *mapperBehavior) Process(evt *event.Event) {
	mapped, err := b.mapper(evt)
	fuse.Trigger(err)
	if mapped != nil {
		fuse.Trigger(b.emitter.Broadcast(mapped))
	}
}

// Recover from an error.
func (b *mapperBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
