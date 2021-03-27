// Tideland Go Together - Cells - Mesh
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package mesh is the runtime package of the Tideland cells event processing.
// It provides types for meshed cells running individual behaviors.
//
// These behaviors are defined based on an interface and can be added to the
// mesh. Here they are running concurrently and can be networked to communicate
// via events. Several useful behaviors are already provided with the behaviors
// package.
//
// New meshes are created with
//
//     msh := mesh.New()
//
// and cells are started with
//
//    msh.Go("foo", NewFooBehavior())
//    msh.Go("bar", NewBarBehavior())
//    msh.Go("baz", NewBazBehavior())
//
// These cells can subscribe each other with
//
//    msh.Subscribe("foo", "bar")
//    msh.Subscribe("foo", "baz")
//
// so that events which are emitted by the cell "foo" will be
// received by the cells "bar" and "baz". Each cell can subscribe
// to multiple other subscribers and even circular subscriptions are
// no problem. But handle with care.
//
// Events from the outside are emitted using
//
//     msh.Emit("foo", mesh.NewEvent("foo", "answer", 42))
//
// In case of many emits to one cell you can get an emitter
// with
//
//     emtr, err := msh.Emitter("foo")
//
// and
//
//     emtrEmit(mesh.NewEvent("foo", "answer", 42))
//
package mesh // import "tideland.dev/go/together/cells/mesh"

// EOF
