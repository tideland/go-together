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
	"net/http"

	"tideland.dev/go/net/httpx"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
	"tideland.dev/go/together/fuse"
)

//--------------------
// HTTP CLIENT BEHAVIOR
//--------------------

// httpClientBehavior performs HTTP requests.
type httpClientBehavior struct {
	id      string
	emitter mesh.Emitter
}

// NewHTTPClientBehavior performs HTTP request and transforms the
// response into emitted payload, depending on the content-type.
func NewHTTPClientBehavior(id string) mesh.Behavior {
	return &httpClientBehavior{
		id: id,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *httpClientBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *httpClientBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *httpClientBehavior) Terminate() error {
	return nil
}

// Process performs the HTTP request.
func (b *httpClientBehavior) Process(evt *event.Event) {
	switch evt.Topic() {
	case TopicHTTPGet:
		fuse.Trigger(b.processGet(evt))
	}
}

// Recover from an error.
func (b *httpClientBehavior) Recover(err interface{}) error {
	return nil
}

// processGet handles the GET request.
func (b *httpClientBehavior) processGet(evt *event.Event) error {
	id := evt.Payload().At("id").AsString("<none>")
	url := evt.Payload().At("url").AsString("")
	resp, err := http.Get(url)
	b.broadcastReply(TopicHTTPGetReply, id, url, resp, err)
	return nil
}

// broadcastReply reads the reply and broadcasts it.
func (b *httpClientBehavior) broadcastReply(topic, id, url string, resp *http.Response, err error) {
	headerToPayload := func() *event.Payload {
		var hplvs []interface{}
		for key := range resp.Header {
			hplvs = append(hplvs, key, resp.Header.Get(key))
		}
		return event.NewPayload(hplvs...)
	}

	var plvs []interface{}

	plvs = append(plvs, "id", id)
	plvs = append(plvs, "url", url)

	if err != nil {
		plvs = append(plvs, "error", err)
	} else {
		var data interface{}
		err = httpx.UnmarshalBody(resp.Body, resp.Header, &data)
		if err != nil {
			plvs = append(plvs, "status-code", resp.StatusCode)
			plvs = append(plvs, "header", headerToPayload())
			plvs = append(plvs, "error", err)
		} else {
			plvs = append(plvs, "status-code", resp.StatusCode)
			plvs = append(plvs, "header", headerToPayload())
			plvs = append(plvs, "data", data)
		}
	}
	_ = b.emitter.Broadcast(event.New(topic, plvs...))
}

// EOF
