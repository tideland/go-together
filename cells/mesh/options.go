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

import ()

//--------------------
// OPTIONS
//--------------------

// Option is a function able to configure a mesh.
type Option func(msh *Mesh)

// QueueCap allows to configure the default queue capacity
// for cells inside a mesh.
func QueueCap(qc int) Option {
	return func(msh *Mesh) {
		if qc > 1 {
			msh.queueCap = qc
		}
	}
}

// EOF
