// Copyright © 2020. All rights reserved.
// Refactorer, modifier: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

// Copyright (C) 2013-2018 by Maxim Bublis <b@codemonkey.ru>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package ekatyp

import (
	"bytes"
	"database/sql/driver"
	"fmt"
)

type (
	// UUID representation compliant with specification described in RFC 4122.
	UUID [_UUID_SIZE]byte
)

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
const (
	// UUID versions

	UUID_V1 byte = 1
	UUID_V2 byte = 2
	UUID_V3 byte = 3
	UUID_V4 byte = 4
	UUID_V5 byte = 5

	// UUID layout variants.

	UUID_VARIANT_NCS       byte = 0
	UUID_VARIANT_RFC4122   byte = 1
	UUID_VARIANT_MICROSOFT byte = 2
	UUID_VARIANT_FUTURE    byte = 3

	// UUID DCE domains.

	UUID_DOMAIN_PERSON = 0
	UUID_DOMAIN_GROUP  = 1
	UUID_DOMAIN_ORG    = 2
)

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
var (
	// _UUID_NULL is special form of UUID that is specified to have all
	// 128 bits set to zero.
	_UUID_NULL = UUID{}

	// Predefined namespace UUIDs.
	UUID_NAMESPACE_DNS  = UUID_OrPanic(UUID_FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
	UUID_NAMESPACE_URL  = UUID_OrPanic(UUID_FromString("6ba7b811-9dad-11d1-80b4-00c04fd430c8"))
	UUID_NAMESPACE_OID  = UUID_OrPanic(UUID_FromString("6ba7b812-9dad-11d1-80b4-00c04fd430c8"))
	UUID_NAMESPACE_X500 = UUID_OrPanic(UUID_FromString("6ba7b814-9dad-11d1-80b4-00c04fd430c8"))
)

// ---------------------------- UUID COMMON METHODS --------------------------- //
// ---------------------------------------------------------------------------- //

// Equal returns true if u and anotherUuid equals, otherwise returns false.
func (u UUID) Equal(anotherUuid UUID) bool {
	return bytes.Equal(u[:], anotherUuid[:])
}

// IsNil reports whether u is nil or not. Is the same as u.Equal(_UUID_NULL).
func (u UUID) IsNil() bool {
	return u.Equal(_UUID_NULL)
}

// SetNil sets the current u to _UUID_NIL. Returns modified u.
func (u *UUID) SetNil() *UUID {
	*u = _UUID_NULL
	return u
}

// Version returns algorithm version used to generate UUID.
func (u UUID) Version() byte {
	return u[6] >> 4
}

// Variant returns UUID layout variant.
func (u UUID) Variant() byte {
	switch {
	case (u[8] >> 7) == 0x00:
		return UUID_VARIANT_NCS
	case (u[8] >> 6) == 0x02:
		return UUID_VARIANT_RFC4122
	case (u[8] >> 5) == 0x06:
		return UUID_VARIANT_MICROSOFT
	case (u[8] >> 5) == 0x07:
		fallthrough
	default:
		return UUID_VARIANT_FUTURE
	}
}

// Bytes returns bytes slice representation of UUID.
func (u UUID) Bytes() []byte {
	return u[:]
}

// Returns canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u UUID) String() string {
	return string(u.hexEncodeTo(make([]byte, 36)))
}

// SetVersion sets version bits.
func (u *UUID) SetVersion(v byte) {
	u[6] = (u[6] & 0x0f) | (v << 4)
}

// SetVariant sets variant bits.
func (u *UUID) SetVariant(v byte) {
	switch v {
	case UUID_VARIANT_NCS:
		u[8] = u[8]&(0xff>>1) | (0x00 << 7)
	case UUID_VARIANT_RFC4122:
		u[8] = u[8]&(0xff>>2) | (0x02 << 6)
	case UUID_VARIANT_MICROSOFT:
		u[8] = u[8]&(0xff>>3) | (0x06 << 5)
	case UUID_VARIANT_FUTURE:
		fallthrough
	default:
		u[8] = u[8]&(0xff>>3) | (0x07 << 5)
	}
}

// --------------------------- UUID CREATION HELPERS -------------------------- //
// ---------------------------------------------------------------------------- //

// UUID_OrPanic is a helper that wraps a call to a function returning (UUID, error)
// and panics if the error is non-nil.
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_OrPanic(u UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return u
}

// UUID_OrNil is a helper that wraps a call to a function returning (UUID, error)
// and returns _UUID_NULL if the error is non-nil.
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_OrNil(u UUID, err error) UUID {
	if err != nil {
		return _UUID_NULL
	}
	return u
}

// -------------------------- UUID RFC4122 GENERATORS ------------------------- //
// ---------------------------------------------------------------------------- //

// UUID_NewV1 returns UUID based on current timestamp and MAC address.
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_NewV1() (UUID, error) {
	return _UUID_RFC4122_Generator.NewV1()
}

// UUID_NewV2 returns DCE Security UUID based on POSIX UID/GID.
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_NewV2(domain byte) (UUID, error) {
	return _UUID_RFC4122_Generator.NewV2(domain)
}

// UUID_NewV3 returns UUID based on MD5 hash of namespace UUID and name.
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_NewV3(ns UUID, name string) UUID {
	return _UUID_RFC4122_Generator.NewV3(ns, name)
}

// UUID_NewV4 returns random generated UUID.
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_NewV4() (UUID, error) {
	return _UUID_RFC4122_Generator.NewV4()
}

// UUID_NewV5 returns UUID based on SHA-1 hash of namespace UUID and name.
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_NewV5(ns UUID, name string) UUID {
	return _UUID_RFC4122_Generator.NewV5(ns, name)
}

// --------------- UUID RFC4122 GENERATOR'S WRAPPERS OF HELPERS --------------- //
// ---------------------------------------------------------------------------- //

// Next methods are the same as just generators v1/v2/v4 but it panics
// if any error occurred while UUID been generated.

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_NewV1_OrPanic() UUID { return UUID_OrPanic(UUID_NewV1()) }

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_NewV2_OrPanic(domain byte) UUID { return UUID_OrPanic(UUID_NewV2(domain)) }

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_NewV4_OrPanic() UUID { return UUID_OrPanic(UUID_NewV4()) }

// Next methods are the same as just generators v1/v2/v4 but it returns
// a NULL UUID "6ba7b810-9dad-11d1-80b4-00c04fd430c8" if any error is occurred
// while UUID been generated.

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_NewV1_OrNil() UUID { return UUID_OrNil(UUID_NewV1()) }

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_NewV_2OrNil(domain byte) UUID { return UUID_OrNil(UUID_NewV2(domain)) }

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_NewV4_OrNil() UUID { return UUID_OrNil(UUID_NewV4()) }

// Next methods are the same as just generators v1/v2/v4, but it returns
// only one argument - an error and saves generated UUID as output argument
// by the address provided by 'dest' arg.
//
// It's useful when you awaits only one argument for being returned.
// For example in if statement to omit else branch.
//
// WARNING!
// No nil check! 'dest' must be not nil, panic otherwise.

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_NewV1_To(dest *UUID) (err error) {
	*dest, err = UUID_NewV1()
	return
}

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_NewV2_To(dest *UUID, domain byte) (err error) {
	*dest, err = UUID_NewV2(domain)
	return
}

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_NewV4_To(dest *UUID) (err error) {
	*dest, err = UUID_NewV4()
	return
}

// ------------------------------- UUID PARSERS ------------------------------- //
// ---------------------------------------------------------------------------- //

// UUID_FromBytes returns UUID converted from raw byte slice input.
// It will return error if the slice isn't 16 bytes long.
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_FromBytes(input []byte) (u UUID, err error) {
	err = u.UnmarshalBinary(input)
	return
}

// UUID_FromString returns UUID parsed from string input.
// Input is expected in a form accepted by UnmarshalText.
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_FromString(input string) (u UUID, err error) {
	err = u.UnmarshalText([]byte(input))
	return
}

// -------------------- UUID PARSER'S WRAPPERS OF HELPERS --------------------- //
// ---------------------------------------------------------------------------- //

// UUID_FromBytes_OrPanic is the same as UUID_OrPanic(UUID_FromBytes(input)).
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_FromBytes_OrPanic(input []byte) UUID {
	return UUID_OrPanic(UUID_FromBytes(input))
}

// UUID_FromString_OrPanic is the same as UUID_OrPanic(UUID_FromString(input)).
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_FromString_OrPanic(input string) UUID {
	return UUID_OrPanic(UUID_FromString(input))
}

// UUID_FromBytes_OrNil is the same as UUID_OrNil(UUID_FromBytes(input)).
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_FromBytes_OrNil(input []byte) UUID {
	return UUID_OrNil(UUID_FromBytes(input))
}

// UUID_FromString_OrNil is the same as UUID_OrNil(UUID_FromString(input)).
//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
func UUID_FromString_OrNil(input string) UUID {
	return UUID_OrNil(UUID_FromString(input))
}

// ------------------------ UUID TEXT ENCODER/DECODER ------------------------- //
// ---------------------------------------------------------------------------- //

// MarshalText implements the encoding.TextMarshaler interface.
// The encoding is the same as returned by String.
func (u UUID) MarshalText() (text []byte, err error) {
	text = []byte(u.String())
	return
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Following formats are supported:
//   - "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
//   - "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
//   - "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
//   - "6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) UnmarshalText(text []byte) (err error) {
	switch len(text) {
	case 32:
		return u.decodeHashLike(text)
	case 36:
		return u.decodeCanonical(text)
	case 38:
		return u.decodeBraced(text)
	case 41:
		fallthrough
	case 45:
		return u.decodeURN(text)
	default:
		return fmt.Errorf("uuid: incorrect UUID length: %s", text)
	}
}

// ------------------------ UUID JSON ENCODER/DECODER ------------------------- //
// ---------------------------------------------------------------------------- //

// MarshalJSON implements the encoding/json.Marshaler interface.
// Returns a JSON string with encoded UUID with the followed format:
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8". Returns JSON null if u is _UUID_NULL.
func (u UUID) MarshalJSON() ([]byte, error) {
	if u == _UUID_NULL {
		return _UUID_JSON_NULL, nil
	}
	return u.jsonMarshal(), nil
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
// Decodes b as encoded JSON UUID string and saves the result to u.
// Supports all UUID variants that u.UnmarshalText() does support but also
// supports JSON null values.
func (u *UUID) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || bytes.Compare(b, _UUID_JSON_NULL) == 0 {
		return nil
	}
	// JSON contains quotes (") because it's raw JSON data and JSON strings
	// has quotes
	if len(b) < 2 {
		return fmt.Errorf("uuid: incorrect UUID length: %s", string(b))
	}
	return u.UnmarshalText(b[1 : len(b)-1])
}

// ----------------------- UUID BINARY ENCODER/DECODER ------------------------ //
// ---------------------------------------------------------------------------- //

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (u UUID) MarshalBinary() (data []byte, err error) {
	data = u.Bytes()
	return
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// It will return error if the slice isn't 16 bytes long.
func (u *UUID) UnmarshalBinary(data []byte) (err error) {
	if len(data) != _UUID_SIZE {
		err = fmt.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(data))
		return
	}
	copy(u[:], data)
	return
}

// ------------------------- UUID SQL ENCODER/DECODER ------------------------- //
// ---------------------------------------------------------------------------- //

// Value implements the driver.Valuer interface. Supports SQL NULL.
func (u UUID) Value() (driver.Value, error) {
	if u == _UUID_NULL {
		return nil, nil
	}
	return u.String(), nil
}

// Scan implements the sql.Scanner interface.
// A 16-byte slice is handled by UnmarshalBinary, while
// a longer byte slice or a string is handled by UnmarshalText. Supports SQL NULL.
func (u *UUID) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		return nil

	case []byte:
		if len(src) == _UUID_SIZE {
			return u.UnmarshalBinary(src)
		}
		return u.UnmarshalText(src)

	case string:
		return u.UnmarshalText([]byte(src))
	}

	return fmt.Errorf("uuid: cannot convert %T to UUID", src)
}
