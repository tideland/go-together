// Tideland Go Together - Cells - Behaviors
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// COUNTER BEHAVIOR
//--------------------

// Counters analyzes the passed event and returns, which counters
// shall be incremented.
type Counters func(evt *event.Event) []string

// counterBehavior counts events based on the counter function.
type counterBehavior struct {
	id       string
	emitter  mesh.Emitter
	counters Counters
	values   map[string]int
}

// NewCounterBehavior creates a counter behavior based on the passed
// function. This function may increase, decrease, or set the counter
// values. Afterwards the counter values will be emitted. All values
// can be reset with the topic "reset!".
func NewCounterBehavior(id string, counters Counters) mesh.Behavior {
	return &counterBehavior{
		id:       id,
		counters: counters,
		values:   map[string]int{},
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *counterBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *counterBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *counterBehavior) Terminate() error {
	return nil
}

// Process counts the event for the return value of the counter func
// and emits this value.
func (b *counterBehavior) Process(evt *event.Event) error {
	switch evt.Topic() {
	case event.TopicStatus:
		return evt.Payload().Reply(event.NewPayload("counter", b.values))
	case event.TopicReset:
		b.values = map[string]int{}
	default:
		counters := b.counters(evt)
		for _, counter := range counters {
			b.values[counter]++
		}
		return b.emitter.Broadcast(event.New(event.TopicCounted, b.values))
	}
	return nil
}

// Recover from an error.
func (b *counterBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
