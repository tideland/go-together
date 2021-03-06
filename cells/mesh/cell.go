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
	"sync"
)

//--------------------
// CELL
//--------------------

// cell runs a behevior networked with other cells.
type cell struct {
	mu       sync.RWMutex
	ctx      context.Context
	name     string
	mesh     Mesh
	behavior Behavior
	in       *stream
	inCells  map[*cell]struct{}
	out      *streams
	drop     func()
}

// newCell starts a new cell working in the background.
func newCell(ctx context.Context, name string, m Mesh, b Behavior, drop func()) *cell {
	c := &cell{
		ctx:      ctx,
		name:     name,
		mesh:     m,
		behavior: b,
		in:       newStream(16),
		inCells:  make(map[*cell]struct{}),
		out:      newStreams(),
		drop:     drop,
	}
	go c.backend()
	return c
}

// Context implements Cell.
func (c *cell) Context() context.Context {
	return c.ctx
}

// Name implements Cell.
func (c *cell) Name() string {
	return c.name
}

// Mesh implements Cell.
func (c *cell) Mesh() Mesh {
	return nil
}

// subscribeTo adds the cell to the out-streams of the
// given in-cell.
func (c *cell) subscribeTo(inCell *cell) {
	c.mu.Lock()
	defer c.mu.Unlock()
	inCell.out.add(c.in)
	c.inCells[inCell] = struct{}{}
}

// unsubscribeFrom removes the cell from the out-streams of the
// given in-cell.
func (c *cell) unsubscribeFrom(inCell *cell) {
	c.mu.Lock()
	defer c.mu.Unlock()
	inCell.out.remove(c.in)
	delete(c.inCells, inCell)
}

// unsubscribeFromAll removes the subscription from all cells this
// one subscribed to.
func (c *cell) unsubscribFromeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for inCell := range c.inCells {
		inCell.out.remove(c.in)
	}
}

// backend runs as goroutine and cares for the behavior. When it ends
// it will send a notification to all subscribers, unsubscribe from
// them, and then tell the mesh that it's not available anymore.
func (c *cell) backend() {
	defer func() {
		c.out.Emit(NewEvent(TerminationTopic, NameKey, c.name))
		c.unsubscribFromeAll()
		c.drop()
	}()
	c.behavior.Go(c, c.in, c.out)
}

// EOF
