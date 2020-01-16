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

// TestTopicPayloadsBehavior tests the topic/payloads behavior.
func TestTopicPayloadsBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()
	defer msh.Stop()

	tpProcessor := func(topic string, pls []*event.Payload) (*event.Payload, error) {
		total := 0
		for _, pl := range pls {
			value := pl.At("value").AsInt(0)
			total += value
		}
		return event.NewPayload("total", total), nil
	}
	cProcessor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		err := accessor.Do(func(index int, evt *event.Event) error {
			total := evt.Payload().At("total").AsInt(0)
			assert.Range(total, 1, 25)
			return nil
		})
		sigc <- true
		return nil, err
	}

	msh.SpawnCells(
		behaviors.NewTopicPayloadsBehavior("topic-payloads", 5, tpProcessor),
		behaviors.NewCollectorBehavior("collector", 10, cProcessor),
	)
	msh.Subscribe("topic-payloads", "collector")

	topics := []string{"alpha", "beta", "gamma"}
	values := []int{1, 2, 3, 4, 5}

	for i := 0; i < 50; i++ {
		topic := generator.OneStringOf(topics...)
		value := generator.OneIntOf(values...)
		msh.Emit("topic-payloads", event.New(topic, "value", value))
	}

	msh.Emit("collector", event.New(event.TopicProcess))
	assert.Wait(sigc, true, 5*time.Second)
}

// EOF
