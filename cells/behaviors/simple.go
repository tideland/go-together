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
	"tideland.dev/go/together/fuse"
	"tideland.dev/go/trace/logger"
)

//--------------------
// SIMPLE BEHAVIOR
//--------------------

// SimpleProcessor is a function type doing the event processing.
type SimpleProcessor func(emitter mesh.Emitter, evt *event.Event) error

// simpleBehavior is a simple event processor using the processor
// function for its own logic.
type simpleBehavior struct {
	id      string
	emitter mesh.Emitter
	process SimpleProcessor
}

// NewSimpleProcessorBehavior creates a behavior based on the passed function.
// Instead of an own logic and an own state it uses the passed simple processor
// function for the event processing.
func NewSimpleProcessorBehavior(id string, processor SimpleProcessor) mesh.Behavior {
	if processor == nil {
		processor = func(emitter mesh.Emitter, evt *event.Event) error {
			logger.Errorf("simple processor %q used without function to handle event %v", id, evt)
			return nil
		}
	}
	return &simpleBehavior{
		id:      id,
		process: processor,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *simpleBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *simpleBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *simpleBehavior) Terminate() error {
	return nil
}

// Process calls the simple processor function.
func (b *simpleBehavior) Process(evt *event.Event) {
	fuse.Trigger(b.process(b.emitter, evt))
}

// Recover from an error.
func (b *simpleBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
