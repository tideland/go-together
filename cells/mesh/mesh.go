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
	"fmt"
	"sync"
)

//--------------------
// MESH
//--------------------

// Mesh manages a closed network of cells.
type Mesh struct {
	mu    sync.RWMutex
	ctx   context.Context
	cells map[string]*cell
}

// New creates new Mesh instance.
func New(ctx context.Context) *Mesh {
	m := &Mesh{
		ctx:   ctx,
		cells: make(map[string]*cell),
	}
	return m
}

// Go starts a cell using the given behavior.
func (m *Mesh) Go(name string, b Behavior) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cells[name] != nil {
		return fmt.Errorf("cell name %q already used", name)
	}
	m.cells[name] = newCell(m.ctx, name, b)
	return nil
}

// Subscribe subscribes the cell with from name to the cell
// with to name.
func (m *Mesh) Subscribe(fromName, toName string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	fromCell := m.cells[fromName]
	toCell := m.cells[toName]
	if fromCell == nil {
		return fmt.Errorf("from cell %q does not exist", fromName)
	}
	if toCell == nil {
		return fmt.Errorf("to cell %q does not exist", toName)
	}
	fromCell.subscribe(toCell)
	return nil
}

// Raise appends an event to the in-queue of a cell.
func (m *Mesh) Raise(name string, evt *Event) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	raiseCell := m.cells[name]
	if raiseCell == nil {
		return fmt.Errorf("cell %q does not exist", name)
	}
	return raiseCell.in.Append(evt)
}

// EOF
