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
// BEHAVIOR
//--------------------

// Behavior describes what cell implementations must understand.
type Behavior interface {
	// Go will be started as wrapped goroutine. It's the responsible
	// of the implementation to run a select loop, receive incomming
	// events via the input queue, and emit events via the output queue
	// if needed.
	Go(ctx context.Context, name string, in InputStream, out OutputStream)
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
func (sb StatelessBehavior) Go(ctx context.Context, name string, in InputStream, out OutputStream) {
	for {
		select {
		case <-ctx.Done():
			return
		case evt := <-in.Pull():
			if err := sb.sf(evt, out); err != nil {
				out.Emit(NewEvent(ErrorTopic, NameKey, name, MessageKey, err.Error()))
			}
		}
	}
}

// EOF
