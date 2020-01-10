// Tideland Go Together - Limiter
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package limiter // import "tideland.dev/go/together/limiter"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
)

//--------------------
// LIMITER
//--------------------

// Job describes a simple function that can be ran by the Limiter.
type Job func() error

// Limiter allows to run only a defined number of jobs at the same time.
type Limiter struct {
	active chan struct{}
}

// New creates a Limiter instance with the passed job limit.
func New(limit int) *Limiter {
	return &Limiter{
		active: make(chan struct{}, limit),
	}
}

// Do executes the passed job if the limit isn't reached and the context
// contains is active and contains no error.
func (l *Limiter) Do(ctx context.Context, job Job) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case l.active <- struct{}{}:
		defer func() {
			<-l.active
		}()
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return job()
	}
}

// EOF
