// Tideland Go Together - Cells - Mesh - Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"testing"

	"tideland.dev/go/audit/asserts"
)

//--------------------
// TESTS
//--------------------

// TestQueueSimple verifies simple appending and pulling of events
// via a queue.
func TestQueueSimple(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	q := newQueue(16)
	topics := []string{"one", "two", "three", "four", "five"}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case evt := <-q.Pull():
				assert.Contains(evt.Topic(), topics)
			}
		}
	}()

	for i := 0; i < 20; i++ {
		topic := topics[i%len(topics)]
		q.Append(NewEvent(topic, nil))
	}

	cancel()
}

// EOF
