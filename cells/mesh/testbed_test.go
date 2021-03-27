// Tideland Go Together - Cells - Mesh - Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh_test // import "tideland.dev/go/together/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestTestbed verifies the working of the testbed for behavior tests.
func TestTestbed(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	forwarder := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				out.EmitEvent(evt)
			}
		}
	}
	behavior := mesh.BehaviorFunc(forwarder)
	count := 0
	tester := func(evt mesh.Event) bool {
		count++
		if count == 3 {
			// Done.
			return true
		}
		return false
	}

	tb := mesh.NewTestbed(behavior, tester)

	tb.Emit("one")
	tb.Emit("two")
	tb.Emit("three")

	err := tb.Wait(time.Second)
	assert.NoError(err)
	assert.Equal(count, 3)
}

// TestTestbedMesh verifies the Mesh stubbing of the testbed.
func TestTestbedMesh(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	mesher := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				switch evt.Topic() {
				case "go":
					out.EmitEvent(evt)
					err := cell.Mesh().Go("anything", nil)
					assert.ErrorContains(err, "cell name 'anything' already used")
				case "subscribe":
					out.EmitEvent(evt)
					err := cell.Mesh().Subscribe("anything", "anything-else")
					assert.ErrorContains(err, "emitter cell 'anything' does not exist")
				case "unsubscribe":
					out.EmitEvent(evt)
					err := cell.Mesh().Unsubscribe("anything", "anything-else")
					assert.ErrorContains(err, "emitter cell 'anything' does not exist")
				case "emit":
					out.EmitEvent(evt)
					err := cell.Mesh().EmitEvent("anything", evt)
					assert.ErrorContains(err, "cell 'anything' does not exist")
				case "emitter":
					out.EmitEvent(evt)
					emtr, err := cell.Mesh().Emitter("anything")
					assert.ErrorContains(err, "cell 'anything' does not exist")
					assert.Nil(emtr)
				case "done":
					out.EmitEvent(evt)
				}
			}
		}
	}
	behavior := mesh.BehaviorFunc(mesher)
	topics := map[string]bool{}
	tester := func(evt mesh.Event) bool {
		topics[evt.Topic()] = true
		return len(topics) == 6
	}

	tb := mesh.NewTestbed(behavior, tester)

	tb.Emit("go")
	tb.Emit("subscribe")
	tb.Emit("unsubscribe")
	tb.Emit("emit")
	tb.Emit("emitter")
	tb.Emit("done")

	err := tb.Wait(time.Second)
	assert.NoError(err)
}

// EOF
