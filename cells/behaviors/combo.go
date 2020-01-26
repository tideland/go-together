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
// SEQUENCE BEHAVIOR
//--------------------

// ComboCriterion is used by the combo behavior. It has to return
// CriterionDone when a combination is complete, CriterionKeep when it
// is so far okay but not complete, CriterionDropFirst when the first
// event shall be dropped, CriterionDropLast when the last event shall
// be dropped, and CriterionClear when the collected events have
// to be cleared for starting over. In case of CriterionDone it
// additionally has to return a payload which will be emitted.
type ComboCriterion func(accessor event.SinkAccessor) (event.CriterionMatch, *event.Payload)

// comboBehavior implements the combo behavior.
type comboBehavior struct {
	id      string
	emitter mesh.Emitter
	matches ComboCriterion
	sink    *event.Sink
}

// NewComboBehavior creates an event sequence behavior. It checks the
// event stream for a combination of events defined by the criterion. In
// this case an event containing the combination is emitted.
func NewComboBehavior(id string, matcher ComboCriterion) mesh.Behavior {
	return &comboBehavior{
		id:      id,
		matches: matcher,
		sink:    event.NewSink(0),
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *comboBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *comboBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *comboBehavior) Terminate() error {
	return b.sink.Clear()
}

// Process matches events for a combination of criteria.
func (b *comboBehavior) Process(evt *event.Event) {
	switch evt.Topic() {
	case event.TopicReset:
		fuse.Trigger(b.sink.Clear())
	default:
		_, err := b.sink.Push(evt)
		fuse.Trigger(err)
		matches, pl := b.matches(b.sink)
		switch matches {
		case event.CriterionDone:
			// All done, emit and start over.
			fuse.Trigger(b.emitter.Broadcast(event.New(TopicComboComplete, pl)))
			b.sink = event.NewSink(0)
		case event.CriterionKeep:
			// So far ok.
		case event.CriterionDropFirst:
			// First event doesn't match.
			_, err := b.sink.PullFirst()
			fuse.Trigger(err)
		case event.CriterionDropLast:
			// Last event doesn't match.
			_, err := b.sink.PullLast()
			fuse.Trigger(err)
		default:
			// Have to start from beginning.
			fuse.Trigger(b.sink.Clear())
		}
	}
}

// Recover from an error.
func (b *comboBehavior) Recover(err interface{}) error {
	return b.sink.Clear()
}

// EOF
