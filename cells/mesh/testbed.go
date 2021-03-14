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
	"context"
	"errors"
	"time"
)

//--------------------
// TESTBED HELPERS
//--------------------

// testbedCell runs the behavior and provides the needed interfaces.
type testbedCell struct {
	ctx      context.Context
	behavior Behavior
	inc      chan *Event
	outc     chan *Event
}

// goTestbedCell initializes the testbed cell and spawns the goroutine.
func goTestbedCell(ctx context.Context, behavior Behavior) *testbedCell {
	tbc := &testbedCell{
		ctx:      ctx,
		behavior: behavior,
		inc:      make(chan *Event),
		outc:     make(chan *Event),
	}
	go tbc.behavior.Go(tbc, tbc, tbc)
	return tbc
}

// Context imepelements mesh.Cell.
func (tbc *testbedCell) Context() context.Context {
	return tbc.ctx
}

// Name imepelements mesh.Cell and returns a static name.
func (tbc *testbedCell) Name() string {
	return "testbed"
}

// Mesh imepelements mesh.Cell.
func (tbc *testbedCell) Mesh() Mesh {
	// TODO Return Mesh implementation.
	return nil
}

// Pull implements mesh.Receptor.
func (tbc *testbedCell) Pull() <-chan *Event {
	return tbc.inc
}

// Emit implements mesh.Emitter.
func (tbc *testbedCell) Emit(evt *Event) error {
	tbc.outc <- evt
	return nil
}

//--------------------
// TESTBED
//--------------------

// Testbed provides a simple environment for the testing of
// individual behaviors.
type Testbed struct {
	ctx    context.Context
	cancel func()
	donec  chan struct{}
	cell   *testbedCell
}

// NewTestbed creates a testbed for the given behavior. The tester function has
// to test the emitted events. Once the final test is done and a criteria is
// fullfilled the tester function has to return true. This signals the Wait()
// method a positive end. Otherwise a timeout will returned.
func NewTestbed(behavior Behavior, tester func(evt *Event) bool) *Testbed {
	ctx, cancel := context.WithCancel(context.Background())
	tb := &Testbed{
		ctx:    ctx,
		cancel: cancel,
		donec:  make(chan struct{}),
		cell:   goTestbedCell(ctx, behavior),
	}
	go func() {
		for {
			select {
			case <-tb.ctx.Done():
				return
			case evt := <-tb.cell.outc:
				if tester(evt) {
					close(tb.donec)
				}
			}
		}
	}()
	return tb
}

// Emit sends an event to the behavior.
func (tb *Testbed) Emit(evt *Event) {
	tb.cell.inc <- evt
}

// Wait waits until test ends or a timeout.
func (tb *Testbed) Wait(timeout time.Duration) error {
	defer tb.cancel()
	select {
	case <-tb.donec:
		return nil
	case to := <-time.After(timeout):
		return errors.New("timeout after " + to.String())
	}
}

// EOF
