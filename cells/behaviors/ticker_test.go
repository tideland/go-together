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

// TestTickerBehavior tests the ticker behavior.
func TestTickerBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()
	defer msh.Stop()

	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		sigc <- accessor.Len()
		return nil, nil
	}

	msh.SpawnCells(
		behaviors.NewTickerBehavior("ticker", 50*time.Millisecond),
		behaviors.NewCollectorBehavior("collector", 10, processor),
	)
	msh.Subscribe("ticker", "collector")

	time.Sleep(125 * time.Millisecond)

	msh.Emit("collector", event.New(event.TopicProcess))
	assert.Wait(sigc, 2, time.Minute)
}

// EOF
