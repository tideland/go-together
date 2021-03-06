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

// mesh manages a closed network of cells. It implements
// the Mesh interface.
type mesh struct {
	mu       sync.RWMutex
	ctx      context.Context
	cells    map[string]*cell
	emitters map[string]*emitter
}

// New creates new Mesh instance.
func New(ctx context.Context) Mesh {
	m := &mesh{
		ctx:      ctx,
		cells:    make(map[string]*cell),
		emitters: make(map[string]*emitter),
	}
	return m
}

// Go implements Mesh.
func (m *mesh) Go(name string, b Behavior) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cells[name] != nil {
		return fmt.Errorf("cell name %q already used", name)
	}
	m.cells[name] = newCell(m.ctx, name, m, b, func() {
		// Callback for cell to unregister.
		m.mu.Lock()
		defer m.mu.Unlock()
		delete(m.cells, name)
		delete(m.emitters, name)
	})
	return nil
}

// Subscribe implements Mesh.
func (m *mesh) Subscribe(fromName, toName string) error {
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

// Unsubscribe implements Mesh.
func (m *mesh) Unsubscribe(toName, fromName string) error {
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

// Emit implements Mesh.
func (m *mesh) Emit(name string, evt *Event) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	emitCell := m.cells[name]
	if emitCell == nil {
		return fmt.Errorf("cell %q does not exist", name)
	}
	return emitCell.in.Emit(evt)
}

// Emitter implements Mesh.
func (m *mesh) Emitter(name string) (Emitter, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	emitCell := m.cells[name]
	if emitCell == nil {
		return nil, fmt.Errorf("cell %q does not exist", name)
	}
	namedEmitter := m.emitters[name]
	if namedEmitter == nil {
		namedEmitter = &emitter{
			strean: emitCell.in,
		}
		m.emitters[name] = namedEmitter
	}
	return namedEmitter, nil
}

// EOF
