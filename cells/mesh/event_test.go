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
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestEventSimple verifies creation and simple access events.
func TestEventSimple(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evt := mesh.NewEvent("test", "a", 1, "b", "2", 3, false, "d", 12.34, "e")

	assert.Equal(evt.Topic(), "test")
	assert.Equal(evt.HasValue("e"), true)
	assert.Equal(evt.HasValue("f"), false)

	i, ok := evt.IntAt("a")
	assert.OK(ok)
	assert.Equal(i, 1)
	s, ok := evt.StringAt("b")
	assert.OK(ok)
	assert.Equal(s, "2")
	b, ok := evt.BoolAt("3")
	assert.OK(ok)
	assert.Equal(b, false)
	f, ok := evt.Float64At("d")
	assert.OK(ok)
	assert.Equal(f, 12.34)
	b, ok = evt.BoolAt("e")
	assert.OK(ok)
	assert.Equal(b, true)
	_, ok = evt.BoolAt("F")
	assert.False(ok)

	evt = mesh.NewEvent("no-payload")

	assert.Equal(evt.Topic(), "no-payload")
}

// TestEventIterateValues verifies the iteration over all values.
func TestEventIterateValues(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evt := mesh.NewEvent("test", "a", 2, "b", 7, "c", 1, "d", 2)
	sum := 0

	evt.Do(func(key string, value interface{}) {
		i, ok := value.(int)

		assert.True(ok)

		sum += i
	})

	assert.Equal(sum, 12)
}

// TestEventFunc verifies a copyable function as value.
func TestSimpleFuncPayload(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sum := 0
	f := mesh.CopyableFunc(func(arg interface{}) error {
		i, ok := arg.(int)
		assert.OK(ok)
		sum += i
		return nil
	})
	evt := mesh.NewEvent("add", "f", f)

	vf, ok := evt.CopyableAt("f")
	assert.OK(ok)
	cf, ok := mesh.IsCopyableFunc(vf)
	assert.OK(ok)

	err := cf.Exec(1)
	assert.NoError(err)
	err = cf.Exec(2)
	assert.NoError(err)
	assert.Equal(sum, 3)
}

// EOF
