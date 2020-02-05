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

// TestRoundRobinBehavior tests the round robin behavior.
func TestRoundRobinBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	plant := mesh.NewTestPlant(assert, behaviors.NewRoundRobinBehavior("rrb"), 5)
	defer plant.Stop()

	for i := 0; i < 25; i++ {
		plant.Emit(event.New("round"))
	}

	for i := 0; i < 5; i++ {
		plant.AssertLength(i, 5)
	}
}

// EOF
