// Tideland Go Together - Cells - Mesh
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package mesh is the runtime package of the Tideland cells event processing.
// It provides types for meshed cells running individual behaviors.
//
// These behaviors are defined based on an interface and can added to the
// mesh. Here they are running concurrently and can be networked and communicate
// via events. Several useful behaviors are already provided with the behaviors
// package.
//
// New meshes are created with
//
//     msh := mesh.New()
//
// and cells are started with
//
//    msh.SpawnCells(
//        NewFooer("a"),
//        NewBarer("b"),
//        NewBazer("c"),
//    )
//
// These cells can subscribe each other with
//
//    msh.Subscribe("a", "b", "c")
//
// so that events which are emitted by the cell "a" will be
// received by the cells "b" and "c". Each cell can subscribe
// to multiple other subscribers and even circular subscriptions are
// no problem. But handle with care.
//
// Events from the outside are emitted using
//
//     msh.Emit("foo", event.New("foo", "answer", 42))
//
package mesh // import "tideland.dev/go/together/cells/mesh"

// EOF
