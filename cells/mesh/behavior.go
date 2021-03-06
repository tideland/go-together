// Tideland Go Together - Cells - Mesh
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh

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
	Emitter(name string) (*Emitter, error)
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
	Go(cell Cell, in InputStream, out OutputStream)
}

//--------------------
// BEHAVIORS
//--------------------

// StatelessFunc defines a function signature for the stateless
// behavior. This function processes an event by being called.
type StatelessFunc func(evt *Event, out OutputStream) error

// StatelessBehavior is a simple behavior using a function
// to process the received events.
type StatelessBehavior struct {
	sf StatelessFunc
}

// NewStatelessBehavior creates a behavior based on the given
// processing function.
func NewStatelessBehavior(sf StatelessFunc) StatelessBehavior {
	return StatelessBehavior{
		sf: sf,
	}
}

// Go implements Behavior.
func (sb StatelessBehavior) Go(cell Cell, in InputStream, out OutputStream) {
	for {
		select {
		case <-cell.Context().Done():
			return
		case evt := <-in.Pull():
			if err := sb.sf(evt, out); err != nil {
				out.Emit(NewEvent(ErrorTopic, NameKey, cell.Name(), MessageKey, err.Error()))
			}
		}
	}
}

// EOF
