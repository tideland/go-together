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
type CopyableFunc func(arg Copyable) error

// Exec executes the function of the copyable func.
func (cf CopyableFunc) Exec(arg Copyable) error {
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
// EVENT
//--------------------

// Event transports a topic and a payload a cell can process.
type Event struct {
	timestamp time.Time
	topic     string
	kvs       map[string]interface{}
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
		kvs:       make(map[string]interface{}),
	}
	evt.parseKVS(kvs...)
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
	_, ok := evt.kvs[key]
	return ok
}

// StringAt returns a value if it is a string.
func (evt *Event) StringAt(key string) (string, bool) {
	tmp, ok := evt.kvs[key]
	if !ok {
		return "", false
	}
	value, ok := tmp.(string)
	return value, ok
}

// IntAt returns a value if it is an int.
func (evt *Event) IntAt(key string) (int, bool) {
	tmp, ok := evt.kvs[key]
	if !ok {
		return 0, false
	}
	value, ok := tmp.(int)
	return value, ok
}

// Float64At returns a value if it is a float64.
func (evt *Event) Float64At(key string) (float64, bool) {
	tmp, ok := evt.kvs[key]
	if !ok {
		return 0.0, false
	}
	value, ok := tmp.(float64)
	return value, ok
}

// BoolAt returns a value if it is a bool.
func (evt *Event) BoolAt(key string) (bool, bool) {
	tmp, ok := evt.kvs[key]
	if !ok {
		return false, false
	}
	value, ok := tmp.(bool)
	return value, ok
}

// CopyableAt returns a value if it is a Copyable.
func (evt *Event) CopyableAt(key string) (Copyable, bool) {
	tmp, ok := evt.kvs[key]
	if !ok {
		return nil, false
	}
	value, ok := tmp.(Copyable)
	if !ok {
		return nil, false
	}
	return value.Copy(), true
}

// ValueLen returns the number of values.
func (evt *Event) ValueLen() int {
	return len(evt.kvs)
}

// ValuesDo performs the given function for all keys and values.
func (evt *Event) ValuesDo(f func(key string, value interface{})) {
	for key, tmp := range evt.kvs {
		switch value := tmp.(type) {
		case Copyable:
			f(key, value.Copy())
		default:
			f(key, tmp)
		}
	}
}

// String implements fmt.Stringer.
func (evt *Event) String() string {
	var kvss []string
	for key, value := range evt.kvs {
		kvss = append(kvss, fmt.Sprintf("%s:%v", key, value))
	}
	tmpl := "Event{Topic: %v Values:[%s]}"
	return fmt.Sprintf(tmpl, evt.topic, strings.Join(kvss, " "))
}

// parseKVS iterates over the key/value values and adds
// them to the payloads values.
func (evt *Event) parseKVS(kvs ...interface{}) {
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
			evt.kvs[key] = true
			continue
		}
		// Talking about a value.
		switch tmp := kv.(type) {
		case string, int, float64, bool, Copyable:
			evt.kvs[key] = tmp
		case fmt.Stringer:
			evt.kvs[key] = tmp.String()
		default:
			evt.kvs[key] = fmt.Sprintf("%v", kv)
		}
	}
}

// EOF
