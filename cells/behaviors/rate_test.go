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

// TestRateBehavior tests the event rate behavior.
func TestRateBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	generator := generators.New(generators.FixedRand())
	msh := mesh.New()
	defer msh.Stop()

	matcher := func(evt *event.Event) (bool, error) {
		return evt.Topic() == "now", nil
	}
	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		analyzer := event.NewSinkAnalyzer(accessor)
		ok, err := analyzer.Match(func(index int, evt *event.Event) (bool, error) {
			return evt.Topic() == behaviors.TopicRate, nil
		})
		sigc <- ok
		return nil, err
	}
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}

	msh.SpawnCells(
		behaviors.NewRateBehavior("rater", matcher, 100),
		behaviors.NewCollectorBehavior("collector", 10000, processor),
	)
	msh.Subscribe("rater", "collector")

	for i := 0; i < 1000; i++ {
		topic := generator.OneStringOf(topics...)
		msh.Emit("rater", event.New(topic))
		generator.SleepOneOf(0, time.Millisecond, 2*time.Millisecond)
	}

	msh.Emit("collector", event.New(event.TopicProcess))
	assert.Wait(sigc, true, 10*time.Second)
}

// EOF
