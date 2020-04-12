// Tideland Go Together - Fuse - Unit Tests
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package fuse_test

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/fuse"
)

//--------------------
// TESTS
//--------------------

// TestTriggerPanic verifies the raising of a panic containing an
// annotated error.
func TestTriggerPanic(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)

	assert.Panics(func() {
		fuse.Trigger(errors.New("should panic"))
	})
}

// EOF
