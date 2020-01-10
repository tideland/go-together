// Tideland Go Together - Loop
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package loop // import "tideland.dev/go/together/loop"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"strings"
	"time"
)

//--------------------
// REASONS
//--------------------

// Reason stores time and reason of one recovering.
type Reason struct {
	Time   time.Time
	Reason interface{}
}

// String creates a string representation of the reason.
func (r Reason) String() string {
	tbs, err := r.Time.MarshalText()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return fmt.Sprintf("['%v' @ %s]", r.Reason, string(tbs))
}

// Reasons maintains a list of recovery reasons a loop has.
type Reasons []Reason

// MakeReasons creates an initial buffer for collecting
// reasons.
func MakeReasons() Reasons {
	return Reasons{}
}

// Append adds a new reason with timestamp to the collection.
func (rs Reasons) Append(reason interface{}) Reasons {
	return append(rs, Reason{
		Time:   time.Now(),
		Reason: reason,
	})
}

// Len returns the length of the recoverings.
func (rs Reasons) Len() int {
	return len(rs)
}

// Last returns the last appended reason.
func (rs Reasons) Last() Reason {
	return rs[len(rs)-1]
}

// Trim returns the last resons defined by l. This way the
// length can be kept constant.
func (rs Reasons) Trim(l int) Reasons {
	if l >= len(rs) {
		return rs
	}
	return rs[len(rs)-l:]
}

// Frequency checks if a given number of panics happened during
// a given duration.
func (rs Reasons) Frequency(num int, dur time.Duration) bool {
	if len(rs) >= num {
		first := rs[len(rs)-num].Time
		last := rs[len(rs)-1].Time
		return last.Sub(first) <= dur
	}
	return false
}

// String creates a string representation of the reasons.
func (rs Reasons) String() string {
	rss := make([]string, len(rs))
	for i, r := range rs {
		rss[i] = r.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(rss, " / "))
}

// EOF
