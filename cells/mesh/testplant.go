// Tideland Go Together - Cells - Mesh
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package mesh // import "tideland.dev/go/together/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"strconv"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/event"
)

//--------------------
// BEHAV
//--------------------

// emitterStub provides the test object the emitter.
type emitterStub struct {
	tp *TestPlant
}

// Subscribers returns the IDs of the configured subscribers.
func (es *emitterStub) Subscribers() []string {
	var ids []string
	for id := range es.tp.subscribers {
		ids = append(ids, id)
	}
	return ids
}

// Emit tries to emit the event to the subscriber with the given ID.
func (es *emitterStub) Emit(id string, evt *event.Event) error {
	bs, ok := es.tp.subscribers[id]
	if !ok {
		return errors.New("not found")
	}
	bs.Process(evt)
	return nil
}

// Broadcast emits the given event to all subscribers.
func (es *emitterStub) Broadcast(evt *event.Event) error {
	for _, bs := range es.tp.subscribers {
		bs.Process(evt)
	}
	return nil
}

// Self emits the given event back to the cell itself.
func (es *emitterStub) Self(evt *event.Event) {
	go es.tp.behavior.Process(evt)
}

// behaviorStub collects events for later tests.
type behaviorStub struct {
	id   string
	evts []*event.Event
}

// ID returns the identificator of the simulated behavior.
func (bs *behaviorStub) ID() string {
	return bs.id
}

// Init doesn't care for the passed emitter.
func (bs *behaviorStub) Init(emitter Emitter) error {
	return nil
}

// Terminate is not interesting for the stub.
func (bs *behaviorStub) Terminate() error {
	return nil
}

// Process collects the received events.
func (bs *behaviorStub) Process(evt *event.Event) {
	bs.evts = append(bs.evts, evt)
}

// Recover is not called by testplant.
func (bs *behaviorStub) Recover(err interface{}) error {
	return nil
}

//--------------------
// TESTPLANT
//--------------------

// TestPlant provides help to test a behavior
type TestPlant struct {
	assert      *asserts.Asserts
	reset       func()
	behavior    Behavior
	subscribers map[string]*behaviorStub
}

// NewTestPlant creates a test plant for the given behavior and the configured
// number of subscribers.
func NewTestPlant(assert *asserts.Asserts, behavior Behavior, subscribers int) *TestPlant {
	tp := &TestPlant{
		assert:      assert,
		behavior:    behavior,
		subscribers: make(map[string]*behaviorStub),
	}
	for i := 0; i < subscribers; i++ {
		id := "sub-" + strconv.Itoa(i)
		tp.subscribers[id] = &behaviorStub{
			id: id,
		}
	}
	tp.assert.IncrCallstackOffset()
	tp.assert.IncrCallstackOffset()
	es := &emitterStub{tp}
	assert.OK(tp.behavior.Init(es))
	return tp
}

// Emit passes an event to the behavior to test.
func (tp *TestPlant) Emit(evt *event.Event) {
	defer func() {
		if r := recover(); r != nil {
			// TODO Add way to test wanted recoverings.
			tp.assert.OK(tp.behavior.Recover(r))
		}
	}()
	tp.behavior.Process(evt)
}

// Stop terminates the behavior.
func (tp *TestPlant) Stop() {
	tp.assert.OK(tp.behavior.Terminate())
}

// AssertLength tests the length of the collected events of a given subscriber.
func (tp *TestPlant) AssertLength(id string, length int) {
	subscriber, ok := tp.subscribers[id]
	tp.assert.OK(ok, "subscriber not found:", id)
	tp.assert.Length(subscriber.evts, length, "collected event length")
}

// AssertAll tests if all collected events of a given subscriber fullfil
// the given test function.
func (tp *TestPlant) AssertAll(id string, test func(*event.Event) bool) {
	subscriber, ok := tp.subscribers[id]
	tp.assert.OK(ok, "subscriber not found:", id)
	for i, evt := range subscriber.evts {
		tp.assert.OK(test(evt), "test failed at", strconv.Itoa(i))
	}
}

// AssertFind tests if the collected events of a given subscriber contain
// at least one matching the given find function.
func (tp *TestPlant) AssertFind(id string, find func(*event.Event) bool) {
	subscriber, ok := tp.subscribers[id]
	tp.assert.OK(ok, "subscriber not found:", id)
	for _, evt := range subscriber.evts {
		if find(evt) {
			return
		}
	}
	tp.assert.OK(errors.New("event not found"))
}

// EOF
