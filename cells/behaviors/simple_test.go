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

// TestSimpleProcessorBehavior tests the simple processor behavior.
func TestSimpleProcessorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()
	defer msh.Stop()

	topics := []string{}
	var wg sync.WaitGroup
	spf := func(emitter mesh.Emitter, evt *event.Event) error {
		topics = append(topics, evt.Topic())
		wg.Done()
		return nil
	}

	msh.SpawnCells(
		behaviors.NewSimpleProcessorBehavior("simple", spf),
	)

	wg.Add(3)
	msh.Emit("simple", event.New("foo"))
	msh.Emit("simple", event.New("bar"))
	msh.Emit("simple", event.New("baz"))

	wg.Wait()
	assert.Length(topics, 3)
}

// EOF
