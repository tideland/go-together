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
	"errors"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"

	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestEventSinkSimple verifies the standard functions of the event sink.
func TestEventSinkSimple(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Empty limited sink.
	max := 5
	sink := mesh.NewEventSink(max)
	assert.Length(sink, 0)
	first, ok := sink.PeekFirst()
	assert.True(mesh.IsNilEvent(first))
	assert.False(ok)
	last, ok := sink.PeekLast()
	assert.True(mesh.IsNilEvent(last))
	assert.False(ok)
	at, ok := sink.PeekAt(-1)
	assert.True(mesh.IsNilEvent(at))
	assert.False(ok)
	at, ok = sink.PeekAt(1337)
	assert.True(mesh.IsNilEvent(at))
	assert.False(ok)

	// Fill empty sink.
	evts := generateEvents(2 * max)
	for _, evt := range evts {
		l := sink.Push(evt)
		assert.True(l <= max)
	}
	assert.Length(sink, max)
	first, ok = sink.PeekFirst()
	assert.OK(ok)
	assert.Equal(first, evts[len(evts)-max])
	last, ok = sink.PeekLast()
	assert.OK(ok)
	assert.Equal(last, evts[len(evts)-1])

	// Filled limited sink.
	sink = mesh.NewEventSink(max, evts...)
	assert.Length(sink, max)
	first, ok = sink.PeekFirst()
	assert.OK(ok)
	assert.Equal(first, evts[len(evts)-max])
	last, ok = sink.PeekLast()
	assert.OK(ok)
	assert.Equal(last, evts[len(evts)-1])
	first = sink.PullFirst()
	assert.Equal(first, evts[len(evts)-max])
	assert.Length(sink, max-1)

	// Filled unlimited sink.
	sink = mesh.NewEventSink(0, evts...)
	assert.Length(sink, len(evts))
	first, ok = sink.PeekFirst()
	assert.OK(ok)
	assert.Equal(first, evts[0])
	last, ok = sink.PeekLast()
	assert.OK(ok)
	assert.Equal(last, evts[len(evts)-1])
}

// TestEventSinkDo verifies the iterating over a sink.
func TestEventSinkDo(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Do without error.
	evts := generateEvents(20)
	sink := mesh.NewEventSink(0, evts...)
	err := sink.Do(func(i int, evt mesh.Event) error {
		assert.Equal(evt, evts[i])
		return nil
	})
	assert.NoError(err)

	// Do with error.
	err = sink.Do(func(i int, evt mesh.Event) error {
		return errors.New("ouch")
	})
	assert.ErrorContains(err, "ouch")
}

// TestEventSinkFunctions verifies the functions on a sink reader.
func TestEventSinkFunctions(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	topics := []string{"a", "b", "b", "b", "a", "a", "b", "b", "a", "b"}
	sink := mesh.NewEventSink(0, generateTopicEvents(topics)...)

	// Filter events and reutrn error.
	evts, err := mesh.EventSinkFilter(sink, func(i int, evt mesh.Event) (bool, error) {
		return evt.Topic() == "a", nil
	})
	assert.NoError(err)
	assert.Length(evts, 4)
	evts, err = mesh.EventSinkFilter(sink, func(i int, evt mesh.Event) (bool, error) {
		return false, errors.New("ouch")
	})
	assert.ErrorContains(err, "ouch")
	assert.Length(evts, 0)

	// Match and mismatch events.
	ok, err := mesh.EventSinkMatch(sink, func(i int, evt mesh.Event) (bool, error) {
		return evt.Topic() == "a" || evt.Topic() == "b", nil
	})
	assert.NoError(err)
	assert.OK(ok)
	ok, err = mesh.EventSinkMatch(sink, func(i int, evt mesh.Event) (bool, error) {
		return evt.Topic() == "a" || evt.Topic() == "x", nil
	})
	assert.NoError(err)
	assert.False(ok)

	// Fold events.
	inject, err := mesh.NewEvent("counts", make(map[string]int))
	assert.NoError(err)
	facc, err := mesh.EventSinkFold(sink, inject, func(i int, acc, evt mesh.Event) (mesh.Event, error) {
		payload := make(map[string]int)
		err := acc.Payload(&payload)
		assert.NoError(err)
		payload[evt.Topic()]++
		return mesh.NewEvent(acc.Topic(), payload)
	})
	assert.NoError(err)
	payload := make(map[string]int)
	err = facc.Payload(&payload)
	assert.NoError(err)
	assert.Equal(payload["a"], 4)
	assert.Equal(payload["b"], 6)
}

//--------------------
// HELPER
//--------------------

// generateEvents generates a number of events for tests.
func generateEvents(count int) []mesh.Event {
	generator := generators.New(generators.FixedRand())
	topics := generator.Words(count)
	return generateTopicEvents(topics)
}

// generateTopicEvents generates a number of events for tests
// based on topics.
func generateTopicEvents(topics []string) []mesh.Event {
	evts := []mesh.Event{}
	for _, topic := range topics {
		evt, _ := mesh.NewEvent(topic)
		evts = append(evts, evt)
	}
	return evts
}

// EOF
