// Tideland Go Together - Fuse
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package fuse // import "tideland.dev/go/together/fuse"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"strings"
	"time"
)

//--------------------
// REASON
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

//--------------------
// REASONS
//--------------------

// Reasons maintains a list of recovery reasons when working with
// panic recovering. It allows to get statistical information
// and maintenance of the collected reasons to decide if a
// recovery makes sense or a loop should be left with an error.
type Reasons struct {
	reasons []Reason
}

// Append adds a new reason with timestamp to the collection.
func (rs *Reasons) Append(reason interface{}) {
	rs.reasons = append(rs.reasons, Reason{
		Time:   time.Now(),
		Reason: reason,
	})
}

// Trim keeps only the defined number of reasons. This way the
// length can be kept constant.
func (rs *Reasons) Trim(keep int) {
	if keep < len(rs.reasons) {
		rs.reasons = rs.reasons[len(rs.reasons)-keep:]
	}
}

// Len returns the length of the recoverings.
func (rs Reasons) Len() int {
	return len(rs.reasons)
}

// Last returns the last appended reason.
func (rs Reasons) Last() Reason {
	return rs.reasons[len(rs.reasons)-1]
}

// Frequency checks if a given number of panics happened during
// a given duration.
func (rs Reasons) Frequency(num int, dur time.Duration) bool {
	if len(rs.reasons) >= num {
		first := rs.reasons[len(rs.reasons)-num].Time
		last := rs.reasons[len(rs.reasons)-1].Time
		return last.Sub(first) <= dur
	}
	return false
}

// String creates a string representation of the reasons.
func (rs Reasons) String() string {
	rss := make([]string, len(rs.reasons))
	for i, r := range rs.reasons {
		rss[i] = r.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(rss, " / "))
}

// EOF
