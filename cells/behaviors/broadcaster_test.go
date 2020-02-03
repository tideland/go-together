// Tideland Go Together - Cells - Behaviors - Unit Tests
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors_test // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestBroadcasterBehavior tests the broadcast behavior.
func TestBroadcasterBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	plant := mesh.NewTestPlant(assert, behaviors.NewBroadcasterBehavior("bb"), 5)
	defer plant.Stop()

	plant.Emit(event.New("a"))
	plant.Emit(event.New("b"))
	plant.Emit(event.New("c"))

	for idx := 0; idx < 5; idx++ {
		plant.AssertLength(idx, 3)
	}
}

// EOF
