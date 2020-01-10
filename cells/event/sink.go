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
	"sync"
	"time"
)

//--------------------
// CONSTANTS
//--------------------

// CriterionMatch signals, how a criterion matches.
type CriterionMatch int

// List of criterion match signals.
const (
	CriterionDone CriterionMatch = iota + 1
	CriterionKeep
	CriterionDropFirst
	CriterionDropLast
	CriterionClear
)

//--------------------
// EVENT SINK FUNCTION TYPES
//--------------------

// SinkDoer performs an operation on an event.
type SinkDoer func(index int, evt *Event) error

// SinkProcessor can be used as a checker function but also inside of
// behaviors to process the content of an event sink and return a new payload.
type SinkProcessor func(accessor SinkAccessor) (*Payload, error)

// SinkFilter checks if an event matches a criterium.
type SinkFilter func(index int, evt *Event) (bool, error)

// SinkFolder allows to reduce (fold) events.
type SinkFolder func(index int, acc *Payload, evt *Event) (*Payload, error)

//--------------------
// EVENT SINK ACCESSOR
//--------------------

// SinkAccessor can be used to read the events in a sink. It is a
// specialized subfunctionality of the event sink.
type SinkAccessor interface {
	// Len returns the number of stored events.
	Len() int

	// PeekFirst returns the first of the collected events.
	PeekFirst() (*Event, bool)

	// PeekLast returns the last of the collected event datas.
	PeekLast() (*Event, bool)

	// PeekAt returns an event at a given index and true if it
	// exists, otherwise nil and false.
	PeekAt(index int) (*Event, bool)

	// Do iterates over all collected events.
	Do(doer SinkDoer) error
}

//--------------------
// EVENT SINK
//--------------------

// Sink stores a number of events ordered by adding them at the end. To
// be used in behaviors for collecting sets of events and operate on them.
type Sink struct {
	mu     sync.RWMutex
	max    int
	events []*Event
	check  SinkProcessor
}

// NewSink creates a sink for events.
func NewSink(max int) *Sink {
	return &Sink{
		max: max,
	}
}

// NewCheckedSink creates a sink for events.
func NewCheckedSink(max int, checker SinkProcessor) *Sink {
	return &Sink{
		max:   max,
		check: checker,
	}
}

// Push adds a new event to the sink.
func (s *Sink) Push(evt *Event) (int, error) {
	s.mu.Lock()
	s.events = append(s.events, evt)
	if s.max > 0 && len(s.events) > s.max {
		s.events = s.events[1:]
	}
	s.mu.Unlock()
	return len(s.events), s.performCheck()
}

// PullFirst returns and removed the first event of the sink.
func (s *Sink) PullFirst() (*Event, error) {
	var evt *Event
	s.mu.Lock()
	if len(s.events) > 0 {
		evt = s.events[0]
		s.events = s.events[1:]
	}
	s.mu.Unlock()
	return evt, s.performCheck()
}

// PullLast returns and removed the last event of the sink.
func (s *Sink) PullLast() (*Event, error) {
	var evt *Event
	s.mu.Lock()
	if len(s.events) > 0 {
		evt = s.events[len(s.events)-1]
		s.events = s.events[:len(s.events)-1]
	}
	s.mu.Unlock()
	return evt, s.performCheck()
}

// Clear removes all collected events.
func (s *Sink) Clear() error {
	s.mu.Lock()
	s.events = nil
	s.mu.Unlock()
	return s.performCheck()
}

// Len implements SinkAccessor.
func (s *Sink) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.events)
}

// PeekFirst implements SinkAccessor.
func (s *Sink) PeekFirst() (*Event, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.events) < 1 {
		return nil, false
	}
	return s.events[0], true
}

// PeekLast implements SinkAccessor.
func (s *Sink) PeekLast() (*Event, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.events) < 1 {
		return nil, false
	}
	return s.events[len(s.events)-1], true
}

// PeekAt implements SinkAccessor.
func (s *Sink) PeekAt(index int) (*Event, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if index < 0 || index > len(s.events)-1 {
		return nil, false
	}
	return s.events[index], true
}

// Do implements SinkAccessor.
func (s *Sink) Do(doer SinkDoer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for index, evt := range s.events {
		if err := doer(index, evt); err != nil {
			return err
		}
	}
	return nil
}

// performCheck calls the checker if configured.
func (s *Sink) performCheck() error {
	if s.check != nil {
		if _, err := s.check(s); err != nil {
			return err
		}
	}
	return nil
}

//--------------------
// EVENT SINK ANALYZER
//--------------------

// SinkAnalyzer is a helpful type to analyze the events collected inside
// an event sink. It's intended to make the life for behavior developers more
// simple.
type SinkAnalyzer struct {
	accessor SinkAccessor
}

// NewSinkAnalyzer creates an analyzer for the given sink accessor.
func NewSinkAnalyzer(accessor SinkAccessor) *SinkAnalyzer {
	return &SinkAnalyzer{
		accessor: accessor,
	}
}

// Filter creates a new accessor containing only the filtered events.
func (sa *SinkAnalyzer) Filter(filter SinkFilter) (SinkAccessor, error) {
	accessor := NewSink(sa.accessor.Len())
	doer := func(index int, evt *Event) error {
		ok, err := filter(index, evt)
		if err != nil {
			accessor = nil
			return err
		}
		if ok {
			accessor.Push(evt)
		}
		return nil
	}
	err := sa.accessor.Do(doer)
	return accessor, err
}

// Match checks if all events match the passed criterion.
func (sa *SinkAnalyzer) Match(matcher SinkFilter) (bool, error) {
	match := true
	doer := func(index int, evt *Event) error {
		ok, err := matcher(index, evt)
		if err != nil {
			match = false
			return err
		}
		match = match && ok
		return nil
	}
	err := sa.accessor.Do(doer)
	return match, err
}

// Fold reduces (folds) the events of the sink.
func (sa *SinkAnalyzer) Fold(inject *Payload, folder SinkFolder) (*Payload, error) {
	acc := inject
	doer := func(index int, evt *Event) error {
		facc, err := folder(index, acc, evt)
		if err != nil {
			acc = nil
			return err
		}
		acc = facc
		return nil
	}
	err := sa.accessor.Do(doer)
	return acc, err
}

// TotalDuration returns the duration between the first and the last event.
func (sa *SinkAnalyzer) TotalDuration() time.Duration {
	first, firstOK := sa.accessor.PeekFirst()
	last, lastOK := sa.accessor.PeekLast()
	if !firstOK || !lastOK {
		return 0
	}
	return last.Timestamp().Sub(first.Timestamp())
}

// MinMaxDuration returns the minimum and maximum durations between events
// in the sink.
func (sa *SinkAnalyzer) MinMaxDuration() (time.Duration, time.Duration) {
	minDuration := sa.TotalDuration()
	maxDuration := 0 * time.Nanosecond
	lastTimestamp := time.Time{}
	doer := func(index int, evt *Event) error {
		if index > 0 {
			duration := evt.Timestamp().Sub(lastTimestamp)
			if duration < minDuration {
				minDuration = duration
			}
			if duration > maxDuration {
				maxDuration = duration
			}
		}
		lastTimestamp = evt.Timestamp()
		return nil
	}
	sa.accessor.Do(doer)
	return minDuration, maxDuration
}

// TopicQuantities returns a map of collected topics and their quantity.
func (sa *SinkAnalyzer) TopicQuantities() map[string]int {
	topics := map[string]int{}
	doer := func(index int, evt *Event) error {
		topics[evt.Topic()] = topics[evt.Topic()] + 1
		return nil
	}
	sa.accessor.Do(doer)
	return topics
}

// TopicFolds reduces the events per topic.
func (sa *SinkAnalyzer) TopicFolds(folder SinkFolder) (map[string]*Payload, error) {
	folds := map[string]*Payload{}
	doer := func(index int, evt *Event) error {
		facc, err := folder(index, folds[evt.Topic()], evt)
		if err != nil {
			folds = nil
			return err
		}
		folds[evt.Topic()] = facc
		return nil
	}
	err := sa.accessor.Do(doer)
	return folds, err
}

// EOF
