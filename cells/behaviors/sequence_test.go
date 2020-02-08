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

// TestSequenceBehavior tests the event sequence behavior.
func TestSequenceBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	criterion := func(accessor event.SinkAccessor) event.CriterionMatch {
		if accessor.Len() < 3 {
			return event.CriterionKeep
		}
		if accessor.Len() > 3 {
			return event.CriterionClear
		}
		// Check if the three events are matching.
		a, _ := accessor.PeekAt(0)
		b, _ := accessor.PeekAt(1)
		now, _ := accessor.PeekAt(2)
		if a.Topic() == "a" && b.Topic() == "b" && now.Topic() == "now" {
			return event.CriterionDone
		}
		return event.CriterionClear
	}
	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		a, _ := accessor.PeekAt(0)
		b, _ := accessor.PeekAt(1)
		now, _ := accessor.PeekAt(2)
		pl := event.NewPayload("0", a.Topic(), "1", b.Topic(), "2", now.Topic())
		return pl, nil
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewSequenceBehavior("sb", criterion, processor), 1)
	defer plant.Stop()

	topics := []string{"a", "b", "c", "d", "now"}
	for i := 0; i < 10000; i++ {
		topic := generator.OneStringOf(topics...)
		plant.Emit(event.New(topic))
	}

	plant.AssertAll(0, func(evt *event.Event) bool {
		return evt.Payload().At("0").AsString("-") == "a" &&
			evt.Payload().At("1").AsString("-") == "b" &&
			evt.Payload().At("2").AsString("-") == "now"
	})
}

// EOF
