// Tideland Go Together - Cells - Mesh - Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/together/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"sync"
	"testing"

	"tideland.dev/go/audit/asserts"
)

//--------------------
// TESTS
//--------------------

// TestStreamSimple verifies simple emitting and pulling of events
// via a stream.
func TestStreamSimple(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	str := newStream()
	topics := []string{"one", "two", "three", "four", "five"}

	var wg sync.WaitGroup

	wg.Add(20)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case evt := <-str.Pull():
				assert.Contains(evt.Topic(), topics)
				wg.Done()
			}
		}
	}()

	for i := 0; i < 20; i++ {
		topic := topics[i%len(topics)]
		err := str.Emit(topic)
		assert.NoError(err)
	}

	wg.Wait()
	cancel()
}

// EOF
