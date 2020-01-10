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
)

//--------------------
// TOPIC/PAYLOADS BEHAVIOR
//--------------------

// TopicPayloadsProcessor processes the collected payloads of a topic and
// returns a new payload to emit.
type TopicPayloadsProcessor func(topic string, payloads []*event.Payload) (*event.Payload, error)

// topicPayloadsBehavior collects and processes payloads by topic.
type topicPayloadsBehavior struct {
	id        string
	emitter   mesh.Emitter
	maximum   int
	collected map[string][]*event.Payload
	process   TopicPayloadsProcessor
}

// NewTopicPayloadsBehavior creates a behavior collecting the payloads
// of events by their topics, processes them, and emits the processed
// result.
func NewTopicPayloadsBehavior(id string, maximum int, processor TopicPayloadsProcessor) mesh.Behavior {
	return &topicPayloadsBehavior{
		id:        id,
		maximum:   maximum,
		collected: make(map[string][]*event.Payload),
		process:   processor,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *topicPayloadsBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *topicPayloadsBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *topicPayloadsBehavior) Terminate() error {
	return nil
}

// Process calls the processor function for the collected payloads
// by the events topic.
func (b *topicPayloadsBehavior) Process(evt *event.Event) error {
	topic := evt.Topic()
	pls := b.collected[topic]
	pls = append(pls, evt.Payload())
	if len(pls) > b.maximum {
		pls = pls[1:]
	}
	b.collected[topic] = pls
	pl, err := b.process(topic, b.collected[topic])
	if err != nil {
		return err
	}
	return b.emitter.Broadcast(event.New(topic, pl))
}

// Recover from an error.
func (b *topicPayloadsBehavior) Recover(err interface{}) error {
	b.collected = make(map[string][]*event.Payload)
	return nil
}

// EOF
