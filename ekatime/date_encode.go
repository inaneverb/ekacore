// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatime

import (
	"errors"
	"strconv"

	"github.com/qioalice/ekago/v2/internal/ekaenc"
)

//goland:noinspection GoSnakeCaseUsage
var (
	_ERR_NIL_DATE_RECEIVER = errors.New("nil ekatime.Date receiver")
	_ERR_NOT_ISO8601_DATE  = errors.New("incorrect ISO8601 date format (must be YYYYMMDD or YYYY-MM-DD)")
	_ERR_BAD_YEAR          = errors.New("year must be in the range [0..4095]")
	_ERR_BAD_MONTH         = errors.New("month must be in the range [1..12]")
	_ERR_BAD_DAY           = errors.New("day must be in the range [1..31]")
	_ERR_BAD_CORRESP_DATE  = errors.New("date must represent valid date (check year, month, day corresponding)")
	_ERR_BAD_JSON_DATE_QUO = errors.New("bad JSON ISO8601 date representation (forgotten quotes?)")
)

// AppendTo generates a string representation of Date and adds it to the b,
// returning a new slice (if it has grown) or the same if there was enough
// space to store 10 bytes (of string representation).
//
// Uses separator as Date's parts separator:
// "YYYY<separator>MM<separator>DD".
//
// Up to 6x faster than fmt.Sprintf().
func (dd Date) AppendTo(b []byte, separator byte) []byte {

	f := func(b []byte) {
		y, m, d := normalizeDate(dd.Split())

		if y >= _YEAR_AS_NUM_STR_MIN && y <= _YEAR_AS_NUM_STR_MAX {
			copy(b[:4], _YEAR_AS_NUM_STR[y-_YEAR_AS_NUM_STR_MIN])
		} else {
			// normalizeDate never return y > 4095, so it's safe
			_ = append(b[:0], strconv.Itoa(int(y))...)
		}

		copy(b[5:7], _DAY_AS_NUM_STR[m-MONTH_JANUARY]) // use day's array, it's ok
		copy(b[8:10], _DAY_AS_NUM_STR[d-1])

		b[4] = separator
		b[7] = separator
	}

	if c, l := cap(b), len(b); c - l >= 10 {
		b = b[:l+10]
		f(b[l:])
		return b
	} else {
		// One more allocation
		b2 := make([]byte, 10)
		f(b2)
		b = append(b, b2...)
		return b
	}
}

// ParseFrom tries to parse b considering with the following format:
// "YYYY<separator>MM<separator>DD", <separator> may be any 1 byte,
// or not presented at all. But if separator was presented between in "YYYY<sep>MM"
// it must be also between other parts and vice-versa.
//
// Read more:
// https://en.wikipedia.org/wiki/ISO_8601
//
// If success, returns nil and saves date into the current Date object.
// Year must be <= 4095, instead it will be overwritten by 4095.
//
// Skips leading spaces, ignores all next data after the day has been scanned.
func (dd *Date) ParseFrom(b []byte) error {

	if dd == nil {
		return _ERR_NIL_DATE_RECEIVER
	}

	var i = 0
	for n := len(b); i < n && b[i] <= ' '; i++ { }

	// Minimum required len: 8 (YYYYMMDD - w/o separators).
	if len(b[i:]) < 8 {
		return _ERR_NOT_ISO8601_DATE
	}

	x1, valid1 := batoi(b[i], b[i+1])
	x2, valid2 := batoi(b[i+2], b[i+3])

	x1 *= 100
	x1 += x2

	if !(valid1 && valid2) || x1 < 0 || x2 < 0 || x1 > 4095 {
		return _ERR_BAD_YEAR
	}

	i += 4
	y := Year(x1)
	wasSeparator := false

	// Skip separator
	if !(b[i] >= '0' && b[i] <= '9') {
		i++
		wasSeparator = true
	}

	x1, valid1 = batoi(b[i], b[i+1])
	m := Month(x1)
	if !valid1 || m < MONTH_JANUARY || m > MONTH_DECEMBER {
		return _ERR_BAD_MONTH
	}

	i += 2
	if wasSeparator {
		i++
	}

	if len(b[i:]) < 2 {
		// Looks like user has incorrect data like:
		// "YYYY-MM-D", "YYYY-MM-"
		return _ERR_NOT_ISO8601_DATE
	}

	x1, valid1 = batoi(b[i], b[i+1])
	if x1 < 1 || x1 > 31 {
		return _ERR_BAD_DAY
	}

	if !IsValidDate(y, m, Day(x1)) {
		return _ERR_BAD_CORRESP_DATE
	}

	*dd = NewDate(y, m, Day(x1))
	return nil
}

// String returns the current Date's string representation in the following format:
// "YYYY/MM/DD".
func (dd Date) String() string {
	return string(dd.AppendTo(make([]byte, 0, 10), '/'))
}

// MarshalJSON encodes the current Date in the following format (quoted)
// "YYYY-MM-DD", and returns it. Always returns nil as error.
//
// JSON null supporting:
// - Writes JSON null if current Date receiver is nil
// - Writes JSON null if current Date == 0.
func (dd *Date) MarshalJSON() ([]byte, error) {

	if dd == nil || dd.ToCmp() == 0 {
		return ekaenc.NULL_JSON_BYTES_SLICE, nil
	}

	b := make([]byte, 12)
	_ = dd.AppendTo(b[1:1:12], '-')

	b[0] = '"'
	b[11] = '"'

	return b, nil
}

// UnmarshalJSON decodes b into the current Date object expecting b contains
// ISO8601 quoted date (only date, not date w/ time) in the one of the following
// formats: "YYYYMMDD", "YYYY-MM-DD".
//
// JSON null supporting:
// - It's ok if there is JSON null and receiver == nil (nothing changes)
// - Zeroes Date if there is JSON null and receiver != nil.
//
// In other cases JSON parsing error or Date.ParseFrom() error is returned.
func (dd *Date) UnmarshalJSON(b []byte) error {

	if ekaenc.IsNullJSON(b) {
		if dd != nil {
			*dd = 0
		}
		return nil
	}

	switch l := len(b); {

	case l != 10 && l != 12:
		// The length must be 10 (8+quotes) or 12 (8+separators+quotes).
		return _ERR_NOT_ISO8601_DATE

	case b[0] != '"' || b[l-1] != '"':
		// Forgotten quotes? Incorrect JSON?
		return _ERR_BAD_JSON_DATE_QUO

	default:
		return dd.ParseFrom(b[1:l-1])
	}
}
