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
	"strings"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestRouterBehavior tests the router behavior.
func TestRouterBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	router := func(evt *event.Event) []string {
		// Topic does not interest in this case.
		return strings.Split(evt.Payload().At("ids").AsString(""), "/")
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewRouterBehavior("rb", router), 5)
	defer plant.Stop()

	plant.Emit(event.New("route-it", "ids", "0/1"))
	plant.Emit(event.New("route-it", "ids", "0/1/2"))
	plant.Emit(event.New("route-it", "ids", "2/3/4"))
	plant.Emit(event.New("route-it", "ids", "unknown"))

	plant.AssertLength(0, 2)
	plant.AssertLength(1, 2)
	plant.AssertLength(2, 2)
	plant.AssertLength(3, 1)
	plant.AssertLength(4, 1)
}

// EOF
