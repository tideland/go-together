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
)

//--------------------
// CONDITION BEHAVIOR
//--------------------

// ConditionTester checks if an event matches a wanted state.
type ConditionTester func(evt *event.Event) bool

// ConditionProcessor handles the matching event.
type ConditionProcessor func(emitter mesh.Emitter, evt *event.Event) error

// conditionBehavior implements the condition behavior.
type conditionBehavior struct {
	id      string
	emitter mesh.Emitter
	test    ConditionTester
	process ConditionProcessor
}

// NewConditionBehavior creates a behavior testing of a cell
// fullfills a given condition. If the test returns true the
// processor is called.
func NewConditionBehavior(id string, tester ConditionTester, processor ConditionProcessor) mesh.Behavior {
	return &conditionBehavior{
		id:      id,
		test:    tester,
		process: processor,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *conditionBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *conditionBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *conditionBehavior) Terminate() error {
	return nil
}

// Process checks the condition.
func (b *conditionBehavior) Process(evt *event.Event) {
	if b.test(evt) {
		fuse.Trigger(b.process(b.emitter, evt))
	}
}

// Recover from an error.
func (b *conditionBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
