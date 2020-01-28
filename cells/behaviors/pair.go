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

	"tideland.dev/go/dsa/identifier"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
	"tideland.dev/go/together/loop"
	"tideland.dev/go/trace/failure"
)

//--------------------
// PAIR BEHAVIOR
//--------------------

// PairCriterion is used by the pair behavior and has to return true, if
// the passed event matches the wanted criterion. While there has been no
// match the first event will be nil and the second the one to test. When
// returning true the found one will be passed as first one so that the
// function can compare.
type PairCriterion func(first, second *event.Event) bool

// pairBehavior checks if events occur in pairs.
type pairBehavior struct {
	id          string
	emitter     mesh.Emitter
	matches     PairCriterion
	timespan    time.Duration
	tickerTopic string
	loop        *loop.Loop
	first       *event.Event
}

// NewPairBehavior creates a behavior checking if two events match a criterion
// defined by the PairCriterion function and the timeout between them is not
// longer than the passed timespan. In case of a positive pair match an according
// event containing both timestamps and both returned datas is emitted. In case
// of a timeout a timeout event is emitted. It's payload is the first timestamp,
// the first data, and the timestamp of the timeout.
func NewPairBehavior(id string, criterion PairCriterion, timespan time.Duration) mesh.Behavior {
	return &pairBehavior{
		id:          id,
		matches:     criterion,
		timespan:    timespan,
		tickerTopic: identifier.NewUUID().String(),
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *pairBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *pairBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	l, err := loop.Go(b.worker)
	if err != nil {
		return failure.Annotate(err, "init pair behavior")
	}
	b.loop = l
	return nil
}

// Terminate the behavior.
func (b *pairBehavior) Terminate() error {
	return b.loop.Stop()
}

// Process evaluates the received events for matching before timeout.
func (b *pairBehavior) Process(evt *event.Event) {
	switch evt.Topic() {
	case b.tickerTopic:
		if b.first == nil {
			return
		}
		if b.first.Timestamp().Add(b.timespan).Before(evt.Timestamp()) {
			// Timeout!
			b.emitTimeout(evt)
		}
	default:
		if b.matches(b.first, evt) {
			if b.first == nil {
				// First match.
				b.first = evt
				return
			}
			// Second match.
			if b.first.Timestamp().Add(b.timespan).Before(evt.Timestamp()) {
				// Timeout!
				b.emitTimeout(evt)
			}
			b.emitPair(evt)
		}
	}
}

// Recover implements the cells.Behavior interface.
func (b *pairBehavior) Recover(err interface{}) error {
	return nil
}

// worker is the sending a cleanup tick to itself.
func (b *pairBehavior) worker(lt loop.Terminator) error {
	ticker := time.NewTicker(b.timespan)
	defer ticker.Stop()
	for {
		select {
		case <-lt.Done():
			return nil
		case <-ticker.C:
			b.emitter.Self(event.New(b.tickerTopic))
		}
	}
}

// emitPair emits the event for a successful pair.
func (b *pairBehavior) emitPair(evt *event.Event) {
	_ = b.emitter.Broadcast(event.New(
		TopicPair,
		"first", b.first,
		"second", evt,
	))
	b.first = nil
}

// emitTimeout emits the event for a pairing timeout.
func (b *pairBehavior) emitTimeout(evt *event.Event) {
	_ = b.emitter.Broadcast(event.New(
		TopicPairTimeout,
		"first", b.first,
		"timeout", evt,
	))
	b.first = nil
}

// EOF
