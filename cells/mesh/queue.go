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
// QUEUE
//--------------------

// queue contains a number of events to be processed by
// a cell behavior.
type queue struct {
	eventc chan *Event
}

// newQueue creates a queue instance with the given buffer size.
func newQueue(size int) *Queue {
	return &queue{
		eventc: make(chan *Event, size),
	}
}

// Pull reads an event out of the queue.
func (q *queue) Pull() <-chan *Event {
	return q.eventc
}

// Append adds an event to the end of the queue. It retries to
// add it to the buffer in case that it's full. The time will
// increase. If it lasts too long, about 5 seconds, a timeout
// error will be returned.
func (q *queue) Append(evt *Event) error {
	wait := 75 * time.Millisecond
	for {
		select {
		case q.eventc <- evt:
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

// queues contains a number of queues to distribute events to.
type queues struct {
	mu     sync.RWMutex
	queues map[*queue]struct{}
}

// newQueue creates a multi-queue instance.
func newQueues() *queues {
	return &queues{
		queues: make(map[*queue]struct{}),
	}
}

// add inserts the given queue to the queues.
func (qs *queues) add(q *queue) {
	qs.mu.Lock()
	defer qs.mu.Unlock()
	qs.queues[q] = struct{}{}
}

// remove deletes the given queue from the queues.
func (qs *queues) remove(q *queue) {
	qs.mu.Lock()
	defer qs.mu.Unlock()
	delete(qs.queues, q)
}

// Append adds an event to the end of all queues.
func (qs *queues) Append(evt *Event) error {
	qs.mu.RLock()
	defer qs.mu.RUnlock()
	for q := range qs.queues {
		if err := q.Append(evt); err != nil {
			return err
		}
	}
	return nil
}

//--------------------
// INPUT AND OUTPUT QUEUE
//--------------------

// InputQueue provides a queue for reading events to process.
type InputQueue interface {
	// Pull reads an event out of the queue.
	Pull() <-chan *Event
}

// OutputQueue provices a queue for writing events
// a subscriber has to process.
type OutputQueue interface {
	// Append adds an event to the end of the queue.
	Append(evt *Event) error
}

// EOF
