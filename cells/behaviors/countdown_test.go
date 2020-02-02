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

// TestCountdownBehavior tests the countdown of events.
func TestCountdownBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	zeroer := func(accessor event.SinkAccessor) (*event.Event, int, error) {
		at := accessor.Len()
		evt := event.New("zero", "after", at)
		return evt, at - 1, nil
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewCountdownBehavior("cb", 5, zeroer), 1)
	defer plant.Stop()

	countdown := func(ct int) {
		for i := 0; i < ct; i++ {
			plant.Emit(event.New("count"))
		}
		plant.AssertFind("sub-0", func(evt *event.Event) bool {
			return evt.Topic() == "zero" && evt.Payload().At("after").AsInt(-1) == ct
		})
	}

	countdown(5)
	countdown(4)
	countdown(3)
	countdown(2)
	countdown(1)
}

// EOF
