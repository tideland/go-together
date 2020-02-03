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
	"tideland.dev/go/audit/generators"
	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestConditionBehavior tests the condition behavior.
func TestConditionBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	size := 1000
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "end"}
	tester := func(evt *event.Event) bool {
		return evt.Topic() == "end"
	}
	processor := func(emitter mesh.Emitter, evt *event.Event) error {
		return emitter.Broadcast(evt)
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewConditionBehavior("cb", tester, processor), 1)
	defer plant.Stop()

	for i := 0; i < size; i++ {
		topic := generator.OneStringOf(topics...)
		plant.Emit(event.New(topic))
	}

	plant.AssertAll(0, func(evt *event.Event) bool {
		return evt.Topic() == "end"
	})
}

// EOF
