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
	"strings"
	"sync"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestMapperBehavior tests the mapping of events.
func TestMapperBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	var wg sync.WaitGroup

	mapper := func(evt *event.Event) (*event.Event, error) {
		text := evt.Payload().At("text").AsString("")
		return event.New(evt.Topic(), "text", strings.ToUpper(text)), nil
	}
	processor := func(emitter mesh.Emitter, evt *event.Event) error {
		wg.Done()
		text := evt.Payload().At("text").AsString("")
		switch evt.Topic() {
		case "a":
			assert.Equal(text, "ABC")
		case "b":
			assert.Equal(text, "DEF")
		case "c":
			assert.Equal(text, "GHI")
		default:
			assert.Fail("mapper didn't work: %s = %s", evt.Topic(), text)
		}
		return nil
	}

	assert.OK(msh.SpawnCells(
		behaviors.NewMapperBehavior("mapper", mapper),
		behaviors.NewSimpleProcessorBehavior("processor", processor),
	))
	assert.OK(msh.Subscribe("mapper", "processor"))

	wg.Add(3)
	assert.OK(msh.Emit("mapper", event.New("a", "text", "abc")))
	assert.OK(msh.Emit("mapper", event.New("b", "text", "def")))
	assert.OK(msh.Emit("mapper", event.New("c", "text", "ghi")))
	wg.Wait()
	assert.OK(msh.Stop())
}

// EOF
