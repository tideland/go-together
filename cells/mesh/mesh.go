// Tideland Go Together - Cells - Mesh
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/together/cells/mesh"

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
		return fmt.Errorf("cell name '%s' already used", name)
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
func (m *mesh) Subscribe(emitterName, receptorName string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	emitterCell := m.cells[emitterName]
	receptorCell := m.cells[receptorName]
	if emitterCell == nil {
		return fmt.Errorf("emitter cell '%s' does not exist", emitterName)
	}
	if receptorCell == nil {
		return fmt.Errorf("receptor cell '%s' does not exist", receptorName)
	}
	receptorCell.subscribeTo(emitterCell)
	return nil
}

// Unsubscribe implements Mesh.
func (m *mesh) Unsubscribe(emitterName, receptorName string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	emitterCell := m.cells[emitterName]
	receptorCell := m.cells[receptorName]
	if emitterCell == nil {
		return fmt.Errorf("emitter cell '%s' does not exist", emitterName)
	}
	if receptorCell == nil {
		return fmt.Errorf("receptor cell '%s' does not exist", receptorName)
	}
	receptorCell.unsubscribeFrom(emitterCell)
	return nil
}

// Emit implements Mesh.
func (m *mesh) Emit(name, topic string, payloads ...interface{}) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return m.EmitEvent(name, evt)
}

// EmitEvent implements Mesh.
func (m *mesh) EmitEvent(name string, evt *Event) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	emitCell := m.cells[name]
	if emitCell == nil {
		return fmt.Errorf("cell '%s' does not exist", name)
	}
	return emitCell.in.EmitEvent(evt)
}

// Emitter implements Mesh.
func (m *mesh) Emitter(name string) (Emitter, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	emitCell := m.cells[name]
	if emitCell == nil {
		return nil, fmt.Errorf("cell '%s' does not exist", name)
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
