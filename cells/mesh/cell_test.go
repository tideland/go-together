// Tideland Go Together - Cells - Mesh - Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/together/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
)

//--------------------
// TESTS
//--------------------

// TestCellSimple provides a simple processing of some
// events.
func TestCellSimple(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	topics := []string{}
	sigc := make(chan interface{})
	collector := func(cell Cell, evt *Event, out Emitter) error {
		topics = append(topics, evt.Topic())
		if len(topics) == 3 {
			close(sigc)
		}
		return nil
	}
	tbCollector := NewRequestBehavior(collector)
	cCollector := newCell(ctx, "collector", meshStub{}, tbCollector, drop)

	cCollector.in.Emit("one")
	cCollector.in.Emit("two")
	cCollector.in.Emit("three")

	assert.WaitClosed(sigc, time.Second)
	assert.Length(topics, 3)
	assert.Equal(strings.Join(topics, " "), "one two three")

	cancel()
}

// TestCellChain provides a chained processing of some
// events.
func TestCellChain(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	topics := []string{}
	sigc := make(chan interface{})
	upcaser := func(cell Cell, evt *Event, out Emitter) error {
		upperTopic := strings.ToUpper(evt.Topic())
		out.Emit(upperTopic)
		return nil
	}
	tbUpcaser := NewRequestBehavior(upcaser)
	cUpcaser := newCell(ctx, "upcaser", meshStub{}, tbUpcaser, drop)
	collector := func(cell Cell, evt *Event, out Emitter) error {
		topics = append(topics, evt.Topic())
		if len(topics) == 3 {
			close(sigc)
		}
		return nil
	}
	tbCollector := NewRequestBehavior(collector)
	cCollector := newCell(ctx, "collector", meshStub{}, tbCollector, drop)
	cCollector.subscribeTo(cUpcaser)

	cUpcaser.in.Emit("one")
	cUpcaser.in.Emit("two")
	cUpcaser.in.Emit("three")

	assert.WaitClosed(sigc, time.Second)
	assert.Length(topics, 3)
	assert.Equal(strings.Join(topics, " "), "ONE TWO THREE")

	cCollector.unsubscribeFrom(cUpcaser)

	cUpcaser.in.Emit("FOUR")
	cUpcaser.in.Emit("FIVE")
	cUpcaser.in.Emit("SIX")

	assert.Length(topics, 3)
	assert.Equal(strings.Join(topics, " "), "ONE TWO THREE")

	cancel()
}

// TestCellAutoUnsubscribe verifies the automatic unsubscription
// and information.
func TestCellAutoUnsubscribe(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	failed := []*Event{}
	collected := []*Event{}
	sigc := make(chan interface{})
	forwarder := func(cell Cell, evt *Event, out Emitter) error {
		return out.EmitEvent(evt)
	}
	cForwarderA := newCell(ctx, "forwarderA", meshStub{}, NewRequestBehavior(forwarder), drop)
	cForwarderB := newCell(ctx, "forwarderB", meshStub{}, NewRequestBehavior(forwarder), drop)
	failer := func(cell Cell, evt *Event, out Emitter) error {
		failed = append(failed, evt)
		if len(failed) == 3 {
			return errors.New("done")
		}
		return out.EmitEvent(evt)
	}
	cFailer := newCell(ctx, "failer", meshStub{}, NewRequestBehavior(failer), drop)
	cFailer.subscribeTo(cForwarderA)
	cFailer.subscribeTo(cForwarderB)
	collector := func(cell Cell, evt *Event, out Emitter) error {
		collected = append(collected, evt)
		if len(collected) == 3 {
			close(sigc)
		}
		return nil
	}
	cCollector := newCell(ctx, "collector", meshStub{}, NewRequestBehavior(collector), drop)
	cCollector.subscribeTo(cFailer)

	cForwarderA.in.Emit("foo")
	cForwarderB.in.Emit("bar")
	cForwarderA.in.Emit("baz")

	assert.WaitClosed(sigc, time.Second)

	cForwarderA.in.Emit("dont-care")
	cForwarderB.in.Emit("dont-care")

	foundc := make(chan interface{})

	for _, evt := range collected {
		if evt.Topic() == TopicError {
			var errpl PayloadCellError
			err := evt.Payload(&errpl)
			assert.NoError(err)
			assert.Equal(errpl.CellName, "failer")
			assert.Equal(errpl.Error, "done")
			close(foundc)
			break
		}
	}

	assert.WaitClosed(foundc, time.Second, "error not found")
	cancel()
}

//--------------------
// STUBS
//--------------------

// meshStub simulates the mesh for the cells.
type meshStub struct{}

func (ms meshStub) Go(name string, b Behavior) error {
	return nil
}

func (ms meshStub) Subscribe(fromName, toName string) error {
	return nil
}

func (ms meshStub) Unsubscribe(toName, fromName string) error {
	return nil
}

func (ms meshStub) Emit(name, topic string, payloads ...interface{}) error {
	return nil
}

func (ms meshStub) EmitEvent(name string, evt *Event) error {
	return nil
}

func (ms meshStub) Emitter(name string) (Emitter, error) {
	return nil, nil
}

// drop simulates the callback to notify the
// mesh of the termination of a cell.
var drop = func() {}

// EOF
