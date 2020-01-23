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
	"time"

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
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()

	trigger := func(emitter mesh.Emitter, evt *event.Event) error {
		if evt.Payload().At("id").AsString("<unknown>") == "nested" {
			assert.OK(emitter.Broadcast(event.New(event.TopicProcess)))
		}
		return nil
	}
	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		assert.OK(accessor.Do(func(index int, evt *event.Event) error {
			assert.Equal(evt.Topic(), behaviors.TopicHTTPGetReply)
			assert.Equal(evt.Payload().At("status-code").AsInt(0), 200)
			switch evt.Payload().At("id").AsString("-") {
			case "simple":
				data := evt.Payload().At("data").AsString("-")
				assert.Equal(data, "Done!")
			case "nested":
				data := evt.Payload().At("data").AsPayload()
				assert.Length(data, 3)
				assert.Equal(data.At("A").AsString("-"), "Foo")
				assert.Equal(data.At("B").AsInt(0), 1234)
				assert.Equal(data.At("C", "0", "E").AsString("-"), "Bar")
				assert.Equal(data.At("C", "0", "F").AsBool(false), true)
				assert.Equal(data.At("C", "0", "G").AsInt(0), 10)
				assert.Equal(data.At("C", "1", "E").AsString("-"), "Baz")
				assert.Equal(data.At("C", "1", "F").AsBool(true), false)
				assert.Equal(data.At("C", "1", "G").AsInt(0), 20)
				assert.Equal(data.At("C", "2", "E").AsString("-"), "Yadda")
				assert.Equal(data.At("C", "2", "F").AsBool(false), true)
				assert.Equal(data.At("C", "2", "G").AsInt(0), 30)
			}
			return nil
		}))
		sigc <- accessor.Len()
		return nil, nil
	}

	assert.OK(msh.SpawnCells(
		behaviors.NewHTTPClientBehavior("client"),
		behaviors.NewSimpleProcessorBehavior("trigger", trigger),
		behaviors.NewCollectorBehavior("collector", 10, processor),
	))
	assert.OK(msh.Subscribe("client", "collector", "trigger"))
	assert.OK(msh.Subscribe("trigger", "collector"))

	assert.OK(msh.Emit("client", event.New(behaviors.TopicHTTPGet, "id", "simple", "url", wa.URL()+"/simple")))
	assert.OK(msh.Emit("client", event.New(behaviors.TopicHTTPGet, "id", "nested", "url", wa.URL()+"/nested")))

	assert.Wait(sigc, 2, time.Second)
	assert.OK(msh.Stop())
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
