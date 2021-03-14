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
				out.Emit(evt)
			}
		}
	}
	behavior := mesh.BehaviorFunc(forwarder)
	count := 0
	tester := func(evt *mesh.Event) bool {
		count++
		if count == 3 {
			// Done.
			return true
		}
		return false
	}

	tb := mesh.NewTestbed(behavior, tester)

	tb.Emit(mesh.NewEvent("one"))
	tb.Emit(mesh.NewEvent("two"))
	tb.Emit(mesh.NewEvent("three"))

	err := tb.Wait(time.Second)
	assert.NoError(err)
	assert.Equal(count, 3)
}

// EOF
