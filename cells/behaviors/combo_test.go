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

// TestComboBehavior tests the event combo behavior. The combination is waiting
// for at least one of the topics "a", "b", "c", and "d" in
func TestComboBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()
	defer msh.Stop()

	matcher := func(accessor event.SinkAccessor) (event.CriterionMatch, *event.Payload) {
		analyzer := event.NewSinkAnalyzer(accessor)
		combo := map[string]int{
			"a": 0,
			"b": 0,
			"c": 0,
			"d": 0,
		}
		matches, err := analyzer.Match(func(index int, evt *event.Event) (bool, error) {
			_, ok := combo[evt.Topic()]
			if ok {
				combo[evt.Topic()]++
			}
			return ok, nil
		})
		// No match or error.
		if err != nil || !matches {
			return event.CriterionDropLast, nil
		}
		// As long as at least one is still zero continue.
		for _, count := range combo {
			if count == 0 {
				return event.CriterionKeep, nil
			}
		}
		// It's done.
		pl := event.NewPayload(
			"a", combo["a"],
			"b", combo["b"],
			"c", combo["c"],
			"d", combo["d"],
		)
		return event.CriterionDone, pl
	}
	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		analyzer := event.NewSinkAnalyzer(accessor)
		ok, err := analyzer.Match(func(index int, evt *event.Event) (bool, error) {
			pl := evt.Payload()
			if len(pl.Keys()) != 4 {
				return false, nil
			}
			for _, key := range pl.Keys() {
				count := pl.At(key).AsInt(0)
				if count == 0 {
					return false, nil
				}
			}
			return true, nil
		})
		sigc <- ok
		return nil, err
	}

	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}

	msh.SpawnCells(
		behaviors.NewComboBehavior("combiner", matcher),
		behaviors.NewCollectorBehavior("collector", 100, processor),
	)
	msh.Subscribe("combiner", "collector")

	for i := 0; i < 1000; i++ {
		topic := generator.OneStringOf(topics...)
		msh.Emit("combiner", event.New(topic))
	}

	msh.Emit("collector", event.New(event.TopicProcess))
	assert.Wait(sigc, true, time.Minute)
}

// EOF
