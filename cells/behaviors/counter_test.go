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
	counters := func(evt *event.Event) []string {
		return evt.Payload().Keys()
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewCounterBehavior("cb", counters), 1)
	defer plant.Stop()

	plant.Emit(event.New("count", "a", 1, "b", 1))
	plant.Emit(event.New("count", "a", 1, "c", 1, "d", 1))
	plant.Emit(event.New("count", "a", 1, "d", 1))

	plant.AssertLength(0, 3)
	plant.AssertLast(0, func(evt *event.Event) bool {
		return evt.Topic() == event.TopicCounted &&
			evt.Payload().At("a").AsInt(-1) == 3 &&
			evt.Payload().At("b").AsInt(-1) == 1 &&
			evt.Payload().At("c").AsInt(-1) == 1 &&
			evt.Payload().At("d").AsInt(-1) == 2
	})

	plant.Reset()
	plant.Emit(event.New(event.TopicReset))
	plant.Emit(event.New("count", "z", "don't care"))

	plant.AssertLast(0, func(evt *event.Event) bool {
		return evt.Topic() == event.TopicCounted &&
			evt.Payload().At("a").AsInt(-1) == -1 &&
			evt.Payload().At("z").AsInt(-1) == 1
	})
}

// EOF
