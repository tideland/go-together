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
	"fmt"
	"time"
)

//--------------------
// TESTBED HELPERS
//--------------------

// testbedMesh implements the Mesh interface.
type testbedMesh struct{}

// Go implements mesh.Mesh and always returns an error.
func (tbm testbedMesh) Go(name string, b Behavior) error {
	return fmt.Errorf("cell name '%s' already used", name)
}

// Subscribe implements mesh.Mesh and always returns an error.
func (tbm testbedMesh) Subscribe(emitterName, receptorName string) error {
	return fmt.Errorf("emitter cell '%s' does not exist", emitterName)
}

// Unsubscribe implements mesh.Mesh and always returns an error.
func (tbm testbedMesh) Unsubscribe(emitterName, receptorName string) error {
	return fmt.Errorf("emitter cell '%s' does not exist", emitterName)
}

// Emit implements mesh.Mesh and always returns an error.
func (tbm testbedMesh) Emit(name, topic string, payloads ...interface{}) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return tbm.EmitEvent(name, evt)
}

// EmitEvent implements mesh.Mesh and always returns an error.
func (tbm testbedMesh) EmitEvent(name string, evt *Event) error {
	return fmt.Errorf("cell '%s' does not exist", name)
}

// Emitter implements mesh.Mesh and always returns an error.
func (tbm testbedMesh) Emitter(name string) (Emitter, error) {
	return nil, fmt.Errorf("cell '%s' does not exist", name)
}

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
	go tbc.backend()
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
	return testbedMesh{}
}

// Pull implements mesh.Receptor.
func (tbc *testbedCell) Pull() <-chan *Event {
	return tbc.inc
}

// Emit implements mesh.Emitter.
func (tbc *testbedCell) Emit(topic string, payloads ...interface{}) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return tbc.EmitEvent(evt)
}

// EmitEvent implements mesh.Emitter.
func (tbc *testbedCell) EmitEvent(evt *Event) error {
	tbc.outc <- evt
	return nil
}

// backend runs the behavior to test.
func (tbc *testbedCell) backend() {
	if err := tbc.behavior.Go(tbc, tbc, tbc); err != nil {
		// Notify subscribers about error.
		tbc.Emit(TopicError, PayloadCellError{
			CellName: tbc.Name(),
			Error:    err.Error(),
		})
	} else {
		// Notify subscribers about termination.
		tbc.Emit(TopicTerminated, PayloadTermination{
			CellName: tbc.Name(),
		})
	}
}

//--------------------
// TESTBED
//--------------------

// Testbed provides a simple environment for the testing of individual behaviors.
// So retrieving the Mesh by the Cell is possible, but using its methods leads to
// errors. Integration tests have to be done differently.
//
// A tester function given when the testbed is started allows to evaluate the
// events emitted by the behavior. As long as the tests aren't done the function
// has to return false. Once returning true for the final tested event
// Testbed.Wait() gets a signal. Otherwise a timeout will be returned to show
// an internal error.
type Testbed struct {
	ctx    context.Context
	cancel func()
	donec  chan struct{}
	cell   *testbedCell
}

// NewTestbed starts a test cell with the given behavior. The tester function
// will be called for each event emitted by the behavior.
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

// Emit creates an event and sends it to the behavior.
func (tb *Testbed) Emit(topic string, payloads ...interface{}) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	tb.EmitEvent(evt)
	return nil
}

// Emit sends an event to the behavior.
func (tb *Testbed) EmitEvent(evt *Event) {
	tb.cell.inc <- evt
}

// Wait waits until test ends or a timeout.
func (tb *Testbed) Wait(timeout time.Duration) error {
	defer tb.cancel()
	now := time.Now()
	select {
	case <-tb.donec:
		return nil
	case to := <-time.After(timeout):
		waited := to.Sub(now)
		return errors.New("timeout after " + waited.String())
	}
}

// EOF
