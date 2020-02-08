// Tideland Go Together - Cells - Event
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package event // import "tideland.dev/go/together/cells/event"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"time"
)

//--------------------
// CONSTANTS
//--------------------

// Standard topics.
const (
	TopicCollected = "collected"
	TopicCounted   = "counted"
	TopicProcess   = "process"
	TopicProcessed = "processed"
	TopicReset     = "reset"
	TopicResult    = "result"
	TopicStatus    = "status"
	TopicTick      = "tick"
)

//--------------------
// EVENT
//--------------------

// Event describes an event of the cells. It contains a topic as well as
// a possible number of key/value pairs as payload.
type Event struct {
	ctx       context.Context
	timestamp time.Time
	topic     string
	payload   *Payload
}

// New creates a new event. The arguments after the topic are taken
// to create a new payload.
func New(topic string, kvs ...interface{}) *Event {
	return WithContext(context.Background(), topic, kvs...)
}

// WithContext creates a new event containing a context allowing to
// cancel it. The arguments after the topic are taken to create a
// new payload.
func WithContext(ctx context.Context, topic string, kvs ...interface{}) *Event {
	return &Event{
		ctx:       ctx,
		timestamp: time.Now(),
		topic:     topic,
		payload:   NewPayload(kvs...),
	}
}

// Context returns the event context.
func (e *Event) Context() context.Context {
	return e.ctx
}

// Done tells if the event is done in the sense of the context. This
// happens if a potential context is canceled or reached timeout or
// deadline.
func (e *Event) Done() bool {
	return e.ctx.Err() != nil
}

// Timestamp returns the event timestamp.
func (e *Event) Timestamp() time.Time {
	return e.timestamp
}

// Topic returns the event topic.
func (e *Event) Topic() string {
	return e.topic
}

// Payload returns the event payload.
func (e *Event) Payload() *Payload {
	return e.payload
}

// String implements fmt.Stringer.
func (e *Event) String() string {
	return fmt.Sprintf("Event{Topic: %v, Payload: %v}", e.topic, e.payload)
}

// EOF
