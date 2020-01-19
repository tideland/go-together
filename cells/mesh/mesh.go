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
	"sync"

	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/trace/failure"
)

//--------------------
// MESH
//--------------------

// Mesh operates a set of interacting cells.
type Mesh struct {
	mu       sync.RWMutex
	cells    cellRegistry
	queueCap int
}

// New creates a new event processing mesh.
func New(opts ...Option) *Mesh {
	msh := &Mesh{
		cells:    cellRegistry{},
		queueCap: 1,
	}
	for _, opt := range opts {
		opt(msh)
	}
	return msh
}

// SpawnCells starts cells running the passed behaviors to work as parts
// of the mesh.
func (msh *Mesh) SpawnCells(behaviors ...Behavior) error {
	msh.mu.Lock()
	defer msh.mu.Unlock()
	// Step one: check IDs.
	var ids []string
	for _, behavior := range behaviors {
		id := behavior.ID()
		if msh.cells.contains(id) {
			ids = append(ids, id)
		}
	}
	if len(ids) != 0 {
		return failure.New("spawn cells: double id(s) %v", ids)
	}
	// Step two: create cells.
	var errs []error
	for _, behavior := range behaviors {
		id := behavior.ID()
		cell, err := newCell(msh, behavior)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		msh.cells.add(id, cell)
	}
	return failure.Annotate(failure.Collect(errs...), "spawn cells")
}

// StopCells terminates the given cells.
func (msh *Mesh) StopCells(ids ...string) error {
	msh.mu.Lock()
	defer msh.mu.Unlock()
	for _, id := range ids {
		if err := msh.cells.unsubscribeFromAll(id); err != nil {
			return err
		}
		if err := msh.cells.remove(id); err != nil {
			return err
		}
	}
	return nil
}

// Cells returns the identifiers of the spawned cells.
func (msh *Mesh) Cells() []string {
	msh.mu.RLock()
	defer msh.mu.RUnlock()
	var ids []string
	for id := range msh.cells {
		ids = append(ids, id)
	}
	return ids
}

// Subscribers retrieves the subscriber IDs of a cell.
func (msh *Mesh) Subscribers(id string) ([]string, error) {
	msh.mu.RLock()
	defer msh.mu.RUnlock()
	// Retrieve all needed cells.
	entry, ok := msh.cells[id]
	if !ok {
		return nil, failure.New("subscribers: %s not found", id)
	}
	return entry.cell.subscribers()
}

// Subscribe connects cells to the given cell.
func (msh *Mesh) Subscribe(id string, subscriberIDs ...string) error {
	msh.mu.Lock()
	defer msh.mu.Unlock()
	return msh.cells.subscribe(id, subscriberIDs)
}

// Unsubscribe disconnect cells from the given cell.
func (msh *Mesh) Unsubscribe(id string, unsubscriberIDs ...string) error {
	msh.mu.Lock()
	defer msh.mu.Unlock()
	return msh.cells.unsubscribe(id, unsubscriberIDs)
}

// Emit sends an event to the given cell.
func (msh *Mesh) Emit(id string, evt *event.Event) error {
	msh.mu.RLock()
	defer msh.mu.RUnlock()
	// Retrieve the needed cell.
	entry, ok := msh.cells[id]
	if !ok {
		return failure.New("emit: %s not found", id)
	}
	return entry.cell.process(evt)
}

// Broadcast sends an event to all cells.
func (msh *Mesh) Broadcast(evt *event.Event) error {
	msh.mu.RLock()
	defer msh.mu.RUnlock()
	cerrs := make([]error, len(msh.cells))
	idx := 0
	// Broadcast.
	for _, entry := range msh.cells {
		cerrs[idx] = entry.cell.process(evt)
		idx++
	}
	// Return collected errors.
	return failure.Collect(cerrs...)
}

// Stop terminates the cells and cleans up.
func (msh *Mesh) Stop() error {
	msh.mu.Lock()
	defer msh.mu.Unlock()
	var errs []error
	// Terminate.
	for _, entry := range msh.cells {
		if err := entry.cell.stop(); err != nil {
			errs = append(errs, err)
		}
	}
	msh.cells = cellRegistry{}
	// Return collected errors.
	return failure.Collect(errs...)
}

// EOF
