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
	"errors"
	"sync"
	"time"
)

//--------------------
// INTERFACES
//--------------------

// Receptor defines the interface to receive events.
type Receptor interface {
	// Pull reads an event out of the input stream.
	Pull() <-chan Event
}

// Emitter defines the interface for emitting events to one
// or more cells.
type Emitter interface {
	// Emit creates a new event and appends it to the output stream.
	Emit(topic string, payloads ...interface{}) error

	// EmitEvent appends the given event to the output stream.
	EmitEvent(evt Event) error
}

//--------------------
// STREAM
//--------------------

// stream manages the flow of events between emitter and receiver.
type stream struct {
	eventc chan Event
}

// newStream creates a stream instance.
func newStream() *stream {
	return &stream{
		eventc: make(chan Event),
	}
}

// Pull reads an event out of the stream.
func (str *stream) Pull() <-chan Event {
	return str.eventc
}

// Emit creates a new event and emits it.
func (str *stream) Emit(topic string, payloads ...interface{}) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return str.EmitEvent(evt)
}

// EmitEvent appends an event to the end of the stream. It retries to
// append it to the buffer in case that it's full. The time will
// increase. If it lasts too long, about 5 seconds, a timeout
// error will be returned.
func (str *stream) EmitEvent(evt Event) error {
	wait := 50 * time.Millisecond
	for {
		select {
		case str.eventc <- evt:
			return nil
		default:
			time.Sleep(wait)
			wait += 50 * time.Millisecond
			if wait > 5*time.Second {
				return errors.New("timeout")
			}
		}
	}
}

//--------------------
// STREAMS
//--------------------

// streams is a set of streans to emit to multiple
// streams at once.
type streams struct {
	mu      sync.RWMutex
	streams map[*stream]struct{}
}

// newStreams creates an empty set of streams.
func newStreams() *streams {
	return &streams{
		streams: make(map[*stream]struct{}),
	}
}

// add add a stream to the set of streams.
func (strs *streams) add(as *stream) {
	strs.mu.Lock()
	defer strs.mu.Unlock()
	strs.streams[as] = struct{}{}
}

// remove deletes a stream from the set of streams.
func (strs *streams) remove(rs *stream) {
	strs.mu.Lock()
	defer strs.mu.Unlock()
	delete(strs.streams, rs)
}

// removeAll deletes all streams from the set of streams.
func (strs *streams) removeAll() {
	strs.mu.Lock()
	defer strs.mu.Unlock()
	strs.streams = make(map[*stream]struct{})
}

// Emit creates a new event and appends it to the end of all
// contained streams.
func (strs *streams) Emit(topic string, payloads ...interface{}) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return strs.EmitEvent(evt)
}

// EmitEvent appends the given event to the end of all contained
// streams.
func (strs *streams) EmitEvent(evt Event) error {
	strs.mu.RLock()
	defer strs.mu.RUnlock()
	for es := range strs.streams {
		if err := es.EmitEvent(evt); err != nil {
			return err
		}
	}
	return nil
}

// EOF
