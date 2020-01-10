// Tideland Go Together - Cells - Mesh - Unit Tests
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
// FUNCTIONS
//--------------------

func GetCellQueueCap(m *Mesh, id string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.cells[id].cell.act.QueueCap()
}

// EOF
