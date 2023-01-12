// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatyp

import (
	"bytes"
	"database/sql/driver"

	"github.com/oklog/ulid/v2"

	"github.com/qioalice/ekago/v4/ekaenc"
	"github.com/qioalice/ekago/v4/ekarand"
)

// Ulid is a Universally Unique Lexicographically Sortable Identifier.
// It's a drop-in replacement of UUID.
// Read more: https://github.com/ulid/spec .
type Ulid ulid.ULID

var ulidNil Ulid

////////////////////////////////////////////////////////////////////////////////

// NewUlidTo creates and saves a new Ulid based on the current time
// and math/rand entropy to the destination.
//
// Thread-safety. Returns ErrIDNilDestination if 'to' is nil.
func NewUlidTo(to *Ulid) error {

	if to == nil {
		return ErrIDNilDestination
	}

	var u, err = ulid.New(ulid.Now(), (*ekarand.MathRandReader)(nil))

	if err == nil {
		*to = Ulid(u)
	}
	return err
}

// NewUlidFromStringTo parses string and saves Ulid to the destination.
// Input is expected in a form accepted by UnmarshalText.
//
// Returns ErrIDNilDestination if 'to' is nil.
func NewUlidFromStringTo(to *Ulid, input string) error {

	if to == nil {
		return ErrIDNilDestination
	}

	var u, err = ulid.ParseStrict(input)

	if err == nil {
		*to = Ulid(u)
	}
	return err
}

////////////////////////////////////////////////////////////////////////////////

// Equal returns true if both of Ulid s are equal, otherwise returns false.
// It returns true if both of them are nil but false if only one.
func (u *Ulid) Equal(other ID) bool {
	if b1, b2 := u == nil, other == nil; b1 || b2 {
		return b1 && b2
	}
	var otherTyped, _ = other.(*Ulid)
	return bytes.Equal(u[:], otherTyped[:])
}

// IsNil reports whether current Ulid is empty (nil).
func (u *Ulid) IsNil() bool {
	return u == nil || bytes.Equal(u[:], ulidNil[:])
}

// SetNil sets the current Ulid to zero Ulid. Does nothing if Ulid is nil.
func (u *Ulid) SetNil() {
	if u != nil {
		*u = ulidNil
	}
}

// Bytes returns bytes slice representation of Ulid. Returns nil if Ulid is nil.
func (u *Ulid) Bytes() []byte {
	if u == nil {
		return nil
	}
	return u[:]
}

// String returns a lexicographically sortable string encoded Ulid
// (26 characters, non-standard base 32) e.g. 01AN4Z07BY79KA1307SR9X4MV3
// Format: tttttttttteeeeeeeeeeeeeeee where t is time and e is entropy.
// Returns "<nil>" if Ulid is nil.
func (u *Ulid) String() string {
	if u == nil {
		return "<nil>"
	}
	return ulid.ULID(*u).String()
}

////////////////////////////////////////////////////////////////////////////////

// MarshalText implements the encoding.TextMarshaler interface by
// returning the string encoded Ulid. Returns "null" if Ulid is nil.
func (u *Ulid) MarshalText() ([]byte, error) {
	if u == nil {
		return ekaenc.NullAsBytesLowerCase(), nil
	}
	return ulid.ULID(*u).MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface by parsing
// the data as string encoded Ulid. Returns ErrIDNilDestination if Ulid is nil.
func (u *Ulid) UnmarshalText(data []byte) error {
	if u == nil {
		return ErrIDNilDestination
	}
	return (*ulid.ULID)(u).UnmarshalText(data)
}

////////////////////////////////////////////////////////////////////////////////

// MarshalJSON implements the json.Marshaler interface.
// Returns a JSON string with encoded Ulid with the followed format:
// "01F58MK33HZR24YGEFXWV619XB". Returns JSON null if Ulid is nil.
func (u *Ulid) MarshalJSON() ([]byte, error) {

	if u == nil {
		return ekaenc.NullAsBytesLowerCase(), nil
	}

	var buf = make([]byte, ulid.EncodedSize+2)
	buf[0], buf[len(buf)-1] = '"', '"'

	if err := ulid.ULID(*u).MarshalTextTo(buf[1 : len(buf)-1]); err != nil {
		return nil, err
	}

	return buf, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// Decodes data as encoded JSON Ulid string and saves the result to u.
// Supports JSON null values.
func (u *Ulid) UnmarshalJSON(data []byte) error {

	switch {
	case u == nil:
		return ErrIDNilDestination

	case len(data) == 0 || ekaenc.IsNullAsBytes(data):
		u.SetNil()
		return nil

	case len(data) != ulid.EncodedSize+2:
		return ulid.ErrDataSize
	}

	return (*ulid.ULID)(u).UnmarshalText(data[1 : len(data)-1])
}

////////////////////////////////////////////////////////////////////////////////

// MarshalBinary implements the encoding.BinaryMarshaler interface by
// returning the Ulid as a byte slice. Returns nil if Ulid is nil.
func (u *Ulid) MarshalBinary() ([]byte, error) {
	if u == nil {
		return nil, nil
	}
	return ulid.ULID(*u).MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface by
// copying the passed data and converting it to an Ulid.
// Returns ErrIDNilDestination if Ulid is nil.
func (u *Ulid) UnmarshalBinary(data []byte) error {
	if u == nil {
		return ErrIDNilDestination
	}
	return (*ulid.ULID)(u).UnmarshalBinary(data)
}

////////////////////////////////////////////////////////////////////////////////

// Value implements the driver.Valuer interface. This returns the value
// represented as a string.
func (u *Ulid) Value() (driver.Value, error) {
	return u.MarshalText()
}

// Scan implements the sql.Scanner interface. It supports scanning
// a string or byte slice. Returns ErrIDNilDestination if Ulid is nil.
func (u *Ulid) Scan(src any) error {
	if u == nil {
		return ErrIDNilDestination
	}
	return (*ulid.ULID)(u).Scan(src)
}

////////////////////////////////////////////////////////////////////////////////

var _ ID = (*Ulid)(nil)
