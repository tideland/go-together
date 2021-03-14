// Tideland Go Together - Cells - Behaviors
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// AGGREGATOR BEHAVIOR
//--------------------

// Aggregator is a function receiving the current aggregated payload
// and event and returns the next aggregated payload.
type Aggregator func(data *mesh.Payload, evt *mesh.Event) (*mesh.Payload, error)

// aggregatorBehavior implements the aggregator behavior.
type aggregatorBehavior struct {
	data      *mesh.Payload
	aggregate Aggregator
}

// NewAggregatorBehavior creates a behavior aggregating the received events
// and emits events with the new aggregate. A "reset!" topic resets the
// aggregate to nil again.
func NewAggregatorBehavior(data *mesh.Payload, aggregator Aggregator) mesh.Behavior {
	return &aggregatorBehavior{
		data:      data,
		aggregate: aggregator,
	}
}

// Go aggregates the event.
func (b *aggregatorBehavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			switch evt.Topic() {
			case TopicReset:
				b.data = nil
				out.Emit(mesh.NewEvent(TopicResetted))
			default:
				data, err := b.aggregate(b.data, evt)
				if err != nil {
					return err
				}
				b.data = data
				if err := out.Emit(mesh.NewEvent(TopicAggregated, b.data)); err != nil {
					return err
				}
			}
		}
	}
}

// EOF
