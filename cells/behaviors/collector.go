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

import "tideland.dev/go/together/cells/mesh"

//--------------------
// COLLECTOR BEHAVIOR
//--------------------

// CollectionProcessorFunc is used to process collected events.
type CollectionProcessorFunc func(r mesh.EventSinkReader) (mesh.Event, error)

// collectorBehavior collects events for processing on demand.
type collectorBehavior struct {
	max     int
	sink    *mesh.EventSink
	process CollectionProcessorFunc
}

// NewCollectorBehavior collects max events. After "process!" topic it processes
// it and emits the result as event. After "reset!" topic the collection is dropped
// to zero.
func NewCollectorBehavior(max int, process CollectionProcessorFunc) mesh.Behavior {
	return &collectorBehavior{
		max:     max,
		sink:    mesh.NewEventSink(max),
		process: process,
	}
}

// Go collects, processes, and resets the collected events.
func (b *collectorBehavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return cell.Context().Err()
		case evt := <-in.Pull():
			switch evt.Topic() {
			case TopicReset:
				b.sink.Clear()
				out.Emit(TopicResetted)
			case TopicProcess:
				pevt, err := b.process(b.sink)
				if err != nil {
					return err
				}
				out.EmitEvent(pevt)
				b.sink.Clear()
			default:
				b.sink.Push(evt)
			}
		}
	}
}

// EOF
