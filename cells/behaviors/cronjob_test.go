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
	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestCronjobBehavior tests the cronjob behavior.
func TestCronjobBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()
	defer assert.NoError(msh.Stop())

	cronjob := func(emitter mesh.Emitter) {
		assert.NoError(emitter.Broadcast(event.New("action")))
	}
	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		sigc <- accessor.Len()
		assert.NoError(accessor.Do(func(index int, evt *event.Event) error {
			assert.Equal(evt.Topic(), "action")
			return nil
		}))
		return nil, nil
	}

	assert.NoError(msh.SpawnCells(
		behaviors.NewCronjobBehavior("cronjob", 50*time.Millisecond, cronjob),
		behaviors.NewCollectorBehavior("collector", 10, processor),
	))
	assert.NoError(msh.Subscribe("cronjob", "collector"))

	time.Sleep(550 * time.Millisecond)

	assert.NoError(msh.Emit("collector", event.New(event.TopicProcess)))
	assert.Wait(sigc, 10, time.Minute)
}

// EOF
