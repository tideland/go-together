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

// TestPairBehavior tests the event pair behavior finding a pair.
func TestPairBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	names := generator.Names(5000)
	matches := func(first, second *event.Event) bool {
		if first == nil {
			// First one.
			return true
		}
		s1 := first.Payload().At("name").AsString("<first>")
		s2 := second.Payload().At("name").AsString("<second>")
		return s1 == s2
	}
	timespan := 5 * time.Millisecond
	plant := mesh.NewTestPlant(assert, behaviors.NewPairBehavior("pb", matches, timespan), 1)
	defer plant.Stop()

	for i := 0; i < 100000; i++ {
		name := generator.OneStringOf(names...)
		plant.Emit(event.New("visitor", "name", name))
	}

	plant.AssertFind(0, func(evt *event.Event) bool {
		return evt.Topic() == behaviors.TopicPair
	})
	plant.AssertFind(0, func(evt *event.Event) bool {
		return evt.Topic() == behaviors.TopicPairTimeout
	})
}

// EOF
