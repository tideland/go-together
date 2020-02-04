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
	"sync"

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
	for i := range es.tp.subscribers {
		ids = append(ids, strconv.Itoa(i))
	}
	return ids
}

// Emit tries to emit the event to the subscriber with the given ID.
func (es *emitterStub) Emit(id string, evt *event.Event) error {
	idx, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	if idx > len(es.tp.subscribers)+1 {
		return errors.New("not found")
	}
	es.tp.subscribers[idx].Process(evt)
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
	go es.tp.Emit(evt)
}

// behaviorStub collects events for later tests.
type behaviorStub struct {
	idx  int
	sink *event.Sink
}

// ID returns the identificator of the simulated behavior.
func (bs *behaviorStub) ID() string {
	return strconv.Itoa(bs.idx)
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
	_, _ = bs.sink.Push(evt)
}

// Recover is not called by testplant.
func (bs *behaviorStub) Recover(err interface{}) error {
	return nil
}

//--------------------
// TESTPLANT
//--------------------

// TestPlant provides help to test behaviors. It is instantiated with an Asserts
// instance, the behavior, and a wanted number of subscribers. Those do get the
// IDs "0" to "N" and simply collect the events emitted by the behavior for later
// test processing.
//
// Then events can be emitted to the behavior so it can do its work. Afterwards
// the AssertsXyz() methods can be used to test the collected events per subscriber.
//
//     plant := mesh.NewTestPlant(assert, NewMyBehavior("myID"), 2)
//     defer plant.Stop()
//
//     plant.Emit(event.New("foo"))
//     plant.Emit(event.New("bar"))
//
//     plant.AssertLength(0, 2)
//     plant.AssertFind(1, func(evt *event.Event) bool {
//         return evt.Topic() == "bar"
//     })
//
// The Stop() at the end ensures the call of the Terminate() method of the behavior.
type TestPlant struct {
	mu          sync.Mutex
	assert      *asserts.Asserts
	behavior    Behavior
	subscribers []*behaviorStub
}

// NewTestPlant creates a test plant for the given behavior and the configured
// number of subscribers.
func NewTestPlant(assert *asserts.Asserts, behavior Behavior, subscribers int) *TestPlant {
	tp := &TestPlant{
		assert:   assert,
		behavior: behavior,
	}
	for i := 0; i < subscribers; i++ {
		tp.subscribers = append(tp.subscribers, &behaviorStub{
			idx:  i,
			sink: event.NewSink(256),
		})
	}
	tp.assert.IncrCallstackOffset()
	tp.assert.IncrCallstackOffset()
	es := &emitterStub{tp}
	assert.OK(tp.behavior.Init(es))
	return tp
}

// Emit passes an event to the behavior to test.
func (tp *TestPlant) Emit(evt *event.Event) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	defer func() {
		if r := recover(); r != nil {
			// TODO Add way to test wanted recoverings.
			tp.assert.OK(tp.behavior.Recover(r))
		}
	}()
	tp.behavior.Process(evt)
}

// Reset clears all collected subscriber events.
func (tp *TestPlant) Reset() {
	for _, subscriber := range tp.subscribers {
		tp.assert.OK(subscriber.sink.Clear())
	}
}

// Stop terminates the behavior and tests its return value.
func (tp *TestPlant) Stop() {
	tp.assert.OK(tp.behavior.Terminate())
}

// AssertLength tests the length of the collected events of a given subscriber.
func (tp *TestPlant) AssertLength(idx int, length int) {
	tp.assert.OK(idx < len(tp.subscribers), "subscriber not found")
	subscriber := tp.subscribers[idx]
	tp.assert.Length(subscriber.sink, length, "collected event length")
}

// AssertAll tests if all collected events of a given subscriber fullfil
// the given test function.
func (tp *TestPlant) AssertAll(idx int, test func(*event.Event) bool) {
	tp.assert.OK(idx < len(tp.subscribers), "subscriber not found")
	subscriber := tp.subscribers[idx]
	tp.assert.OK(subscriber.sink.Do(func(i int, evt *event.Event) error {
		tp.assert.OK(test(evt), "test failed at", strconv.Itoa(i))
		return nil
	}))
}

// AssertFind tests if the collected events of a given subscriber contain
// at least one matching the given match function.
func (tp *TestPlant) AssertFind(idx int, matches func(*event.Event) bool) {
	tp.assert.OK(idx < len(tp.subscribers), "subscriber not found")
	subscriber := tp.subscribers[idx]
	found := false
	tp.assert.OK(subscriber.sink.Do(func(i int, evt *event.Event) error {
		found = found || matches(evt)
		return nil
	}))
	tp.assert.OK(found, "event not found")
}

// AssertNone tests if the collected events of a given subscriber contain
// no one matching the given match function.
func (tp *TestPlant) AssertNone(idx int, matches func(*event.Event) bool) {
	tp.assert.OK(idx < len(tp.subscribers), "subscriber not found")
	subscriber := tp.subscribers[idx]
	found := false
	tp.assert.OK(subscriber.sink.Do(func(i int, evt *event.Event) error {
		found = found || matches(evt)
		return nil
	}))
	tp.assert.False(found, "event found")
}

// AssertFirst tests if the first of the collected events of a given subscriber
// fullfils a given test function.
func (tp *TestPlant) AssertFirst(idx int, test func(*event.Event) bool) {
	tp.assert.OK(idx < len(tp.subscribers), "subscriber not found")
	subscriber := tp.subscribers[idx]
	evt, ok := subscriber.sink.PeekFirst()
	tp.assert.OK(ok, "event not found")
	tp.assert.OK(test(evt), "test failed")
}

// AssertLast tests if the last of the collected events of a given subscriber
// fullfils a given test function.
func (tp *TestPlant) AssertLast(idx int, test func(*event.Event) bool) {
	tp.assert.OK(idx < len(tp.subscribers), "subscriber not found")
	subscriber := tp.subscribers[idx]
	evt, ok := subscriber.sink.PeekLast()
	tp.assert.OK(ok, "event not found")
	tp.assert.OK(test(evt), "test failed")
}

// EOF
