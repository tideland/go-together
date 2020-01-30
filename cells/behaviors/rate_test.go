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

// TestRateBehavior tests the event rate behavior.
func TestRateBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}
	matcher := func(evt *event.Event) (bool, error) {
		return evt.Topic() == "now", nil
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewRateBehavior("rater", matcher, 100), 1)
	defer plant.Stop()

	for i := 0; i < 1000; i++ {
		topic := generator.OneStringOf(topics...)
		plant.Emit(event.New(topic))
	}
	plant.AssertAll("sub-0", func(evt *event.Event) bool {
		return evt.Topic() == behaviors.TopicRate
	})
}

// EOF
