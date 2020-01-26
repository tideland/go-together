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

// TestPairBehavior tests the event pair behavior.
func TestPairBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	generator := generators.New(generators.FixedRand())
	msh := mesh.New()

	matchCount := make(map[string]int)
	names := generator.Names(5000)
	waitForName := generator.OneStringOf(names...)
	matches := func(evt *event.Event, pl *event.Payload) (*event.Payload, bool) {
		// Wait for a name visiting twice during timout.
		name := evt.Payload().At("name").AsString("<none>")
		if name == waitForName {
			// Hit!
			return event.NewPayload("name", name), true
		}
		return pl, false
	}
	processor := func(emitter mesh.Emitter, evt *event.Event) error {
		matchCount[evt.Topic()]++
		if len(matchCount) == 2 {
			sigc <- true
		}
		return nil
	}
	timespan := 10 * time.Millisecond

	assert.OK(msh.SpawnCells(
		behaviors.NewPairBehavior("pairer", matches, timespan),
		behaviors.NewSimpleProcessorBehavior("processor", processor),
	))
	assert.OK(msh.Subscribe("pairer", "processor"))

	go func() {
		for i := 0; i < 100000; i++ {
			name := generator.OneStringOf(names...)
			_ = msh.Emit("pairer", event.New("visitor", "name", name))
		}
	}()

	assert.Wait(sigc, true, 30*time.Second)
	assert.OK(matchCount[behaviors.TopicPair] > 0)
	assert.OK(matchCount[behaviors.TopicPairTimeout] > 0)
	assert.OK(msh.Stop())
}

// EOF
