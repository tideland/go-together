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

// TestFilterBehavior tests the filter behavior.
func TestFilterBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	selects := 0
	excludes := 0
	msh := mesh.New()
	defer msh.Stop()

	filter := func(evt *event.Event) (bool, error) {
		payload := evt.Payload().At("test").AsString("")
		return evt.Topic() == payload, nil
	}
	selectConditioner := func(evt *event.Event) bool {
		selects++
		return selects == 2
	}
	excludesConditioner := func(evt *event.Event) bool {
		excludes++
		return excludes == 4
	}
	processor := func(emitter mesh.Emitter, evt *event.Event) error {
		sigc <- true
		return nil
	}

	msh.SpawnCells(
		behaviors.NewSelectFilterBehavior("select", filter),
		behaviors.NewConditionBehavior("selects", selectConditioner, processor),
		behaviors.NewExcludeFilterBehavior("exclude", filter),
		behaviors.NewConditionBehavior("excludes", excludesConditioner, processor),
	)
	msh.Subscribe("select", "selects")
	msh.Subscribe("exclude", "excludes")

	msh.Emit("select", event.New("a", "test", "a"))
	msh.Emit("select", event.New("a", "test", "b"))
	msh.Emit("select", event.New("a", "test", "c"))
	msh.Emit("select", event.New("a", "test", "d"))
	msh.Emit("select", event.New("b", "test", "b"))
	msh.Emit("select", event.New("b", "test", "a"))

	assert.Wait(sigc, true, time.Second)

	msh.Emit("exclude", event.New("a", "test", "a"))
	msh.Emit("exclude", event.New("a", "test", "b"))
	msh.Emit("exclude", event.New("a", "test", "c"))
	msh.Emit("exclude", event.New("a", "test", "d"))
	msh.Emit("exclude", event.New("b", "test", "b"))
	msh.Emit("exclude", event.New("b", "test", "a"))

	assert.Wait(sigc, true, time.Second)
}

// EOF
