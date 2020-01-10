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
// SEQUENCE BEHAVIOR
//--------------------

// SequenceCriterion is used by the sequence behavior. It has to return
// CriterionDone when a sequence is complete, CriterionKeep when it is
// so far okay but not complete, and CriterionClear when the sequence
// doesn't match and has to be cleared.
type SequenceCriterion func(accessor event.SinkAccessor) event.CriterionMatch

// sequenceBehavior implements the sequence behavior.
type sequenceBehavior struct {
	id      string
	emitter mesh.Emitter
	matches SequenceCriterion
	process event.SinkProcessor
	sink    *event.Sink
}

// NewSequenceBehavior creates an event sequence behavior. It checks the
// event stream for a sequence defined by the criterion. In this case an
// event containing the sequence is emitted.
func NewSequenceBehavior(id string, criterion SequenceCriterion, processor event.SinkProcessor) mesh.Behavior {
	return &sequenceBehavior{
		id:      id,
		matches: criterion,
		process: processor,
		sink:    event.NewSink(0),
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *sequenceBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *sequenceBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *sequenceBehavior) Terminate() error {
	b.sink.Clear()
	return nil
}

// Process ...
func (b *sequenceBehavior) Process(evt *event.Event) error {
	switch evt.Topic() {
	case event.TopicReset:
		b.sink.Clear()
	default:
		b.sink.Push(evt)
		matches := b.matches(b.sink)
		switch matches {
		case event.CriterionDone:
			// All done, process and start over.
			pl, err := b.process(b.sink)
			if err != nil {
				return err
			}
			b.emitter.Broadcast(event.New(TopicSequence, pl))
			b.sink = event.NewSink(0)
		case event.CriterionKeep:
			// So far ok.
		default:
			// Have to start from beginning.
			b.sink.Clear()
		}
	}
	return nil
}

// Recover implements the cells.Behavior interface.
func (b *sequenceBehavior) Recover(err interface{}) error {
	b.sink.Clear()
	return nil
}

// EOF
