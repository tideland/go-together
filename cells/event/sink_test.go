// Tideland Go Together - Cells - Event - Unit Tests
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package event_test // import "tideland.dev/go/together/cells/event"

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"sort"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/generators"
	"tideland.dev/go/together/cells/event"
)

//--------------------
// TESTS
//--------------------

// TestSink tests the simple event sink.
func TestSink(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	checkTopic := func(evt *event.Event) {
		assert.Contents(evt.Topic(), topicData)
	}

	// Empty sink.
	sink := event.NewSink(0)
	first, ok := sink.PeekFirst()
	assert.Nil(first)
	assert.False(ok)
	last, ok := sink.PeekLast()
	assert.Nil(last)
	assert.False(ok)
	at, ok := sink.PeekAt(-1)
	assert.Nil(at)
	assert.False(ok)
	at, ok = sink.PeekAt(4711)
	assert.Nil(at)
	assert.False(ok)

	// Limited number of events.
	sink = event.NewSink(5)
	addEvents(assert, 10, sink)
	assert.Length(sink, 5)

	// Unlimited number of events.
	sink = event.NewSink(0)
	addEvents(assert, 10, sink)
	assert.Length(sink, 10)

	first, ok = sink.PeekFirst()
	assert.True(ok)
	checkTopic(first)
	last, ok = sink.PeekLast()
	assert.True(ok)
	checkTopic(last)

	for i := 0; i < sink.Len(); i++ {
		at, ok = sink.PeekAt(i)
		assert.True(ok)
		checkTopic(at)
	}
}

// TestSinkIteration tests the event sink iteration.
func TestSinkIteration(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sink := event.NewSink(0)

	addEvents(assert, 10, sink)
	assert.Length(sink, 10)

	err := sink.Do(func(index int, evt *event.Event) error {
		assert.Contents(evt.Topic(), topicData)
		v := evt.Payload().At("bool").AsBool(false)
		assert.True(v)
		return nil
	})
	assert.NoError(err)

	ok, err := event.NewSinkAnalyzer(sink).Match(func(index int, evt *event.Event) (bool, error) {
		topicOK := contains(evt.Topic(), topicData)
		payloadOK := evt.Payload().At("bool").AsBool(false)
		return topicOK && payloadOK, nil
	})
	assert.NoError(err)
	assert.True(ok)
}

// TestSinkIterationError tests the event sink iteration error.
func TestSinkIterationError(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sink := event.NewSink(0)

	addEvents(assert, 10, sink)

	err := sink.Do(func(index int, event *event.Event) error {
		return errors.New("ouch")
	})
	assert.ErrorMatch(err, "ouch")
	ok, err := event.NewSinkAnalyzer(sink).Match(func(index int, event *event.Event) (bool, error) {
		// The bool true won't be passed to outside.
		return true, errors.New("ouch")
	})
	assert.False(ok)
	assert.ErrorMatch(err, "ouch")
}

// TestCheckedSink tests the checking of new events.
func TestCheckedSink(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	payloadc := asserts.MakeWaitChan()
	donec := asserts.MakeWaitChan()
	count := 0
	wanted := []string{"f", "c", "c"}
	checker := func(events event.SinkAccessor) (*event.Payload, error) {
		count++
		defer func() {
			if count == 100 {
				donec <- struct{}{}
			}
		}()
		if events.Len() < len(wanted) {
			return nil, nil
		}
		ok, err := event.NewSinkAnalyzer(events).Match(func(index int, evt *event.Event) (bool, error) {
			return evt.Topic() == wanted[index], nil
		})
		if err != nil {
			return nil, err
		}
		if ok {
			first, _ := events.PeekFirst()
			last, _ := events.PeekLast()
			payload := last.Timestamp().Sub(first.Timestamp())
			payloadc <- payload
		}
		return nil, nil
	}
	sink := event.NewCheckedSink(3, checker)

	go addEvents(assert, 100, sink)

	for {
		select {
		case payload := <-payloadc:
			d, ok := payload.(time.Duration)
			assert.True(ok)
			assert.True(d > 0)
		case <-donec:
			return
		case <-time.After(5 * time.Second):
			assert.Fail()
		}
	}
}

// TestSinkAnalyzer tests analyzing an event sink.
func TestSinkAnalyzer(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sink := event.NewSink(0)
	analyzer := event.NewSinkAnalyzer(sink)

	addEvents(assert, 100, sink)

	// Check filtering.
	threechecker := func(index int, event *event.Event) (bool, error) {
		if event.Topic() == "three" {
			return true, nil
		}
		return false, nil
	}
	ferrchecker := func(index int, event *event.Event) (bool, error) {
		return false, errors.New("ouch")
	}
	filterSink, err := analyzer.Filter(threechecker)
	assert.NoError(err)
	assert.True(filterSink.Len() < sink.Len(), "less events with topic 'three' than total number")
	filterSink, err = analyzer.Filter(ferrchecker)
	assert.ErrorMatch(err, "ouch", "error is returned correctly")

	// Check matching.
	filterSink, err = analyzer.Filter(threechecker)
	assert.NoError(err)
	ok, err := event.NewSinkAnalyzer(filterSink).Match(threechecker)
	assert.NoError(err)
	assert.True(ok, "all events in filterSink do have topic 'three'")
	ok, err = event.NewSinkAnalyzer(filterSink).Match(ferrchecker)
	assert.ErrorMatch(err, "ouch", "error is returned correctly")

	// Check folding.
	testCount := 0
	threefolder := func(index int, acc *event.Payload, evt *event.Event) (*event.Payload, error) {
		if evt.Topic() == "three" {
			testCount++
			count := acc.At("count").AsInt(0)
			return event.NewPayload("count", count+1), nil
		}
		return acc, nil
	}
	ferrfolder := func(index int, acc *event.Payload, evt *event.Event) (*event.Payload, error) {
		return nil, errors.New("ouch")
	}
	pl := event.NewPayload("count", 0)
	pl, err = analyzer.Fold(pl, threefolder)
	assert.NoError(err)
	count := pl.At("count").AsInt(0)
	assert.Equal(count, testCount, "accumulator has been updated correctly")
	pl, err = analyzer.Fold(pl, ferrfolder)
	assert.ErrorMatch(err, "ouch", "error is returned correctly")

	// Check total duration.
	dsink := event.NewSink(0)
	danalyzer := event.NewSinkAnalyzer(dsink)
	duration := danalyzer.TotalDuration()
	assert.Equal(duration, 0*time.Nanosecond, "empty sink has no duration")

	addEvents(assert, 1, dsink)
	duration = danalyzer.TotalDuration()
	assert.Equal(duration, 0*time.Nanosecond, "sink containing one event has no duration")

	addEvents(assert, 1, dsink)
	first, ok := dsink.PeekFirst()
	assert.True(ok)
	last, ok := dsink.PeekLast()
	assert.True(ok)
	duration = danalyzer.TotalDuration()
	assert.Equal(duration, last.Timestamp().Sub(first.Timestamp()), "total duration calculated correctly")

	// Check minimum/maximum duration.
	durations := []time.Duration{}
	timestamp := time.Time{}
	sink.Do(func(index int, event *event.Event) error {
		if index == 0 {
			timestamp = event.Timestamp()
			return nil
		}
		duration = event.Timestamp().Sub(timestamp)
		durations = append(durations, duration)
		timestamp = event.Timestamp()
		return nil
	})
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})
	dmin := durations[0]
	dmax := durations[len(durations)-1]
	min, max := analyzer.MinMaxDuration()
	assert.Equal(min, dmin, "minimum duration is correct")
	assert.Equal(max, dmax, "maximum duration is correct")

	// Check topic quantities.
	quantities := analyzer.TopicQuantities()
	assert.Length(quantities, len(topicData))
	for topic, quantity := range quantities {
		assert.Contents(topic, topicData, "topic is one of the topics")
		assert.Range(quantity, 1, 100, "quantity is in range")
	}

	tfolder := func(index int, acc *event.Payload, evt *event.Event) (*event.Payload, error) {
		if acc == nil {
			return event.NewPayload(), nil
		}
		count := acc.At("count").AsInt(1)
		return acc.Clone("count", count+1), nil
	}
	terrfolder := func(index int, acc *event.Payload, evt *event.Event) (*event.Payload, error) {
		return nil, errors.New("ouch")
	}
	payloads, err := analyzer.TopicFolds(tfolder)
	assert.NoError(err)
	assert.Length(quantities, len(topicData))
	for topic, payload := range payloads {
		assert.Contents(topic, topicData, "topic is one of the topics")
		count := payload.At("count").AsInt(-1)
		assert.Range(count, 1, 100, "quantity is in range")
	}
	payloads, err = analyzer.TopicFolds(terrfolder)
	assert.ErrorMatch(err, "ouch", "error is returned correctly")
}

//--------------------
// HELPER
//--------------------

// topicData contains the test topics.
var topicData = []string{"one", "two", "three", "four", "five"}

// payloadData contains the payload data.
var payloadData = []interface{}{"string", "foo", "bool", true, "int", 1}

// addEvents adds a number of events to a sink.
func addEvents(assert *asserts.Asserts, count int, sink *event.Sink) {
	generator := generators.New(generators.FixedRand())
	for i := 0; i < count; i++ {
		topic := generator.OneStringOf(topicData...)
		sleep := generator.Duration(2*time.Millisecond, 4*time.Millisecond)
		event := event.New(topic, payloadData...)
		n, err := sink.Push(event)
		assert.NoError(err)
		assert.True(n > 0)
		time.Sleep(sleep)
	}
}

// contains checks if test is content of set.
func contains(test string, set []string) bool {
	for _, content := range set {
		if test == content {
			return true
		}
	}
	return false
}

// EOF
