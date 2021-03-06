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

import ()

//--------------------
// EMITTER
//--------------------

// emitter allows the continuous emitting of events to a cell
// without having to resolve the cell name each time.
type emitter struct {
	strean *stream
}

// Emit implements Emitter.
func (e *emitter) Emit(evt *Event) error {
	return e.strean.Emit(evt)
}

// EOF
