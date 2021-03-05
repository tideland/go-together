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
	mu      sync.RWMutex
	ctx     context.Context
	cells   map[string]*cell
	emitter map[string]*Emitter
}

// New creates new Mesh instance.
func New(ctx context.Context) *Mesh {
	m := &Mesh{
		ctx:     ctx,
		cells:   make(map[string]*cell),
		emitter: make(map[string]*Emitter),
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
	m.cells[name] = newCell(m.ctx, m, name, b)
	return nil
}

// drop removes a cell and only can be done by itself
// notifying the mesh that it ends its work.
func (m *Mesh) drop(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.cells, name)
	delete(m.emitter, name)
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
	fromCell.subscribeTo(toCell)
	return nil
}

// Unsubscribe unsubscribes the cell with to name from the cell
// with from name.
func (m *Mesh) Unsubscribe(toName, fromName string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	toCell := m.cells[toName]
	fromCell := m.cells[fromName]
	if toCell == nil {
		return fmt.Errorf("to cell %q does not exist", toName)
	}
	if fromCell == nil {
		return fmt.Errorf("from cell %q does not exist", fromName)
	}
	toCell.unsubscribeFrom(fromCell)
	return nil
}

// Emit raises an event to the named cell.
func (m *Mesh) Emit(name string, evt *Event) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	emitCell := m.cells[name]
	if emitCell == nil {
		return fmt.Errorf("cell %q does not exist", name)
	}
	return emitCell.in.Emit(evt)
}

// Emitter returns a static emitter for the named cell.
func (m *Mesh) Emitter(name string) (*Emitter, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	emitCell := m.cells[name]
	if emitCell == nil {
		return nil, fmt.Errorf("cell %q does not exist", name)
	}
	m.emitter[name] = (*Emitter)(emitCell.in)
	return m.emitter[name], nil
}

// EOF
