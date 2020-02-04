// Tideland Go Together - Cells - Behaviors - Unit Tests
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors_test // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestTickerBehavior tests the ticker behavior.
func TestTickerBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	plant := mesh.NewTestPlant(assert, behaviors.NewTickerBehavior("ticker", 50*time.Millisecond), 1)
	defer plant.Stop()

	time.Sleep(125 * time.Millisecond)

	plant.AssertLength(0, 2)
}

// EOF
