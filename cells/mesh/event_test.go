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

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestEventSimple verifies creation and simple access events.
func TestEventSimple(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	payloadIn := []string{"a", "b", "c"}
	evt, err := mesh.NewEvent("")
	assert.ErrorContains(err, "event needs topic")
	assert.True(mesh.IsNilEvent(evt))

	evt, err = mesh.NewEvent("test-a")

	assert.NoError(err)
	assert.Equal(evt.Topic(), "test-a")
	assert.False(evt.HasPayload())

	evt, err = mesh.NewEvent("test-b", payloadIn)
	payloadOut := []string{}

	assert.NoError(err)
	assert.True(evt.HasPayload())

	err = evt.Payload(&payloadOut)

	assert.NoError(err)
	assert.Length(payloadOut, 3)
}

// EOF
