// Tideland Go Together - Cells - Behaviors
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
	"tideland.dev/go/trace/logger"
)

//--------------------
// LOGGER BEHAVIOR
//--------------------

// loggerBehavior is a behaior for the logging of events.
type loggerBehavior struct {
	id      string
	emitter mesh.Emitter
}

// NewLoggerBehavior creates a logging behavior. It logs emitted
// events with info level.
func NewLoggerBehavior(id string) mesh.Behavior {
	return &loggerBehavior{
		id: id,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *loggerBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *loggerBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *loggerBehavior) Terminate() error {
	return nil
}

// Process logs the event at info level.
func (b *loggerBehavior) Process(evt *event.Event) error {
	logger.Infof("(%s) logging event %v", b.id, evt)
	return nil
}

// Recover from an error. Can't even log, it's a logging problem.
func (b *loggerBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
