// Tideland Go Together - Cells - Mesh - Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh

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
	collector := func(evt *Event, out OutputStream) error {
		topics = append(topics, evt.Topic())
		if len(topics) == 3 {
			close(sigc)
		}
		return nil
	}
	tbCollector := NewStatelessBehavior(collector)
	cCollector := newCell(ctx, ownerStub{}, "collector", tbCollector)

	cCollector.in.Emit(NewEvent("one"))
	cCollector.in.Emit(NewEvent("two"))
	cCollector.in.Emit(NewEvent("three"))

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
	upcaser := func(evt *Event, out OutputStream) error {
		upperTopic := strings.ToUpper(evt.Topic())
		out.Emit(NewEvent(upperTopic))
		return nil
	}
	tbUpcaser := NewStatelessBehavior(upcaser)
	cUpcaser := newCell(ctx, ownerStub{}, "upcaser", tbUpcaser)
	collector := func(evt *Event, out OutputStream) error {
		topics = append(topics, evt.Topic())
		if len(topics) == 3 {
			close(sigc)
		}
		return nil
	}
	tbCollector := NewStatelessBehavior(collector)
	cCollector := newCell(ctx, ownerStub{}, "collector", tbCollector)
	cCollector.subscribeTo(cUpcaser)

	cUpcaser.in.Emit(NewEvent("one"))
	cUpcaser.in.Emit(NewEvent("two"))
	cUpcaser.in.Emit(NewEvent("three"))

	assert.WaitClosed(sigc, time.Second)
	assert.Length(topics, 3)
	assert.Equal(strings.Join(topics, " "), "ONE TWO THREE")

	cCollector.unsubscribeFrom(cUpcaser)

	cUpcaser.in.Emit(NewEvent("FOUR"))
	cUpcaser.in.Emit(NewEvent("FIVE"))
	cUpcaser.in.Emit(NewEvent("SIX"))

	assert.Length(topics, 3)
	assert.Equal(strings.Join(topics, " "), "ONE TWO THREE")

	cancel()
}

// TestCellAutoUnsubscribe verifies the automatic unsubscription
// and information.
func TestCellAutoUnsubscribe(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	events := []*Event{}
	sigc := make(chan interface{})
	forwarder := func(evt *Event, out OutputStream) error {
		return out.Emit(evt)
	}
	cForwarderA := newCell(ctx, ownerStub{}, "forwarderA", NewStatelessBehavior(forwarder))
	cForwarderB := newCell(ctx, ownerStub{}, "forwarderB", NewStatelessBehavior(forwarder))
	failer := func(evt *Event, out OutputStream) error {
		if evt.Topic() == "fail" {
			msg, _ := evt.StringAt("message")
			return errors.New(msg)
		}
		return out.Emit(evt)
	}
	cFailer := newCell(ctx, ownerStub{}, "failer", NewStatelessBehavior(failer))
	cFailer.subscribeTo(cForwarderA)
	cFailer.subscribeTo(cForwarderB)
	collector := func(evt *Event, out OutputStream) error {
		events = append(events, evt)
		if evt.Topic() == ErrorTopic {
			close(sigc)
		}
		return nil
	}
	cCollector := newCell(ctx, ownerStub{}, "collector", NewStatelessBehavior(collector))
	cCollector.subscribeTo(cFailer)

	cForwarderA.in.Emit(NewEvent("foo"))
	cForwarderB.in.Emit(NewEvent("bar"))
	cForwarderA.in.Emit(NewEvent("fail", "message", "ouch"))

	assert.WaitClosed(sigc, time.Second)
	cForwarderA.in.Emit(NewEvent("dont-care"))
	cForwarderB.in.Emit(NewEvent("dont-care"))

	i := len(events)
	assert.True(i < 4)
	for i = 0; i < len(events); i++ {
		if events[i].Topic() == ErrorTopic {
			break
		}
	}
	name, _ := events[i].StringAt(NameKey)
	assert.Equal(name, "failer")
	message, _ := events[i].StringAt(MessageKey)
	assert.Equal(message, "ouch")

	cancel()
}

//--------------------
// STUBS
//--------------------

// ownerStub simulates the mesh.
type ownerStub struct{}

func (d ownerStub) drop(name string) {}

// EOF
