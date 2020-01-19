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
	"time"

	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// RATE BEHAVIOR
//--------------------

// RateWindowCriterion is used by the rate window behavior and has to return
// true, if the passed event matches a criterion for rate window measuring.
type RateWindowCriterion func(evt *event.Event) (bool, error)

// rateWindowBehavior implements the rate window behavior.
type rateWindowBehavior struct {
	id       string
	emitter  mesh.Emitter
	sink     *event.Sink
	matches  RateWindowCriterion
	count    int
	duration time.Duration
	process  event.SinkProcessor
}

// NewRateWindowBehavior creates an event rate window behavior. It checks
// if an event matches the passed criterion. If count events match during
// duration the process function is called. Its returned payload is
// emitted as new event with topic "rate-window". A received "reset" as
// topic resets the collected matches.
func NewRateWindowBehavior(
	id string,
	matches RateWindowCriterion,
	count int,
	duration time.Duration,
	process event.SinkProcessor) mesh.Behavior {
	return &rateWindowBehavior{
		id:       id,
		sink:     event.NewSink(count),
		matches:  matches,
		count:    count,
		duration: duration,
		process:  process,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *rateWindowBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *rateWindowBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *rateWindowBehavior) Terminate() error {
	return nil
}

// Process implements the cells.Behavior interface.
func (b *rateWindowBehavior) Process(evt *event.Event) error {
	switch evt.Topic() {
	case event.TopicReset:
		return b.sink.Clear()
	default:
		ok, err := b.matches(evt)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if _, err = b.sink.Push(evt); err != nil {
			return err
		}
		if b.sink.Len() == b.count {
			// Got enough matches, check duration.
			first, _ := b.sink.PeekFirst()
			last, _ := b.sink.PeekLast()
			difference := last.Timestamp().Sub(first.Timestamp())
			if difference <= b.duration {
				// We've got a burst!
				pl, err := b.process(b.sink)
				if err != nil {
					return err
				}
				if err = b.emitter.Broadcast(event.New(TopicRateWindow, pl)); err != nil {
					return err
				}
			}
			if _, err := b.sink.PullFirst(); err != nil {
				return err
			}
		}
		return nil
	}
}

// Recover implements the cells.Behavior interface.
func (b *rateWindowBehavior) Recover(err interface{}) error {
	b.sink = event.NewSink(b.count)
	return nil
}

// EOF
