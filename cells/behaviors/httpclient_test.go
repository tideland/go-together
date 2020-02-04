// Tideland Go Together - Cells - Behaviors - Unit Tests
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors_test // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/json"
	"net/http"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/environments"
	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestHTTPClientBehaviorGet tests the HTTP client behavior, here
// the GET method.
func TestHTTPClientBehaviorGet(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := initWebAsserter(assert)
	plant := mesh.NewTestPlant(assert, behaviors.NewHTTPClientBehavior("cb"), 1)
	defer plant.Stop()

	plant.Emit(event.New(behaviors.TopicHTTPGet, "id", "simple", "url", wa.URL()+"/simple"))
	plant.Emit(event.New(behaviors.TopicHTTPGet, "id", "nested", "url", wa.URL()+"/nested"))

	plant.AssertFirst(0, func(evt *event.Event) bool {
		return evt.Topic() == behaviors.TopicHTTPGetReply &&
			evt.Payload().At("status-code").AsInt(0) == 200 &&
			evt.Payload().At("data").AsString("-") == "Done!"
	})
	plant.AssertLast(0, func(evt *event.Event) bool {
		data := evt.Payload().At("data").AsPayload()
		return evt.Topic() == behaviors.TopicHTTPGetReply &&
			evt.Payload().At("status-code").AsInt(0) == 200 &&
			data.At("A").AsString("-") == "Foo" &&
			data.At("B").AsInt(0) == 1234 &&
			data.At("C", "0", "E").AsString("-") == "Bar" &&
			data.At("C", "0", "F").AsBool(false) &&
			data.At("C", "0", "G").AsInt(0) == 10
	})

	// assert.Equal(data.At("C", "1", "E").AsString("-"), "Baz")
	// assert.Equal(data.At("C", "1", "F").AsBool(true), false)
	// assert.Equal(data.At("C", "1", "G").AsInt(0), 20)
	// assert.Equal(data.At("C", "2", "E").AsString("-"), "Yadda")
	// assert.Equal(data.At("C", "2", "F").AsBool(false), true)
	// assert.Equal(data.At("C", "2", "G").AsInt(0), 30)
}

//--------------------
// TESTS
//--------------------

type inner struct {
	E string
	F bool
	G int
}

type outer struct {
	A string
	B int
	C []inner
}

func initSimpleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(environments.HeaderContentType, environments.ContentTypePlain)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Done!"))
	}
}

func initNestedHandler() http.HandlerFunc {
	v := outer{
		A: "Foo",
		B: 1234,
		C: []inner{
			{
				E: "Bar",
				F: true,
				G: 10,
			},
			{
				E: "Baz",
				F: false,
				G: 20,
			},
			{
				E: "Yadda",
				F: true,
				G: 30,
			},
		},
	}
	b, _ := json.Marshal(v)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(environments.HeaderContentType, environments.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

func initWebAsserter(assert *asserts.Asserts) *environments.WebAsserter {
	wa := environments.NewWebAsserter(assert)

	wa.Handle("/simple", initSimpleHandler())
	wa.Handle("/nested", initNestedHandler())

	return wa
}

// EOF
