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
	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestRoundRobinBehavior tests the round robin behavior.
func TestRoundRobinBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()
	defer msh.Stop()

	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		sigc <- accessor.Len()
		return nil, nil
	}

	msh.SpawnCells(
		behaviors.NewRoundRobinBehavior("round-robin"),
		behaviors.NewCollectorBehavior("round-robin-1", 10, processor),
		behaviors.NewCollectorBehavior("round-robin-2", 10, processor),
		behaviors.NewCollectorBehavior("round-robin-3", 10, processor),
		behaviors.NewCollectorBehavior("round-robin-4", 10, processor),
		behaviors.NewCollectorBehavior("round-robin-5", 10, processor),
	)
	msh.Subscribe("round-robin", "round-robin-1", "round-robin-2", "round-robin-3", "round-robin-4", "round-robin-5")

	for i := 0; i < 25; i++ {
		err := msh.Emit("round-robin", event.New("round"))
		assert.Nil(err)
	}

	time.Sleep(50 * time.Millisecond)

	for i := 1; i < 6; i++ {
		id := fmt.Sprintf("round-robin-%d", i)
		msh.Emit(id, event.New(event.TopicProcess))
		assert.Wait(sigc, 5, time.Second)
	}
}

// EOF
