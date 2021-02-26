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
	"sync"
)

//--------------------
// MESH
//--------------------

// Mesh manages a closed network of cells.
type Mesh struct {
	mu sync.RWMutex
}

// New creates new Mesh instance.
func New() *Mesh {
	m := &Mesh{}
	return m
}

// EOF
