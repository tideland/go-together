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
	defer assert.NoError(msh.Stop())

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

	assert.NoError(msh.SpawnCells(
		behaviors.NewSelectFilterBehavior("select", filter),
		behaviors.NewConditionBehavior("selects", selectConditioner, processor),
		behaviors.NewExcludeFilterBehavior("exclude", filter),
		behaviors.NewConditionBehavior("excludes", excludesConditioner, processor),
	))
	assert.NoError(msh.Subscribe("select", "selects"))
	assert.NoError(msh.Subscribe("exclude", "excludes"))

	data := [][2]string{
		{"a", "a"},
		{"a", "b"},
		{"a", "c"},
		{"a", "d"},
		{"b", "b"},
		{"b", "a"},
	}
	for _, d := range data {
		assert.NoError(msh.Emit("select", event.New(d[0], "test", d[1])))
	}
	assert.Wait(sigc, true, time.Second)

	for _, d := range data {
		assert.NoError(msh.Emit("exclude", event.New(d[0], "test", d[1])))
	}
	assert.Wait(sigc, true, time.Second)
}

// EOF
