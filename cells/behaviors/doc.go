// Tideland Go Together - Cells - Behaviors
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package behaviors provides several generic and always useful
// standard behaviors for the Tideland Go Library Cells. They are
// simply created by calling NewXyzBehavior(). Their configuration
// is done using constructor arguments. Some of them take functions
// implementations of interfaces to control their processing. These
// behaviors are:
//
// Aggregator aggregates events and emits each aggregated value.
//
// Broadcaster simply emits received events to all subscribers.
//
// Callback calls a number of passed functions for each received event.
//
// Collector collects events which can be processed on demand.
//
// Combo waits for a user-defined combination of events.
//
// Condition tests events for conditions using a tester function
// and calls a processor then.
//
// Countdown counts a defined number of events down and processes
// the collected ones at zero.
//
// Counter counts events, the counters can be retrieved.
//
// Cronjob performs a given function with an emitter every given duration.
//
// Evaluator evaluates events based on a user-defined function which
// returns a rating.
//
// Filter emits received events based on a user-defined filter.
//
// Finite State Machine allows to build finite state machines for events.
//
// Logger logs received events with level INFO.
//
// Mapper maps received events based on a user-defined function to
// new events.
//
// Mesh Router allows to create a list of cell IDs where the received
// event is then routed to.
//
// Once calls the once function only for the first event it receives.
//
// Pair checks if the event stream contains two matching ones based on a
// user-based criterion in a given timespan.
//
// Rate measures times between a number of criterion fitting events and
// emits the result.
//
// Rate Window checks if a number of events in a given timespan matches
// a given criterion.
//
// Round Robin distribtes the received events round robin to the subscribed
// cells.
//
// Router allows to create a list of subscriber cell IDs where the received
// event is then routed to.
//
// Sequence checks the event stream for a defined sequence of events
// discovered by a user-defined criterion.
//
// Simple Processor allows to not implement a behavior but only use
// one function for event processing.
//
// Ticker emits tick events in a defined interval.
//
// Topic/Payloads collects the payloads of events by topics, processes
// those, and emits the processed result.
package behaviors // import "tideland.dev/go/together/cells/behaviors"

// EOF
