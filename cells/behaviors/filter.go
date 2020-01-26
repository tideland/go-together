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
// FILTER BEHAVIOR
//--------------------

// filterMode describes if the filter works selecting or excluding.
type filterMode int

// Flags for the filter.
const (
	selectFilter filterMode = iota
	excludeFilter
)

// Filter is a function type checking if an event shall be filtered.
type Filter func(evt *event.Event) (bool, error)

// filterBehavior is a simple repeater using the filter
// function to check if an event shall be selected or excluded
// for re-emitting.
type filterBehavior struct {
	id      string
	emitter mesh.Emitter
	mode    filterMode
	matches Filter
}

// NewSelectFilterBehavior creates a filter behavior based on the passed function.
// It re-emits every received event for which the filter function returns true.
func NewSelectFilterBehavior(id string, matches Filter) mesh.Behavior {
	return &filterBehavior{
		id:      id,
		mode:    selectFilter,
		matches: matches,
	}
}

// NewExcludeFilterBehavior creates a filter behavior based on the passed function.
// It re-emits every received event for which the filter function returns false.
func NewExcludeFilterBehavior(id string, matches Filter) mesh.Behavior {
	return &filterBehavior{
		id:      id,
		mode:    excludeFilter,
		matches: matches,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *filterBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *filterBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *filterBehavior) Terminate() error {
	return nil
}

// Process emits the event when the filter func returns true and the
// mode is select or it returns false and the mode is exclude.
func (b *filterBehavior) Process(evt *event.Event) {
	ok, err := b.matches(evt)
	fuse.Trigger(err)
	switch b.mode {
	case selectFilter:
		// Select those who match.
		if ok {
			fuse.Trigger(b.emitter.Broadcast(evt))
		}
	case excludeFilter:
		// Select those who don't match.
		if !ok {
			fuse.Trigger(b.emitter.Broadcast(evt))
		}
	}
}

// Recover from an error.
func (b *filterBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
