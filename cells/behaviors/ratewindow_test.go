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
	"fmt"
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

// TestRateWindowBehavior tests the event rate window behavior.
func TestRateWindowBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	generator := generators.New(generators.FixedRand())
	msh := mesh.New()
	defer msh.Stop()

	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "bang"}
	duration := 50 * time.Millisecond
	matcher := func(evt *event.Event) (bool, error) {
		match := evt.Topic() == "bang"
		return match, nil
	}
	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		first, _ := accessor.PeekFirst()
		last, _ := accessor.PeekLast()
		difference := last.Timestamp().Sub(first.Timestamp())
		return event.NewPayload("difference", difference), nil
	}
	oncer := func(emitter mesh.Emitter, evt *event.Event) error {
		difference := evt.Payload().At("difference").AsDuration(0)
		assert.True(difference < duration)
		assert.Equal(evt.Topic(), behaviors.TopicRateWindow)
		sigc <- difference
		return nil
	}

	msh.SpawnCells(
		behaviors.NewRateWindowBehavior("windower", matcher, 5, duration, processor),
		behaviors.NewOnceBehavior("oncer", oncer),
	)
	msh.Subscribe("windower", "oncer")

	for i := 0; i < 250; i++ {
		topic := generator.OneStringOf(topics...)
		msh.Emit("windower", event.New(topic))
		time.Sleep(time.Millisecond)
	}

	assert.WaitTested(sigc, func(v interface{}) error {
		difference := v.(time.Duration)
		if difference > 50*time.Millisecond {
			return fmt.Errorf("diff %v", difference)
		}
		return nil
	}, 5*time.Second)
}

// EOF
