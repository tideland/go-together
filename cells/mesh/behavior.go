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

// EventProcessingFunc defines a function signature for the simple
// event processing behavior. This function processes every event.
type EventProcessingFunc func(evt *Event, out OutputStream)

// EventProcessingBehavior is a simple behavior using a function
// to process the received events.
type EventProcessingBehavior struct {
	epf EventProcessingFunc
}

// NewEventProcessingBehavior creates a behavior based on the given
// processing function.
func NewEventProcessingBehavior(epf EventProcessingFunc) EventProcessingBehavior {
	return EventProcessingBehavior{
		epf: epf,
	}
}

// Go implements Behavior.
func (epb EventProcessingBehavior) Go(ctx context.Context, name string, in InputStream, out OutputStream) {
	for {
		select {
		case <-ctx.Done():
			return
		case evt := <-in.Pull():
			epb.epf(evt, out)
		}
	}
}

// EOF
