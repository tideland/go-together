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
	"time"

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
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()

	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		eventInfos := []string{}
		assert.OK(accessor.Do(func(index int, evt *event.Event) error {
			eventInfos = append(eventInfos, evt.Topic())
			return nil
		}))
		sigc <- eventInfos
		return nil, nil
	}

	lockA := lockMachine{"a", 0}
	lockB := lockMachine{"b", 0}

	assert.OK(msh.SpawnCells(
		behaviors.NewFSMBehavior("lock-a", behaviors.FSMStatus{"locked", lockA.Locked, nil}),
		behaviors.NewFSMBehavior("lock-b", behaviors.FSMStatus{"locked", lockB.Locked, nil}),
		newRestorerBehavior("restorer"),
		behaviors.NewCollectorBehavior("collector-a", 10, processor),
		behaviors.NewCollectorBehavior("collector-b", 10, processor),
	))
	assert.OK(msh.Subscribe("lock-a", "restorer", "collector-a"))
	assert.OK(msh.Subscribe("lock-b", "restorer", "collector-b"))

	// 1st run: emit not enough and press button.
	assert.OK(msh.Emit("lock-a", event.New("coin", "cents", 20)))
	assert.OK(msh.Emit("lock-a", event.New("coin", "cents", 20)))
	assert.OK(msh.Emit("lock-a", event.New("coin", "cents", 20)))
	assert.OK(msh.Emit("lock-a", event.New("info")))
	assert.OK(msh.Emit("lock-a", event.New("press-button")))
	assert.OK(msh.Emit("lock-a", event.New("check-cents")))
	assert.OK(msh.Emit("restorer", event.New("grab")))

	time.Sleep(100 * time.Millisecond)

	assert.OK(msh.Emit("collector-a", event.New(event.TopicProcess)))
	assert.OK(msh.Emit("collector-a", event.New(event.TopicReset)))

	assert.Wait(sigc, []string{"status", "coins-dropped", "cents-checked"}, time.Second)

	// 2nd run: unlock the lock and lock it again.
	assert.OK(msh.Emit("lock-a", event.New("coin", "cents", 50)))
	assert.OK(msh.Emit("lock-a", event.New("coin", "cents", 20)))
	assert.OK(msh.Emit("lock-a", event.New("coin", "cents", 50)))
	assert.OK(msh.Emit("lock-a", event.New("info")))
	assert.OK(msh.Emit("lock-a", event.New("press-button")))

	time.Sleep(100 * time.Millisecond)

	assert.OK(msh.Emit("collector-a", event.New(event.TopicProcess)))
	assert.OK(msh.Emit("collector-a", event.New(event.TopicReset)))

	assert.Wait(sigc, []string{"unlocked", "status", "coins-dropped"}, time.Second)

	// 3rd run: put a plastic chip in the lock.
	assert.OK(msh.Emit("lock-a", event.New("plastic-chip")))

	time.Sleep(100 * time.Millisecond)

	assert.OK(msh.Emit("collector-a", event.New(event.TopicProcess)))
	assert.OK(msh.Emit("collector-a", event.New(event.TopicReset)))

	assert.Wait(sigc, []string{"dunno"}, time.Second)

	// 4th run: try a bad action.
	assert.OK(msh.Emit("lock-b", event.New("screwdriver")))

	time.Sleep(100 * time.Millisecond)

	assert.OK(msh.Emit("collector-b", event.New(event.TopicProcess)))
	assert.OK(msh.Emit("collector-b", event.New(event.TopicReset)))

	assert.Wait(sigc, []string{"error"}, time.Second)
	assert.OK(msh.Stop())
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
		))
		return behaviors.FSMStatus{evt.Topic(), nil, fmt.Errorf("don't try to break me")}
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

// restorerBehavior for test.
type restorerBehavior struct {
	id      string
	emitter mesh.Emitter
	cents   int
}

func newRestorerBehavior(id string) mesh.Behavior {
	return &restorerBehavior{
		id:    id,
		cents: 0,
	}
}

func (b *restorerBehavior) ID() string {
	return b.id
}

func (b *restorerBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

func (b *restorerBehavior) Terminate() error {
	return nil
}

func (b *restorerBehavior) Process(evt *event.Event) {
	switch evt.Topic() {
	case "grab-coins":
		_ = b.emitter.Broadcast(event.New("cents", "cents", b.cents))
		b.cents = 0
	case "drop-coins":
		b.cents += payloadCents(evt)
	}
}

func (b *restorerBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
