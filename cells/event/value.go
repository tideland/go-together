// Tideland Go Together - Cells - Event
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package event // import "tideland.dev/go/together/cells/event"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"strconv"
	"time"
)

//--------------------
// CONSTANTS
//--------------------

// Defaults.
const (
	DefaultKey   = "value"
	DefaultValue = true
)

// Different formats for the parsing of strings
// into times.
var timeFormats = []string{
	"Mon Jan 2 15:04:05 -0700 MST 2006",
	"2006-01-02 15:04:05.999999999 -0700 MST",
	time.ANSIC,
	time.Kitchen,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RubyDate,
	time.Stamp,
	time.StampMicro,
	time.StampMilli,
	time.StampNano,
	time.UnixDate,
}

// For time converts.
const (
	maxInt = int64(^uint(0) >> 1)
)

//--------------------
// VALUE
//--------------------

// Value contains one payload value.
type Value struct {
	raw interface{}
	err error
}

// IsDefined returns true if this value is defined.
func (v *Value) IsDefined() bool {
	return v.raw != nil
}

// IsUndefined returns true if this value is undefined.
func (v *Value) IsUndefined() bool {
	return v.raw == nil
}

// IsPayload returns true if the value is a payload.
func (v *Value) IsPayload() bool {
	if v.raw == nil {
		return false
	}
	_, ok := v.raw.(*Payload)
	return ok
}

// AsString returns the value as string, dv is taken as default value.
func (v *Value) AsString(dv string) string {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.raw.(type) {
	case string:
		return tv
	case int:
		return strconv.Itoa(tv)
	case float64:
		return strconv.FormatFloat(tv, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(tv)
	case time.Time:
		return tv.Format(time.RFC3339Nano)
	case time.Duration:
		return tv.String()
	case *Payload:
		return tv.String()
	}
	return dv
}

// AsInt returns the value as int, dv is taken as default value.
func (v *Value) AsInt(dv int) int {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.raw.(type) {
	case string:
		i, err := strconv.Atoi(tv)
		if err != nil {
			return dv
		}
		return i
	case int:
		return tv
	case float64:
		return int(tv)
	case bool:
		if tv {
			return 1
		}
		return 0
	case time.Time:
		ns := tv.UnixNano()
		if ns > maxInt {
			return dv
		}
		return int(ns)
	case time.Duration:
		ns := tv.Nanoseconds()
		if ns > maxInt {
			return dv
		}
		return int(ns)
	}
	return dv
}

// AsFloat64 returns the value as float64, dv is taken as default value.
func (v *Value) AsFloat64(dv float64) float64 {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.raw.(type) {
	case string:
		f, err := strconv.ParseFloat(tv, 64)
		if err != nil {
			return dv
		}
		return f
	case int:
		return float64(tv)
	case float64:
		return tv
	case bool:
		if tv {
			return 1.0
		}
		return 0.0
	case time.Time:
		ns := tv.UnixNano()
		return float64(ns)
	case time.Duration:
		ns := tv.Nanoseconds()
		return float64(ns)
	}
	return dv
}

// AsBool returns the value as bool, dv is taken as default value.
func (v *Value) AsBool(dv bool) bool {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.raw.(type) {
	case string:
		b, err := strconv.ParseBool(tv)
		if err != nil {
			return dv
		}
		return b
	case int:
		return tv == 1
	case float64:
		return tv == 1.0
	case bool:
		return tv
	case time.Time:
		return tv.UnixNano() > 0
	case time.Duration:
		return tv.Nanoseconds() > 0
	}
	return dv
}

// AsTime returns the value as time, dv is taken as default value.
func (v *Value) AsTime(dv time.Time) time.Time {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.raw.(type) {
	case string:
		for _, timeFormat := range timeFormats {
			t, err := time.Parse(timeFormat, tv)
			if err == nil {
				return t
			}
		}
		return dv
	case int:
		return time.Time{}.Add(time.Duration(tv))
	case float64:
		d := int64(tv)
		return time.Time{}.Add(time.Duration(d))
	case bool:
		if tv {
			return time.Time{}.Add(1)
		}
		return time.Time{}
	case time.Time:
		return tv
	case time.Duration:
		return time.Time{}.Add(tv)
	}
	return dv
}

// AsDuration returns the value as duration, dv is taken as default value.
func (v *Value) AsDuration(dv time.Duration) time.Duration {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.raw.(type) {
	case string:
		d, err := time.ParseDuration(tv)
		if err == nil {
			return d
		}
		return dv
	case int:
		return time.Duration(tv)
	case float64:
		d := int64(tv)
		return time.Duration(d)
	case bool:
		if tv {
			return 1
		}
		return 0
	case time.Time:
		return time.Duration(tv.UnixNano())
	case time.Duration:
		return tv
	}
	return dv
}

// AsPayload returns the value as payload.
func (v *Value) AsPayload() *Payload {
	if v.IsUndefined() {
		return NewPayload()
	}
	switch tv := v.raw.(type) {
	case *Payload:
		return tv
	default:
		return &Payload{
			values: map[string]interface{}{
				DefaultKey: tv,
			},
		}
	}
}

// AsPayloadChan returns the value as payload channel.
func (v *Value) AsPayloadChan() PayloadChan {
	if v.IsUndefined() {
		return nil
	}
	tv, ok := v.raw.(PayloadChan)
	if !ok {
		return nil
	}
	return tv
}

// String implements fmt.Stringer.
func (v *Value) String() string {
	return fmt.Sprintf("%v", v.raw)
}

// Error implements error.
func (v *Value) Error() string {
	return fmt.Sprintf("%v", v.err)
}

// EOF
