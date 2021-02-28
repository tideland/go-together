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
	"sync"
	"testing"

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

	var wg sync.WaitGroup

	wg.Add(3)

	tb := &testBehavior{
		f: func(evt *Event, out OutputStream) {
			t := evt.Topic()
			assert.Contains(t, []string{"one", "two", "three"})
			wg.Done()
		},
	}
	c := newCell(ctx, "test", tb)

	c.in.Emit(NewEvent("one", nil))
	c.in.Emit(NewEvent("two", nil))
	c.in.Emit(NewEvent("three", nil))

	wg.Wait()
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
