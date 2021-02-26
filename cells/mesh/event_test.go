// Tideland Go Together - Cells - Mesh - Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh_test

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

// TestSimpleKeyValuePayload verifies creation and simple access of a key/value payload.
func TestSimpleKeyValuePayload(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	var p mesh.Payload = mesh.NewKeyValuePayload("a", 1, "b", "2", 3, false, "d", 12.34, "e")

	kvp, ok := mesh.IsKeyValuePayload(p)
	assert.OK(ok)

	assert.Equal(kvp.At("a"), 1)
	assert.Equal(kvp.At("b"), "2")
	assert.Equal(kvp.At("3"), false)
	assert.Equal(kvp.At("d"), 12.34)
	assert.Equal(kvp.At("e"), true)
	assert.Equal(kvp.Has("e"), true)
	assert.Equal(kvp.At("f"), nil)
	assert.Equal(kvp.Has("f"), false)

	i, ok := kvp.IntAt("a")
	assert.OK(ok)
	assert.Equal(i, 1)
	s, ok := kvp.StringAt("b")
	assert.OK(ok)
	assert.Equal(s, "2")
	b, ok := kvp.BoolAt("3")
	assert.OK(ok)
	assert.Equal(b, false)
	f, ok := kvp.Float64At("d")
	assert.OK(ok)
	assert.Equal(f, 12.34)
	b, ok = kvp.BoolAt("e")
	assert.OK(ok)
	assert.Equal(b, true)
	_, ok = kvp.BoolAt("F")
	assert.False(ok)
}

// TestDoKeyValuePayload verifies the iteration over all key/value payload data.
func TestDoKeyValuePayload(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	kvp := mesh.NewKeyValuePayload("a", 2, "b", 7, "c", 1, "d", 2)
	sum := 0

	kvp.Do(func(key string, value interface{}) {
		i, ok := value.(int)

		assert.True(ok)

		sum += i
	})

	assert.Equal(sum, 12)
}

// TestNestedKeyValuePayload verifies the deep copy of a nested payload.
func TestNestedKeyValuePayload(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bottom := mesh.NewKeyValuePayload("ba", 1, "bb", true)
	mid := mesh.NewKeyValuePayload("ma", 2, "mb", true, "mc", bottom)
	top := mesh.NewKeyValuePayload("ta", 3, "tb", true, "tc", mid)

	ctp := top.DeepCopy()
	ct, ok := mesh.IsKeyValuePayload(ctp)
	assert.OK(ok)
	cmp, ok := ct.PayloadAt("tc")
	assert.OK(ok)
	cm, ok := mesh.IsKeyValuePayload(cmp)
	assert.OK(ok)
	bmp, ok := cm.PayloadAt("mc")
	assert.OK(ok)
	bm, ok := mesh.IsKeyValuePayload(bmp)
	assert.OK(ok)
	assert.True(bm.Has("ba"))
	assert.True(bm.Has("bb"))
}

// TestSimpleFuncPayload verifies the executable function payload.
func TestSimpleFuncPayload(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sum := 0

	var p mesh.Payload = mesh.NewFuncPayload(func(arg mesh.Payload) error {
		kvp, ok := mesh.IsKeyValuePayload(arg)
		assert.OK(ok)

		a, ok := kvp.IntAt("a")
		assert.OK(ok)
		b, ok := kvp.IntAt("b")
		assert.OK(ok)

		sum = a + b

		return nil
	})
	fp, ok := mesh.IsFuncPayload(p)
	assert.OK(ok)

	input := mesh.NewKeyValuePayload("a", 6, "b", 4)
	err := fp.Exec(input)
	assert.NoError(err)
	assert.Equal(sum, 10)

	cp := fp.DeepCopy()
	cfp, ok := mesh.IsFuncPayload(cp)
	assert.OK(ok)

	input = mesh.NewKeyValuePayload("a", 1, "b", 2)
	err = cfp.Exec(input)
	assert.NoError(err)
	assert.Equal(sum, 3)
}

// TestEvent verifies event creation and access.
func TestEvent(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evt := mesh.NewEvent("test", mesh.NewKeyValuePayload("a", 1, "b", 2))

	assert.Equal(evt.Topic(), "test")

	p := evt.Payload()
	kvp, ok := mesh.IsKeyValuePayload(p)
	assert.OK(ok)
	a, ok := kvp.IntAt("a")
	assert.OK(ok)
	assert.Equal(a, 1)
	b, ok := kvp.IntAt("b")
	assert.OK(ok)
	assert.Equal(b, 2)
}

// EOF
