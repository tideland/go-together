// Tideland Go Together - Actor
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package actor // import "tideland.dev/go/together/actor"

//--------------------
// IMPORTS
//--------------------

import (
	"context"

	"tideland.dev/go/together/loop"
	"tideland.dev/go/together/notifier"
)

//--------------------
// OPTIONS
//--------------------

// Option defines the signature of an option setting function.
type Option func(act *Actor) error

// WithContext allows to pass a context for cancellation or timeout.
func WithContext(ctx context.Context) Option {
	return func(act *Actor) error {
		act.options = append(act.options, loop.WithContext(ctx))
		return nil
	}
}

// WithQueueCap defines the channel capacity for actions sent to an Actor.
func WithQueueCap(c int) Option {
	return func(act *Actor) error {
		if c < 1 {
			c = 1
		}
		act.actionC = make(chan Action, c)
		return nil
	}
}

// WithRecoverer defines the panic handler of an actor.
func WithRecoverer(rcvr Recoverer) Option {
	return func(act *Actor) error {
		act.options = append(act.options, loop.WithRecoverer(loop.Recoverer(rcvr)))
		return nil
	}
}

// WithSignalbox add a notifier to make external monitors aware of
// the Actors internal status.
func WithSignalbox(signalbox *notifier.Signalbox) Option {
	return func(act *Actor) error {
		act.options = append(act.options, loop.WithSignalbox(signalbox))
		return nil
	}
}

// WithFinalizer sets a function for finalizing the
// work of a Loop.
func WithFinalizer(finalizer loop.Finalizer) Option {
	return func(act *Actor) error {
		act.options = append(act.options, loop.WithFinalizer(finalizer))
		return nil
	}
}

// EOF
