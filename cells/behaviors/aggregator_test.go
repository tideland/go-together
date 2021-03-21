// Tideland Go Together - Cells - Behaviors
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors_test // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"strconv"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestAggregatorBehavior tests the aggregator behavior. Aggreting will
// here simply be the counting of events, the mitting of this value will
// be tested.
func TestAggregatorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	count := 50
	aggregator := func(aggregate interface{}, evt *mesh.Event) (interface{}, error) {
		words := aggregate.(map[string]bool)
		words[evt.Topic()] = true
		return words, nil
	}
	behavior := behaviors.NewAggregatorBehavior(map[string]bool{}, aggregator)
	tester := func(evt *mesh.Event) bool {
		switch evt.Topic() {
		case behaviors.TopicResetted:
			return true
		case behaviors.TopicAggregated:
			var words map[string]bool
			err := evt.Payload(&words)
			assert.NoError(err)
			assert.Length(words, count)
		}
		return false
	}
	tb := mesh.NewTestbed(behavior, tester)

	// Run the test.
	for i := 0; i < count; i++ {
		topic := strconv.Itoa(i)
		tb.Emit(topic)
	}
	tb.Emit(behaviors.TopicAggregate)
	tb.Emit(behaviors.TopicReset)

	err := tb.Wait(time.Second)
	assert.NoError(err)
}

// EOF
