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
	"strings"
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

// TestAggregatorBehavior tests the aggregator behavior. Scenario
// is simply to concatenate the random topics to the passed in topics in
// the payload at "topic".
func TestAggregatorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	count := 50
	aggregate := func(pl *event.Payload, evt *event.Event) (*event.Payload, error) {
		topic := evt.Topic()
		topics := pl.At("topics").AsString("")
		topics += "/" + topic
		return event.NewPayload("topics", topics), nil
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewAggregatorBehavior("ab", event.NewPayload(), aggregate), 1)
	defer plant.Stop()

	for i := 0; i < count; i++ {
		topic := generator.Word()
		plant.Emit(event.New(topic))
	}

	pl, plc := event.NewReplyPayload()
	plant.Emit(event.New(event.TopicStatus, pl))

	select {
	case pl := <-plc:
		topics := pl.At("topics").AsString("")
		splitted := strings.Split(topics, "/")
		assert.Length(splitted, count+1)
	case <-time.After(5 * time.Second):
		assert.Fail("timeout")
	}

	plant.AssertLength(0, count)
}

// EOF
