// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"errors"

	"github.com/qioalice/ekago/v4/ekaenc"
)

//goland:noinspection GoSnakeCaseUsage
var (
	_ERR_NIL_TIME_RECEIVER = errors.New("nil ekatime.Date receiver")
	_ERR_NOT_ISO8601_TIME  = errors.New("incorrect ISO8601 time format (must be hhmmss, hh:mm:ss or without second)")
	_ERR_BAD_HOUR          = errors.New("hours must be in the range [0..23]")
	_ERR_BAD_MINUTE        = errors.New("minutes must be in the range [0..59]")
	_ERR_BAD_SECOND        = errors.New("seconds must be in the range [0..59]")
	_ERR_BAD_CORRESP_TIME  = errors.New("time must represent valid time (24h format)")
	_ERR_BAD_JSON_TIME_QUO = errors.New("bad JSON ISO8601 time representation (forgotten quotes?)")
)

// AppendTo generates a string representation of Time and adds it to the b,
// returning a new slice (if it has grown) or the same if there was enough
// space to store 8 bytes (of string representation).
//
// Uses separator as Time's parts separator:
// "hh<separator>mm<separator>ss".
//
// Up to 6x faster than fmt.Sprintf().
func (t Time) AppendTo(b []byte, separator byte) []byte {

	f := func(b []byte) {
		hh, mm, ss := normalizeTime(t.Split())

		copy(b[0:2], _TIME_PART_AS_NUM_STR[hh])
		copy(b[3:5], _TIME_PART_AS_NUM_STR[mm])
		copy(b[6:8], _TIME_PART_AS_NUM_STR[ss])

		b[2] = separator
		b[5] = separator
	}

	if c, l := cap(b), len(b); c-l >= 8 {
		b = b[:l+8]
		f(b[l:])
		return b
	} else {
		// One more allocation
		b2 := make([]byte, 8)
		f(b2)
		b = append(b, b2...)
		return b
	}
}

// ParseFrom tries to parse b considering with the following format:
// "hh<separator>mm<separator>ss", <separator> may be any 1 byte,
// or not presented at all. But if separator was presented between in "hh<sep>mm"
// it must be also between other parts and vice-versa.
//
// Read more:
// https://en.wikipedia.org/wiki/ISO_8601
//
// Skips leading spaces, ignores all next data after the day has been scanned.
func (t *Time) ParseFrom(b []byte) error {

	if t == nil {
		return _ERR_NIL_TIME_RECEIVER
	}

	var i = 0
	for n := len(b); i < n && b[i] <= ' '; i++ {
	}

	// Minimum required len: 4 (hhmm - w/o separators, w/o seconds).
	if len(b[i:]) < 4 {
		return _ERR_NOT_ISO8601_TIME
	}

	x, valid := batoi(b[i], b[i+1])
	if !valid || x < 0 || x > 23 {
		return _ERR_BAD_HOUR
	}

	i += 2
	hh := Hour(x)
	wasSeparator := false

	// Skip separator
	if !(b[i] >= '0' && b[i] <= '9') {
		i++
		wasSeparator = true
	}

	// At this code point, len(b) may == 1. Check it.
	if len(b[i:]) == 1 {
		return _ERR_NOT_ISO8601_TIME
	}

	x, valid = batoi(b[i], b[i+1])
	if !valid || x < 0 || x > 59 {
		return _ERR_BAD_MINUTE
	}

	i += 2
	mm := Minute(x)
	ss := Second(0)

	// At this code point user may provide "hhmm" w/o seconds.
	// Check whether seconds are provided.
	if l := len(b[i:]); l > 0 {
		// We need 2 symbols if there was no separator, or 3 symbols otherwise.
		if (l == 1 && !wasSeparator) || (l == 2 && wasSeparator) {
			return _ERR_NOT_ISO8601_TIME
		}
		if wasSeparator {
			i++
		}
		x, valid = batoi(b[i], b[i+1])
		if !valid || x < 0 || x > 59 {
			return _ERR_BAD_SECOND
		}
		ss = Second(x)
	}

	if !IsValidTime(hh, mm, ss) {
		return _ERR_BAD_CORRESP_TIME
	}

	*t = NewTime(hh, mm, ss)
	return nil
}

// String returns the current Time's string representation in the following format:
// "hh:mm:ss".
func (t Time) String() string {
	return string(t.AppendTo(make([]byte, 0, 8), ':'))
}

// MarshalJSON encodes the current Time in the following format (quoted)
// "hh:mm:ss", and returns it. Always returns nil as error.
//
// JSON null supporting:
// - Writes JSON null if current Time receiver == nil,
// - Writes JSON null if current Time == 0.
func (t *Time) MarshalJSON() ([]byte, error) {

	if t == nil || *t == 0 {
		return ekaenc.NullAsBytesLowerCase(), nil
	}

	b := make([]byte, 10)
	_ = t.AppendTo(b[1:1:10], ':')

	b[0] = '"'
	b[9] = '"'

	return b, nil
}

// UnmarshalJSON decodes b into the current Time object expecting b contains
// ISO8601 quoted time (only time, not time w/ date) in the one of the following
// formats: "hhmm", "hh:mm", "hhmmss", "hh:mm:ss".
//
// JSON null supporting:
// - It's ok if there is JSON null and receiver == nil (nothing changes)
// - Zeroes Time if there is JSON null and receiver != nil.
//
// In other cases JSON parsing error or Time.ParseFrom() error is returned.
func (t *Time) UnmarshalJSON(b []byte) error {

	if ekaenc.IsNullAsBytes(b) {
		if t != nil {
			*t = 0
		}
		return nil
	}

	switch l := len(b); {

	case !(l >= 6 && l <= 8) && l != 10:
		// The length must be:
		// - 6: "hhmm",
		// - 7: "hh:mm",
		// - 8: "hhmmss"
		// - 10: "hh:mm:ss"
		return _ERR_NOT_ISO8601_TIME

	case b[0] != '"' || b[l-1] != '"':
		// Forgotten quotes? Incorrect JSON?
		return _ERR_BAD_JSON_TIME_QUO

	default:
		return t.ParseFrom(b[1 : l-1])
	}
}
