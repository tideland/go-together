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

// RateCriterion is used by the rate behavior and has to return true, if
// the passed event matches a criterion for rate measuring.
type RateCriterion func(evt *event.Event) (bool, error)

// rateBehavior calculates the average rate of events matching a criterion.
type rateBehavior struct {
	id        string
	emitter   mesh.Emitter
	matches   RateCriterion
	count     int
	last      time.Time
	durations []time.Duration
}

// NewRateBehavior creates an even rate measuiring behavior. Each time the
// criterion function returns true for a received event the duration between
// this and the last one is calculated and emitted together with the timestamp.
// Additionally a moving average, lowest, and highest duration is calculated
// and emitted too. A "reset!" as topic resets the stored values.
func NewRateBehavior(id string, matches RateCriterion, count int) mesh.Behavior {
	return &rateBehavior{
		id:        id,
		matches:   matches,
		count:     count,
		last:      time.Now(),
		durations: []time.Duration{},
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *rateBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *rateBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *rateBehavior) Terminate() error {
	return nil
}

// Process calculates the rate of matching events.
func (b *rateBehavior) Process(evt *event.Event) error {
	switch evt.Topic() {
	case event.TopicReset:
		b.last = time.Now()
		b.durations = []time.Duration{}
	default:
		ok, err := b.matches(evt)
		if err != nil {
			return err
		}
		if ok {
			current := evt.Timestamp()
			duration := current.Sub(b.last)
			b.last = current
			b.durations = append(b.durations, duration)
			if len(b.durations) > b.count {
				b.durations = b.durations[1:]
			}
			total := 0 * time.Nanosecond
			low := 0x7FFFFFFFFFFFFFFF * time.Nanosecond
			high := 0 * time.Nanosecond
			for _, d := range b.durations {
				total += d
				if d < low {
					low = d
				}
				if d > high {
					high = d
				}
			}
			avg := total / time.Duration(len(b.durations))
			return b.emitter.Broadcast(event.New(
				TopicRate,
				"time", current,
				"duration", duration,
				"high", high,
				"low", low,
				"average", avg,
			))
		}
	}
	return nil
}

// Recover implements the cells.Behavior interface.
func (b *rateBehavior) Recover(err interface{}) error {
	b.last = time.Now()
	b.durations = []time.Duration{}
	return nil
}

// EOF
