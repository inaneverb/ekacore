// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatyp

import (
	"bytes"
	"database/sql/driver"

	"github.com/qioalice/ekago/v3/ekarand"

	"github.com/oklog/ulid/v2"
)

type (
	// ULID is a Universally Unique Lexicographically Sortable Identifier.
	// It's a drop-in replacement of UUID.
	// Read more: https://github.com/ulid/spec .
	ULID ulid.ULID
)

// ---------------------------- UUID COMMON METHODS --------------------------- //
// ---------------------------------------------------------------------------- //

// Equal returns true if both of ULID s are equal, otherwise returns false.
func (u ULID) Equal(anotherUlid ULID) bool {
	return bytes.Equal(u[:], anotherUlid[:])
}

// IsNil reports whether current ULID is empty (nil).
func (u ULID) IsNil() bool {
	return u == ULID(_UUID_NULL)
}

// SetNil sets the current ULID to zero ULID. Returns modified ULID.
func (u *ULID) SetNil() *ULID {
	*u = ULID(_UUID_NULL)
	return u
}

// Bytes returns bytes slice representation of ULID.
func (u ULID) Bytes() []byte {
	return u[:]
}

// String returns a lexicographically sortable string encoded ULID
// (26 characters, non-standard base 32) e.g. 01AN4Z07BY79KA1307SR9X4MV3
// Format: tttttttttteeeeeeeeeeeeeeee where t is time and e is entropy.
func (u ULID) String() string {
	return ulid.ULID(u).String()
}

// --------------------------- UUID CREATION HELPERS -------------------------- //
// ---------------------------------------------------------------------------- //

// ULID_OrPanic is a helper that wraps a call to a function returning (ULID, error)
// and panics if the error is non-nil.
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func ULID_OrPanic(u ULID, err error) ULID {
	if err != nil {
		panic(err)
	}
	return u
}

// ULID_OrNil is a helper that wraps a call to a function returning (ULID, error)
// and returns zero ULID if the error is non-nil.
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func ULID_OrNil(u ULID, err error) ULID {
	if err != nil {
		return ULID(_UUID_NULL)
	}
	return u
}

// ------------------------------ ULID GENERATORS ----------------------------- //
// ---------------------------------------------------------------------------- //

// ULID_New() returns an new ULID based on the current time and math/rand entropy.
// Thread-safety.
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func ULID_New() (ULID, error) {
	u, err := ulid.New(ulid.Now(), (*ekarand.MathRandReader)(nil))
	return ULID(u), err
}

// ------------------- ULID GENERATOR'S WRAPPERS OF HELPERS ------------------- //
// ---------------------------------------------------------------------------- //

// Next methods are the same as just generators but it panics
// if any error occurred while ULID been generated.

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func ULID_New_OrPanic() ULID {
	return ULID_OrPanic(ULID_New())
}

// Next methods are the same as just generators but it returns
// a zero ULID if any error is occurred while UUID been generated.

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func ULID_New_OrNil() ULID {
	return ULID_OrNil(ULID_New())
}

// Next methods are the same as just generators but it returns
// only one argument - an error and saves generated ULID as output argument
// by the address provided by 'dest' arg.
//
// It's useful when you awaits only one argument for being returned.
// For example in if statement to omit else branch.
//
// WARNING!
// No nil check! 'dest' must be not nil, panic otherwise.

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func ULID_New_To(dest *ULID) (err error) {
	*dest, err = ULID_New()
	return
}

// ------------------------------- ULID PARSERS ------------------------------- //
// ---------------------------------------------------------------------------- //

// ULID_FromString returns ULID parsed from string input.
// Input is expected in a form accepted by UnmarshalText.
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func ULID_FromString(input string) (ULID, error) {
	u, err := ulid.ParseStrict(input)
	return ULID(u), err
}

// -------------------- ULID PARSER'S WRAPPERS OF HELPERS --------------------- //
// ---------------------------------------------------------------------------- //

// ULID_FromString_OrPanic is the same as ULID_OrPanic(ULID_FromString(input)).
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func ULID_FromString_OrPanic(input string) ULID {
	return ULID_OrPanic(ULID_FromString(input))
}

// ULID_FromString_OrNil is the same as ULID_OrNil(UUID_FromString(input)).
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func ULID_FromString_OrNil(input string) ULID {
	return ULID_OrNil(ULID_FromString(input))
}

// ------------------------ UUID TEXT ENCODER/DECODER ------------------------- //
// ---------------------------------------------------------------------------- //

// MarshalText implements the encoding.TextMarshaler interface by
// returning the string encoded ULID.
func (u ULID) MarshalText() ([]byte, error) {
	return ulid.ULID(u).MarshalText()
}

// MarshalTextTo writes the ULID as a string to the given buffer.
// ErrBufferSize is returned when the len(dst) != 26.
func (u ULID) MarshalTextTo(dest []byte) error {
	return ulid.ULID(u).MarshalTextTo(dest)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface by
// parsing the data as string encoded ULID.
//
// ErrDataSize is returned if the len(v) is different from an encoded
// ULID's length. Invalid encodings produce undefined ULIDs.
func (u *ULID) UnmarshalText(data []byte) error {
	return (*ulid.ULID)(u).UnmarshalText(data)
}

// ------------------------ UUID JSON ENCODER/DECODER ------------------------- //
// ---------------------------------------------------------------------------- //

// MarshalJSON implements the encoding/json.Marshaler interface.
// Returns a JSON string with encoded ULID with the followed format:
// "01F58MK33HZR24YGEFXWV619XB". Returns JSON null if ULID is nil.
func (u ULID) MarshalJSON() ([]byte, error) {

	if u.IsNil() {
		return _UUID_JSON_NULL, nil
	}

	buf := make([]byte, ulid.EncodedSize+2)
	buf[0], buf[len(buf)-1] = '"', '"'

	if err := ulid.ULID(u).MarshalTextTo(buf[1 : len(buf)-1]); err != nil {
		return nil, err
	}

	return buf, nil
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
// Decodes data as encoded JSON ULID string and saves the result to u.
// Supports JSON null values.
func (u *ULID) UnmarshalJSON(data []byte) error {

	if len(data) == 0 || bytes.Compare(data, _UUID_JSON_NULL) == 0 {
		return nil
	}

	if len(data) != ulid.EncodedSize+2 {
		return ulid.ErrDataSize
	}

	return (*ulid.ULID)(u).UnmarshalText(data[1 : len(data)-1])
}

// ----------------------- UUID BINARY ENCODER/DECODER ------------------------ //
// ---------------------------------------------------------------------------- //

// MarshalBinary implements the encoding.BinaryMarshaler interface by
// returning the ULID as a byte slice.
func (u ULID) MarshalBinary() ([]byte, error) {
	return ulid.ULID(u).MarshalBinary()
}

// MarshalBinaryTo writes the binary encoding of the ULID to the given buffer.
// ErrBufferSize is returned when the len(dst) != 16.
func (u ULID) MarshalBinaryTo(dst []byte) error {
	return ulid.ULID(u).MarshalBinaryTo(dst)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface by
// copying the passed data and converting it to an ULID. ErrDataSize is
// returned if the data length is different from ULID length.
func (u *ULID) UnmarshalBinary(data []byte) error {
	return (*ulid.ULID)(u).UnmarshalBinary(data)
}

// ------------------------- UUID SQL ENCODER/DECODER ------------------------- //
// ---------------------------------------------------------------------------- //

// Value implements the sql/driver.Valuer interface. This returns the value
// represented as a string.
func (u ULID) Value() (driver.Value, error) {
	return ulid.ULID(u).MarshalText()
}

// Scan implements the sql.Scanner interface. It supports scanning
// a string or byte slice.
func (u *ULID) Scan(src interface{}) error {
	return (*ulid.ULID)(u).Scan(src)
}
