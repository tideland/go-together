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
	"time"

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
	sigc := asserts.MakeMultiWaitChan(size)
	msh := mesh.New()
	defer assert.NoError(msh.Stop())

	tester := func(evt *event.Event) bool {
		return evt.Topic() == "end"
	}
	processor := func(emitter mesh.Emitter, evt *event.Event) error {
		sigc <- evt.Topic()
		return nil
	}

	assert.NoError(msh.SpawnCells(behaviors.NewConditionBehavior("condition", tester, processor)))

	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "end"}

	for i := 0; i < size; i++ {
		topic := generator.OneStringOf(topics...)
		assert.NoError(msh.Emit("condition", event.New(topic)))
	}

	assert.Wait(sigc, "end", time.Second)
}

// EOF
