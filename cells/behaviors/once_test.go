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

// TestOnceBehavior tests the once behavior.
func TestOnceBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()
	defer msh.Stop()

	oneTimer := func(emitter mesh.Emitter, evt *event.Event) error {
		topic := evt.Topic()
		sigc <- topic
		return emitter.Broadcast(event.New(topic + "/" + topic))
	}
	msh.SpawnCells(
		behaviors.NewOnceBehavior("first", oneTimer),
		behaviors.NewOnceBehavior("second", oneTimer),
	)
	msh.Subscribe("first", "second")

	msh.Emit("first", event.New("foo"))
	msh.Emit("first", event.New("bar"))

	assert.Wait(sigc, "foo", time.Second)
	assert.Wait(sigc, "foo/foo", time.Second)
}

// EOF
