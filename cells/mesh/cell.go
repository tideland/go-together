// Tideland Go Together - Cells - Mesh
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh

//--------------------
// IMPORT
//--------------------

import (
	"context"
)

//--------------------
// CELL
//--------------------

// cell runs a behevior networked with other cells.
type cell struct {
	ctx          context.Context
	name         string
	behavior     Behavior
	subscribedTo *registry
	in           *stream
	out          *registry
}

// newCell starts a new cell working in the background.
func newCell(ctx context.Context, name string, b Behavior) *cell {
	c := &cell{
		ctx:          ctx,
		name:         name,
		behavior:     b,
		subscribedTo: newRegistry(),
		in:           newStream(16),
		out:          newRegistry(),
	}
	go c.backend()
	return c
}

// subscribe adds the in queue of this cell to the out cells of the
// given cell.
func (c *cell) subscribe(sc *cell) {
	c.subscribedTo.add(sc)
	sc.out.add(c)
}

// unsubscribeAll removes the subscription from all cells this
// one subscribed to.
func (c *cell) unsubscribeAll() {
	c.subscribedTo.removeAll(c)
}

// backend runs as goroutine and cares for the behavior. When it ends
// then all subscriptions are unsubscribed and subscribers get a
// notification.
func (c *cell) backend() {
	defer func() {
		c.unsubscribeAll()
		c.out.Emit(NewEvent(TerminationTopic, NameKey, c.name))
	}()
	c.behavior.Go(c.ctx, c.name, c.in, c.out)
}

// EOF
