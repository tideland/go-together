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

// TestMapperBehavior tests the mapping of events.
func TestMapperBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	mapper := func(evt *event.Event) (*event.Event, error) {
		text := evt.Payload().At("text").AsString("")
		return event.New(evt.Topic(), "text", strings.ToUpper(text)), nil
	}
	plant := mesh.NewTestPlant(assert, behaviors.NewMapperBehavior("mb", mapper), 1)
	defer plant.Stop()

	plant.Emit(event.New("a", "text", "abc"))
	plant.Emit(event.New("b", "text", "def"))
	plant.Emit(event.New("c", "text", "ghi"))

	plant.AssertFind(0, func(evt *event.Event) bool {
		return evt.Topic() == "a" && evt.Payload().At("text").AsString("-") == "ABC"
	})
	plant.AssertFind(0, func(evt *event.Event) bool {
		return evt.Topic() == "b" && evt.Payload().At("text").AsString("-") == "DEF"
	})
	plant.AssertFind(0, func(evt *event.Event) bool {
		return evt.Topic() == "c" && evt.Payload().At("text").AsString("-") == "GHI"
	})
}

// EOF
