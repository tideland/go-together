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

// AggregatorFunc is a function receiving the current aggregated payload
// and event and returns the next aggregated payload.
type AggregatorFunc func(aggregate interface{}, evt *mesh.Event) (interface{}, error)

// aggregatorBehavior implements the aggregator behavior.
type aggregatorBehavior struct {
	initial    interface{}
	aggregate  interface{}
	aggregator AggregatorFunc
}

// NewAggregatorBehavior creates a behavior aggregating the received events
// and emits events with the new aggregate. A "reset!" topic resets the
// aggregate to nil again.
func NewAggregatorBehavior(aggregate interface{}, aggregator AggregatorFunc) mesh.Behavior {
	return &aggregatorBehavior{
		initial:    aggregate,
		aggregate:  aggregate,
		aggregator: aggregator,
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
			case TopicAggregate:
				// Request aggregated.
				if err := out.Emit(TopicAggregated, b.aggregate); err != nil {
					return err
				}
			case TopicReset:
				// Reset to initial value.
				b.aggregate = b.initial
				if err := out.Emit(TopicResetted); err != nil {
					return err
				}
			default:
				// Aggregate the event.
				aggregate, err := b.aggregator(b.aggregate, evt)
				if err != nil {
					return err
				}
				b.aggregate = aggregate
			}
		}
	}
}

// EOF
