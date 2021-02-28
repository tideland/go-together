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

// streams contains a number of streams to distribute events to.
type streams struct {
	mu      sync.RWMutex
	streams map[*stream]struct{}
}

// newStreams creates a multi-stream instance.
func newStreams() *streams {
	return &streams{
		streams: make(map[*stream]struct{}),
	}
}

// add inserts the given stream to the streams.
func (strs *streams) add(str *stream) {
	strs.mu.Lock()
	defer strs.mu.Unlock()
	strs.streams[str] = struct{}{}
}

// remove deletes the given stream from the streams.
func (strs *streams) remove(str *stream) {
	strs.mu.Lock()
	defer strs.mu.Unlock()
	delete(strs.streams, str)
}

// Emit appends an event to the end of all streams.
func (strs *streams) Emit(evt *Event) error {
	strs.mu.RLock()
	defer strs.mu.RUnlock()
	for str := range strs.streams {
		if err := str.Emit(evt); err != nil {
			return err
		}
	}
	return nil
}

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

// EOF
