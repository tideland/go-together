// Tideland Go Together - Cells - Mesh - Unit Tests
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package mesh_test // import "tideland.dev/go/together/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestEmptyTestPlant verifies the instantiation of the test plant.
func TestEmptyTestPlant(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	plant := mesh.NewTestPlant(assert, NewTestBehavior("foo"), 0)

	plant.Emit(event.New("set", "a", 1))
}

// TestSubscribingTestPlant verifies the test plant with two subscribers.
// After setting two values the length is broadcasted and one value is
// emitted to subscriber "1"
func TestSubscribingTestPlant(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	plant := mesh.NewTestPlant(assert, NewTestBehavior("foo"), 2)

	plant.Emit(event.New("set", "a", 1))
	plant.Emit(event.New("set", "b", 2))
	plant.Emit(event.New("length"))
	plant.Emit(event.New("emit", "to", "1", "value", 12345))

	plant.AssertLength(0, 1)
	plant.AssertLength(1, 2)
	plant.AssertAll(0, func(evt *event.Event) bool {
		return evt.Topic() == "set" && evt.Payload().At("length").AsInt(-1) == 2
	})
	plant.AssertFind(1, func(evt *event.Event) bool {
		return evt.Topic() == "set" && evt.Payload().At("value").AsInt(-1) == 12345
	})
	plant.AssertFirst(1, func(evt *event.Event) bool {
		return evt.Payload().At("length").AsInt(-1) == 2
	})
	plant.AssertLast(1, func(evt *event.Event) bool {
		return evt.Payload().At("value").AsInt(-1) == 12345
	})
	plant.AssertNone(1, func(evt *event.Event) bool {
		return evt.Topic() == "not here"
	})

	plant.Reset()
	plant.AssertLength(0, 0)
	plant.AssertLength(1, 0)
}

// EOF
