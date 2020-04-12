// Tideland Go Together - Fuse
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package fuse // import "tideland.dev/go/together/fuse"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/trace/failure"
	"tideland.dev/go/trace/location"
)

//--------------------
// TRIGGER
//--------------------

// Trigger raises an annotated panic in case of an error.
func Trigger(err error) {
	if err != nil {
		code := location.At(2).Code("PANIC")
		panic(failure.Annotate(err, code))
	}
}

// EOF
