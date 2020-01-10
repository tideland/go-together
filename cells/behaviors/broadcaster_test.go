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

// TestBroadcasterBehavior tests the broadcast behavior.
func TestBroadcasterBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()
	defer msh.Stop()

	mktester := func() behaviors.ConditionTester {
		counter := 0
		return func(evt *event.Event) bool {
			counter++
			return counter == 3
		}
	}
	processor := func(emitter mesh.Emitter, evt *event.Event) error {
		sigc <- true
		return nil
	}

	msh.SpawnCells(
		behaviors.NewBroadcasterBehavior("broadcast"),
		behaviors.NewConditionBehavior("test-a", mktester(), processor),
		behaviors.NewConditionBehavior("test-b", mktester(), processor),
	)
	msh.Subscribe("broadcast", "test-a", "test-b")

	msh.Emit("broadcast", event.New("test"))
	msh.Emit("broadcast", event.New("test"))
	msh.Emit("broadcast", event.New("test"))

	assert.Wait(sigc, true, time.Second)
	assert.Wait(sigc, true, time.Second)
}

// EOF
