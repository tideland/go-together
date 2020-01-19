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
// is simply to count the lengths of the random topic until it
// reached the value 100.
func TestAggregatorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	msh := mesh.New()
	defer assert.NoError(msh.Stop())

	aggregate := func(pl *event.Payload, evt *event.Event) (*event.Payload, error) {
		topic := evt.Topic()
		topics := pl.At("topics").AsString("")
		topics += "/" + topic
		return event.NewPayload("topics", topics), nil
	}

	assert.NoError(msh.SpawnCells(behaviors.NewAggregatorBehavior("aggregator", event.NewPayload(), aggregate)))

	for i := 0; i < 50; i++ {
		topic := generator.Word()
		assert.NoError(msh.Emit("aggregator", event.New(topic)))
	}

	pl, plc := event.NewReplyPayload()
	evt := event.New(event.TopicStatus, pl)

	assert.NoError(msh.Emit("aggregator", evt))

	select {
	case pl := <-plc:
		topics := pl.At("topics").AsString("")
		splitted := strings.Split(topics, "/")
		assert.Length(splitted, 51)
	case <-time.After(5 * time.Second):
		assert.Fail("timeout")
	}
}

// EOF
