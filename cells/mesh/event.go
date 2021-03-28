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
	payload   json.RawMessage
}

// nilEvent is returned in case of errors.
var nilEvent = Event{topic: TopicNil}

// IsNilEvent checks if an event is the nil event.
func IsNilEvent(evt Event) bool {
	return evt.topic == TopicNil
}

// NewEvent creates a new Event based on a topic. The payloads are optional.
func NewEvent(topic string, payloads ...interface{}) (Event, error) {
	evt := nilEvent
	if topic == "" {
		return evt, fmt.Errorf("event needs topic")
	}
	evt.timestamp = time.Now().UTC()
	evt.topic = topic
	// Check if the only value is a payload.
	switch len(payloads) {
	case 0:
		return evt, nil
	case 1:
		bs, err := json.Marshal(payloads[0])
		if err != nil {
			return evt, fmt.Errorf("cannot marshal payload: %v", err)
		}
		evt.payload = bs
	default:
		bs, err := json.Marshal(payloads)
		if err != nil {
			return evt, fmt.Errorf("cannot marshal payload: %v", err)
		}
		evt.payload = bs
	}
	return evt, nil
}

// Timestamp returns the event timestamp.
func (evt Event) Timestamp() time.Time {
	return evt.timestamp
}

// Topic returns the event topic.
func (evt Event) Topic() string {
	return evt.topic
}

// HasPayload checks if the event contains a payload.
func (evt Event) HasPayload() bool {
	return evt.payload != nil
}

// Payload unmarshals the payload of the event.
func (evt Event) Payload(payload interface{}) error {
	if evt.payload == nil {
		return fmt.Errorf("Event contains no payload")
	}
	err := json.Unmarshal(evt.payload, payload)
	if err != nil {
		return fmt.Errorf("cannont unmarshal payload: %v", err)
	}
	return nil
}

// String implements fmt.Stringer.
func (evt Event) String() string {
	return fmt.Sprintf(
		"Event{Timestamp:%s Topic:%v Payload:%v}",
		evt.timestamp.Format(time.RFC3339Nano), evt.topic, string(evt.payload),
	)
}

// MarshalJSON implements the custom JSON marshaling of the event.
func (evt Event) MarshalJSON() ([]byte, error) {
	tmp := struct {
		Timestamp time.Time       `json:"timestamp"`
		Topic     string          `json:"topic"`
		Payload   json.RawMessage `json:"payload,omitempty"`
	}{
		Timestamp: evt.timestamp,
		Topic:     evt.topic,
		Payload:   evt.payload,
	}
	return json.Marshal(tmp)
}

// UnmarshalJSON implements the custom JSON unmarshaling of the event.
func (evt *Event) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Timestamp time.Time       `json:"timestamp"`
		Topic     string          `json:"topic"`
		Payload   json.RawMessage `json:"payload,omitempty"`
	}{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	evt.timestamp = tmp.Timestamp
	evt.topic = tmp.Topic
	evt.payload = tmp.Payload
	return nil
}

// EOF
