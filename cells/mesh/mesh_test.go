// Tideland Go Together - Cells - Mesh - Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestNewMesh verifies the simple creation of a mesh.
func TestNewMesh(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	msh := mesh.New(ctx)

	assert.NotNil(msh)

	cancel()
}

// TestMeshGo verifies the starting of a cell via mesh.
func TestMeshGo(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan interface{})
	behaviorFunc := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		sigc <- cell.Name()
		return nil
	}
	msh := mesh.New(ctx)

	msh.Go("testing", mesh.BehaviorFunc(behaviorFunc))

	assert.Wait(sigc, "testing", time.Second)

	cancel()
}

// TestMeshEmit verifies the emitting of events to one cell.
func TestMeshEmit(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan interface{})
	behaviorFunc := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		i := 0
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				i++
				if evt.Topic() == "get-i" {
					sigc <- i
				}
			}
		}
	}
	msh := mesh.New(ctx)
	err := msh.Emit("testing", mesh.NewEvent("one"))
	assert.ErrorContains(err, "cell \"testing\" does not exist")

	msh.Go("testing", mesh.BehaviorFunc(behaviorFunc))

	err = msh.Emit("testing", mesh.NewEvent("one"))
	assert.NoError(err)

	msh.Emit("testing", mesh.NewEvent("two"))
	msh.Emit("testing", mesh.NewEvent("three"))
	msh.Emit("testing", mesh.NewEvent("get-i"))

	assert.Wait(sigc, 4, time.Second)

	cancel()
}

// TestMeshEmitter verifies the emitting of events to one cell using an emitter.
func TestMeshEmitter(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan interface{})
	behaviorFunc := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		i := 0
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				i++
				if evt.Topic() == "get-i" {
					sigc <- i
				}
			}
		}
	}
	msh := mesh.New(ctx)
	emtr, err := msh.Emitter("testing")
	assert.ErrorContains(err, "cell \"testing\" does not exist")

	msh.Go("testing", mesh.BehaviorFunc(behaviorFunc))
	emtr, err = msh.Emitter("testing")
	assert.NoError(err)

	emtr.Emit(mesh.NewEvent("one"))
	emtr.Emit(mesh.NewEvent("two"))
	emtr.Emit(mesh.NewEvent("three"))
	emtr.Emit(mesh.NewEvent("get-i"))

	assert.Wait(sigc, 4, time.Second)

	cancel()
}

// EOF
