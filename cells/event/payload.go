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
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"tideland.dev/go/trace/failure"
)

//--------------------
// PAYLOAD CHANNEL
//--------------------

// PayloadChan is intended to be sent with an event as payload
// so that a behavior can use it to answer a request.
type PayloadChan chan *Payload

// Wait waits for a returned payload until receiving or timeout.
func (plc PayloadChan) Wait(timeout time.Duration) (*Payload, error) {
	select {
	case pl, ok := <-plc:
		if !ok {
			return nil, failure.New("payload channel has been closed")
		}
		return pl, nil
	case <-time.After(timeout):
		return nil, failure.New("no returned payload until timeout")
	}
}

//--------------------
// PAYLOAD
//--------------------

// Payload contains key/value pairs of data.
type Payload struct {
	values map[string]interface{}
	replyc PayloadChan
}

// NewPayload creates a payload with the given pairs of keys
// and values. In case of payloads or maps as keys thos will
// be merged, in case of arrays or slices those will be merged
// with the index as string. In case of those as values they
// will all be nested payload values.
func NewPayload(kvs ...interface{}) *Payload {
	pl := &Payload{
		values: map[string]interface{}{},
	}
	pl.setKeyValues(kvs...)
	return pl
}

// NewReplyPayload creates a new payload like NewPayload() but
// also a reply channel.
func NewReplyPayload(kvs ...interface{}) (*Payload, PayloadChan) {
	pl := &Payload{
		values: map[string]interface{}{},
		replyc: make(PayloadChan, 1),
	}
	pl.setKeyValues(kvs...)
	return pl, pl.replyc
}

// Keys returns the keys of the payload.
func (pl *Payload) Keys() []string {
	var keys []string
	for key := range pl.values {
		keys = append(keys, key)
	}
	return keys
}

// At returns the value at the given key. This value may
// be empty.
func (pl *Payload) At(keys ...string) *Value {
	accessError := func(msg string, vs ...interface{}) *Value {
		return &Value{
			err: failure.New(msg, vs...),
		}
	}
	// Analyse values at keys.
	switch len(keys) {
	case 0:
		return accessError("no key passed")
	case 1:
		if v, ok := pl.values[keys[0]]; ok {
			return &Value{
				raw: v,
			}
		}
		return accessError("no payload value at key %q", keys[0])
	default:
		if _, ok := pl.values[keys[0]]; ok {
			if npl, ok := pl.values[keys[0]].(*Payload); ok {
				return npl.At(keys[1:]...)
			}
			return accessError("value at key %q is no payload", keys[0])
		}
		return accessError("no payload value at key %q", keys[0])
	}
}

// Do performs a function for all key/value pairs.
func (pl *Payload) Do(f func(key string, value *Value) error) error {
	var errs []error
	for key, rawValue := range pl.values {
		value := &Value{
			raw: rawValue,
		}
		errs = append(errs, f(key, value))
	}
	return failure.Collect(errs...)
}

// Reply allows the receiver of a payload to reply via a channel.
func (pl *Payload) Reply(rpl *Payload) error {
	if pl.replyc == nil {
		return failure.New("payload contains no reply channel")
	}
	select {
	case pl.replyc <- rpl:
		return nil
	default:
		return failure.New("payload reply channel is closed")
	}
}

// Clone creates a new payload with the content of the current one and
// applies the given changes.
func (pl *Payload) Clone(kvs ...interface{}) *Payload {
	plc := &Payload{
		values: map[string]interface{}{},
	}
	for key, value := range pl.values {
		plc.values[key] = value
	}
	plc.setKeyValues(kvs...)
	return plc
}

// Len returns the number of values of the payload.
func (pl *Payload) Len() int {
	return len(pl.values)
}

// String implements fmt.Stringer.
func (pl *Payload) String() string {
	var kvs []string
	for key, value := range pl.values {
		kvs = append(kvs, fmt.Sprintf("%s:%v", key, value))
	}
	return fmt.Sprintf("Payload[%s]", strings.Join(kvs, " "))
}

// setKeyValues iterates over the key/value values and adds
// them to the payloads values.
func (pl *Payload) setKeyValues(kvs ...interface{}) {
	var key string
	for i, kv := range kvs {
		if i%2 == 0 {
			// Talking about a key.
			if plk, ok := kv.(*Payload); ok {
				// A payload, merge values.
				pl.mergeMap(plk.values)
				pl.replyc = plk.replyc
				continue
			}
			switch reflect.TypeOf(kv).Kind() {
			case reflect.Map:
				// A map, merge it.
				pl.mergeMap(kv)
			case reflect.Array, reflect.Slice:
				// An iteratable, merge it.
				pl.mergeIteratable(kv)
			default:
				// Any other key.
				key = fmt.Sprintf("%v", kv)
				pl.values[key] = DefaultValue
			}
			continue
		}
		// Talking about a value.
		switch reflect.TypeOf(kv).Kind() {
		case reflect.Map, reflect.Array, reflect.Slice:
			// Nest it.
			pl.values[key] = NewPayload(kv)
		default:
			// Add it.
			pl.values[key] = kv
		}
	}
}

// mergeMap maps any map into the own values.
func (pl *Payload) mergeMap(m interface{}) {
	var kvs []interface{}
	iter := reflect.ValueOf(m).MapRange()
	for iter.Next() {
		key := iter.Key().Interface()
		value := iter.Value().Interface()
		kvs = append(kvs, key, value)
	}
	pl.setKeyValues(kvs...)
}

// mergeIteratable maps arrays and slices into the
// own values.
func (pl *Payload) mergeIteratable(i interface{}) {
	var kvs []interface{}
	v := reflect.ValueOf(i)
	l := v.Len()
	for i := 0; i < l; i++ {
		key := strconv.Itoa(i)
		value := v.Index(i).Interface()
		kvs = append(kvs, key, value)
	}
	pl.setKeyValues(kvs...)
}

// EOF
