// Tideland Go Together - Loop
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package loop // import "tideland.dev/go/together/loop"

//--------------------
// IMPORTS
//--------------------

import (
	"context"

	"tideland.dev/go/together/notifier"
	"tideland.dev/go/trace/failure"
)

//--------------------
// OPTIONS
//--------------------

// Option defines the signature of an option setting function.
type Option func(loop *Loop) error

// WithContext allows to pass a context for cancellation or timeout.
func WithContext(ctx context.Context) Option {
	return func(loop *Loop) error {
		if ctx == nil {
			return failure.New("invalid loop option: context is nil")
		}
		loop.closer.Add(ctx.Done())
		return nil
	}
}

// WithRecoverer defines the panic handler of a loop.
func WithRecoverer(recoverer Recoverer) Option {
	return func(loop *Loop) error {
		if recoverer == nil {
			return failure.New("invalid loop option: recoverer is nil")
		}
		loop.recoverer = recoverer
		return nil
	}
}

// WithSignalbox adds a signalbox to make external monitors aware of
// the Loop internal status.
func WithSignalbox(signalbox *notifier.Signalbox) Option {
	return func(loop *Loop) error {
		if signalbox == nil {
			return failure.New("invalid loop option: signalbox is nil")
		}
		loop.bundle.Add(signalbox)
		return nil
	}
}

// WithFinalizer sets a function for finalizing the
// work of a Loop.
func WithFinalizer(finalizer Finalizer) Option {
	return func(loop *Loop) error {
		if finalizer == nil {
			return failure.New("invalid loop option: finalizer is nil")
		}
		loop.finalizer = finalizer
		return nil
	}
}

// EOF
