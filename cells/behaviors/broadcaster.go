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
// BROARDCASTER BEHAVIOR
//--------------------

// broadcasterBehavior implements the broadcaster behavior.
type broadcasterBehavior struct{}

// NewBroadcasterBeehavior creates a behavior simply emitting pulled events.
func NewBroadcasterBehavior() mesh.Behavior {
	return &broadcasterBehavior{}
}

// Go aggregates the event.
func (b *broadcasterBehavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			if err := out.EmitEvent(evt); err != nil {
				return err
			}
		}
	}
}

// EOF
