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

// TestCallbackBehavior tests the callback behavior.
func TestCallbackBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()
	defer msh.Stop()

	cbdA := []string{}
	cbfA := func(emitter mesh.Emitter, evt *event.Event) error {
		cbdA = append(cbdA, evt.Topic())
		return nil
	}
	cbdB := 0
	cbfB := func(emitter mesh.Emitter, evt *event.Event) error {
		cbdB++
		return nil
	}
	sigc := asserts.MakeWaitChan()
	cbfC := func(emitter mesh.Emitter, evt *event.Event) error {
		if evt.Topic() == "baz" {
			sigc <- true
		}
		return nil
	}

	msh.SpawnCells(behaviors.NewCallbackBehavior("callback", cbfA, cbfB, cbfC))

	msh.Emit("callback", event.New("foo"))
	msh.Emit("callback", event.New("bar"))
	msh.Emit("callback", event.New("baz"))

	assert.Wait(sigc, true, time.Second)
	assert.Equal(cbdA, []string{"foo", "bar", "baz"})
	assert.Equal(cbdB, 3)
}

// EOF
