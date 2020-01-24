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

// TestSequenceBehavior tests the event sequence behavior.
func TestSequenceBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	generator := generators.New(generators.FixedRand())
	msh := mesh.New()

	sequence := []string{"a", "b", "now"}
	sequencer := func(accessor event.SinkAccessor) event.CriterionMatch {
		analyzer := event.NewSinkAnalyzer(accessor)
		matcher := func(index int, evt *event.Event) (bool, error) {
			ok := evt.Topic() == sequence[index]
			return ok, nil
		}
		matches, err := analyzer.Match(matcher)
		if err != nil || !matches {
			return event.CriterionClear
		}
		if accessor.Len() == len(sequence) {
			return event.CriterionDone
		}
		return event.CriterionKeep
	}
	analyzer := func(accessor event.SinkAccessor) (*event.Payload, error) {
		first, ok := accessor.PeekFirst()
		assert.OK(ok)
		return first.Payload(), nil
	}
	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		var indexes []int
		err := accessor.Do(func(_ int, evt *event.Event) error {
			index := evt.Payload().At("index").AsInt(-1)
			indexes = append(indexes, index)
			return nil
		})
		assert.Nil(err)
		sigc <- indexes
		return nil, nil
	}
	topics := []string{"a", "b", "c", "d", "now"}

	assert.OK(msh.SpawnCells(
		behaviors.NewSequenceBehavior("sequencer", sequencer, analyzer),
		behaviors.NewCollectorBehavior("collector", 100, processor),
	))
	assert.OK(msh.Subscribe("sequencer", "collector"))

	for i := 0; i < 1000; i++ {
		topic := generator.OneStringOf(topics...)
		assert.OK(msh.Emit("sequencer", event.New(topic, "index", i)))
		generator.SleepOneOf(0, 1*time.Millisecond, 2*time.Millisecond)
	}

	assert.OK(msh.Emit("collector", event.New(event.TopicProcess)))
	assert.Wait(sigc, []int{155, 269, 287, 298, 523, 888}, time.Minute)
	assert.OK(msh.Stop())
}

// EOF
