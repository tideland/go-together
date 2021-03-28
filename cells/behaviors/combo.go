// Tideland Go Together - Cells - Behaviors
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"

	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// COMBO BEHAVIOR
//--------------------

// ComboCriterionFunc is used by the combo behavior. It has to return
// CriterionDone when a combination is complete, CriterionKeep when it
// is so far okay but not yet complete, CriterionDropFirst when the first
// event shall be dropped, CriterionDropLast when the last event shall
// be dropped, and CriterionClear when the collected events have to be
// cleared for starting over. In case of CriterionDone it additionally
// has to return a payload which will be emitted.
type ComboCriterionFunc func(r mesh.EventSinkReader) (CriterionMatch, interface{})

// comboBehavior implements the combo behavior.
type comboBehavior struct {
	matches ComboCriterionFunc
	sink    *mesh.EventSink
}

// NewComboBehavior creates a behavior checking an event stream for a
// combination of events defined by a criterion. In case of a matching
// situation an according event is emitted.
func NewComboBehavior(matcher ComboCriterionFunc) mesh.Behavior {
	return &comboBehavior{
		matches: matcher,
		sink:    mesh.NewEventSink(0),
	}
}

// Go collects the events and runs the criterion matcher to check them.
func (b *comboBehavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			switch evt.Topic() {
			case TopicReset:
				b.sink.Clear()
			default:
				b.sink.Push(evt)
				matches, pl := b.matches(b.sink)
				switch matches {
				case CriterionDone:
					out.Emit(TopicCriterionDone, pl)
				case CriterionKeep:
				case CriterionDropFirst:
					b.sink.PullFirst()
				case CriterionDropLast:
					b.sink.PullLast()
				default:
					return fmt.Errorf("invalid criterion matcher result: %v", matches)
				}
			}
		}
	}
}

// EOF
