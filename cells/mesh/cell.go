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
	ctx      context.Context
	name     string
	behavior Behavior
	in       *stream
	out      *streams
}

// newCell starts a new cell working in the background.
func newCell(ctx context.Context, name string, b Behavior) *cell {
	c := &cell{
		ctx:      ctx,
		name:     name,
		behavior: b,
		in:       newStream(16),
		out:      newStreams(),
	}
	go c.behavior.Go(c.ctx, c.name, c.in, c.out)
	return c
}

// subscribe adds the in queue of this cell to the out cells of the
// given cell.
func (c *cell) subscribe(to *cell) {
	to.out.add(c.in)
}

// EOF
