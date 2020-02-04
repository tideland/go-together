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

// TestRateWindowBehavior tests the event rate window behavior.
func TestRateWindowBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "bang"}
	matcher := func(evt *event.Event) (bool, error) {
		// Signal when topic is "bang".
		return evt.Topic() == "bang", nil
	}
	duration := 50 * time.Millisecond
	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		// Got signals with matching rate, return info about it.
		first, _ := accessor.PeekFirst()
		last, _ := accessor.PeekLast()
		difference := last.Timestamp().Sub(first.Timestamp())
		return event.NewPayload("difference", difference), nil
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewRateWindowBehavior("rwb", matcher, 5, duration, processor), 1)
	defer plant.Stop()

	for i := 0; i < 250; i++ {
		topic := generator.OneStringOf(topics...)
		plant.Emit(event.New(topic))
		time.Sleep(time.Millisecond)
	}

	plant.AssertAll(0, func(evt *event.Event) bool {
		return evt.Topic() == behaviors.TopicRateWindow && evt.Payload().At("difference").AsDuration(2*duration) <= duration
	})

	plant.Reset()
	plant.Emit(event.New(event.TopicReset))

	for i := 0; i < 5; i++ {
		topic := generator.OneStringOf(topics...)
		plant.Emit(event.New(topic))
		time.Sleep(duration)
	}

	plant.AssertNone(0, func(evt *event.Event) bool {
		return evt.Topic() == behaviors.TopicRateWindow
	})
}

// EOF
