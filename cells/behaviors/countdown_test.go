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
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()
	defer assert.NoError(msh.Stop())

	zeroer := func(accessor event.SinkAccessor) (*event.Event, int, error) {
		at := accessor.Len()
		evt := event.New("zero", at)
		return evt, at - 1, nil
	}
	tester := func(evt *event.Event) bool {
		return evt.Topic() == "zero"
	}
	processor := func(emitter mesh.Emitter, evt *event.Event) error {
		sigc <- evt.Topic()
		return nil
	}

	assert.NoError(msh.SpawnCells(
		behaviors.NewCountdownBehavior("countdowner", 5, zeroer),
		behaviors.NewConditionBehavior("conditioner", tester, processor),
	))
	assert.NoError(msh.Subscribe("countdowner", "conditioner"))

	countdown := func(ct int) {
		for i := 0; i < ct; i++ {
			err := msh.Emit("countdowner", event.New("count"))
			assert.Nil(err)
		}
		assert.Wait(sigc, "zero", time.Second)
	}

	countdown(5)
	countdown(4)
	countdown(3)
	countdown(2)
	countdown(1)
}

// EOF
