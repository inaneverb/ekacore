// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
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
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net"
	"sync"
	"time"

	"github.com/qioalice/ekago/v2/ekasys"
)

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
type (
	// _UUID_Generator provides interface for generating UUIDs.
	_UUID_Generator interface {
		NewV1() (UUID, error)
		NewV2(domain byte) (UUID, error)
		NewV3(ns UUID, name string) UUID
		NewV4() (UUID, error)
		NewV5(ns UUID, name string) UUID
	}

	// Default generator implementation.
	_T_UUID_RFC4122_Generator struct {
		storageMutex sync.Mutex

		hwAddr    [6]byte
		hwAddrErr error

		clockSequence    uint16
		clockSequenceErr error

		rand io.Reader

		lastTime uint64
	}
)

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
const (
	_UUID_SIZE = 16 // size of a UUID in bytes.

	// Difference in 100-nanosecond intervals between
	// UUID epoch (October 15, 1582) and Unix epoch (January 1, 1970).
	_UUID_EPOCH_START = 122192928000000000
)

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
var (
	_UUID_RFC4122_Generator = newRFC4122Generator(rand.Reader)

	// String parse helpers.
	_UUID_URN_Prefix = []byte("urn:uuid:")
	_UUID_ByteGroups = []int{8, 4, 4, 4, 12}

	_UUID_JSON_NULL = []byte("null")
)

// hexEncodeTo encodes UUID as hex to dest with canonical string representation:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (w/o quotes). Requires: len(dest) >= 36.
// Returns dest.
func (u UUID) hexEncodeTo(dest []byte) []byte {

	dest[8] = '-'
	dest[13] = '-'
	dest[18] = '-'
	dest[23] = '-'

	hex.Encode(dest[0:8], u[0:4])
	hex.Encode(dest[9:13], u[4:6])
	hex.Encode(dest[14:18], u[6:8])
	hex.Encode(dest[19:23], u[8:10])
	hex.Encode(dest[24:], u[10:])

	return dest
}

// jsonMarshal returns canonical string representation of UUID with double quotes:
// "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx".
func (u UUID) jsonMarshal() []byte {

	buf := make([]byte, 38)

	buf[0] = '"'
	buf[37] = '"'

	u.hexEncodeTo(buf[1:37])

	return buf
}

// decodeCanonical decodes UUID string in format
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8".
func (u *UUID) decodeCanonical(t []byte) (err error) {
	if t[8] != '-' || t[13] != '-' || t[18] != '-' || t[23] != '-' {
		return fmt.Errorf("uuid: incorrect UUID format %s", t)
	}

	src := t[:]
	dst := u[:]

	for i, byteGroup := range _UUID_ByteGroups {
		if i > 0 {
			src = src[1:] // skip dash
		}
		_, err = hex.Decode(dst[:byteGroup/2], src[:byteGroup])
		if err != nil {
			return
		}
		src = src[byteGroup:]
		dst = dst[byteGroup/2:]
	}

	return
}

// decodeHashLike decodes UUID string in format
// "6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodeHashLike(t []byte) (err error) {
	src := t[:]
	dst := u[:]

	if _, err = hex.Decode(dst, src); err != nil {
		return err
	}
	return
}

// decodeBraced decodes UUID string in format
// "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}" or in format
// "{6ba7b8109dad11d180b400c04fd430c8}".
func (u *UUID) decodeBraced(t []byte) (err error) {
	l := len(t)

	if t[0] != '{' || t[l-1] != '}' {
		return fmt.Errorf("uuid: incorrect UUID format %s", t)
	}

	return u.decodePlain(t[1 : l-1])
}

// decodeURN decodes UUID string in format
// "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8" or in format
// "urn:uuid:6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodeURN(t []byte) (err error) {
	total := len(t)

	urn_uuid_prefix := t[:9]

	if !bytes.Equal(urn_uuid_prefix, _UUID_URN_Prefix) {
		return fmt.Errorf("uuid: incorrect UUID format: %s", t)
	}

	return u.decodePlain(t[9:total])
}

// decodePlain decodes UUID string in canonical format
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8" or in hash-like format
// "6ba7b8109dad11d180b400c04fd430c8".
func (u *UUID) decodePlain(t []byte) (err error) {
	switch len(t) {
	case 32:
		return u.decodeHashLike(t)
	case 36:
		return u.decodeCanonical(t)
	default:
		return fmt.Errorf("uuid: incorrrect UUID length: %s", t)
	}
}

// NewV1 returns UUID based on current timestamp and MAC address.
func (g *_T_UUID_RFC4122_Generator) NewV1() (UUID, error) {
	u := UUID{}

	timeNow, clockSeq, err := g.getClockSequence()
	if err != nil {
		return _UUID_NULL, err
	}
	binary.BigEndian.PutUint32(u[0:], uint32(timeNow))
	binary.BigEndian.PutUint16(u[4:], uint16(timeNow>>32))
	binary.BigEndian.PutUint16(u[6:], uint16(timeNow>>48))
	binary.BigEndian.PutUint16(u[8:], clockSeq)

	if g.hwAddrErr != nil {
		return _UUID_NULL, err
	}
	copy(u[10:], g.hwAddr[:])

	u.SetVersion(UUID_V1)
	u.SetVariant(UUID_VARIANT_RFC4122)

	return u, nil
}

// NewV2 returns DCE Security UUID based on POSIX UID/GID.
func (g *_T_UUID_RFC4122_Generator) NewV2(domain byte) (UUID, error) {
	u, err := g.NewV1()
	if err != nil {
		return _UUID_NULL, err
	}

	switch domain {
	case UUID_DOMAIN_PERSON:
		binary.BigEndian.PutUint32(u[:], ekasys.PosixCachedUid())
	case UUID_DOMAIN_GROUP:
		binary.BigEndian.PutUint32(u[:], ekasys.PosixCachedGid())
	}

	u[9] = domain

	u.SetVersion(UUID_V2)
	u.SetVariant(UUID_VARIANT_RFC4122)

	return u, nil
}

// NewV3 returns UUID based on MD5 hash of namespace UUID and name.
func (g *_T_UUID_RFC4122_Generator) NewV3(ns UUID, name string) UUID {
	u := g.newFromHash(md5.New(), ns, name)
	u.SetVersion(UUID_V3)
	u.SetVariant(UUID_VARIANT_RFC4122)

	return u
}

// NewV4 returns random generated UUID.
func (g *_T_UUID_RFC4122_Generator) NewV4() (UUID, error) {
	u := UUID{}
	if _, err := io.ReadFull(g.rand, u[:]); err != nil {
		return _UUID_NULL, err
	}
	u.SetVersion(UUID_V4)
	u.SetVariant(UUID_VARIANT_RFC4122)

	return u, nil
}

// NewV5 returns UUID based on SHA-1 hash of namespace UUID and name.
func (g *_T_UUID_RFC4122_Generator) NewV5(ns UUID, name string) UUID {
	u := g.newFromHash(sha1.New(), ns, name)
	u.SetVersion(UUID_V5)
	u.SetVariant(UUID_VARIANT_RFC4122)

	return u
}

// Returns epoch and clock sequence.
func (g *_T_UUID_RFC4122_Generator) getClockSequence() (uint64, uint16, error) {

	if g.clockSequenceErr != nil {
		return 0, 0, g.clockSequenceErr
	}

	g.storageMutex.Lock()
	defer g.storageMutex.Unlock()

	timeNow := _UUID_EPOCH_START + uint64(time.Now().UnixNano()/100)
	// Clock didn't change since last UUID generation.
	// Should increase clock sequence.
	if timeNow <= g.lastTime {
		g.clockSequence++
	}
	g.lastTime = timeNow

	return timeNow, g.clockSequence, nil
}

// Returns UUID based on hashing of namespace UUID and name.
func (_ *_T_UUID_RFC4122_Generator) newFromHash(h hash.Hash, ns UUID, name string) UUID {
	u := UUID{}
	h.Write(ns[:])
	h.Write([]byte(name))
	copy(u[:], h.Sum(nil))

	return u
}

// newRFC4122Generator creates, initializes and returns an RFC4122 UUID generator.
func newRFC4122Generator(randReader io.Reader) _UUID_Generator {

	var (
		g      _T_UUID_RFC4122_Generator
		ifaces []net.Interface
	)

	g.rand = randReader

	ifaces, g.hwAddrErr = net.Interfaces()
	// no need to check err because if err != nil, ifaces is empty

	for _, iface := range ifaces {
		if len(iface.HardwareAddr) >= 6 {
			copy(g.hwAddr[:], iface.HardwareAddr)
		}
	}

	// do not overwrite g.hwAddrErr
	// but update it if g.hwAddr has not been set after the loop above
	if g.hwAddrErr == nil && bytes.Equal(g.hwAddr[:], make([]byte, 6)) {
		g.hwAddrErr = fmt.Errorf("uuid: no HW address found")
	}

	// Initialize g.hwAddr randomly in case of real network interfaces absence.
	if g.hwAddrErr != nil {
		if _, err := io.ReadFull(g.rand, g.hwAddr[:]); err == nil {
			g.hwAddr[0] |= 0x01 // Set multicast bit as recommended by RFC 4122
			g.hwAddrErr = nil   // Ignore all previous errors
		} else {
			g.hwAddrErr = fmt.Errorf("uuid: no HW address found and failed to mock it: %s", err.Error())
		}
	}

	buf := make([]byte, 2)
	if _, err := io.ReadFull(g.rand, buf); err == nil {
		g.clockSequence = binary.BigEndian.Uint16(buf)
	} else {
		g.clockSequenceErr = fmt.Errorf("uuid: failed to generate clock sequence: %s", err.Error())
	}

	return &g
}
