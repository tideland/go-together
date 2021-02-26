// Tideland Go Together - Cells - Mesh
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"strings"
	"time"
)

//--------------------
// PAYLOAD
//--------------------

// Payload defines the interface for event payloads which are only defined
// by the property to be deep copyable.
type Payload interface {
	// DeepCopy has to be implemented by the concrete payload to ensure
	// that only copies will be transported.
	DeepCopy() Payload
}

// emptyPayload is used if an event with payload nil is created.
type emptyPayload struct{}

// DeepCopy implements Payload.
func (ep emptyPayload) DeepCopy() Payload {
	return ep
}

// KeyValuePayload implements payload and contains a number of key/value
// pairs where the values having the types string, int, float64, bool, or
// Payload.
type KeyValuePayload struct {
	keyValues map[string]interface{}
}

// IsKeyValuePayload checks if a payload us a key/value payload.
func IsKeyValuePayload(p Payload) (*KeyValuePayload, bool) {
	kvp, ok := p.(*KeyValuePayload)
	return kvp, ok
}

// NewKeyValuePayload creates a key/value payload by parsing the passed
// keys and values. Those are alternating keys and values. Keys will
// be converted to strings if needed, also invalid values. A final key
// will be set to true.
func NewKeyValuePayload(kvs ...interface{}) *KeyValuePayload {
	kvp := &KeyValuePayload{
		keyValues: map[string]interface{}{},
	}
	kvp.setKeyValues(kvs...)
	return kvp
}

// Has check if the key/value payload contains the passed key.
func (kvp *KeyValuePayload) Has(key string) bool {
	_, ok := kvp.keyValues[key]
	return ok
}

// At returns the value of the key. If there's no value nil will be
// returned.
func (kvp *KeyValuePayload) At(key string) interface{} {
	tmp, ok := kvp.keyValues[key]
	if !ok {
		return nil
	}
	switch value := tmp.(type) {
	case Payload:
		return value.DeepCopy()
	default:
		return tmp
	}
}

// StringAt returns a value if it is a string.
func (kvp *KeyValuePayload) StringAt(key string) (string, bool) {
	tmp, ok := kvp.keyValues[key]
	if !ok {
		return "", false
	}
	value, ok := tmp.(string)
	return value, ok
}

// IntAt returns a value if it is an int.
func (kvp *KeyValuePayload) IntAt(key string) (int, bool) {
	tmp, ok := kvp.keyValues[key]
	if !ok {
		return 0, false
	}
	value, ok := tmp.(int)
	return value, ok
}

// Float64At returns a value if it is a float64.
func (kvp *KeyValuePayload) Float64At(key string) (float64, bool) {
	tmp, ok := kvp.keyValues[key]
	if !ok {
		return 0.0, false
	}
	value, ok := tmp.(float64)
	return value, ok
}

// BoolAt returns a value if it is a bool.
func (kvp *KeyValuePayload) BoolAt(key string) (bool, bool) {
	tmp, ok := kvp.keyValues[key]
	if !ok {
		return false, false
	}
	value, ok := tmp.(bool)
	return value, ok
}

// PayloadAt returns a value if it is a payload.
func (kvp *KeyValuePayload) PayloadAt(key string) (Payload, bool) {
	tmp, ok := kvp.keyValues[key]
	if !ok {
		return nil, false
	}
	value, ok := tmp.(Payload)
	return value, ok
}

// Do performs the given function for all keys and values.
func (kvp *KeyValuePayload) Do(f func(key string, value interface{})) {
	for key, tmp := range kvp.keyValues {
		switch value := tmp.(type) {
		case Payload:
			f(key, value.DeepCopy())
		default:
			f(key, tmp)
		}
	}
}

// Len returns the number of values of the payload.
func (kvp *KeyValuePayload) Len() int {
	return len(kvp.keyValues)
}

// String implements fmt.Stringer.
func (kvp *KeyValuePayload) String() string {
	var kvs []string
	for key, value := range kvp.keyValues {
		kvs = append(kvs, fmt.Sprintf("%s:%v", key, value))
	}
	return fmt.Sprintf("Payload[%s]", strings.Join(kvs, " "))
}

// DeepCopy implements Payload.
func (kvp *KeyValuePayload) DeepCopy() Payload {
	dc := &KeyValuePayload{
		keyValues: map[string]interface{}{},
	}
	for key, tmp := range kvp.keyValues {
		switch value := tmp.(type) {
		case Payload:
			dc.keyValues[key] = value.DeepCopy()
		default:
			dc.keyValues[key] = value
		}
	}
	return dc
}

// setKeyValues iterates over the key/value values and adds
// them to the payloads values.
func (kvp *KeyValuePayload) setKeyValues(kvs ...interface{}) {
	var key string
	for i, kv := range kvs {
		if i%2 == 0 {
			// Talking about a key.
			switch tmp := kv.(type) {
			case string:
				key = tmp
			case fmt.Stringer:
				key = tmp.String()
			default:
				key = fmt.Sprintf("%v", kv)
			}
			// Preset with a default key.
			kvp.keyValues[key] = true
			continue
		}
		// Talking about a value.
		switch tmp := kv.(type) {
		case string, int, float64, bool, Payload:
			kvp.keyValues[key] = tmp
		case fmt.Stringer:
			kvp.keyValues[key] = tmp.String()
		default:
			kvp.keyValues[key] = fmt.Sprintf("%v", kv)
		}
	}
}

// FuncPayload allows to pass a function as payload inside an event.
type FuncPayload func(arg Payload) error

// IsFuncPayload checks if a payload is a function payload.
func IsFuncPayload(p Payload) (FuncPayload, bool) {
	fp, ok := p.(FuncPayload)
	return fp, ok
}

// NewFuncPayload creates a function payload out of the given function.
func NewFuncPayload(f func(arg Payload) error) FuncPayload {
	return FuncPayload(f)
}

// Exec executes the function of the fuction payload.
func (fp FuncPayload) Exec(arg Payload) error {
	return fp(arg)
}

// DeepCopy implements Payload.
func (fp FuncPayload) DeepCopy() Payload {
	return fp
}

//--------------------
// EVENT
//--------------------

// Event transports a topic and a payload a cell can process.
type Event struct {
	timestamp time.Time
	topic     string
	payload   Payload
}

// NewEvent creates a new event based on topic and payload. The latter
// will be deep copied and so can be reused for other events too.
func NewEvent(topic string, payload Payload) *Event {
	var pc Payload
	if payload == nil {
		pc = emptyPayload{}
	} else {
		pc = payload.DeepCopy()
	}
	return &Event{
		timestamp: time.Now().UTC(),
		topic:     topic,
		payload:   pc,
	}
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
func (e *Event) Payload() Payload {
	return e.payload
}

// String implements fmt.Stringer.
func (e *Event) String() string {
	return fmt.Sprintf("Event{Topic: %v, Payload: %v}", e.topic, e.payload)
}

// EOF
