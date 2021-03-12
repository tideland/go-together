// Tideland Go Together - Cells - Mesh
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/together/cells/mesh"

//--------------------
// IMPORT
//--------------------

import (
	"context"
)

//--------------------
// MESH
//--------------------

// Mesh describes the interface to a mesh of a cell from the
// perspective of a behavior.
type Mesh interface {
	// Go starts a cell using the given behavior.
	Go(name string, b Behavior) error

	// Subscribe subscribes the cell with from name to the cell
	// with to name.
	Subscribe(fromName, toName string) error

	// Unsubscribe unsubscribes the cell with to name from the cell
	// with from name.
	Unsubscribe(toName, fromName string) error

	// Emit raises an event to the named cell.
	Emit(name string, evt *Event) error

	// Emitter returns a static emitter for the named cell.
	Emitter(name string) (Emitter, error)
}

//--------------------
// CELL
//--------------------

// Cell describes the interface to a cell from the perspective
// of a behavior.
type Cell interface {
	// Context returns the context of mesh and cell.
	Context() context.Context

	// Name returns the name of the deployed cell running the
	// behavior.
	Name() string

	// Mesh returns the mesh of the cell.
	Mesh() Mesh
}

//--------------------
// BEHAVIOR
//--------------------

// Behavior describes what cell implementations must understand.
type Behavior interface {
	// Go will be started as wrapped goroutine. It's the responsible
	// of the implementation to run a select loop, receive incomming
	// events via the input queue, and emit events via the output queue
	// if needed.
	//
	//     for {
	//         select {
	//         case <-cell.Context().Done():
	//             return nil
	//         case evt := <-in.Pull():
	//             ...
	//             out.Emit(mesh.NewVent("my-topic", "x", 12345))
	//         }
	//     }
	//
	Go(cell Cell, in Receptor, out Emitter) error
}

//--------------------
// BEHAVIORS
//--------------------

// BehaviorFunc simplifies implementation of a behavior when only
// one function is needed. It can be deployed via
//
//     myMesh.Go("my-name", BehaviorFunc(myFunc))
type BehaviorFunc func(cell Cell, in Receptor, out Emitter) error

// Go implements Behavior.
func (bf BehaviorFunc) Go(cell Cell, in Receptor, out Emitter) error {
	return bf(cell, in, out)
}

// RequestFunc defines a function signature for the request
// behavior. It is called per received event.
type RequestFunc func(cell Cell, evt *Event, out Emitter) error

// RequestBehavior is a simple behavior using a function
// to process the received events.
type RequestBehavior struct {
	rf RequestFunc
}

// NewRequestBehavior creates a behavior based on the given
// processing function.
func NewRequestBehavior(rf RequestFunc) RequestBehavior {
	return RequestBehavior{
		rf: rf,
	}
}

// Go implements Behavior.
func (rb RequestBehavior) Go(cell Cell, in Receptor, out Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			if err := rb.rf(cell, evt, out); err != nil {
				return err
			}
		}
	}
}

// EOF
