// Copyright © 2020-2023. All rights reserved.
// Refactorer, modifier: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatyp

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/gofrs/uuid"

	"github.com/qioalice/ekago/v4/ekaenc"
	"github.com/qioalice/ekago/v4/ekastr"
)

// Uuid is a Universally Unique IDentifier.
// Read more: https://en.wikipedia.org/wiki/Universally_unique_identifier
type Uuid uuid.UUID

var uuidNil Uuid

var (
	ErrUuidUnsupportedVersion = errors.New("uuid: unsupported version")
)

const (
	// UUID versions

	UuidV1 byte = 1
	UuidV2 byte = 2
	UuidV3 byte = 3
	UuidV4 byte = 4
	UuidV5 byte = 5
	UuidV6 byte = 6
	UuidV7 byte = 7

	// UUID layout variants.

	UuidVariantNcs       byte = 0
	UuidVariantRfc4122   byte = 1
	UuidVariantMicrosoft byte = 2
	UuidVariantFuture    byte = 3
)

////////////////////////////////////////////////////////////////////////////////

// NewUuidTo creates and saves a new Uuid with provided version.
// Only v1, v4, v6, v7 are supported. Returns ErrIDNilDestination if 'to' is nil.
//
// WARNING!
// Returns ErrUuidUnsupportedVersion if any of UuidV2, UuidV3, UuidV5 is passed
// as a version. Remember, only v1, v4, v6 and v7 are supported!
func NewUuidTo(to *Uuid, v byte) error {

	if to == nil {
		return ErrIDNilDestination
	}

	var u uuid.UUID
	var err error

	switch v {
	case UuidV1:
		u, err = uuid.NewV1()
	case UuidV4:
		u, err = uuid.NewV4()
	case UuidV6:
		u, err = uuid.NewV6()
	case UuidV7:
		u, err = uuid.NewV7()
	default:
		return ErrUuidUnsupportedVersion
	}

	if err != nil {
		return err
	}

	*to = Uuid(u)
	return nil
}

// NewUuidFromStringTo parses string and saves Uuid to the destination.
// Input is expected in a form accepted by UnmarshalText.
//
// Returns ErrIDNilDestination if 'to' is nil.
func NewUuidFromStringTo(to *Uuid, input string) error {
	return to.UnmarshalText(ekastr.ToBytes(input))
}

////////////////////////////////////////////////////////////////////////////////

// Version returns algorithm version used to generate Uuid.
// Returns 0 if Uuid is nil.
func (u *Uuid) Version() byte {
	if u == nil {
		return 0
	}
	return u[6] >> 4
}

// SetVersion sets version bits. Does nothing if Uuid is nil.
func (u *Uuid) SetVersion(v byte) {
	if u != nil {
		u[6] = (u[6] & 0x0f) | (v << 4)
	}
}

// Variant returns Uuid layout variant.
// Returns UuidVariantFuture if Uuid is nil.
func (u *Uuid) Variant() byte {
	switch {
	case u == nil:
		fallthrough
	case (u[8] >> 5) == 0x07:
		fallthrough
	default:
		return UuidVariantFuture
	case (u[8] >> 7) == 0x00:
		return UuidVariantNcs
	case (u[8] >> 6) == 0x02:
		return UuidVariantRfc4122
	case (u[8] >> 5) == 0x06:
		return UuidVariantMicrosoft
	}
}

// SetVariant sets variant bits. Does nothing if Uuid is nil.
func (u *Uuid) SetVariant(v byte) {
	if u == nil {
		return
	}
	switch v {
	case UuidVariantNcs:
		u[8] = u[8]&(0xff>>1) | (0x00 << 7)
	case UuidVariantRfc4122:
		u[8] = u[8]&(0xff>>2) | (0x02 << 6)
	case UuidVariantMicrosoft:
		u[8] = u[8]&(0xff>>3) | (0x06 << 5)
	case UuidVariantFuture:
		fallthrough
	default:
		u[8] = u[8]&(0xff>>3) | (0x07 << 5)
	}
}

////////////////////////////////////////////////////////////////////////////////

// Equal returns true if both of Uuid s are equal, otherwise returns false.
// It returns true if both of them are nil but false if only one.
func (u *Uuid) Equal(other ID) bool {
	if b1, b2 := u == nil, other == nil; b1 || b2 {
		return b1 && b2
	}
	var otherTyped, _ = other.(*Uuid)
	return bytes.Equal(u[:], otherTyped[:])
}

// IsNil reports whether current Uuid is empty (nil).
func (u *Uuid) IsNil() bool {
	return u == nil || bytes.Equal(u[:], uuidNil[:])
}

// SetNil sets the current Uuid to zero Uuid. Does nothing if Uuid is nil.
func (u *Uuid) SetNil() {
	if u != nil {
		*u = uuidNil
	}
}

// Bytes returns bytes slice representation of Uuid. Returns nil if Uuid is nil.
func (u *Uuid) Bytes() []byte {
	if u == nil {
		return nil
	}
	return u[:]
}

// String returns canonical string representation of Uuid:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx. Returns "<nil>" if Uuid is nil.
func (u *Uuid) String() string {
	if u == nil {
		return "<nil>"
	}
	return uuid.UUID(*u).String()
}

////////////////////////////////////////////////////////////////////////////////

// MarshalText implements the encoding.TextMarshaler interface by
// returning the string encoded Uuid. Returns "null" if Uuid is nil.
func (u *Uuid) MarshalText() ([]byte, error) {
	if u == nil {
		return ekaenc.NullAsBytesLowerCase(), nil
	}
	return ekastr.ToBytes(uuid.UUID(*u).String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface by parsing
// the data as string encoded Uuid. Returns ErrIDNilDestination if Uuid is nil.
func (u *Uuid) UnmarshalText(text []byte) error {
	if u == nil {
		return ErrIDNilDestination
	}
	return (*uuid.UUID)(u).UnmarshalText(text)
}

////////////////////////////////////////////////////////////////////////////////

// MarshalJSON implements the json.Marshaler interface.
// Returns a JSON string with encoded Uuid with the followed format:
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8". Returns JSON null if Uuid is nil.
func (u *Uuid) MarshalJSON() ([]byte, error) {

	if u == nil {
		return ekaenc.NullAsBytesLowerCase(), nil
	}

	var data, err = u.MarshalText()
	if err != nil {
		return nil, err
	}

	var buf = make([]byte, len(data)+2)
	buf[0], buf[len(buf)-1] = '"', '"'

	copy(buf[1:len(buf)-1], data)
	return buf, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// Decodes data as encoded JSON Uuid string and saves the result to u.
// Supports all Uuid variants that UnmarshalText() does support but also
// supports JSON null values.
func (u *Uuid) UnmarshalJSON(data []byte) error {

	switch {
	case u == nil:
		return ErrIDNilDestination

	case len(data) == 0 || ekaenc.IsNullAsBytes(data):
		u.SetNil()
		return nil
	}

	return (*uuid.UUID)(u).UnmarshalText(data[1 : len(data)-1])
}

////////////////////////////////////////////////////////////////////////////////

// MarshalBinary implements the encoding.BinaryMarshaler interface by
// returning the Uuid as a byte slice. Returns nil if Uuid is nil.
func (u *Uuid) MarshalBinary() ([]byte, error) {

	if u == nil {
		return nil, nil
	}

	var buf = make([]byte, len(u))

	copy(buf, u[:])
	return buf, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface by
// copying the passed data and converting it to an Uuid.
// Returns ErrIDNilDestination if Uuid is nil.
func (u *Uuid) UnmarshalBinary(data []byte) error {

	switch {
	case u == nil:
		return ErrIDNilDestination

	case len(data) != len(u):
		return errors.New("uuid: incorrect data length for binary unmarshal")
	}

	copy(u[:], data)
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// Value implements the driver.Valuer interface. This returns the value
// represented as a string.
func (u *Uuid) Value() (driver.Value, error) {
	return u.MarshalText()
}

// Scan implements the sql.Scanner interface. It supports scanning
// a string or byte slice. Returns ErrIDNilDestination if Uuid is nil.
func (u *Uuid) Scan(src any) error {

	if u == nil {
		return ErrIDNilDestination
	}

	return (*uuid.UUID)(u).Scan(src)
}

////////////////////////////////////////////////////////////////////////////////

var _ ID = (*Uuid)(nil)
