// Tideland Go Together - Cells - Event - Unit Tests
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package event_test // import "tideland.dev/go/together/cells/event"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/event"
)

//--------------------
// TESTS
//--------------------

// TestTopicOnly verifies creation of an event with only a topic.
func TestTopicOnly(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evt := event.New("test")

	assert.Equal(evt.Topic(), "test")
	vfoo := evt.Payload().At("foo")
	assert.ErrorMatch(vfoo, ".*no payload value at key \"foo\".*")
}

// TestDone verifies a done event based on a context.
func TestDone(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evt := event.New("test")

	assert.False(evt.Done())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	evt = event.WithContext(ctx, "test")

	assert.False(evt.Done())
	cancel()
	assert.True(evt.Done())

	ctx, cancel = context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	evt = event.WithContext(ctx, "test")

	assert.False(evt.Done())
	time.Sleep(100 * time.Millisecond)
	assert.True(evt.Done())
}

// TestKeyValues verifies creation of an event with a topic
// and key/value pairs.
func TestKeyValues(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evt := event.New("topic", "a", 1, "b", "2")

	assert.Equal(evt.Topic(), "topic")
	va := evt.Payload().At("a").AsInt(0)
	assert.Equal(va, 1)
	vb := evt.Payload().At("b").AsInt(0)
	assert.Equal(vb, 2)
}

// TestWithPayload verifies creation of an event with a topic
// and an external created payload.
func TestWithPayload(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	pl := event.NewPayload("a", 1, "b", "2")
	evt := event.New("topic", pl)

	assert.Equal(evt.Topic(), "topic")
	va := evt.Payload().At("a").AsInt(0)
	assert.Equal(va, 1)
	vb := evt.Payload().At("b").AsString("<none>")
	assert.Equal(vb, "2")
}

// TestDefaultValue verifies creation of an event with a topic
// and a valueless final key.
func TestDefaultValue(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evt := event.New("topic", "a")

	assert.Equal(evt.Topic(), "topic")
	va := evt.Payload().At("a").AsBool(false)
	assert.True(va)

	evt = event.New("topic", "a", 1, "b", 2, "c")

	assert.Equal(evt.Topic(), "topic")
	vc := evt.Payload().At("c").AsBool(false)
	assert.True(vc)
}

// EOF
