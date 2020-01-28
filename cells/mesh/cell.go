// Tideland Go Together - Cells - Mesh
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package mesh // import "tideland.dev/go/together/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/together/actor"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/fuse"
	"tideland.dev/go/trace/failure"
)

//--------------------
// CELL
//--------------------

// cell runs a behavior for the processing of events and emitting of
// resulting events.
type cell struct {
	msh             *Mesh
	behavior        Behavior
	queueCap        int
	subscribedCells map[string]*cell
	act             *actor.Actor
}

// newCell creates a new cell running the given behavior in a goroutine.
func newCell(msh *Mesh, behavior Behavior) (*cell, error) {
	c := &cell{
		msh:             msh,
		behavior:        behavior,
		queueCap:        msh.queueCap,
		subscribedCells: map[string]*cell{},
	}
	// Initialize the behavior with the freshly created cell as emitter.
	err := c.behavior.Init(c)
	if err != nil {
		return nil, failure.Annotate(err, "init cell %q", behavior.ID())
	}
	// Configure the cell if wanted.
	if configurator, ok := behavior.(Configurator); ok {
		configurator.Configure(c)
	}
	// Start the backend.
	act, err := actor.Go(
		actor.WithQueueCap(c.queueCap),
		actor.WithRecoverer(behavior.Recover),
		actor.WithFinalizer(c.finalize),
	)
	if err != nil {
		return nil, failure.Annotate(err, "init cell %q", behavior.ID())
	}
	c.act = act
	return c, nil
}

// Subscribers is part of the emitter interface and returns the
// the IDs of the subscriber cells.
func (c *cell) Subscribers() []string {
	var subscriberIDs []string
	for subscriberID := range c.subscribedCells {
		subscriberIDs = append(subscriberIDs, subscriberID)
	}
	return subscriberIDs
}

// Emit is part of Emitter interface and emits the given event
// to the given subscriber if it exists.
func (c *cell) Emit(id string, evt *event.Event) error {
	subscriber, ok := c.subscribedCells[id]
	if !ok {
		return failure.New("emit: cell %q is no subscriber", id)
	}
	return subscriber.process(evt)
}

// Broadcast is part of Emitter interface and emits the given
// event to all subscribers.
func (c *cell) Broadcast(evt *event.Event) error {
	var serrs []error
	for _, subscriber := range c.subscribedCells {
		serrs = append(serrs, subscriber.process(evt))
	}
	return failure.Collect(serrs...)
}

// Self is part of Emitter interface and emits the given event
// back to the cell itself.
func (c *cell) Self(evt *event.Event) {
	fuse.Trigger(c.msh.Emit(c.behavior.ID(), evt))
}

// SetQueueCap is part of the Configurable interface and allows
// a behavior to set the queue capacity.
func (c *cell) SetQueueCap(qc int) {
	if qc > 1 {
		c.queueCap = qc
	}
}

// subscribers returns the subscriber IDs of the cell.
func (c *cell) subscribers() ([]string, error) {
	var subscriberIDs []string
	if aerr := c.act.DoSync(func() {
		subscriberIDs = c.Subscribers()
	}); aerr != nil {
		return nil, failure.Annotate(aerr, "subscribers of cell %q", c.behavior.ID())
	}
	return subscriberIDs, nil
}

// subscribe adds cells to the subscribers of this cell.
func (c *cell) subscribe(subscribers []*cell) error {
	if aerr := c.act.DoSync(func() {
		for _, subscriber := range subscribers {
			c.subscribedCells[subscriber.behavior.ID()] = subscriber
		}
	}); aerr != nil {
		return failure.Annotate(aerr, "subscribe cell %q", c.behavior.ID())
	}
	return nil
}

// unsubscribe removes cells from the subscribers of this cell.
func (c *cell) unsubscribe(subscribers []*cell) error {
	if aerr := c.act.DoSync(func() {
		for _, subscriber := range subscribers {
			delete(c.subscribedCells, subscriber.behavior.ID())
		}
	}); aerr != nil {
		return failure.Annotate(aerr, "unsubscribe cell %q", c.behavior.ID())
	}
	return nil
}

// process lets the cell behavior process the event asynchronously.
func (c *cell) process(evt *event.Event) error {
	if aerr := c.act.DoAsync(func() {
		if evt.Done() {
			return
		}
		c.behavior.Process(evt)
	}); aerr != nil {
		return failure.Annotate(aerr, "processing cell %q", c.behavior.ID())
	}
	return nil
}

// Finalize implements the loop.Finalizer to perform termination
// when the actor stops.
func (c *cell) finalize(err error) error {
	terr := c.behavior.Terminate()
	c.behavior = &terminatedBehavior{c.behavior.ID()}
	c.subscribedCells = map[string]*cell{}
	return failure.Collect(err, terr)
}

// stop tells the actor to stop with finalizing for termination
// of the behavior.
func (c *cell) stop() error {
	if aerr := c.act.Stop(); aerr != nil {
		return failure.Annotate(aerr, "stopping cell %q", c.behavior.ID())
	}
	return nil
}

//--------------------
// TERMINATED BEHAVIOR
//--------------------

// terminatedBehavior will be used by a cell after shutting down.
type terminatedBehavior struct {
	id string
}

func (db *terminatedBehavior) ID() string {
	return db.id
}

func (db *terminatedBehavior) Init(emitter Emitter) error {
	return nil
}

func (db *terminatedBehavior) Terminate() error {
	return nil
}

func (db *terminatedBehavior) Process(evt *event.Event) {}

func (db *terminatedBehavior) Recover(r interface{}) error {
	return nil
}

// EOF
