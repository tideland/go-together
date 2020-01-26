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
// MESH ROUTER BEHAVIOR
//--------------------

// meshRouterBehavior check for each received event which cell will
// get it based on the router function.
type meshRouterBehavior struct {
	id      string
	emitter mesh.Emitter
	routeTo Router
}

// NewMeshRouterBehavior creates a mesh router behavior using the passed function
// to determine to which cells the received event shall be re-emitted.
func NewMeshRouterBehavior(id string, router Router) mesh.Behavior {
	return &meshRouterBehavior{
		id:      id,
		routeTo: router,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *meshRouterBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *meshRouterBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *meshRouterBehavior) Terminate() error {
	return nil
}

// Process emits the event to those ids returned by the router function.
func (b *meshRouterBehavior) Process(evt *event.Event) {
	ids := b.routeTo(evt)
	for _, id := range ids {
		fuse.Trigger(b.emitter.Mesh().Emit(id, evt))
	}
}

// Recover from an error.
func (b *meshRouterBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
