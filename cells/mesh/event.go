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
	"fmt"
	"strings"
	"time"
)

//--------------------
// COPYABLE
//--------------------

// Copyable defines the interface for event values which are only
// defined by being copyable.
type Copyable interface {
	// Copy has to be implemented by the concrete value to ensure
	// that only copies will be transported.
	Copy() Copyable
}

// CopyableFunc allows to pass a function as value inside an event.
type CopyableFunc func(arg interface{}) error

// Exec executes the function of the copyable func.
func (cf CopyableFunc) Exec(arg interface{}) error {
	return cf(arg)
}

// Copy implements Copyable.
func (cf CopyableFunc) Copy() Copyable {
	return cf
}

// IsCopyableFunc checks if an arument is a copyable function.
func IsCopyableFunc(v interface{}) (CopyableFunc, bool) {
	cf, ok := v.(CopyableFunc)
	return cf, ok
}

//--------------------
// PAYLOAD
//--------------------

// Payload contains a number of key/value pairs and implements
// the Copyable interface.
type Payload struct {
	kvs map[string]interface{}
}

// NewPayload creates a payload with the given pairs of keys
// and values. In case of payloads or maps as keys thos will
// be merged, in case of arrays or slices those will be merged
// with the index as string. In case of those as values they
// will all be nested payload values.
func NewPayload(kvs ...interface{}) *Payload {
	pl := &Payload{
		kvs: map[string]interface{}{},
	}
	pl.parseKVS(kvs...)
	return pl
}

// HasValue checks if the payload values contain one with
// the given key.
func (pl *Payload) HasValue(key string) bool {
	_, ok := pl.kvs[key]
	return ok
}

// StringAt returns a value if it is a string.
func (pl *Payload) StringAt(key string) (string, bool) {
	tmp, ok := pl.kvs[key]
	if !ok {
		return "", false
	}
	value, ok := tmp.(string)
	return value, ok
}

// IntAt returns a value if it is an int.
func (pl *Payload) IntAt(key string) (int, bool) {
	tmp, ok := pl.kvs[key]
	if !ok {
		return 0, false
	}
	value, ok := tmp.(int)
	return value, ok
}

// Float64At returns a value if it is a float64.
func (pl *Payload) Float64At(key string) (float64, bool) {
	tmp, ok := pl.kvs[key]
	if !ok {
		return 0.0, false
	}
	value, ok := tmp.(float64)
	return value, ok
}

// BoolAt returns a value if it is a bool.
func (pl *Payload) BoolAt(key string) (bool, bool) {
	tmp, ok := pl.kvs[key]
	if !ok {
		return false, false
	}
	value, ok := tmp.(bool)
	return value, ok
}

// CopyableAt returns a value if it is a Copyable.
func (pl *Payload) CopyableAt(key string) (Copyable, bool) {
	tmp, ok := pl.kvs[key]
	if !ok {
		return nil, false
	}
	value, ok := tmp.(Copyable)
	if !ok {
		return nil, false
	}
	return value.Copy(), true
}

// PayloadAt returns a value if it is a Payload.
func (pl *Payload) PayloadAt(key string) (*Payload, bool) {
	tmp, ok := pl.kvs[key]
	if !ok {
		return nil, false
	}
	value, ok := tmp.(*Payload)
	if !ok {
		return nil, false
	}
	return value.Copy().(*Payload), true
}

// Len returns the number of values.
func (pl *Payload) Len() int {
	return len(pl.kvs)
}

// Do performs the given function for all keys and values.
func (pl *Payload) Do(f func(key string, value interface{})) {
	for key, tmp := range pl.kvs {
		switch value := tmp.(type) {
		case Copyable:
			f(key, value.Copy())
		default:
			f(key, tmp)
		}
	}
}

// String implements fmt.Stringer.
func (pl *Payload) String() string {
	var kvss []string
	for key, value := range pl.kvs {
		kvss = append(kvss, fmt.Sprintf("%s:%v", key, value))
	}
	tmpl := "Payload{Values:[%s]}"
	return fmt.Sprintf(tmpl, strings.Join(kvss, " "))
}

// Copy implements Copyable.
func (pl *Payload) Copy() Copyable {
	plc := &Payload{
		kvs: make(map[string]interface{}),
	}
	for k, v := range pl.kvs {
		vc, ok := v.(Copyable)
		if ok {
			plc.kvs[k] = vc.Copy()
		} else {
			plc.kvs[k] = v
		}
	}
	return plc
}

// parseKVS iterates over the key/value values and adds
// them to the payloads values.
func (pl *Payload) parseKVS(kvs ...interface{}) {
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
			pl.kvs[key] = true
			continue
		}
		// Talking about a value.
		switch tmp := kv.(type) {
		case string, int, float64, bool, Copyable:
			pl.kvs[key] = tmp
		case fmt.Stringer:
			pl.kvs[key] = tmp.String()
		default:
			pl.kvs[key] = fmt.Sprintf("%v", kv)
		}
	}
}

//--------------------
// EVENT
//--------------------

// Event transports a topic and a payload a cell can process.
type Event struct {
	timestamp time.Time
	topic     string
	payload   *Payload
}

// NewEvent creates a new event based on a topic. Additionally
// a set of keys and values can be added. The given values will
// be interpreted as alternating keys and values. Keys will
// be converted to strings if needed, also invalid values. A
// final key without a value will be set to true.
func NewEvent(topic string, kvs ...interface{}) *Event {
	evt := &Event{
		timestamp: time.Now().UTC(),
		topic:     topic,
		payload:   NewPayload(kvs...),
	}
	return evt
}

// Timestamp returns the event timestamp.
func (evt *Event) Timestamp() time.Time {
	return evt.timestamp
}

// Topic returns the event topic.
func (evt *Event) Topic() string {
	return evt.topic
}

// HasValue checks if the event values contain one with
// the given key.
func (evt *Event) HasValue(key string) bool {
	return evt.payload.HasValue(key)
}

// StringAt returns a value if it is a string.
func (evt *Event) StringAt(key string) (string, bool) {
	return evt.payload.StringAt(key)
}

// IntAt returns a value if it is an int.
func (evt *Event) IntAt(key string) (int, bool) {
	return evt.payload.IntAt(key)
}

// Float64At returns a value if it is a float64.
func (evt *Event) Float64At(key string) (float64, bool) {
	return evt.payload.Float64At(key)
}

// BoolAt returns a value if it is a bool.
func (evt *Event) BoolAt(key string) (bool, bool) {
	return evt.payload.BoolAt(key)
}

// CopyableAt returns a value if it is a Copyable.
func (evt *Event) CopyableAt(key string) (Copyable, bool) {
	return evt.payload.CopyableAt(key)
}

// PayloadAt returns a value if it is a Payload.
func (evt *Event) PayloadAt(key string) (*Payload, bool) {
	return evt.PayloadAt(key)
}

// Len returns the number of values.
func (evt *Event) Len() int {
	return evt.payload.Len()
}

// Do performs the given function for all keys and values.
func (evt *Event) Do(f func(key string, value interface{})) {
	evt.payload.Do(f)
}

// String implements fmt.Stringer.
func (evt *Event) String() string {
	return fmt.Sprintf("Event{Topic:%v %v}", evt.topic, evt.payload)
}

// EOF
