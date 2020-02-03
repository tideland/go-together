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

// TestCallbackBehavior tests the callback behavior.
func TestCallbackBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	cbfA := func(emitter mesh.Emitter, evt *event.Event) error {
		return emitter.Emit("0", evt)
	}
	cbfB := func(emitter mesh.Emitter, evt *event.Event) error {
		return emitter.Emit("1", evt)
	}
	cbfC := func(emitter mesh.Emitter, evt *event.Event) error {
		if err := emitter.Emit("0", evt); err != nil {
			return err
		}
		return emitter.Emit("1", evt)
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewCallbackBehavior("cb", cbfA, cbfB, cbfC), 2)
	defer plant.Stop()

	plant.Emit(event.New("foo"))
	plant.Emit(event.New("bar"))
	plant.Emit(event.New("baz"))

	plant.AssertLength(0, 6)
	plant.AssertLength(1, 6)
}

// EOF
