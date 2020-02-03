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

// TestCollectorBehavior tests the collector behavior.
func TestCollectorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		return event.NewPayload("length", accessor.Len()), nil
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewCollectorBehavior("cb", 10, processor), 1)
	defer plant.Stop()

	// Don't care for words, we collect maximally 10 events.
	for _, word := range generator.Words(25) {
		plant.Emit(event.New("collect", word))
	}

	plant.Emit(event.New(event.TopicProcess))
	plant.AssertLength(0, 26)
	plant.AssertFind(0, func(evt *event.Event) bool {
		return evt.Topic() == event.TopicResult && evt.Payload().At("length").AsInt(-1) == 10
	})

	plant.Emit(event.New(event.TopicReset))
	plant.Emit(event.New(event.TopicProcess))
	plant.AssertFind(0, func(evt *event.Event) bool {
		return evt.Topic() == event.TopicResult && evt.Payload().At("length").AsInt(-1) == 0
	})
}

// EOF
