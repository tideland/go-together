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

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestEvaluatorBehavior verifies the evaluating of events.
func TestEvaluatorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evaluator := func(evt *event.Event) (float64, error) {
		f, err := strconv.ParseFloat(evt.Topic(), 64)
		return f, err
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewEvaluatorBehavior("eb", evaluator), 1)
	defer plant.Stop()

	// Standard evaluating.
	topics := []string{"2", "1", "1", "1", "3", "2", "3", "1", "3", "9"}
	for _, topic := range topics {
		plant.Emit(event.New(topic))
	}
	plant.AssertFind(0, func(evt *event.Event) bool {
		return evt.Topic() == behaviors.TopicEvaluation &&
			evt.Payload().At("count").AsInt(0) == 10 &&
			evt.Payload().At("min-rating").AsFloat64(0.0) == 1.0 &&
			evt.Payload().At("max-rating").AsFloat64(0.0) == 9.0 &&
			evt.Payload().At("avg-rating").AsFloat64(0.0) == 2.6 &&
			evt.Payload().At("med-rating").AsFloat64(0.0) == 2.0
	})

	// Reset and check with only one value.
	plant.Emit(event.New(event.TopicReset))
	plant.Emit(event.New("1234"))

	plant.AssertFind(0, func(evt *event.Event) bool {
		return evt.Topic() == behaviors.TopicEvaluation &&
			evt.Payload().At("count").AsInt(0) == 1 &&
			evt.Payload().At("min-rating").AsFloat64(0.0) == 1234.0 &&
			evt.Payload().At("max-rating").AsFloat64(0.0) == 1234.0 &&
			evt.Payload().At("avg-rating").AsFloat64(0.0) == 1234.0 &&
			evt.Payload().At("med-rating").AsFloat64(0.0) == 1234.0
	})

	// Crash evaluating.
	plant.Emit(event.New(event.TopicReset))
	topics = []string{"2", "1", "3", "4", "crash", "1", "2", "1", "2", "1"}
	for _, topic := range topics {
		plant.Emit(event.New(topic))
	}

	plant.AssertFind(0, func(evt *event.Event) bool {
		return evt.Topic() == behaviors.TopicEvaluation &&
			evt.Payload().At("count").AsInt(0) == 5 &&
			evt.Payload().At("min-rating").AsFloat64(0.0) == 1.0 &&
			evt.Payload().At("max-rating").AsFloat64(0.0) == 2.0 &&
			evt.Payload().At("avg-rating").AsFloat64(0.0) == 1.4 &&
			evt.Payload().At("med-rating").AsFloat64(0.0) == 1.0
	})
}

// TestMovingEvaluatorBehavior tests the evaluator behavior with a
// moving number of values.
func TestMovingEvaluatorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evaluator := func(evt *event.Event) (float64, error) {
		f, err := strconv.ParseFloat(evt.Topic(), 64)
		return f, err
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewMovingEvaluatorBehavior("evaluator", evaluator, 5), 1)
	defer plant.Stop()

	// Standard evaluating.
	topics := []string{"1", "2", "1", "1", "9", "2", "3", "1", "3", "2"}
	for _, topic := range topics {
		plant.Emit(event.New(topic))
	}

	plant.AssertFind(0, func(evt *event.Event) bool {
		return evt.Topic() == behaviors.TopicEvaluation &&
			evt.Payload().At("count").AsInt(0) == 5 &&
			evt.Payload().At("min-rating").AsFloat64(0.0) == 1.0 &&
			evt.Payload().At("max-rating").AsFloat64(0.0) == 3.0 &&
			evt.Payload().At("avg-rating").AsFloat64(0.0) == 2.2 &&
			evt.Payload().At("med-rating").AsFloat64(0.0) == 2.0
	})
}

// EOF
