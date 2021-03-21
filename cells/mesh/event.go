// Tideland Go Together - Cells - Mesh
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/together/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/json"
	"fmt"
	"time"
)

//--------------------
// EVENT
//--------------------

// Event transports a topic and a payload a cell can process. The
// payload is anything marshalled into JSON and will be unmarshalled
// when a receiving cell accesses it.
type Event struct {
	timestamp time.Time
	topic     string
	payload   []byte
}

// NewEvent creates a new event based on a topic. The payloads are optional.
func NewEvent(topic string, payloads ...interface{}) (*Event, error) {
	evt := &Event{
		timestamp: time.Now().UTC(),
		topic:     topic,
	}
	// Check if the only value is a payload.
	switch len(payloads) {
	case 0:
		return evt, nil
	case 1:
		bs, err := json.Marshal(payloads[0])
		if err != nil {
			return nil, fmt.Errorf("cannot marshal payload: %v", err)
		}
		evt.payload = bs
	default:
		bs, err := json.Marshal(payloads)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal payload: %v", err)
		}
		evt.payload = bs
	}
	return evt, nil
}

// Timestamp returns the event timestamp.
func (evt *Event) Timestamp() time.Time {
	return evt.timestamp
}

// Topic returns the event topic.
func (evt *Event) Topic() string {
	return evt.topic
}

// HasPayload checks if the event contains a payload.
func (evt *Event) HasPayload() bool {
	return evt.payload != nil
}

// Payload unmarshals the payload of the event.
func (evt *Event) Payload(payload interface{}) error {
	if evt.payload == nil {
		return fmt.Errorf("event contains no payload")
	}
	err := json.Unmarshal(evt.payload, payload)
	if err != nil {
		return fmt.Errorf("cannont unmarshal payload: %v", err)
	}
	return nil
}

// String implements fmt.Stringer.
func (evt *Event) String() string {
	return fmt.Sprintf("Event{Topic:%v Payload:%v}", evt.topic, string(evt.payload))
}

// EOF
