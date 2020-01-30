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
	filter := func(evt *event.Event) (bool, error) {
		payload := evt.Payload().At("test").AsString("")
		return evt.Topic() == payload, nil
	}

	plant := mesh.NewTestPlant(assert, behaviors.NewSelectFilterBehavior("sfb", filter), 1)
	data := [][2]string{
		{"a", "a"},
		{"a", "b"},
		{"a", "c"},
		{"a", "d"},
		{"b", "b"},
		{"b", "a"},
	}
	for _, d := range data {
		plant.Emit(event.New(d[0], "test", d[1]))
	}
	plant.AssertLength("sub-0", 2)
	plant.Stop()

	plant = mesh.NewTestPlant(assert, behaviors.NewExcludeFilterBehavior("efb", filter), 1)
	data = [][2]string{
		{"a", "a"},
		{"a", "b"},
		{"a", "c"},
		{"a", "d"},
		{"b", "b"},
		{"b", "a"},
	}
	for _, d := range data {
		plant.Emit(event.New(d[0], "test", d[1]))
	}
	plant.AssertLength("sub-0", 4)
	plant.Stop()
}

// EOF
