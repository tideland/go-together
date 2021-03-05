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
	"errors"
	"sync"
	"time"
)

//--------------------
// INPUT AND OUTPUT STREAMS
//--------------------

// InputStream provides a stream for reading events to process.
type InputStream interface {
	// Pull reads an event out of the stream.
	Pull() <-chan *Event
}

// OutputStream provices a stream for emitting events
// the subscribers have to process.
type OutputStream interface {
	// Emit appends an event to the end of all streams.
	Emit(evt *Event) error
}

//--------------------
// STREAM
//--------------------

// stream manages the flow of events between emitter and receiver.
type stream struct {
	eventc chan *Event
}

// newStream creates a stream instance with the given buffer size.
func newStream(size int) *stream {
	return &stream{
		eventc: make(chan *Event, size),
	}
}

// Pull reads an event out of the stream.
func (str *stream) Pull() <-chan *Event {
	return str.eventc
}

// Emit appends an event to the end of the stream. It retries to
// append it to the buffer in case that it's full. The time will
// increase. If it lasts too long, about 5 seconds, a timeout
// error will be returned.
func (str *stream) Emit(evt *Event) error {
	wait := 75 * time.Millisecond
	for {
		select {
		case str.eventc <- evt:
			return nil
		default:
			time.Sleep(wait)
			wait *= 2
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
func (s *streams) add(as *stream) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streams[as] = struct{}{}
}

// remove deletes a stream from the set of streams.
func (s *streams) remove(rs *stream) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.streams, rs)
}

// removeAll deletes all streams from the set of streams.
func (s *streams) removeAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streams = make(map[*stream]struct{})
}

// Emit implements OutputStream emitting an event to all streams.
func (s *streams) Emit(evt *Event) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for es := range s.streams {
		if err := es.Emit(evt); err != nil {
			return err
		}
	}
	return nil
}

// EOF
