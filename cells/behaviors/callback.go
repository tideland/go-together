// Tideland Go Together - Cells - Behaviors
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// CALLBACK BEHAVIOR
//--------------------

// CallbackFunc is a function called by the behavior when it receives an event.
type CallbackFunc func(evt mesh.Event, out mesh.Emitter) error

// callbackBehavior implements the aggregator behavior.
type callbackBehavior struct {
	callbacks []CallbackFunc
}

// NewCallbackBehavior creates a behavior aggregating the received events
// and emits events with the new aggregate. A "reset!" topic resets the
// aggregate to nil again.
func NewCallbackBehavior(callbacks ...CallbackFunc) mesh.Behavior {
	return &callbackBehavior{
		callbacks: callbacks,
	}
}

// Go calls the callbacks for any event.
func (b *callbackBehavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			for _, callback := range b.callbacks {
				if err := callback(evt, out); err != nil {
					return err
				}
			}
		}
	}
}

// EOF
