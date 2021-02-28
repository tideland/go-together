// Tideland Go Together - Cells - Mesh - Tests
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
	"context"
	"strings"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
)

//--------------------
// TESTS
//--------------------

// TestCellSimple provides a simple processing of some
// events.
func TestCellSimple(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	topics := []string{}
	sigc := make(chan interface{})
	collector := func(evt *Event, out OutputStream) {
		topics = append(topics, evt.Topic())
		if len(topics) == 3 {
			close(sigc)
		}
	}
	tbCollector := &testBehavior{f: collector}
	cCollector := newCell(ctx, "collector", tbCollector)

	cCollector.in.Emit(NewEvent("one"))
	cCollector.in.Emit(NewEvent("two"))
	cCollector.in.Emit(NewEvent("three"))

	assert.WaitClosed(sigc, time.Second)
	assert.Length(topics, 3)
	assert.Equal(strings.Join(topics, " "), "one two three")

	cancel()
}

// TestCellChain provides a chained processing of some
// events.
func TestCellChain(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	topics := []string{}
	sigc := make(chan interface{})
	upcaser := func(evt *Event, out OutputStream) {
		upperTopic := strings.ToUpper(evt.Topic())
		out.Emit(NewEvent(upperTopic))
	}
	tbUpcaser := &testBehavior{f: upcaser}
	cUpcaser := newCell(ctx, "upcaser", tbUpcaser)
	collector := func(evt *Event, out OutputStream) {
		topics = append(topics, evt.Topic())
		if len(topics) == 3 {
			close(sigc)
		}
	}
	tbCollector := &testBehavior{f: collector}
	cCollector := newCell(ctx, "collector", tbCollector)

	cCollector.subscribe(cUpcaser)

	cUpcaser.in.Emit(NewEvent("one"))
	cUpcaser.in.Emit(NewEvent("two"))
	cUpcaser.in.Emit(NewEvent("three"))

	assert.WaitClosed(sigc, time.Second)
	assert.Length(topics, 3)
	assert.Equal(strings.Join(topics, " "), "ONE TWO THREE")

	cancel()
}

//--------------------
// HELPER
//--------------------

type testBehavior struct {
	f func(evt *Event, out OutputStream)
}

func (tb *testBehavior) Go(ctx context.Context, name string, in InputStream, out OutputStream) {
	for {
		select {
		case <-ctx.Done():
			return
		case evt := <-in.Pull():
			tb.f(evt, out)
		}
	}
}

// EOF
