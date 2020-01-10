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

// TestCounterBehavior tests the counting of events.
func TestCounterBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()
	defer msh.Stop()

	counters := func(evt *event.Event) []string {
		return evt.Payload().Keys()
	}
	tester := func(evt *event.Event) bool {
		a := evt.Payload().At("a").AsInt(0)
		b := evt.Payload().At("b").AsInt(0)
		c := evt.Payload().At("c").AsInt(0)
		d := evt.Payload().At("d").AsInt(0)
		return a == 3 && b == 1 && c == 1 && d == 2
	}
	processor := func(emitter mesh.Emitter, evt *event.Event) error {
		sigc <- true
		return nil
	}

	msh.SpawnCells(
		behaviors.NewCounterBehavior("counter", counters),
		behaviors.NewConditionBehavior("conditioner", tester, processor),
	)
	msh.Subscribe("counter", "conditioner")

	msh.Emit("counter", event.New("count", "a", 1, "b", 1))
	msh.Emit("counter", event.New("count", "a", 1, "c", 1, "d", 1))
	msh.Emit("counter", event.New("count", "a", 1, "d", 1))

	assert.Wait(sigc, true, time.Second)
}

// EOF
