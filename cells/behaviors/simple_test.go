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

// TestSimpleProcessorBehavior tests the simple processor behavior.
func TestSimpleProcessorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	spf := func(emitter mesh.Emitter, evt *event.Event) error {
		return emitter.Broadcast(evt)
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewSimpleProcessorBehavior("simple", spf), 1)
	defer plant.Stop()

	plant.Emit(event.New("foo"))
	plant.Emit(event.New("bar"))
	plant.Emit(event.New("baz"))

	plant.AssertLength("sub-0", 3)
}

// EOF
