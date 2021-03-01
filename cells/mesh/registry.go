// Tideland Go Together - Cells - Mesh
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
	"sync"
)

//--------------------
// REGISTRY
//--------------------

type registry struct {
	mu    sync.RWMutex
	cells map[*cell]struct{}
}

// newRegistry creates a cell registry.
func newRegistry() *registry {
	return &registry{
		cells: make(map[*cell]struct{}),
	}
}

// add registers a cell.
func (r *registry) add(ac *cell) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cells[ac] = struct{}{}
}

// remove unregisters a cell.
func (r *registry) remove(rc *cell) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.cells, rc)
}

// removeAll unregisters a cell from the output of all registered cells.
func (r *registry) removeAll(rc *cell) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for c := range r.cells {
		rc.out.remove(c)
	}
}

// Emit implements OutputStream emitting an event to all
// cells.
func (r *registry) Emit(evt *Event) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for c := range r.cells {
		if err := c.in.Emit(evt); err != nil {
			return err
		}
	}
	return nil

}

// EOF
