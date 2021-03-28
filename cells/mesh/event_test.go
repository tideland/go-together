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
	"encoding/json"
	"testing"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestEventSimple verifies events without payloads.
func TestEventSimple(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	evt, err := mesh.NewEvent("")
	assert.ErrorContains(err, "event needs topic")
	assert.True(mesh.IsNilEvent(evt))

	evt, err = mesh.NewEvent("test")
	assert.NoError(err)
	assert.Equal(evt.Topic(), "test")
	assert.False(evt.HasPayload())
}

// TestEventPayload verifies events with one or more payloads.
func TestEventPayload(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	payloadIn := []string{"a", "b", "c"}
	payloadOutA := []string{}
	evt, err := mesh.NewEvent("test", payloadIn)
	assert.NoError(err)
	assert.Equal(evt.Topic(), "test")
	assert.True(evt.HasPayload())
	err = evt.Payload(&payloadOutA)
	assert.NoError(err)
	assert.Length(payloadOutA, 3)
	assert.Equal(payloadOutA, payloadIn)

	payloadOutB := []int{}
	evt, err = mesh.NewEvent("test", 1, 2, 3, 4, 5)
	assert.NoError(err)
	assert.Equal(evt.Topic(), "test")
	assert.True(evt.HasPayload())
	err = evt.Payload(&payloadOutB)
	assert.NoError(err)
	assert.Length(payloadOutB, 5)
	assert.Equal(payloadOutB, []int{1, 2, 3, 4, 5})
}

// TestEventMarshaling verifies the event marshaling and unmarshaling.
func TestEventMarshaling(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	evtIn, err := mesh.NewEvent("test")
	assert.NoError(err)
	data, err := json.Marshal(evtIn)
	assert.NoError(err)

	evtOut := mesh.Event{}
	err = json.Unmarshal(data, &evtOut)
	assert.NoError(err)
	assert.Equal(evtOut, evtIn)

	plEvtA, err := mesh.NewEvent("payload-a")
	assert.NoError(err)
	plEvtB, err := mesh.NewEvent("payload-b")
	assert.NoError(err)
	plEvtC, err := mesh.NewEvent("payload-c")
	assert.NoError(err)

	evtIn, err = mesh.NewEvent("test", plEvtA, plEvtB, plEvtC)
	assert.NoError(err)
	data, err = json.Marshal(evtIn)
	assert.NoError(err)

	evtOut = mesh.Event{}
	err = json.Unmarshal(data, &evtOut)
	assert.NoError(err)
	assert.Equal(evtOut, evtIn)
	pl := []mesh.Event{}
	err = evtOut.Payload(&pl)
	assert.NoError(err)
	assert.Equal(pl[0], plEvtA)
	assert.Equal(pl[1], plEvtB)
	assert.Equal(pl[2], plEvtC)
}

// EOF
