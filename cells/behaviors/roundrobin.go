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
	"sort"

	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// ROUND ROBIN BEHAVIOR
//--------------------

// roundRobinBehavior emit the received events round robin to its
// subscribers in a very simple way.
type roundRobinBehavior struct {
	id      string
	emitter mesh.Emitter
	current int
}

// NewRoundRobinBehavior creates a behavior emitting the received events to
// its subscribers in a very simple way. Subscriptions or unsubscriptions
// during runtime may influence the order.
func NewRoundRobinBehavior(id string) mesh.Behavior {
	return &roundRobinBehavior{
		id: id,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *roundRobinBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *roundRobinBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *roundRobinBehavior) Terminate() error {
	return nil
}

// Process emits the event round robin to the subscribers.
func (b *roundRobinBehavior) Process(evt *event.Event) error {
	subscribers := b.emitter.Subscribers()
	ls := len(subscribers)
	if ls == 0 {
		return nil
	}
	if b.current >= ls {
		b.current = 0
	}
	sort.Strings(subscribers)
	b.emitter.Emit(subscribers[b.current], evt)
	b.current++
	return nil
}

// Recover from an error.
func (b *roundRobinBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
