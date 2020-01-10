// Tideland Go Together - Cells - Event
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package event provides the event and its payload emitted by
// the user of a mesh and the contained cells. It is read-only
// to avoid modifications while multiple behaviors are using it
// concurrently. Accessors and cloners make live for behavior
// developers as well easier as the sink and its methods for
// analyzing.
package event // import "tideland.dev/go/together/cells/event"

// EOF
