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
	"strconv"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestEvaluatorBehavior tests the evaluator behavior.
func TestEvaluatorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()
	defer assert.NoError(msh.Stop())

	evaluator := func(evt *event.Event) (float64, error) {
		i, err := strconv.Atoi(evt.Topic())
		return float64(i), err
	}
	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		evt, ok := accessor.PeekLast()
		assert.True(ok)
		sigc <- evt
		return nil, nil
	}

	assert.NoError(msh.SpawnCells(
		behaviors.NewEvaluatorBehavior("evaluator", evaluator),
		behaviors.NewCollectorBehavior("collector", 1000, processor),
	))
	assert.NoError(msh.Subscribe("evaluator", "collector"))

	// Standard evaluating.
	topics := []string{"2", "1", "1", "1", "3", "2", "3", "1", "3", "9"}
	for _, topic := range topics {
		assert.NoError(msh.Emit("evaluator", event.New(topic)))
	}
	time.Sleep(100 * time.Millisecond)

	assert.NoError(msh.Emit("collector", event.New(event.TopicProcess)))
	assert.NoError(msh.Emit("collector", event.New(event.TopicReset)))

	assert.WaitTested(sigc, func(value interface{}) error {
		evt, ok := value.(*event.Event)
		assert.True(ok)
		assert.Equal(evt.Payload().At("count").AsInt(0), 10)
		assert.Equal(evt.Payload().At("min-rating").AsFloat64(0.0), 1.0)
		assert.Equal(evt.Payload().At("max-rating").AsFloat64(0.0), 9.0)
		assert.Equal(evt.Payload().At("avg-rating").AsFloat64(0.0), 2.6)
		assert.Equal(evt.Payload().At("med-rating").AsFloat64(0.0), 2.0)
		return nil
	}, time.Second)

	// Reset and check with only one value.
	assert.NoError(msh.Emit("evaluator", event.New(event.TopicReset)))
	assert.NoError(msh.Emit("evaluator", event.New("1234")))
	time.Sleep(100 * time.Millisecond)

	assert.NoError(msh.Emit("collector", event.New(event.TopicProcess)))
	assert.NoError(msh.Emit("collector", event.New(event.TopicReset)))

	assert.WaitTested(sigc, func(value interface{}) error {
		evt, ok := value.(*event.Event)
		assert.True(ok)
		assert.Equal(evt.Payload().At("count").AsInt(0), 1)
		assert.Equal(evt.Payload().At("min-rating").AsFloat64(0.0), 1234.0)
		assert.Equal(evt.Payload().At("max-rating").AsFloat64(0.0), 1234.0)
		assert.Equal(evt.Payload().At("avg-rating").AsFloat64(0.0), 1234.0)
		assert.Equal(evt.Payload().At("med-rating").AsFloat64(0.0), 1234.0)
		return nil
	}, time.Second)

	// Crash evaluating.
	topics = []string{"2", "1", "3", "4", "crash", "1", "2", "1", "2", "1"}
	for _, topic := range topics {
		assert.NoError(msh.Emit("evaluator", event.New(topic)))
	}
	time.Sleep(100 * time.Millisecond)

	assert.NoError(msh.Emit("collector", event.New(event.TopicProcess)))
	assert.NoError(msh.Emit("collector", event.New(event.TopicReset)))

	assert.WaitTested(sigc, func(value interface{}) error {
		evt, ok := value.(*event.Event)
		assert.True(ok)
		assert.Equal(evt.Payload().At("count").AsInt(0), 5)
		assert.Equal(evt.Payload().At("min-rating").AsFloat64(0.0), 1.0)
		assert.Equal(evt.Payload().At("max-rating").AsFloat64(0.0), 2.0)
		assert.Equal(evt.Payload().At("avg-rating").AsFloat64(0.0), 1.4)
		assert.Equal(evt.Payload().At("med-rating").AsFloat64(0.0), 1.0)
		return nil
	}, time.Second)
}

// TestLimitedEvaluatorBehavior tests the limited evaluator behavior.
func TestLimitedEvaluatorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()
	defer assert.NoError(msh.Stop())

	evaluator := func(evt *event.Event) (float64, error) {
		i, err := strconv.Atoi(evt.Topic())
		return float64(i), err
	}
	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		evt, ok := accessor.PeekLast()
		assert.True(ok)
		sigc <- evt
		return nil, nil
	}

	assert.NoError(msh.SpawnCells(
		behaviors.NewMovingEvaluatorBehavior("evaluator", evaluator, 5),
		behaviors.NewCollectorBehavior("collector", 1000, processor),
	))
	assert.NoError(msh.Subscribe("evaluator", "collector"))

	// Standard evaluating.
	topics := []string{"1", "2", "1", "1", "9", "2", "3", "1", "3", "2"}
	for _, topic := range topics {
		assert.NoError(msh.Emit("evaluator", event.New(topic)))
	}
	time.Sleep(100 * time.Millisecond)

	assert.NoError(msh.Emit("collector", event.New(event.TopicProcess)))
	assert.NoError(msh.Emit("collector", event.New(event.TopicReset)))

	assert.WaitTested(sigc, func(value interface{}) error {
		evt, ok := value.(*event.Event)
		assert.True(ok)
		assert.Equal(evt.Payload().At("count").AsInt(0), 5)
		assert.Equal(evt.Payload().At("min-rating").AsFloat64(0.0), 1.0)
		assert.Equal(evt.Payload().At("max-rating").AsFloat64(0.0), 3.0)
		assert.Equal(evt.Payload().At("avg-rating").AsFloat64(0.0), 2.2)
		assert.Equal(evt.Payload().At("med-rating").AsFloat64(0.0), 2.0)
		return nil
	}, time.Second)
}

// EOF
