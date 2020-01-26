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
	"tideland.dev/go/together/fuse"
)

//--------------------
// AGGREGATOR BEHAVIOR
//--------------------

// Aggregator is a function receiving the current aggregated payload
// and event and returns the next aggregated payload.
type Aggregator func(p *event.Payload, evt *event.Event) (*event.Payload, error)

// aggregatorBehavior implements the aggregator behavior.
type aggregatorBehavior struct {
	id        string
	emitter   mesh.Emitter
	payload   *event.Payload
	aggregate Aggregator
}

// NewAggregatorBehavior creates a behavior aggregating the received events
// and emits events with the new aggregate. A "reset!" topic resets the
// aggregate to nil again.
func NewAggregatorBehavior(id string, pl *event.Payload, aggregator Aggregator) mesh.Behavior {
	return &aggregatorBehavior{
		id:        id,
		payload:   pl,
		aggregate: aggregator,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *aggregatorBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *aggregatorBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *aggregatorBehavior) Terminate() error {
	return nil
}

// Process aggregates the event.
func (b *aggregatorBehavior) Process(evt *event.Event) {
	switch evt.Topic() {
	case event.TopicStatus:
		fuse.Trigger(evt.Payload().Reply(b.payload))
	case event.TopicReset:
		b.payload = nil
	default:
		pl, err := b.aggregate(b.payload, evt)
		fuse.Trigger(err)
		b.payload = pl
		fuse.Trigger(b.emitter.Broadcast(event.New(TopicAggregated, pl)))
	}
}

// Recover from an error.
func (b *aggregatorBehavior) Recover(err interface{}) error {
	println("recovering of aggregator called")
	b.payload = event.NewPayload()
	return nil
}

// EOF
