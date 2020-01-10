// Tideland Go Together - Cells - Mesh - Unit Tests
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package mesh_test // import "tideland.dev/go/together/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"errors"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// CONSTANTS
//--------------------

const waitTimeout = time.Second

//--------------------
// TESTS
//--------------------

// TestSpawnCells verifies starting the mesh, spawning some
// cells, and stops the mesh.
func TestSpawnCells(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	err := msh.SpawnCells(
		NewTestBehavior("foo"),
		NewTestBehavior("bar"),
		NewTestBehavior("baz"),
	)
	assert.NoError(err)

	ids := msh.Cells()
	assert.Length(ids, 3)
	assert.Contains(ids, "foo")
	assert.Contains(ids, "bar")
	assert.Contains(ids, "baz")

	err = msh.Stop()
	assert.NoError(err)
}

// TestSpawnDoubleCells verifies starting the mesh, spawning double
// cells, and checking the returned error.
func TestSpawnDoubleCells(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	err := msh.SpawnCells(
		NewTestBehavior("foo"),
	)
	assert.NoError(err)

	err = msh.SpawnCells(
		NewTestBehavior("foo"),
		NewTestBehavior("bar"),
	)
	assert.ErrorContains(err, "spawn cells: double id(s) [foo]")

	ids := msh.Cells()
	assert.Length(ids, 1)

	err = msh.Stop()
	assert.NoError(err)
}

// TestSpawnErrorCells verifies starting the mesh, spawning cell
// returning an error during init, and checking the returned error.
func TestSpawnErrorCells(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	err := msh.SpawnCells(
		NewTestBehavior("foo"),
		NewTestBehavior("bar"),
	)
	assert.NoError(err)

	err = msh.SpawnCells(
		NewTestBehavior("crash"),
		NewTestBehavior("baz"),
		NewTestBehavior("boom"),
	)
	assert.ErrorMatch(err, ".*spawn cells.*crashing.*exploding")

	ids := msh.Cells()
	assert.Length(ids, 3)

	err = msh.Stop()
	assert.NoError(err)
}

// TestSpawnConfigureCells verifies starting the mesh, spawning a self
// configurating cell, and stops the mesh.
func TestSpawnConfigureCells(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	err := msh.SpawnCells(
		NewTestBehavior("foo"),
		&TestConfigureBehavior{NewTestBehavior("bar"), 16},
		&TestConfigureBehavior{NewTestBehavior("baz"), 64},
		&TestConfigureBehavior{NewTestBehavior("invalid"), -128},
	)
	assert.NoError(err)
	assert.Equal(mesh.GetCellQueueCap(msh, "foo"), 1)
	assert.Equal(mesh.GetCellQueueCap(msh, "bar"), 16)
	assert.Equal(mesh.GetCellQueueCap(msh, "baz"), 64)
	assert.Equal(mesh.GetCellQueueCap(msh, "invalid"), 1)

	err = msh.Stop()
	assert.NoError(err)
}

// TestConfigureMesh verifies starting a configured mesh, spawning cells,
// and stops the mesh.
func TestConfigureMesh(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New(mesh.QueueCap(8))

	err := msh.SpawnCells(
		NewTestBehavior("foo"),
		&TestConfigureBehavior{NewTestBehavior("bar"), 2},
		&TestConfigureBehavior{NewTestBehavior("baz"), 16},
		&TestConfigureBehavior{NewTestBehavior("invalid"), -128},
	)
	assert.NoError(err)
	assert.Equal(mesh.GetCellQueueCap(msh, "foo"), 8)
	assert.Equal(mesh.GetCellQueueCap(msh, "bar"), 2)
	assert.Equal(mesh.GetCellQueueCap(msh, "baz"), 16)
	assert.Equal(mesh.GetCellQueueCap(msh, "invalid"), 8)

	err = msh.Stop()
	assert.NoError(err)

	msh = mesh.New(mesh.QueueCap(-16))

	err = msh.SpawnCells(
		NewTestBehavior("foo"),
		&TestConfigureBehavior{NewTestBehavior("bar"), 2},
		&TestConfigureBehavior{NewTestBehavior("baz"), 16},
		&TestConfigureBehavior{NewTestBehavior("invalid"), -128},
	)
	assert.NoError(err)
	assert.Equal(mesh.GetCellQueueCap(msh, "foo"), 1)
	assert.Equal(mesh.GetCellQueueCap(msh, "bar"), 2)
	assert.Equal(mesh.GetCellQueueCap(msh, "baz"), 16)
	assert.Equal(mesh.GetCellQueueCap(msh, "invalid"), 1)

	err = msh.Stop()
	assert.NoError(err)
}

// TestStopCells verifies stopping some cells.
func TestStopCells(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	// Initial environment with subscriptions.
	err := msh.SpawnCells(
		NewTestBehavior("foo"),
		NewTestBehavior("bar"),
		NewTestBehavior("baz"),
	)
	assert.NoError(err)

	ids := msh.Cells()
	assert.Length(ids, 3)
	assert.Contains(ids, "foo")
	assert.Contains(ids, "bar")
	assert.Contains(ids, "baz")

	msh.Subscribe("foo", "bar", "baz")

	fooS, err := msh.Subscribers("foo")
	assert.NoError(err)
	assert.Length(fooS, 2)
	assert.Contains(fooS, "bar")
	assert.Contains(fooS, "baz")

	// Stopping shall unsubscribe too.
	err = msh.StopCells("baz")

	ids = msh.Cells()
	assert.Length(ids, 2)
	assert.Contains(ids, "foo")
	assert.Contains(ids, "bar")

	fooS, err = msh.Subscribers("foo")
	assert.NoError(err)
	assert.Length(fooS, 1)
	assert.Contains(fooS, "bar")

	err = msh.Stop()
	assert.NoError(err)
}

// TestTermination verifies calling the termination method.
func TestTermination(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	// Initial environment with subscriptions.
	err := msh.SpawnCells(
		NewTestBehavior("bang"),
	)
	assert.NoError(err)

	err = msh.Stop()
	assert.ErrorMatch(err, ".*breaking.*")
}

// TestEmitEvents verifies emitting some events to a node.
func TestEmitEvents(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	err := msh.SpawnCells(
		NewTestBehavior("foo"),
	)
	assert.NoError(err)

	msh.Emit("foo", event.New("set", "a", 1))
	msh.Emit("foo", event.New("set", "b", 2))
	msh.Emit("foo", event.New("set", "c", 3))

	pl, plc := event.NewReplyPayload()

	msh.Emit("foo", event.New("send", pl))

	plr, err := plc.Wait(waitTimeout)

	assert.NoError(err)
	assert.Equal(plr.At("a").AsInt(0), 1)
	assert.Equal(plr.At("b").AsInt(0), 2)
	assert.Equal(plr.At("c").AsInt(0), 3)

	err = msh.Stop()
	assert.NoError(err)
}

// TestEmitContextEvents verifies emitting some events with a context
// to a node. Some of those will timeout.
func TestEmitContextEvents(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	err := msh.SpawnCells(
		NewTestBehavior("foo"),
	)
	assert.NoError(err)

	ctxA := context.Background()
	ctxB, cancel := context.WithTimeout(ctxA, 5*time.Millisecond)
	defer cancel()

	msh.Emit("foo", event.WithContext(ctxA, "set", "a", 5))
	msh.Emit("foo", event.WithContext(ctxA, "set", "b", 5))

	time.Sleep(20 * time.Millisecond)

	msh.Emit("foo", event.WithContext(ctxB, "set", "b", 10))

	pl, plc := event.NewReplyPayload()

	msh.Emit("foo", event.New("send", pl))

	plr, err := plc.Wait(waitTimeout)

	assert.NoError(err)
	assert.Equal(plr.At("a").AsInt(0), 5)
	assert.Equal(plr.At("b").AsInt(0), 5)

	err = msh.Stop()
	assert.NoError(err)
}

// TestBroadcastEvents verifies broadcasting some events to a node.
func TestBroadcastEvents(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()
	assertData := func(id string) {
		pl, plc := event.NewReplyPayload()

		msh.Emit(id, event.New("send", pl))

		plr, err := plc.Wait(waitTimeout)

		assert.NoError(err)
		assert.Equal(plr.At("a").AsInt(0), 1)
		assert.Equal(plr.At("b").AsInt(0), 2)
		assert.Equal(plr.At("c").AsInt(0), 3)
	}

	err := msh.SpawnCells(
		NewTestBehavior("foo"),
		NewTestBehavior("bar"),
		NewTestBehavior("baz"),
	)
	assert.NoError(err)

	msh.Broadcast(event.New("set", "a", 1, "b", 2, "c", 3))

	assertData("foo")
	assertData("bar")
	assertData("baz")

	err = msh.Stop()
	assert.NoError(err)
}

// TestBehaviorEmit verifies the emitting to individual subscribers.
func TestBehaviorEmit(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	err := msh.SpawnCells(
		NewTestBehavior("foo"),
		NewTestBehavior("bar"),
		NewTestBehavior("baz"),
	)
	assert.NoError(err)

	msh.Subscribe("foo", "bar", "baz")

	msh.Emit("foo", event.New("emit", "to", "bar", "value", 1234))
	msh.Emit("foo", event.New("emit", "to", "baz", "value", 4321))

	assertSend := func(id string, value int) {
		pl, plc := event.NewReplyPayload()
		msh.Emit(id, event.New("send", pl))
		plr, err := plc.Wait(waitTimeout)
		assert.NoError(err)
		assert.Equal(plr.At("value").AsInt(0), value)
	}

	waitEvents(assert, msh, "foo")
	assertSend("bar", 1234)
	assertSend("baz", 4321)
}

// TestSubscribe verifies the subscription of cells.
func TestSubscribe(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	err := msh.SpawnCells(
		NewTestBehavior("foo"),
		NewTestBehavior("bar"),
		NewTestBehavior("baz"),
	)
	assert.NoError(err)

	msh.Subscribe("foo", "bar", "baz")

	// Directly ask mesh.
	fooS, err := msh.Subscribers("foo")
	assert.NoError(err)
	assert.Length(fooS, 2)
	assert.Contains(fooS, "bar")
	assert.Contains(fooS, "baz")

	// Send event to store subscribers
	msh.Emit("foo", event.New("subscribers"))
	pl, plc := event.NewReplyPayload()
	msh.Emit("foo", event.New("send", pl))
	plr, err := plc.Wait(waitTimeout)
	assert.NoError(err)
	assert.Equal(plr.At("bar").AsInt(0), 1)
	assert.Equal(plr.At("baz").AsInt(0), 1)

	// Set additional values and let emit length.
	msh.Emit("foo", event.New("set", "a", 1, "b", 1234))
	msh.Emit("foo", event.New("length"))
	waitEvents(assert, msh, "foo")

	// Ask bar for received length.
	msh.Emit("bar", event.New("send", pl))
	plr, err = plc.Wait(waitTimeout)
	assert.NoError(err)
	assert.Equal(plr.At("length").AsInt(0), 4)

	err = msh.Stop()
	assert.NoError(err)
}

// TestUnsubscribe verifies the unsubscription of cells.
func TestUnsubscribe(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	err := msh.SpawnCells(
		NewTestBehavior("foo"),
		NewTestBehavior("bar"),
		NewTestBehavior("baz"),
	)
	assert.NoError(err)

	// Subscribe bar and baz, test both.
	msh.Subscribe("foo", "bar", "baz")

	fooS, err := msh.Subscribers("foo")
	assert.NoError(err)
	assert.Length(fooS, 2)
	assert.Contains(fooS, "bar")
	assert.Contains(fooS, "baz")

	// Unsubscribe baz.
	msh.Unsubscribe("foo", "baz")

	fooS, err = msh.Subscribers("foo")
	assert.NoError(err)
	assert.Length(fooS, 1)
	assert.Contains(fooS, "bar")

	err = msh.Stop()
	assert.NoError(err)
}

// TestInvalidSubscriptions verifies the invalid (un)subscriptions of cells.
func TestInvalidSubscriptions(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	err := msh.SpawnCells(
		NewTestBehavior("foo"),
		NewTestBehavior("bar"),
	)
	assert.NoError(err)

	err = msh.Subscribe("foo", "bar", "baz")
	assert.ErrorMatch(err, ".*cannot find cell.*")

	err = msh.Subscribe("foo", "bar")
	assert.NoError(err)

	err = msh.Unsubscribe("foo", "bar", "baz")
	assert.ErrorMatch(err, ".*cannot find cell.*")

	err = msh.Unsubscribe("foo", "bar")
	assert.NoError(err)

	err = msh.Unsubscribe("foo", "bar")
	assert.NoError(err)

	err = msh.Stop()
	assert.NoError(err)
}

// TestSubscriberIDs verifies the retrieval of subscriber IDs.
func TestSubscriberIDs(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	err := msh.SpawnCells(
		NewTestBehavior("foo"),
		NewTestBehavior("bar"),
		NewTestBehavior("baz"),
	)
	assert.NoError(err)

	err = msh.Subscribe("foo", "bar", "baz")
	assert.NoError(err)

	subscriberIDs, err := msh.Subscribers("foo")
	assert.NoError(err)
	assert.Length(subscriberIDs, 2)

	subscriberIDs, err = msh.Subscribers("bar")
	assert.NoError(err)
	assert.Length(subscriberIDs, 0)

	err = msh.Unsubscribe("foo", "baz")
	assert.NoError(err)

	subscriberIDs, err = msh.Subscribers("foo")
	assert.NoError(err)
	assert.Length(subscriberIDs, 1)

	err = msh.Stop()
	assert.NoError(err)
}

//--------------------
// HELPERS
//--------------------

func waitEvents(assert *asserts.Asserts, msh *mesh.Mesh, id string) {
	pl, plc := event.NewReplyPayload()
	msh.Emit(id, event.New("send", pl))
	_, err := plc.Wait(waitTimeout)
	assert.NoError(err)
}

type TestBehavior struct {
	id      string
	emitter mesh.Emitter
	datas   map[string]int
}

func NewTestBehavior(id string) *TestBehavior {
	return &TestBehavior{
		id:    id,
		datas: make(map[string]int),
	}
}

func (tb *TestBehavior) ID() string {
	return tb.id
}

func (tb *TestBehavior) Init(emitter mesh.Emitter) error {
	switch tb.id {
	case "crash":
		return errors.New("crashing")
	case "boom":
		return errors.New("exploding")
	}
	tb.emitter = emitter
	return nil
}

func (tb *TestBehavior) Terminate() error {
	if tb.id == "bang" {
		return errors.New("breaking")
	}
	return nil
}

func (tb *TestBehavior) Process(evt *event.Event) error {
	switch evt.Topic() {
	case "set":
		return evt.Payload().Do(func(key string, value *event.Value) error {
			tb.datas[key] = value.AsInt(-1)
			return nil
		})
	case "emit":
		to := evt.Payload().At("to").AsString("<unknown>")
		value := evt.Payload().At("value").AsInt(-1)
		tb.emitter.Emit(to, event.New("set", "value", value))
	case "subscribers":
		ids := tb.emitter.Subscribers()
		for _, id := range ids {
			tb.datas[id] = 1
		}
	case "length":
		tb.emitter.Broadcast(event.New("set", "length", len(tb.datas)))
	case "send":
		evt.Payload().Reply(event.NewPayload(tb.datas))
	case "clear":
		tb.datas = make(map[string]int)
	}
	return nil
}

func (tb *TestBehavior) Recover(r interface{}) error {
	return nil
}

type TestConfigureBehavior struct {
	*TestBehavior
	queueCap int
}

func (tcb *TestConfigureBehavior) Configure(c mesh.Configurable) {
	c.SetQueueCap(tcb.queueCap)
}

// EOF
