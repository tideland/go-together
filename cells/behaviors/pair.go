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
// PAIR BEHAVIOR
//--------------------

// PairCriterion is used by the pair behavior and has to return true, if
// the passed event matches the wanted criterion. The given payload first
// is an empty one, later the returned payload. This allows the matching
// routine to maintain a state. In case of a timeout it will be reset.
type PairCriterion func(evt *event.Event, pl *event.Payload) (*event.Payload, bool)

// pairBehavior checks if events occur in pairs.
type pairBehavior struct {
	id       string
	emitter  mesh.Emitter
	matches  PairCriterion
	timespan time.Duration
	payload  *event.Payload
	matched  *time.Time
}

// NewPairBehavior creates a behavior checking if two events match a criterion
// defined by the PairCriterion function and the timeout between them is not
// longer than the passed timespan. In case of a positive pair match an according
// event containing both timestamps and both returned datas is emitted. In case
// of a timeout a timeout event is emitted. It's payload is the first timestamp,
// the first data, and the timestamp of the timeout.
func NewPairBehavior(id string, criterion PairCriterion, timespan time.Duration) mesh.Behavior {
	return &pairBehavior{
		id:       id,
		matches:  criterion,
		timespan: timespan,
		payload:  event.NewPayload(),
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *pairBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *pairBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *pairBehavior) Terminate() error {
	return nil
}

// Process evaluates the received events for matching before timeout.
func (b *pairBehavior) Process(evt *event.Event) error {
	ok := false
	timestamp := evt.Timestamp()
	b.payload, ok = b.matches(evt, b.payload)
	if b.payload == nil {
		// Criterion reset payload.
		b.payload = event.NewPayload()
	}
	if !ok {
		// Nothing to see, go on.
		return nil
	}
	if b.matched == nil {
		// First match.
		b.matched = &timestamp
		return nil
	}
	// Second match, check for timeout.
	timeout := b.matched.Add(b.timespan)
	if timeout.After(timestamp) {
		// Event in time.
		b.emitPair(timestamp)
	} else {
		// Sorry, too late.
		b.emitTimeout(timeout)
	}
	return nil
}

// Recover implements the cells.Behavior interface.
func (b *pairBehavior) Recover(err interface{}) error {
	return nil
}

// emitPair emits the event for a successful pair.
func (b *pairBehavior) emitPair(timestamp time.Time) {
	b.emitter.Broadcast(event.New(
		TopicPair,
		"first", *b.matched,
		"second", timestamp,
		"payload", b.payload,
	))
	b.payload = event.NewPayload()
	b.matched = nil
}

// emitTimeout emits the event for a pairing timeout.
func (b *pairBehavior) emitTimeout(timeout time.Time) {
	b.emitter.Broadcast(event.New(
		TopicPairTimeout,
		"payload", b.payload,
		"timeout", timeout,
	))
	b.payload = event.NewPayload()
	b.matched = nil
}

// EOF
