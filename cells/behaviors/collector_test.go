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
	"tideland.dev/go/audit/generators"
	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestCollectorBehavior tests the collector behavior.
func TestCollectorBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	generator := generators.New(generators.FixedRand())
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()
	defer msh.Stop()

	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		sigc <- accessor.Len()
		return nil, nil
	}

	msh.SpawnCells(behaviors.NewCollectorBehavior("collector", 10, processor))

	// Don't care for words, we collect maximally 10 events.
	for _, word := range generator.Words(25) {
		msh.Emit("collector", event.New("collect", word))
	}

	msh.Emit("collector", event.New(event.TopicProcess))
	assert.Wait(sigc, 10, time.Second)

	msh.Emit("collector", event.New(event.TopicReset))

	msh.Emit("collector", event.New(event.TopicProcess))
	assert.Wait(sigc, 0, time.Second)
}

// EOF
