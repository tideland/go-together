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
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"
	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestAggregatorBehavior tests the aggregator behavior. Scenario
// is simply to concatenate the random topics to the passed in topics in
// the payload at "topic".
func TestAggregatorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	count := 50
	aggregator := func(data *mesh.Payload, evt *mesh.Event) (*mesh.Payload, error) {
		counted, err := data.IntAt("counted")
		assert.NoError(err)
		counted++
		return mesh.NewPayload("counted", counted), nil
	}
	pl := mesh.NewPayload("counted", 0)
	behavior := behaviors.NewAggregatorBehavior(pl, aggregator)
	tester := func(evt *mesh.Event) bool {
		switch evt.Topic() {
		case behaviors.TopicResetted:
			return true
		case behaviors.TopicAggregated:
			counted, err := evt.IntAt("counted")
			assert.NoError(err)
			assert.True(counted <= count)
		}
		return false
	}
	tb := mesh.NewTestbed(behavior, tester)

	for i := 0; i < count; i++ {
		topic := generator.Word()
		tb.Emit(mesh.NewEvent(topic))
	}
	tb.Emit(mesh.NewEvent(behaviors.TopicReset))

	err := tb.Wait(time.Second)
	assert.NoError(err)
}

// EOF
