// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"errors"

	"github.com/qioalice/ekago/v3/internal/ekaenc"
)

//goland:noinspection GoSnakeCaseUsage
var (
	_ERR_NIL_TIMESTAMP_RECEIVER  = errors.New("nil ekatime.Timestamp receiver")
	_ERR_NOT_ISO8601_TIMESTAMP   = errors.New("incorrect ISO8601 timestamp format (must be YYYY-MM-DDThh:mm:ss)")
	_ERR_BAD_TIMESTAMP_SEPARATOR = errors.New("ISO8601 require using 'T' as date time separator")
	_ERR_BAD_JSON_TIMESTAMP_QUO  = errors.New("bad JSON ISO8601 timestamp representation (forgotten quotes?)")
)

// AppendTo generates a string representation of Timestamp and adds it to the b,
// returning a new slice (if it has grown) or the same if there was enough
// space to store 19 bytes (of string representation).
//
// Uses separator as Date and Time's parts separator:
// "YYYY<separatorDate>MM<separatorDate>DD hh<separatorTime>mm<separatorTime>ss".
//
// Up to 10x faster than fmt.Sprintf().
func (ts Timestamp) AppendTo(b []byte, separatorDate, separatorTime byte) []byte {

	b = NewDate(dateFromUnix(ts)).AppendTo(b, separatorDate)
	b = append(b, ' ')
	b = NewTime(timeFromUnix(ts)).AppendTo(b, separatorTime)

	return b
}

// ParseFrom tries to parse b considering with the following format:
// "YYYY<sep>MM<sep>DD<reqSep>hh<sep>mm".
// - <sep> may be any 1 byte, or not presented at all.
//   But if separator was presented between in "YYYY<sep>MM" it must be also
//   between other parts and vice-versa.
// - <reqSep> required separator: space or 'T' char.
//
// Read more:
// https://en.wikipedia.org/wiki/ISO_8601
//
// If success, returns nil and saves date into the current Timestamp object.
// Year must be <= 4095, instead it will be overwritten by 4095.
//
// Skips leading spaces, ignores all next data after the day has been scanned.
func (ts *Timestamp) ParseFrom(b []byte) error {


	if ts == nil {
		return _ERR_NIL_TIMESTAMP_RECEIVER
	}

	var (
		d Date
		t Time
	)

	// Regardless Date.ParseFrom() has the loop that skips spaces,
	// we need to know we must start parsing Time from.
	// So, skip spaces manually.

	var i = 0
	for n := len(b); i < n && b[i] <= ' '; i++ { }

	if err := d.ParseFrom(b[i:]); err != nil {
		return err
	}

	// Date.ParseFrom() finished w/ no errors.
	// So, it was be either 8 chars (YYYYMMDD) or 10 (YYYY-MM-DD).
	//
	// YYYYMMDD
	// YYYY-MM-DD
	//        ^    we need to check this char.
	// If it's a digit, it was 6 char format. Otherwise - 8.

	i += 7
	if b[i] >= '0' && b[i] <= '9' {
		// Was 8 char format.
		i += 1
	} else {
		// Was 10 char format.
		i += 3
	}

	// It's OK if time is not presented.
	b = b[i:]
	i = 0
	if len(b) > 0 {
		// ISO8601 requires a separator between date and time. And it must be 'T' char.
		// We allows also a space. Why not?
		if b[i] != ' ' && b[i] != 't' && b[i] != 'T' {
			return _ERR_BAD_TIMESTAMP_SEPARATOR
		}
		i++
		if err := t.ParseFrom(b[i:]); err != nil {
			return err
		}
	}

	*ts = d.WithTime(t.Split())
	return nil
}

// String returns the current Timestamp's human-readable string representation
// in the following format: "YYYY/MM/DD hh:mm:ss".
func (ts Timestamp) String() string {
	return string(ts.AppendTo(make([]byte, 0, 19), '/', ':'))
}

// MarshalJSON encodes the current Time in the following format (quoted)
// "YYYY-MM-DDThh:mm:ss", and returns it. Always returns nil as error.
//
// JSON null supporting:
// - Writes JSON null if current Timestamp receiver == nil.
// - Writes JSON null if current Timestamp == 0 (*).
//
// -----
//
// (*): Yes, it's kinda weird that 01 Jan 1970 00:00:00 is JSON null,
// but let's be honest. When the last time you'll need to marshal EXACTLY
// this date and it's not a null-like value?
// And, of course, you always may use 01 Jan 1970 00:00:01 and it will be marshalled
// correctly.
func (ts *Timestamp) MarshalJSON() ([]byte, error) {

	if ts == nil || *ts == 0 {
		return ekaenc.NULL_JSON_BYTES_SLICE, nil
	}

	// Date: 10 chars (YYYY-MM-DD)
	// Clock: 8 chars (hh:mm:ss)
	// Quotes: 2 chars ("")
	// Date clock separator: 1 char (T)
	// Summary: 21 char.
	b := make([]byte, 21)

	_ = ts.Date().AppendTo(b[1:1:20], '-')
	_ = ts.Time().AppendTo(b[12:12:20], ':')

	b[0] = '"'
	b[11] = 'T'
	b[20] = '"'

	return b, nil
}

// UnmarshalJSON decodes b into the current Timestamp object expecting b contains
// ISO8601 quoted date with time in the one of the following formats:
//   "YYYY-MM-DDThh:mm:ss" (recommended),
//   "YYYYMMDDThh:mm:ss", "YYYY-MM-DDThhmmss", "YYYYMMDDThhmmss"
//
// JSON null supporting:
// - It's ok if there is JSON null and receiver == nil (nothing changes)
// - Zeroes Timestamp if there is JSON null and receiver != nil
//   (yes, sets to 01 Jan 1970 00:00:00).
//
// In other cases JSON parsing error or Date.ParseFrom() error is returned.
func (ts *Timestamp) UnmarshalJSON(b []byte) error {

	if ekaenc.IsNullJSON(b) {
		if ts != nil {
			*ts = 0
		}
		return nil
	}

	// Length variants:
	// Date length variants: 10, 8
	// Clock length variants: 4, 5, 6, 8,
	// Quotes: 2,
	// Date time separator: 1.
	//
	// So, summary variants:
	// - 15 (8+4+2+1),
	// - 16 (8+5+2+1),
	// - 17 (8+6+2+1, 10+4+2+1),
	// - 18 (10+5+2+1),
	// - 19 (8+8+2+1, 10+6+2+1),
	// - 21 (10+8+2+1).

	switch l := len(b); {

	case !(l >= 15 && l <= 19) && l != 21:
		return _ERR_NOT_ISO8601_TIMESTAMP

	case b[0] != '"' && b[l-1] != '"':
		return _ERR_BAD_JSON_TIMESTAMP_QUO

	default:
		return ts.ParseFrom(b[1:l-1])
	}
}
