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
	"crypto/rand"
	"fmt"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

func TestFromBytes(t *testing.T) {

	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b1 := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	u1, err := UUID_FromBytes(b1)

	require.NoError(t, err)
	require.Equal(t, u, u1)

	b2 := []byte{}
	_, err = UUID_FromBytes(b2)

	require.Error(t, err)
}

func BenchmarkFromBytes(b *testing.B) {
	b.ReportAllocs()
	bytes := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	for i := 0; i < b.N; i++ {
		_, _ = UUID_FromBytes(bytes)
	}
}

func TestMarshalBinary(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b1 := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	b2, err := u.MarshalBinary()

	require.NoError(t, err)
	require.Equal(t, b1, b2)
}

func BenchmarkMarshalBinary(b *testing.B) {
	b.ReportAllocs()
	u, _ := UUID_NewV4()
	for i := 0; i < b.N; i++ {
		_, _ = u.MarshalBinary()
	}
}

func TestUnmarshalBinary(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b1 := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	u1 := UUID{}
	err := u1.UnmarshalBinary(b1)

	require.NoError(t, err)
	require.Equal(t, u, u1)

	b2 := []byte{}
	u2 := UUID{}
	err = u2.UnmarshalBinary(b2)

	require.Error(t, err)
}

func TestFromString(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	s1 := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	s2 := "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}"
	s3 := "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	s4 := "6ba7b8109dad11d180b400c04fd430c8"
	s5 := "urn:uuid:6ba7b8109dad11d180b400c04fd430c8"

	_, err := UUID_FromString("")
	require.Error(t, err)

	u1, err := UUID_FromString(s1)
	require.NoError(t, err)
	require.Equal(t, u, u1)

	u2, err := UUID_FromString(s2)
	require.NoError(t, err)
	require.Equal(t, u, u2)

	u3, err := UUID_FromString(s3)
	require.NoError(t, err)
	require.Equal(t, u, u3)

	u4, err := UUID_FromString(s4)
	require.NoError(t, err)
	require.Equal(t, u, u4)

	u5, err := UUID_FromString(s5)
	require.NoError(t, err)
	require.Equal(t, u, u5)
}

func BenchmarkFromString(b *testing.B) {
	b.ReportAllocs()
	str := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	for i := 0; i < b.N; i++ {
		_, _ = UUID_FromString(str)
	}
}

func BenchmarkFromStringUrn(b *testing.B) {
	b.ReportAllocs()
	str := "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	for i := 0; i < b.N; i++ {
		_, _ = UUID_FromString(str)
	}
}

func BenchmarkFromStringWithBrackets(b *testing.B) {
	b.ReportAllocs()
	str := "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}"
	for i := 0; i < b.N; i++ {
		_, _ = UUID_FromString(str)
	}
}

func TestFromStringShort(t *testing.T) {
	// Invalid 35-character UUID string
	s1 := "6ba7b810-9dad-11d1-80b4-00c04fd430c"

	for i := len(s1); i >= 0; i-- {
		_, err := UUID_FromString(s1[:i])
		require.Error(t, err)
	}
}

func TestFromStringLong(t *testing.T) {
	// Invalid 37+ character UUID string
	strings := []string{
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8=",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
		"{6ba7b810-9dad-11d1-80b4-00c04fd430c8}f",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c800c04fd430c8",
	}

	for _, str := range strings {
		_, err := UUID_FromString(str)
		require.Error(t, err)
	}
}

func TestFromStringInvalid(t *testing.T) {
	// Invalid UUID string formats
	strings := []string{
		"6ba7b8109dad11d180b400c04fd430c86ba7b8109dad11d180b400c04fd430c8",
		"urn:uuid:{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
		"uuid:urn:6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"uuid:urn:6ba7b8109dad11d180b400c04fd430c8",
		"6ba7b8109-dad-11d1-80b4-00c04fd430c8",
		"6ba7b810-9dad1-1d1-80b4-00c04fd430c8",
		"6ba7b810-9dad-11d18-0b4-00c04fd430c8",
		"6ba7b810-9dad-11d1-80b40-0c04fd430c8",
		"6ba7b810+9dad+11d1+80b4+00c04fd430c8",
		"(6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
		"{6ba7b810-9dad-11d1-80b4-00c04fd430c8>",
		"zba7b810-9dad-11d1-80b4-00c04fd430c8",
		"6ba7b810-9dad11d180b400c04fd430c8",
		"6ba7b8109dad-11d180b400c04fd430c8",
		"6ba7b8109dad11d1-80b400c04fd430c8",
		"6ba7b8109dad11d180b4-00c04fd430c8",
	}

	for _, str := range strings {
		_, err := UUID_FromString(str)
		require.Error(t, err)
	}
}

func TestFromStringOrNil(t *testing.T) {
	require.Equal(t, _UUID_NULL, UUID_FromString_OrNil(""))
}

func TestFromBytesOrNil(t *testing.T) {
	require.Equal(t, _UUID_NULL, UUID_FromBytes_OrNil([]byte{}))
}

func TestMarshalText(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b1 := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	b2, err := u.MarshalText()
	require.NoError(t, err)
	require.Equal(t, b1, b2)
}

func BenchmarkMarshalText(b *testing.B) {
	b.ReportAllocs()
	u, _ := UUID_NewV4()
	for i := 0; i < b.N; i++ {
		_, _ = u.MarshalText()
	}
}

func TestUnmarshalText(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b1 := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	u1 := UUID{}
	err := u1.UnmarshalText(b1)
	require.NoError(t, err)
	require.Equal(t, u, u1)

	b2 := []byte("")
	u2 := UUID{}
	err = u2.UnmarshalText(b2)
	require.Error(t, err)
}

func BenchmarkUnmarshalText(b *testing.B) {
	b.ReportAllocs()
	bytes := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	u := UUID{}
	for i := 0; i < b.N; i++ {
		_ = u.UnmarshalText(bytes)
	}
}

func BenchmarkMarshalToString(b *testing.B) {
	b.ReportAllocs()
	u, _ := UUID_NewV4()
	for i := 0; i < b.N; i++ {
		_ = u.String()
	}
}

type faultyReader struct {
	callsNum   int
	readToFail int // Read call number to fail
}

func (r *faultyReader) Read(dest []byte) (int, error) {
	r.callsNum++
	if (r.callsNum - 1) == r.readToFail {
		return 0, fmt.Errorf("io: reader is faulty")
	}
	return rand.Read(dest)
}

func TestNewV1(t *testing.T) {
	u1, err := UUID_NewV1()
	require.NoError(t, err)
	require.Equal(t, UUID_V1, u1.Version())
	require.Equal(t, UUID_VARIANT_RFC4122, u1.Variant())

	u2, err := UUID_NewV1()
	require.NoError(t, err)
	require.Equal(t, UUID_V1, u2.Version())
	require.Equal(t, UUID_VARIANT_RFC4122, u2.Variant())

	require.NotEqual(t, u1, u2)
}

func TestNewV1FaultyRand(t *testing.T) {
	g := newRFC4122Generator(new(faultyReader))
	u1, err := g.NewV1()
	require.Error(t, err)
	require.Equal(t, _UUID_NULL, u1)
}

func BenchmarkNewV1(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = UUID_NewV1()
	}
}

func TestNewV2(t *testing.T) {
	u1, err := UUID_NewV2(UUID_DOMAIN_PERSON)
	require.NoError(t, err)
	require.Equal(t, UUID_V2, u1.Version())
	require.Equal(t, UUID_VARIANT_RFC4122, u1.Variant())

	u2, err := UUID_NewV2(UUID_DOMAIN_GROUP)
	require.NoError(t, err)
	require.Equal(t, UUID_V2, u2.Version())
	require.Equal(t, UUID_VARIANT_RFC4122, u2.Variant())

	u3, err := UUID_NewV2(UUID_DOMAIN_ORG)
	require.NoError(t, err)
	require.Equal(t, UUID_V2, u3.Version())
	require.Equal(t, UUID_VARIANT_RFC4122, u3.Variant())
}

func TestNewV2FaultyRand(t *testing.T) {
	g := newRFC4122Generator(new(faultyReader))
	u1, err := g.NewV2(UUID_DOMAIN_PERSON)
	require.Error(t, err)
	require.Equal(t, _UUID_NULL, u1)
}

func BenchmarkNewV2(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = UUID_NewV2(UUID_DOMAIN_PERSON)
	}
}

func TestNewV3(t *testing.T) {
	u1 := UUID_NewV3(UUID_NAMESPACE_DNS, "www.example.com")
	require.Equal(t, UUID_V3, u1.Version())
	require.Equal(t, UUID_VARIANT_RFC4122, u1.Variant())
	require.Equal(t, "5df41881-3aed-3515-88a7-2f4a814cf09e", u1.String())

	u2 := UUID_NewV3(UUID_NAMESPACE_DNS, "example.com")
	require.NotEqual(t, u2, u1)

	u3 := UUID_NewV3(UUID_NAMESPACE_DNS, "example.com")
	require.Equal(t, u3, u2)

	u4 := UUID_NewV3(UUID_NAMESPACE_URL, "example.com")
	require.NotEqual(t, u4, u3)
}

func BenchmarkNewV3(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = UUID_NewV3(UUID_NAMESPACE_DNS, "www.example.com")
	}
}

func TestNewV4(t *testing.T) {
	u1, err := UUID_NewV4()
	require.NoError(t, err)
	require.Equal(t, UUID_V4, u1.Version())
	require.Equal(t, UUID_VARIANT_RFC4122, u1.Variant())

	u2, err := UUID_NewV4()
	require.NoError(t, err)
	require.Equal(t, UUID_V4, u2.Version())
	require.Equal(t, UUID_VARIANT_RFC4122, u2.Variant())

	require.NotEqual(t, u2, u1)
}

func TestNewV4FaultyRand(t *testing.T) {
	g := newRFC4122Generator(new(faultyReader))
	u1, err := g.NewV4()
	require.Error(t, err)
	require.Equal(t, _UUID_NULL, u1)
}

func TestNewV4PartialRead(t *testing.T) {
	g := newRFC4122Generator(iotest.OneByteReader(rand.Reader))
	u1, err := g.NewV4()
	zeros := bytes.Count(u1.Bytes(), []byte{0})
	mostlyZeros := zeros >= 10

	require.NoError(t, err)
	require.False(t, mostlyZeros)
}

func BenchmarkNewV4(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = UUID_NewV4()
	}
}

func TestNewV5(t *testing.T) {
	u1 := UUID_NewV5(UUID_NAMESPACE_DNS, "www.example.com")
	require.Equal(t, UUID_V5, u1.Version())
	require.Equal(t, UUID_VARIANT_RFC4122, u1.Variant())
	require.Equal(t, "2ed6657d-e927-568b-95e1-2665a8aea6a2", u1.String())

	u2 := UUID_NewV5(UUID_NAMESPACE_DNS, "example.com")
	require.NotEqual(t, u2, u1)

	u3 := UUID_NewV5(UUID_NAMESPACE_DNS, "example.com")
	require.Equal(t, u3, u2)

	u4 := UUID_NewV5(UUID_NAMESPACE_URL, "example.com")
	require.NotEqual(t, u4, u3)
}

func BenchmarkNewV5(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = UUID_NewV5(UUID_NAMESPACE_DNS, "www.example.com")
	}
}

func TestValue(t *testing.T) {
	u, err := UUID_FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	require.NoError(t, err)

	val, err := u.Value()
	require.NoError(t, err)
	require.Equal(t, u.String(), val)
}

func TestValueNil(t *testing.T) {
	u := UUID{}

	val, err := u.Value()
	require.NoError(t, err)
	require.Equal(t, nil, val)
}

func TestScanBinary(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b1 := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	u1 := UUID{}
	err := u1.Scan(b1)
	require.NoError(t, err)
	require.Equal(t, u, u1)

	b2 := []byte{}
	u2 := UUID{}

	err = u2.Scan(b2)
	require.Error(t, err)
}

func TestScanString(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	s1 := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

	u1 := UUID{}
	err := u1.Scan(s1)
	require.NoError(t, err)
	require.Equal(t, u, u1)

	s2 := ""
	u2 := UUID{}

	err = u2.Scan(s2)
	require.Error(t, err)
}

func TestScanText(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	b1 := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	u1 := UUID{}
	err := u1.Scan(b1)
	require.NoError(t, err)
	require.Equal(t, u, u1)

	b2 := []byte("")
	u2 := UUID{}
	err = u2.Scan(b2)
	require.Error(t, err)
}

func TestScanUnsupported(t *testing.T) {
	u := UUID{}

	err := u.Scan(true)
	require.Error(t, err)
}

func TestScanNil(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	err := u.Scan(nil)
	require.NoError(t, err)
}

func TestBytes(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	bytes1 := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	require.Equal(t, u.Bytes(), bytes1)
}

func TestString(t *testing.T) {
	require.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c8", UUID_NAMESPACE_DNS.String())
}

func TestEqual(t *testing.T) {
	require.NotEqual(t, UUID_NAMESPACE_DNS, UUID_NAMESPACE_URL)
	require.Equal(t, UUID_NAMESPACE_DNS, UUID_NAMESPACE_DNS)
}

func TestVersion(t *testing.T) {
	u := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	require.Equal(t, UUID_V1, u.Version())
}

func TestSetVersion(t *testing.T) {
	u := UUID{}
	u.SetVersion(4)
	require.Equal(t, UUID_V4, u.Version())
}

func TestVariant(t *testing.T) {
	u1 := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	require.Equal(t, UUID_VARIANT_NCS, u1.Variant())

	u2 := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	require.Equal(t, UUID_VARIANT_RFC4122, u2.Variant())

	u3 := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	require.Equal(t, UUID_VARIANT_MICROSOFT, u3.Variant())

	u4 := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xe0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	require.Equal(t, UUID_VARIANT_FUTURE, u4.Variant())
}

func TestSetVariant(t *testing.T) {
	u := UUID{}

	u.SetVariant(UUID_VARIANT_NCS)
	require.Equal(t, UUID_VARIANT_NCS, u.Variant())

	u.SetVariant(UUID_VARIANT_RFC4122)
	require.Equal(t, UUID_VARIANT_RFC4122, u.Variant())

	u.SetVariant(UUID_VARIANT_MICROSOFT)
	require.Equal(t, UUID_VARIANT_MICROSOFT, u.Variant())

	u.SetVariant(UUID_VARIANT_FUTURE)
	require.Equal(t, UUID_VARIANT_FUTURE, u.Variant())
}

func TestMust(t *testing.T) {
	require.Panics(t, func() {
		UUID_OrPanic(func() (UUID, error) {
			return _UUID_NULL, fmt.Errorf("uuid: expected error")
		}())
	})
}
