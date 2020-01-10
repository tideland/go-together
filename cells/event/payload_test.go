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
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/event"
)

//--------------------
// TESTS
//--------------------

// TestSimplePayload verifies creation of a payload with key/value pairs.
func TestSimplePayload(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	now := time.Now()
	pl := event.NewPayload("a", 1, "b", "2", "c", now, "d")

	assert.Equal(pl.At("a").AsInt(0), 1)
	assert.Equal(pl.At("b").AsInt(0), 2)
	assert.Equal(pl.At("b").AsString("0"), "2")
	assert.Equal(pl.At("c").AsTime(time.Time{}), now)
	assert.True(pl.At("d").AsBool(false))
	assert.True(pl.At("d").IsDefined())
	assert.True(pl.At("e").IsUndefined())
}

// TestPayloadDefaults verifies retrieving default values from payloads.
func TestPayloadDefaults(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	now := time.Now()
	pl := event.NewPayload()

	assert.Equal(pl.At("a").AsString("foo"), "foo")
	assert.Equal(pl.At("a").AsInt(1234), 1234)
	assert.Equal(pl.At("a").AsFloat64(12.34), 12.34)
	assert.Equal(pl.At("a").AsBool(true), true)
	assert.Equal(pl.At("a").AsTime(now), now)
	assert.Equal(pl.At("a").AsDuration(time.Second), time.Second)
}

// TestNestedPayloads verifies retrieving values from nested payloads.
func TestNestedPayloads(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	plc := event.NewPayload("ca", 100, "cb", 200)
	plb := event.NewPayload("ba", 10, "bb", plc)
	pla := event.NewPayload("aa", 1, "ab", plb)

	assert.Equal(pla.At("aa").AsInt(0), 1)
	assert.Equal(pla.At("ab", "ba").AsInt(0), 10)
	assert.Equal(pla.At("ab", "bb", "ca").AsInt(0), 100)
	assert.Equal(pla.At("ab", "bb", "cb").AsInt(0), 200)
}

// TestPayloadValue verifies merging of payloads as values.
func TestPayloadValue(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	pla := event.NewPayload("a", 1, "b", 2)
	plb := event.NewPayload(pla)

	assert.Equal(plb.At("a").AsInt(0), 1)
	assert.Equal(plb.At("b").AsInt(0), 2)
}

// TestMapKeys tests the treating of map keys.
func TestMapKeys(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ma := map[int]string{
		1: "foo",
		2: "bar",
	}
	mb := map[string]interface{}{
		"a": 12.34,
		"b": ma,
	}
	pl := event.NewPayload("c", true, mb)

	assert.Equal(pl.At("a").AsFloat64(0.0), 12.34)
	assert.Equal(pl.At("b", "1").AsString("-"), "foo")
	assert.Equal(pl.At("b", "2").AsString("-"), "bar")
	assert.True(pl.At("c").AsBool(false))
}

// TestIteratableKeys tests the treating of iteratable keys.
func TestIteratableKeys(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	kas := []string{"a", "b", "c"}
	kbs := [3]int{1, 2, 3}
	pl := event.NewPayload(
		"a", event.NewPayload(kas),
		"b", event.NewPayload(kbs),
	)

	assert.Equal(pl.At("a", "0").AsString("-"), "a")
	assert.Equal(pl.At("a", "1").AsString("-"), "b")
	assert.Equal(pl.At("a", "2").AsString("-"), "c")

	assert.Equal(pl.At("b", "0").AsInt(0), 1)
	assert.Equal(pl.At("b", "1").AsInt(0), 2)
	assert.Equal(pl.At("b", "2").AsInt(0), 3)
}

// TestMapValues tests the treating of map values as nested
// payloads.
func TestMapValues(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ma := map[int]string{
		1: "foo",
		2: "bar",
	}
	mb := map[string]interface{}{
		"a": 12.34,
		"b": ma,
	}
	pl := event.NewPayload("a", 1, "b", mb)

	assert.Equal(pl.At("a").AsInt(0), 1)
	assert.Equal(pl.At("b", "a").AsFloat64(0.0), 12.34)
	assert.Equal(pl.At("b", "b", "1").AsString(""), "foo")
	assert.Equal(pl.At("b", "b", "2").AsString(""), "bar")
}

// TestIteratableValues tests the treating of array and slice values
// as nested payloads.
func TestIteratableValues(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	vas := []string{"a", "b", "c"}
	vbs := [3]int{1, 2, 3}
	pl := event.NewPayload("a", vas, "b", vbs)

	assert.Equal(pl.At("a", "0").AsString("-"), "a")
	assert.Equal(pl.At("a", "1").AsString("-"), "b")
	assert.Equal(pl.At("a", "2").AsString("-"), "c")

	assert.Equal(pl.At("b", "0").AsInt(0), 1)
	assert.Equal(pl.At("b", "1").AsInt(0), 2)
	assert.Equal(pl.At("b", "2").AsInt(0), 3)
}

// TestPayloadClone verifies the cloning of payloads together with
// modification of individual ones.
func TestPayloadClone(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	pla := event.NewPayload("a", 1, "b", "two", "c", 3.0)

	assert.Equal(pla.At("a").AsInt(0), 1)
	assert.Equal(pla.At("b").AsString("zero"), "two")
	assert.Equal(pla.At("c").AsFloat64(0.0), 3.0)

	plb := pla.Clone("a", "4711", "d", "foo")

	assert.Equal(plb.At("a").AsInt(0), 4711)
	assert.Equal(plb.At("b").AsString("zero"), "two")
	assert.Equal(plb.At("c").AsFloat64(0.0), 3.0)
	assert.Equal(plb.At("d").AsString("bar"), "foo")

	assert.Length(plb.Keys(), 4)
}

// EOF
