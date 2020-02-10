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

// TestTopicPayloadsBehavior tests the topic/payloads behavior.
func TestTopicPayloadsBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	processor := func(topic string, pls []*event.Payload) (*event.Payload, error) {
		total := 0
		for _, pl := range pls {
			value := pl.At("value").AsInt(0)
			total += value
		}
		return event.NewPayload("total", total), nil
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewTopicPayloadsBehavior("tpp", 5, processor), 1)
	defer plant.Stop()

	topics := []string{"alpha", "beta", "gamma"}
	values := []int{1, 2, 3, 4, 5}

	for i := 0; i < 50; i++ {
		topic := generator.OneStringOf(topics...)
		value := generator.OneIntOf(values...)
		plant.Emit(event.New(topic, "value", value))
	}

	plant.AssertAll(0, func(evt *event.Event) bool {
		topic := evt.Topic()
		topicOK := topic == "alpha" || topic == "beta" || topic == "gamma"
		total := evt.Payload().At("total").AsInt(-1)
		totalOK := total >= 1 && total <= 250
		return topicOK && totalOK
	})
}

// EOF
