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
	"fmt"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestFSMBehavior tests the finite state machine behavior.
func TestFSMBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	lock := lockMachine{"one", 0}
	plant := mesh.NewTestPlant(assert, behaviors.NewFSMBehavior("fsmb", behaviors.FSMStatus{"locked", lock.Locked, nil}), 1)
	defer plant.Stop()

	// 1st run: emit not enough and press button.
	plant.Emit(event.New("coin", "cents", 20))
	plant.Emit(event.New("coin", "cents", 20))
	plant.Emit(event.New("coin", "cents", 20))
	plant.Emit(event.New("info"))
	plant.Emit(event.New("press-button"))
	plant.Emit(event.New("check-cents"))

	plant.AssertFind(0, func(evt *event.Event) bool {
		return evt.Topic() == "status" &&
			evt.Payload().At("status").AsString("-") == "locked" &&
			evt.Payload().At("cents").AsInt(-1) == 60
	})
	plant.AssertLast(0, func(evt *event.Event) bool {
		return evt.Topic() == "cents-checked" && evt.Payload().At("cents").AsInt(-1) == 0
	})
	plant.Reset()

	// 2nd run: unlock the lock and lock it again.
	plant.Emit(event.New("coin", "cents", 50))
	plant.Emit(event.New("coin", "cents", 20))
	plant.Emit(event.New("coin", "cents", 50))
	plant.Emit(event.New("info"))
	plant.Emit(event.New("press-button"))

	plant.AssertFind(0, func(evt *event.Event) bool {
		return evt.Topic() == "unlocked"
	})
	plant.AssertLast(0, func(evt *event.Event) bool {
		return evt.Topic() == "coins-dropped" && evt.Payload().At("cents").AsInt(-1) == 20
	})
	plant.Reset()

	// 3rd run: put a plastic chip in the lock.
	plant.Emit(event.New("plastic-chip"))

	plant.AssertLast(0, func(evt *event.Event) bool {
		return evt.Topic() == "dunno"
	})
	plant.Reset()

	// 4th run: try a bad action.
	plant.Emit(event.New("screwdriver"))

	plant.AssertLast(0, func(evt *event.Event) bool {
		return evt.Topic() == "error" && evt.Payload().At("message").AsString("-") == "don't try to break me"
	})
}

//--------------------
// HELPERS
//--------------------

// cents retrieves the cents out of the payload of an event.
func payloadCents(evt *event.Event) int {
	return evt.Payload().At("cents").AsInt(-1)
}

// lockMachine will be unlocked if enough money is inserted. After
// that it can be locked again.
type lockMachine struct {
	id    string
	cents int
}

// Locked represents the locked state receiving coins.
func (m *lockMachine) Locked(emitter mesh.Emitter, evt *event.Event) behaviors.FSMStatus {
	switch evt.Topic() {
	case "check-cents":
		_ = emitter.Broadcast(event.New(
			"cents-checked",
			"id", m.id,
			"cents", m.cents,
		))
	case "info":
		_ = emitter.Broadcast(event.New(
			"status",
			"id", m.id,
			"status", "locked",
			"cents", m.cents,
		))
	case "coin":
		cents := payloadCents(evt)
		if cents < 1 {
			return behaviors.FSMStatus{"locked-error", nil, fmt.Errorf("do not insert buttons")}
		}
		m.cents += cents
		if m.cents > 100 {
			m.cents -= 100
			_ = emitter.Broadcast(event.New(
				"unlocked",
				"id", m.id,
				"status", "unlocked",
			))
			return behaviors.FSMStatus{"unlocked", m.Unlocked, nil}
		}
	case "press-button":
		if m.cents > 0 {
			_ = emitter.Broadcast(event.New(
				"coins-dropped",
				"id", m.id,
				"cents", m.cents,
			))
			m.cents = 0
		}
	case "screwdriver":
		_ = emitter.Broadcast(event.New(
			"error",
			"id", m.id,
			"message", "don't try to break me",
		))
		return behaviors.FSMStatus{"locked-error", nil, fmt.Errorf("don't try to break me")}
	default:
		_ = emitter.Broadcast(event.New(
			"dunno",
			"id", m.id,
		))
	}
	return behaviors.FSMStatus{"locked", m.Locked, nil}
}

// Unlocked represents the unlocked state receiving coins.
func (m *lockMachine) Unlocked(emitter mesh.Emitter, evt *event.Event) behaviors.FSMStatus {
	switch evt.Topic() {
	case "check-cents":
		_ = emitter.Broadcast(event.New(
			"cents-checked",
			"id", m.id,
			"cents", m.cents,
		))
	case "info":
		_ = emitter.Broadcast(event.New(
			"status",
			"id", m.id,
			"status", "unlocked",
			"cents", m.cents,
		))
	case "coin":
		cents := payloadCents(evt)
		_ = emitter.Broadcast(event.New(
			"coins-returned",
			"id", m.id,
			"cents", cents,
		))
	case "press-button":
		if m.cents > 0 {
			_ = emitter.Broadcast(event.New(
				"coins-dropped",
				"id", m.id,
				"cents", m.cents,
			))
			m.cents = 0
		}
		return behaviors.FSMStatus{"locked", m.Locked, nil}
	default:
		_ = emitter.Broadcast(event.New(
			"dunno",
			"id", m.id,
		))
	}
	return behaviors.FSMStatus{"unlocked", m.Unlocked, nil}
}

// EOF
