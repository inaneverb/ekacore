// Copyright © 2013-2023. All rights reserved.
// Author: Maxim Bublis <b@codemonkey.ru>
// Modifier: Ilya Yuryevich <iyuryevich@pm.me, https://github.com/qioalice>
// License: https://opensource.org/licenses/MIT

package ekatyp

// This file provides implementations of the Universally Unique Identifier
// (Uuid), as specified in RFC-4122 and the Peabody RFC Draft (revision 03).
//
// RFC-4122[1] provides the specification for versions 1, 3, 4, and 5. The
// Peabody Uuid RFC Draft[2] provides the specification for the new k-sortable
// UUIDs, versions 6 and 7.
//
// DCE 1.1[3] provides the specification for version 2, but version 2 support
// was removed from this package in v4 due to some concerns with the
// specification itself. Reading the spec, it seems that it would result in
// generating UUIDs that aren't such unique. In having read the spec it seemed
// that our implementation did not meet the spec. It also seems to be at-odds
// with RFC 4122, meaning we would need quite a bit of special code to support
// it. Lastly, there were no Version 2 implementations that we could find to
// ensure we were understanding the specification correctly.
//
// [1] https://tools.ietf.org/html/rfc4122
// [2] https://datatracker.ietf.org/doc/html/draft-peabody-dispatch-new-uuid-format-03
// [3] http://pubs.opengroup.org/onlinepubs/9696989899/chap5.htm#tagcjh_08_02_01_01

// Forked from here: https://github.com/gofrs/uuid

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"net"
	"sync"
	"time"

	"github.com/inaneverb/ekacore/ekaenc/v4"
	"github.com/inaneverb/ekacore/ekaext/v4"
	"github.com/inaneverb/ekacore/ekarand/v4"
	"github.com/inaneverb/ekacore/ekastr/v4"
	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

// CHANGELOG.
// Changelog describes what changes were made comparing to original package.
// - Added ErrUuidInvalidSize.
// - Added Uuid.MarshalTextTo(), Uuid.MarshalBinaryTo().
// - Changed signature of newUuidGenerator():
//   - Added an ability to pass arguments;
//   - Allowed arguments: io.Reader (source of rand data), HW addr func source.
// - Removed NewGenWithHWAF(), newUuidGenerator() just does the job.
// - Renamed types, constants.
// - Removed NullUUID (use Uuid instead).
// - All methods of Uuid are nil-safe and has a pointer-based receiver.
// - Removed Generator interface.

type (
	// Uuid is an array type to represent the value of Uuid
	// as defined in RFC-4122.
	Uuid [_UUID_SIZE_BINARY]byte
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	UUID_V1 byte = 1 // (date-time and MAC address)
	UUID_V2 byte = 2 // (date-time and MAC address, DCE security version)
	UUID_V3 byte = 3 // (namespace name-based)
	UUID_V4 byte = 4 // (random)
	UUID_V5 byte = 5 // (namespace name-based)
	UUID_V6 byte = 6 // (k-sortable timestamp and random, compatible with v1)
	UUID_V7 byte = 7 // (k-sortable timestamp and random)
	UUID_V8 byte = 8 // (k-sortable timestamp, random and user's data)

	UUID_VARIANT_NCS       byte = 0
	UUID_VARIANT_RFC4122   byte = 1
	UUID_VARIANT_MICROSOFT byte = 2
	UUID_VARIANT_FUTURE    byte = 3
)

////////////////////////////////////////////////////////////////////////////////

type (
	// _UuidGenerator is a reference Uuid generator based on the specifications
	// laid out in RFC-4122 and DCE 1.1: Authentication and Security Services.
	_UuidGenerator struct {
		storageMutex  sync.Mutex
		rand          io.Reader
		readFunc      _UuidReaderFunc
		epochFunc     _UuidEpochFunc
		lastTime      uint64
		clockSequence uint16
		hardwareAddr  [6]byte
	}

	// _UuidEpochFunc is type used to get time.Time objects.
	_UuidEpochFunc = func() time.Time

	// _UuidHardwareAddrFunc is type used to provide hardware (MAC) addresses.
	_UuidHardwareAddrFunc = func() (net.HardwareAddr, error)

	_UuidReaderFunc = func([]byte) (int, error)

	// _UuidTimestamp is the count of 100-nanosecond intervals
	// since 00:00:00.00, 15 October 1582 within a UUID_V1.
	_UuidTimestamp uint64
)

//goland:noinspection GoSnakeCaseUsage
const (
	_UUID_SIZE_BINARY = 16 // Size of an Uuid in bytes (and binary format).
	_UUID_SIZE_TEXT   = 36 // Size of a text representation of Uuid in bytes.

	// Amount of 100 nsec in 1 sec.
	_UUID_100NS_IN_SECOND = 10000000

	// Difference in 100-nanosecond intervals between
	// Uuid epoch (October 15, 1582) and Unix epoch (January 1, 1970).
	_UUID_EPOCH_START = 122192928000000000
)

var (
	// uuidGeneratorDefault is generator that uses ekarand.CryptoRandReader()
	// and uuidDefaultHardwareAddrFunc as a hardware address getter.
	uuidGeneratorDefault = newUuidGenerator(nil, rand.Reader)

	// uuidGeneratorFast is the same as uuidGeneratorDefault,
	// but uses ekarand.FastRandReader() instead.
	uuidGeneratorFast = newUuidGenerator(nil, ekarand.GetFastRandReader())

	// uuidNil is a empty (zero) but not nil Uuid.
	uuidNil = Uuid{}

	// uuidIndexes32 is the bytes' indexes of payload in 32-bytes hash uuid
	// (without hyphens) text representation (6ba7b8109dad11d180b400c04fd430c8).
	uuidIndexes32 = [_UUID_SIZE_BINARY]byte{
		0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30,
	}

	// uuidIndexes36 is the bytes' indexes of payload in 36-bytes
	// canonical uuid (with hyphens) text representation
	// (6ba7b810-9dad-11d1-80b4-00c04fd430c8).
	uuidIndexes36 = [_UUID_SIZE_BINARY]byte{
		0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34,
	}

	// uuidRType represents rtype of Uuid. Lightweight than reflect package.
	uuidRType = ekaunsafe.RTypeOf(Uuid{})
)

var (
	ErrUuidUnsupportedVersion = errors.New("uuid: unsupported version")
	ErrUuidInvalidFormat      = errors.New("uuid: invalid format")
)

////////////////////////////////////////////////////////////////////////////////

// NewUuidTo creates and saves a new Uuid with provided version.
// Only UUID_V1, UUID_V4, UUID_V6, UUID_V7 are supported.
//
// Returns ErrUuidUnsupportedVersion if incorrect version is requested.
// Returns ErrIDNilDestination if 'to' is nil.
func NewUuidTo(to *Uuid, v byte) error {
	return newUuidTo(to, v, uuidGeneratorDefault, nil, nil, 0, 0)
}

// NewUuidFastTo is the same as NewUuidTo() (and same rules are applied here),
// but uses custom CSPRNG to generate UUIDs. Thus, it's almost 10x faster.
func NewUuidFastTo(to *Uuid, v byte) error {
	return newUuidTo(to, v, uuidGeneratorFast, nil, nil, 0, 0)
}

// NewUuidWithNamespaceTo allows to create UUID_V3, UUID_V5 versions,
// which are based on provided namespace 'ns' and custom user's data.
// An empty or nil user's data is allowed and won't lead to error.
func NewUuidWithNamespaceTo(to *Uuid, v byte, ns Uuid, data []byte) error {
	return newUuidTo(to, v, uuidGeneratorDefault, &ns, data, 0, 0)
}

// NewUuidWithCustomTo allows to create UUID_V8 version, that is based on
// UNIX timestamp, user's data and random data.
//
// You may specify how many bits from 'data' you want to add to the Uuid,
// using 'n'. There are some rules:
//
//   - 'n' == 0, up to 32 bits are used from 'data',
//     depends on amount of non-zero bytes right-to-left.
//
//   - 'n' > 32, all 32 bits are used.
//
//   - 'n' ∈ [1..32], exactly the number of bits that are used,
//     right-to-left.
//
// Take a look where your 'data' bits will be placed:
//
//	  0                   1                   2                   3
//	 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                           unix_ts_ms                          |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|          unix_ts_ms           |  ver  |        rand           |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|var|                         rand                              |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//	|                         rand & data                           |
//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
// The last 4 bytes of Uuid is where your 'data' may be added.
// All examples, that are listed below shows only them (last 4 bytes),
// and skips first 12 ones.
//
// BitsCount: 14
// Data: 0000 0001 | 0000 0011 | 1100 1110 | 1111 0011
// Uuid: ???? ???? | ???? ???? | ??00 1100 | 1111 0011
//
// If 'n' == 0:
//
// Bytes from 'data' counted right-to-left, and is used if:
// - it's not zero and all next bytes are zero,
// - it's zero, but some next bytes (or all of them) is not zero.
//
// Example:
//
// Data: 0000 0000 | 0010 0000 | 1011 1101 | 1000 0000 (0x0020BD80).
// Uuid: ???? ???? | 0010 0000 | 1011 1101 | 1000 0000.
// Explanation: 4th byte is zero, and ignored (a random data is placed there).
//
// Data: 0000 0000 | 0010 0000 | 0000 0000 | 1000 0000 (0x00200080).
// Uuid: ???? ???? | 0010 0000 | 0000 0000 | 1000 0000.
// Explanation: The same as above, but 2nd byte is zero. They're also used,
// because despite 2nd is zero, 3rd IS NOT. So, we should also count 2nd one.
//
// Data: 0000 0000 | 0000 0000 | 0000 0000 | 1000 0000 (0x00000080).
// Uuid: ???? ???? | ???? ???? | ???? ???? | 1000 0000.
// Explanation: The same as above, but now 2nd and 3rd are zero.
// The same for 4th. So, all of them are ignored.
func NewUuidWithCustomTo(to *Uuid, data uint32, n uint8) error {
	return newUuidTo(to, UUID_V8, uuidGeneratorDefault, nil, nil, data, n)
}

// NewUuidFastWithCustomTo is the same as NewUuidWithCustomTo(),
// but uses custom CSPRNG to generate UUIDs. Thus, it's almost 10x faster.
func NewUuidFastWithCustomTo(to *Uuid, data uint32, n uint8) error {
	return newUuidTo(to, UUID_V8, uuidGeneratorFast, nil, nil, data, n)
}

// NewUuidFromStringTo parses string and saves Uuid to the destination.
// Input is expected in a form accepted by Uuid.UnmarshalText().
//
// Returns ErrIDNilDestination if 'to' is nil.
func NewUuidFromStringTo(to *Uuid, input string) error {
	return to.UnmarshalText(ekastr.ToBytes(input))
}

////////////////////////////////////////////////////////////////////////////////

// UuidExtractTimestamp extracts UNIX timestamp as time.Time from Uuid,
// but only if it's UUID_V1, UUID_V6, UUID_V7, UUID_V8.
// Returns zero time.Time otherwise, if 'u' is nil, empty or incorrect.
//
// WARNING!
// Use this function for UUID_V8 with caution! You should make sure,
// that UNIX timestamp was added, and there's exactly that data, not yours
// or random one.
func UuidExtractTimestamp(u *Uuid) time.Time {

	// We don't need nil check here, cause u.Version() returns 0 for nil Uuid.
	// Thus, it will be default switch-case branch, that returns time.Time{}.

	switch u.Version() {
	case UUID_V1:
		var low = uint64(binary.BigEndian.Uint32(u[0:4]))
		var mid = uint64(binary.BigEndian.Uint16(u[4:6]))
		var hi = uint64(binary.BigEndian.Uint16(u[6:8]) & 0xFFF)
		return _UuidTimestamp(low + (mid << 32) + (hi << 48)).ToStd()

	case UUID_V6:
		var hi = uint64(binary.BigEndian.Uint32(u[0:4]))
		var mid = uint64(binary.BigEndian.Uint16(u[4:6]))
		var low = uint64(binary.BigEndian.Uint16(u[6:8]) & 0xFFF)
		return _UuidTimestamp(low + (mid << 12) + (hi << 28)).ToStd()

	case UUID_V7, UUID_V8:
		var ms = (uint64(u[0]) << 40) | (uint64(u[1]) << 32) |
			(uint64(u[2]) << 24) | (uint64(u[3]) << 16) | (uint64(u[4]) << 8) |
			uint64(u[5])
		return time.UnixMilli(int64(ms))

	default:
		return time.Time{}
	}
}

// UuidExtractCustomData extracts user's custom data as uint32 from Uuid,
// but only if it's UUID_V8.
// Returns 0 otherwise, if 'u' is nil, empty or incorrect.
//
// WARNING!
// It's your responsibility to maintain value of 'n'. You should pass the same,
// as was used for the NewUuidWithCustomTo() or NewUuidFastWithCustomTo(),
// during Uuid creation.
func UuidExtractCustomData(u *Uuid, n uint8) uint32 {

	// We don't need nil check here, cause u.Version() returns 0 for nil Uuid.
	// Thus, it will be default switch-case branch, that returns 0.

	switch u.Version() {
	case UUID_V8:
		return u.extractData(n)

	default:
		return 0
	}
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
		u[6] = (u[6] & 0x0F) | (v << 4)
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
		return UUID_VARIANT_FUTURE
	case (u[8] >> 7) == 0x00:
		return UUID_VARIANT_NCS
	case (u[8] >> 6) == 0x02:
		return UUID_VARIANT_RFC4122
	case (u[8] >> 5) == 0x06:
		return UUID_VARIANT_MICROSOFT
	}
}

// SetVariant sets variant bits. Does nothing if Uuid is nil.
func (u *Uuid) SetVariant(v byte) {
	if u == nil {
		return
	}
	switch v {
	case UUID_VARIANT_NCS:
		u[8] = u[8]&(0xFF>>1) | (0x00 << 7)
	case UUID_VARIANT_RFC4122:
		u[8] = u[8]&(0xFF>>2) | (0x02 << 6)
	case UUID_VARIANT_MICROSOFT:
		u[8] = u[8]&(0xFF>>3) | (0x06 << 5)
	case UUID_VARIANT_FUTURE:
		fallthrough
	default:
		u[8] = u[8]&(0xFF>>3) | (0x07 << 5)
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
	var data, _ = u.MarshalText()
	return ekastr.FromBytes(data)
}

////////////////////////////////////////////////////////////////////////////////

// Format implements fmt.Formatter for Uuid values.
//
// Rules for printf verbs:
//
//   - 'x', 'X' outputs only the hex digits of the Uuid,
//     using lower case for 'x' and upper case for 'X'.
//
//   - 'v', '+v', 's', 'q' outputs the canonical RFC-4122 string representation.
//
//   - 'S' outputs the RFC-4122 format, but with capital hex digits.
//
//   - '#v' outputs the "Go syntax" representation, which is a [16]byte.
//
//   - All other verbs not handled directly by the fmt package (like '%p')
//     are unsupported and will return "%!verb(uuid.Uuid=value)"
//     as recommended by the fmt package.
func (u *Uuid) Format(f fmt.State, c rune) {

	var b []byte

	switch {
	case u == nil:
		return

	case c == 'v' && f.Flag('#'):
		_, _ = fmt.Fprintf(f, "%#v", [_UUID_SIZE_BINARY]byte(*u))
		return

	case c == 'x' || c == 'X':
		b = make([]byte, 32)
		hex.Encode(b, u[:])

	case c == 'v' || c == 's' || c == 'S':
		b = make([]byte, _UUID_SIZE_TEXT)
		_ = u.MarshalTextTo(b)

	case c == 'q':
		b = make([]byte, _UUID_SIZE_TEXT+2)
		_ = u.MarshalTextTo(b[1:37])
		b[0], b[37] = '"', '"'

	default:
		b = make([]byte, _UUID_SIZE_TEXT+15)
		copy(b, ekastr.ToBytes("%!_(uuid.Uuid="))
		b[2], b[50] = byte(c), ')'
		_ = u.MarshalTextTo(b[14:50])
	}

	if c == 'X' || c == 'S' {
		ekastr.ToUpper(b)
	}

	_, _ = f.Write(b)
}

////////////////////////////////////////////////////////////////////////////////

// MarshalText implements the encoding.TextMarshaler interface.
// The encoding is the same as returned by the String() method.
func (u *Uuid) MarshalText() ([]byte, error) {

	if u == nil {
		return ekaenc.NullAsBytesLowerCase(), nil
	}

	var buf = make([]byte, _UUID_SIZE_TEXT)
	var err = u.MarshalTextTo(buf)

	return ekaext.If(err == nil, buf, nil), err
}

// MarshalTextTo does the same as MarshalText() does,
// but writes result to the 'b'. Returns an error, if len(b) != 36.
func (u *Uuid) MarshalTextTo(b []byte) error {
	switch {
	case (u == nil && len(b) < 4) || (u != nil && len(b) != _UUID_SIZE_TEXT):
		return ErrIDNotEnoughBuffer
	case u == nil:
		copy(b, ekaenc.NullAsStringLowerCase())
	default:
		var ti = &uuidIndexes36
		b[8], b[13], b[18], b[23] = '-', '-', '-', '-'
		for i, n := 0, len(ti); i < n; i++ {
			b[ti[i]], b[ti[i]+1] = ekaenc.EncodeHexFull(u[i])
		}
	}
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Following formats are supported:
//
//	"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
//	"{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
//	"urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
//	"6ba7b8109dad11d180b400c04fd430c8"
//	"{6ba7b8109dad11d180b400c04fd430c8}",
//	"urn:uuid:6ba7b8109dad11d180b400c04fd430c8"
//
// ABNF for supported Uuid text representation follows:
//
//	URN := 'urn'
//	Uuid-NID := 'uuid'
//
//	hexdig := '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' |
//	          'a' | 'b' | 'c' | 'd' | 'e' | 'f' |
//	          'A' | 'B' | 'C' | 'D' | 'E' | 'F'
//
//	hexoct := hexdig hexdig
//	2hexoct := hexoct hexoct
//	4hexoct := 2hexoct 2hexoct
//	6hexoct := 4hexoct 2hexoct
//	12hexoct := 6hexoct 6hexoct
//
//	hashlike := 12hexoct
//	canonical := 4hexoct '-' 2hexoct '-' 2hexoct '-' 6hexoct
//
//	plain := canonical | hashlike
//	uuid := canonical | hashlike | braced | urn
//
//	braced := '{' plain '}' | '{' hashlike  '}'
//	urn := URN ':' Uuid-NID ':' plain
func (u *Uuid) UnmarshalText(b []byte) error {

	if u == nil {
		return ErrIDNilDestination
	}

	var n = len(b)

	switch n {
	case 32: // hash
	case 36: // canonical
	case 34, 38:
		if b[0] != '{' || b[n-1] != '}' {
			return fmt.Errorf("uuid: incorrect format in string %q", b)
		}
		b = b[1 : n-1]
	case 41, 45:
		if string(b[:9]) != "urn:uuid:" {
			return fmt.Errorf("uuid: incorrect format in string %q", b[:9])
		}
		b = b[9:]
	default:
		return fmt.Errorf("uuid: incorrect length %d in string %q", len(b), b)
	}

	if n == 36 && b[8]&b[13]&b[18]&b[23]&'-' != b[8]|b[13]|b[18]|b[23]|'-' {
		return fmt.Errorf("uuid: incorrect format in string %q", b)
	}

	var ti = ekaext.If(n == 36, &uuidIndexes36, &uuidIndexes32)

	for i, n := 0, len(ti); i < n; i++ {
		var x, y = b[ti[i]], b[ti[i]+1]
		if u[i] = ekaenc.DecodeHexFull(x, y); u[i] == 0 && x|y != 0 {
			return ErrUuidInvalidFormat
		}
	}

	return nil
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
func (u *Uuid) UnmarshalJSON(b []byte) error {

	switch n := len(b); {
	case u == nil:
		return ErrIDNilDestination

	case n == 0 || ekaenc.IsNullAsBytes(b):
		u.SetNil()
		return nil

	case n < 2:
		return ErrUuidInvalidFormat

	default:
		return u.UnmarshalText(b[1 : n-1])
	}
}

////////////////////////////////////////////////////////////////////////////////

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (u *Uuid) MarshalBinary() ([]byte, error) {

	if u == nil {
		return nil, nil
	}

	var buf = make([]byte, _UUID_SIZE_BINARY)

	copy(buf, u[:])
	return buf, nil
}

// MarshalBinaryTo does the same as MarshalBinary() does,
// but writes result to the 'b'. Returns an error, if len(b) != 16.
func (u *Uuid) MarshalBinaryTo(b []byte) error {

	switch {
	case u != nil && len(b) != _UUID_SIZE_BINARY:
		return ErrIDNotEnoughBuffer

	case u != nil:
		copy(b, u[:])
	}

	return nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// It will return an error if the slice isn't 16 bytes long.
func (u *Uuid) UnmarshalBinary(data []byte) error {

	switch {
	case u == nil:
		return ErrIDNilDestination

	case len(data) != _UUID_SIZE_BINARY:
		return errors.New("uuid: incorrect data length for binary unmarshal")
	}

	copy(u[:], data)
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// Value implements the driver.Valuer interface.
func (u *Uuid) Value() (driver.Value, error) {
	return u.String(), nil
}

// Scan implements the sql.Scanner interface.
// A 16-byte slice will be handled by UnmarshalBinary, while
// a longer byte slice or a string will be handled by UnmarshalText.
func (u *Uuid) Scan(src any) error {

	if u == nil {
		return ErrIDNilDestination
	}

	switch i := ekaunsafe.UnpackInterface(src); {
	case i.Type == uuidRType:

	case i.Type == ekaunsafe.RTypeBytes():
		var src, ub, ut = *(*[]byte)(i.Word), u.UnmarshalBinary, u.UnmarshalText
		return ekaext.If(len(src) == _UUID_SIZE_BINARY, ub, ut)(src)

	case i.Type == ekaunsafe.RTypeString():
		return u.UnmarshalText(ekastr.ToBytes(*(*string)(i.Word)))
	}

	return fmt.Errorf("uuid: cannot convert %T to Uuid", src)
}

////////////////////////////////////////////////////////////////////////////////

// TextLen returns how many bytes will take a text representation of Uuid.
// Text representation may be obtained by MarshalText(), MarshalTextTo()
// or String().
func (*Uuid) TextLen() int { return _UUID_SIZE_TEXT }

// BinaryLen returns how many bytes will take a binary representation of Uuid.
// Text representation may be obtained by MarshalBinary(), MarshalBinaryTo().
func (*Uuid) BinaryLen() int { return _UUID_SIZE_BINARY }

////////////////////////////////////////////////////////////////////////////////

// insertData inserts 'n' bits from 'x' to the Uuid.
func (u *Uuid) insertData(x uint32, n uint8) {

	if n == 0 {
		return
	}

	n = ((n - 1) & 0x1F) + 1 // Clamp upper limit to 32; now n ∈ [1..32].

	var m uint32 = 0xFFFFFFFF >> (0x20 - n) // Prepare mask.

	x &= m // Clear unused bits.

	u[12] = u[12]&byte((^m)>>24) | byte(x>>24)
	u[13] = u[13]&byte((^m)>>16) | byte(x>>16)
	u[14] = u[14]&byte((^m)>>8) | byte(x>>8)
	u[15] = u[15]&byte(^m) | byte(x)
}

// extractData extracts 'n' bits from Uuid and returns them as uint32.
func (u *Uuid) extractData(n uint8) (r uint32) {

	if n == 0 {
		return 0
	}

	n = ((n - 1) & 0x1F) + 1 // Clamp upper limit to 32; now n ∈ [1..32].

	r |= uint32(u[12]) << 24
	r |= uint32(u[13]) << 16
	r |= uint32(u[14]) << 8
	r |= uint32(u[15])

	return r & (0xFFFFFFFF >> (0x20 - n))
}

////////////////////////////////////////////////////////////////////////////////

// newUuidTo is a generic Uuid constructor with custom _UuidGenerator.
func newUuidTo(
	to *Uuid, v byte, gen *_UuidGenerator,
	ns *Uuid, name []byte, data uint32, n uint8) error {

	switch {
	case to == nil:
		return ErrIDNilDestination
	case ns == nil && (v == UUID_V3 || v == UUID_V5):
		return errors.New("uuid: BUG: v3/v5 is requested, but namespace is nil")
	}

	var u Uuid
	var err error

	switch v {
	case UUID_V1:
		gen.NewV1(&u)
	case UUID_V3, UUID_V5:
		gen.NewV3V5(&u, *ns, name, v)
	case UUID_V4:
		err = gen.NewV4(&u)
	case UUID_V6:
		err = gen.NewV6(&u)
	case UUID_V7, UUID_V8:
		err = gen.NewV7V8(&u, data, v, n)
	default:
		return ErrUuidUnsupportedVersion
	}

	if err != nil {
		return err
	}

	*to = u
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// newUuidGenerator returns a new instance of _UuidGenerator with some defaults.
// You can overwrite them, providing specific arguments.
//
// Rules:
//   - Pass io.Reader to use it as a custom rand data source (CSPRNG).
//   - Pass func() (net.HardwareAddr, error) to use it
//     as a custom hardware address getter
//   - All other types are prohibited and ignored.
//   - If you pass more than 1 same types, the last one will be used.
func newUuidGenerator(hwAddrFunc _UuidHardwareAddrFunc, rnd io.Reader) *_UuidGenerator {

	hwAddrFunc = ekaext.If(hwAddrFunc == nil, uuidDefaultHardwareAddrFunc, hwAddrFunc)
	rnd = ekaext.If(rnd == nil, rand.Reader, rnd)

	var gen = _UuidGenerator{epochFunc: time.Now, rand: rnd}

	var buf = make([]byte, 16)
	_, _ = gen.rand.Read(buf)
	gen.clockSequence = binary.BigEndian.Uint16(buf[:2])

	if hardwareAddr, err := hwAddrFunc(); err == nil {
		copy(gen.hardwareAddr[:], hardwareAddr)
	} else {
		copy(gen.hardwareAddr[:], buf[2:])
		gen.hardwareAddr[0] |= 0x01
	}

	return &gen
}

////////////////////////////////////////////////////////////////////////////////

// NewV1 returns an Uuid based on the current timestamp and MAC address.
func (g *_UuidGenerator) NewV1(to *Uuid) {

	var timeNow, clockSeq = g.getClockSequence()

	binary.BigEndian.PutUint32(to[0:], uint32(timeNow))
	binary.BigEndian.PutUint16(to[4:], uint16(timeNow>>32))
	binary.BigEndian.PutUint16(to[6:], uint16(timeNow>>48))
	binary.BigEndian.PutUint16(to[8:], clockSeq)

	copy(to[10:], g.hardwareAddr[:])

	to.SetVersion(UUID_V1)
	to.SetVariant(UUID_VARIANT_RFC4122)
}

// NewV3V5 returns an Uuid based on the MD5 or SHA1 hash
// of the namespace Uuid and name.
func (g *_UuidGenerator) NewV3V5(to *Uuid, ns Uuid, name []byte, v byte) {

	var h = ekaext.If(v == UUID_V3, md5.New, sha1.New)()
	g.newFromHash(to, h, ns, name)

	to.SetVersion(v)
	to.SetVariant(UUID_VARIANT_RFC4122)
}

// NewV4 returns a randomly generated Uuid.
func (g *_UuidGenerator) NewV4(to *Uuid) error {

	// go1.19:
	// io.ReadFull() or io.ReadAtLeast() performs 1 *UNNECESSARY* allocation.
	// I don't know why, and that's why I use just rand.Read() here.

	var q, err = g.rand.Read(to[:])
	if err = g.errRandRead(_UUID_SIZE_BINARY, q, err); err != nil {
		return err
	}

	to.SetVersion(UUID_V4)
	to.SetVariant(UUID_VARIANT_RFC4122)

	return nil
}

type FReader func([]byte) (int, error)

func (f FReader) Read(b []byte) (int, error) {
	return f(b)
}

func NewUuidV4() (Uuid, error) {

	var u Uuid

	var q, err = uuidGeneratorFast.readFunc(u[:])
	if err = uuidGeneratorFast.errRandRead(_UUID_SIZE_BINARY, q, err); err != nil {
		return uuidNil, err
	}

	u.SetVersion(UUID_V4)
	u.SetVariant(UUID_VARIANT_RFC4122)

	return u, nil
}

// NewV6 returns a k-sortable Uuid based on a timestamp and 48 bits of
// pseudorandom data. The timestamp in a V6 Uuid is the same as V1,
// with the bit order being adjusted to allow the Uuid to be k-sortable.
func (g *_UuidGenerator) NewV6(to *Uuid) error {

	// go1.19:
	// io.ReadFull() or io.ReadAtLeast() performs 1 *UNNECESSARY* allocation.
	// I don't know why, and that's why I use just rand.Read() here.

	var q, err = g.rand.Read(to[10:])
	if err = g.errRandRead(_UUID_SIZE_BINARY-10, q, err); err != nil {
		return err
	}

	var timeNow, clockSeq = g.getClockSequence()

	binary.BigEndian.PutUint32(to[0:], uint32(timeNow>>28))   // time_high
	binary.BigEndian.PutUint16(to[4:], uint16(timeNow>>12))   // time_mid
	binary.BigEndian.PutUint16(to[6:], uint16(timeNow&0xFFF)) // time_low
	binary.BigEndian.PutUint16(to[8:], clockSeq&0x3FFF)       // clk_seq_hi_res

	to.SetVersion(UUID_V6)
	to.SetVariant(UUID_VARIANT_RFC4122)

	return nil
}

// NewV7V8 returns a k-sortable Uuid based on
// the current millisecond precision UNIX epoch
// and 74 bits that are shared with pseudorandom data and user's data.
func (g *_UuidGenerator) NewV7V8(to *Uuid, x uint32, v byte, n uint8) error {

	var tn = g.epochFunc()
	var ms = uint64(tn.Unix())*1e3 + uint64(tn.Nanosecond())/1e6

	to[0], to[1], to[2] = byte(ms>>40), byte(ms>>32), byte(ms>>24)
	to[3], to[4], to[5] = byte(ms>>16), byte(ms>>8), byte(ms)

	// go1.19:
	// io.ReadFull() or io.ReadAtLeast() performs 1 *UNNECESSARY* allocation.
	// I don't know why, and that's why I use just rand.Read() here.

	var q, err = g.rand.Read(to[6:])
	if err = g.errRandRead(_UUID_SIZE_BINARY-6, q, err); err != nil {
		return err
	}

	to.insertData(x, n)
	to.SetVersion(v)
	to.SetVariant(UUID_VARIANT_RFC4122)

	return nil
}

////////////////////////////////////////////////////////////////////////////////

// getClockSequence returns the epoch and clock sequence for V1 and V6 UUIDs.
func (g *_UuidGenerator) getClockSequence() (uint64, uint16) {

	g.storageMutex.Lock()
	defer g.storageMutex.Unlock()

	var timeNow = g.getEpoch()

	if timeNow <= g.lastTime {
		// Clock didn't change since last Uuid generation.
		// Should increase clock sequence.
		g.clockSequence++
	}

	g.lastTime = timeNow
	return timeNow, g.clockSequence
}

// Returns the difference between Uuid epoch (October 15, 1582)
// and current time in 100-nanosecond intervals.
func (g *_UuidGenerator) getEpoch() uint64 {
	return _UUID_EPOCH_START + uint64(g.epochFunc().UnixNano()/100)
}

// newFromHash returns the Uuid based on the hashing of the namespace
// from the provided 'ns' Uuid and name.
func (g *_UuidGenerator) newFromHash(
	to *Uuid, h hash.Hash, ns Uuid, name []byte) {

	h.Write(ns[:])
	h.Write(name)

	copy(to[:], h.Sum(nil))
}

// errRandRead check provided arguments and if it looks like an exception,
// returns an error about reading from _UuidGenerator.rand.
//
// Returned error means that reading was unsuccessful, or got fewer bytes
// than required. io.EOF or io.ErrUnexpectedEOF treated as ok, if required
// amount of bytes were read.
func (g *_UuidGenerator) errRandRead(require, got int, orig error) error {

	var isEOF = errors.Is(orig, io.EOF) || errors.Is(orig, io.ErrUnexpectedEOF)

	switch {
	case isEOF && got == _UUID_SIZE_BINARY:
		orig = nil
	case orig != nil:
		// orig will be returned as is; nothing to do here
	case got != require:
		const D = "uuid: BUG: only %d of %d random bytes were read"
		return fmt.Errorf(D, got, require)
	}

	return orig
}

////////////////////////////////////////////////////////////////////////////////

// ToStd returns the UTC time.Time representation of a _UuidTimestamp.
func (t _UuidTimestamp) ToStd() time.Time {
	var secs = int64(t) / _UUID_100NS_IN_SECOND
	var nsecs = 100 * (int64(t) % _UUID_100NS_IN_SECOND)
	return time.Unix(secs-(_UUID_EPOCH_START/_UUID_100NS_IN_SECOND), nsecs)
}

////////////////////////////////////////////////////////////////////////////////

// uuidDefaultHardwareAddrFunc is default net.HardwareAddr getter.
func uuidDefaultHardwareAddrFunc() (net.HardwareAddr, error) {

	var netIFaces, err = net.Interfaces()
	if err != nil {
		return net.HardwareAddr{}, err
	}

	for i, n := 0, len(netIFaces); i < n; i++ {
		if len(netIFaces[i].HardwareAddr) >= 6 {
			return netIFaces[i].HardwareAddr, nil
		}
	}

	return net.HardwareAddr{}, errors.New("uuid: no hardware address found")
}

////////////////////////////////////////////////////////////////////////////////

var _ ID = (*Uuid)(nil)
