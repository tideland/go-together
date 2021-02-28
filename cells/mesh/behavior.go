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
)

//--------------------
// BEHAVIOR
//--------------------

// Behavior describes what cell implementations must understand.
type Behavior interface {
	// Go will be started as wrapped goroutine. It's the responsible
	// of the implementation to run a select loop, receive incomming
	// events via the input queue, and emit events via the output queue
	// if needed.
	Go(ctx context.Context, name string, in InputStream, out OutputStream)
}

// EOF
