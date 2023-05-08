// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/inaneverb/ekacore/ekamath/v4"
	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

//goland:noinspection GoRedundantConversion,GoBoolExpressions
var (
	tda = []any{
		// DO NOT CHANGE AN ORDER OF EXISTED ELEMENTS!
		// DO NOT CHANGE AN EXISTED ELEMENTS AT ALL!
		// ONLY ADD A NEW ONES IF YOU NEED.
		bool(0 == 0),           // 0
		byte(0),                // 1
		rune(0),                // 2
		int(0),                 // 3
		int8(0),                // 4
		int16(0),               // 5
		int32(0),               // 6
		int64(0),               // 7
		uint(0),                // 8
		uint8(0),               // 9
		uint16(0),              // 10
		uint32(0),              // 11
		uint64(0),              // 12
		float32(0),             // 13
		float64(0),             // 14
		string(""),             // 15
		[]string(nil),          // 16
		[]byte(nil),            // 17
		[][]byte(nil),          // 18
		map[string]string(nil), // 19
		map[string]any(nil),    // 20
		complex64(0),           // 21
		complex128(0),          // 22
		uintptr(0),             // 23
		unsafe.Pointer(nil),    // 24
		time.Time{},            // 25
		time.Duration(0),       // 26
		bytes.Buffer{},         // 27
		(*bytes.Buffer)(nil),   // 28
	}
	td1 = []struct {
		f  func() uintptr
		eq uint64
	}{
		{f: ekaunsafe.RTypeBool, eq: 1 << 0},
		{f: ekaunsafe.RTypeByte, eq: 1<<1 | 1<<9},
		{f: ekaunsafe.RTypeRune, eq: 1<<2 | 1<<6},
		{f: ekaunsafe.RTypeInt, eq: 1 << 3},
		{f: ekaunsafe.RTypeInt8, eq: 1 << 4},
		{f: ekaunsafe.RTypeInt16, eq: 1 << 5},
		{f: ekaunsafe.RTypeInt32, eq: 1<<6 | 1<<2},
		{f: ekaunsafe.RTypeInt64, eq: 1 << 7},
		{f: ekaunsafe.RTypeUint, eq: 1 << 8},
		{f: ekaunsafe.RTypeUint8, eq: 1<<9 | 1<<1},
		{f: ekaunsafe.RTypeUint16, eq: 1 << 10},
		{f: ekaunsafe.RTypeUint32, eq: 1 << 11},
		{f: ekaunsafe.RTypeUint64, eq: 1 << 12},
		{f: ekaunsafe.RTypeFloat32, eq: 1 << 13},
		{f: ekaunsafe.RTypeFloat64, eq: 1 << 14},
		{f: ekaunsafe.RTypeString, eq: 1 << 15},
		{f: ekaunsafe.RTypeStringArray, eq: 1 << 16},
		{f: ekaunsafe.RTypeBytes, eq: 1 << 17},
		{f: ekaunsafe.RTypeBytesArray, eq: 1 << 18},
		{f: ekaunsafe.RTypeMapStringString, eq: 1 << 19},
		{f: ekaunsafe.RTypeMapStringAny, eq: 1 << 20},
		{f: ekaunsafe.RTypeComplex64, eq: 1 << 21},
		{f: ekaunsafe.RTypeComplex128, eq: 1 << 22},
		{f: ekaunsafe.RTypeUintptr, eq: 1 << 23},
		{f: ekaunsafe.RTypeUnsafePointer, eq: 1 << 24},
		{f: ekaunsafe.RTypeTimeTime, eq: 1 << 25},
		{f: ekaunsafe.RTypeTimeDuration, eq: 1 << 26},
		{f: ekaunsafe.RTypeBytesBuffer, eq: 1 << 27},
		{f: ekaunsafe.RTypeBytesBufferPtr, eq: 1 << 28},
	}
	td2 = []struct {
		f  func(uintptr) bool
		eq uint64
	}{
		{f: ekaunsafe.RTypeIsNumericAny, eq: 4095 << 1},   // [1..12] as idx
		{f: ekaunsafe.RTypeIsRealAny, eq: 16383 << 1},     // [1..14] as idx
		{f: ekaunsafe.RTypeIsIntAny, eq: 31<<3 | 1<<2},    // [3..7,2] as idx
		{f: ekaunsafe.RTypeIsIntFixed, eq: 15<<4 | 1<<2},  // [4..7,2] as idx
		{f: ekaunsafe.RTypeIsUintAny, eq: 31<<8 | 1<<1},   // [8..12,1] as idx
		{f: ekaunsafe.RTypeIsUintFixed, eq: 15<<9 | 1<<1}, // [9..12,1] as idx
		{f: ekaunsafe.RTypeIsFloatAny, eq: 3 << 13},       // [13..14] as idx
		{f: ekaunsafe.RTypeIsComplexAny, eq: 3 << 21},     // [21..22] as idx
		{f: ekaunsafe.RTypeIsStringLike, eq: 5 << 15},     // [15,17] as idx
	}
	pt = []struct {
		f    func() uintptr
		z    uintptr
		name string
	}{
		{f: ekaunsafe.RTypeBool, z: ekaunsafe.RTypeOf(tda[0]), name: "RTypeBool"},
		{f: ekaunsafe.RTypeByte, z: ekaunsafe.RTypeOf(tda[1]), name: "RTypeByte"},
		{f: ekaunsafe.RTypeRune, z: ekaunsafe.RTypeOf(tda[2]), name: "RTypeRune"},
		{f: ekaunsafe.RTypeInt, z: ekaunsafe.RTypeOf(tda[3]), name: "RTypeInt"},
		{f: ekaunsafe.RTypeInt8, z: ekaunsafe.RTypeOf(tda[4]), name: "RTypeInt8"},
		{f: ekaunsafe.RTypeInt16, z: ekaunsafe.RTypeOf(tda[5]), name: "RTypeInt16"},
		{f: ekaunsafe.RTypeInt32, z: ekaunsafe.RTypeOf(tda[6]), name: "RTypeInt32"},
		{f: ekaunsafe.RTypeInt64, z: ekaunsafe.RTypeOf(tda[7]), name: "RTypeInt64"},
		{f: ekaunsafe.RTypeUint, z: ekaunsafe.RTypeOf(tda[8]), name: "RTypeUint"},
		{f: ekaunsafe.RTypeUint8, z: ekaunsafe.RTypeOf(tda[9]), name: "RTypeUint8"},
		{f: ekaunsafe.RTypeUint16, z: ekaunsafe.RTypeOf(tda[10]), name: "RTypeUint16"},
		{f: ekaunsafe.RTypeUint32, z: ekaunsafe.RTypeOf(tda[11]), name: "RTypeUint32"},
		{f: ekaunsafe.RTypeUint64, z: ekaunsafe.RTypeOf(tda[12]), name: "RTypeUint64"},
		{f: ekaunsafe.RTypeFloat32, z: ekaunsafe.RTypeOf(tda[13]), name: "RTypeFloat32"},
		{f: ekaunsafe.RTypeFloat64, z: ekaunsafe.RTypeOf(tda[14]), name: "RTypeFloat64"},
		{f: ekaunsafe.RTypeString, z: ekaunsafe.RTypeOf(tda[15]), name: "RTypeString"},
		{f: ekaunsafe.RTypeStringArray, z: ekaunsafe.RTypeOf(tda[16]), name: "RTypeStringArray"},
		{f: ekaunsafe.RTypeBytes, z: ekaunsafe.RTypeOf(tda[17]), name: "RTypeBytes"},
		{f: ekaunsafe.RTypeBytesArray, z: ekaunsafe.RTypeOf(tda[18]), name: "RTypeBytesArray"},
		{f: ekaunsafe.RTypeMapStringString, z: ekaunsafe.RTypeOf(tda[19]), name: "RTypeMapStringString"},
		{f: ekaunsafe.RTypeMapStringAny, z: ekaunsafe.RTypeOf(tda[20]), name: "RTypeMapStringInterface"},
		{f: ekaunsafe.RTypeComplex64, z: ekaunsafe.RTypeOf(tda[21]), name: "RTypeComplex64"},
		{f: ekaunsafe.RTypeComplex128, z: ekaunsafe.RTypeOf(tda[22]), name: "RTypeComplex128"},
		{f: ekaunsafe.RTypeUintptr, z: ekaunsafe.RTypeOf(tda[23]), name: "RTypeUintptr"},
		{f: ekaunsafe.RTypeUnsafePointer, z: ekaunsafe.RTypeOf(tda[24]), name: "RTypeUnsafePointer"},
		{f: ekaunsafe.RTypeTimeTime, z: ekaunsafe.RTypeOf(tda[25]), name: "RTypeTimeTime"},
		{f: ekaunsafe.RTypeTimeDuration, z: ekaunsafe.RTypeOf(tda[26]), name: "RTypeTimeDuration"},
		{f: ekaunsafe.RTypeBytesBuffer, z: ekaunsafe.RTypeOf(tda[27]), name: "RTypeBytesBuffer"},
		{f: ekaunsafe.RTypeBytesBufferPtr, z: ekaunsafe.RTypeOf(tda[28]), name: "RTypeBytesBufferPtr"},
	}
)

func testRType(t *testing.T, tdIdx uint8) {
	var eqIdx = td1[tdIdx].eq
	for i, z, n := 0, uint64(1), ekamath.Min(64, len(td1)); i < n; i++ {
		if eqIdx&z > 0 {
			assert.Equal(t, ekaunsafe.RTypeOf(tda[i]), td1[tdIdx].f())
		} else {
			assert.NotEqual(t, ekaunsafe.RTypeOf(tda[i]), td1[tdIdx].f())
		}
		z <<= 1
	}
}

func testRTypeIs(t *testing.T, tdIdx uint8) {
	var eqIdx = td2[tdIdx].eq
	for i, z, n := 0, uint64(1), ekamath.Min(64, len(td1)); i < n; i++ {
		if eqIdx&z > 0 {
			assert.True(t, td2[tdIdx].f(ekaunsafe.RTypeOf(tda[i])))
		} else {
			assert.False(t, td2[tdIdx].f(ekaunsafe.RTypeOf(tda[i])))
		}
		z <<= 1
	}
}

func TestRTypePrintTable(t *testing.T) {
	var maxWidth = 0
	for i, n := 0, len(pt); i < n; i++ {
		if l := len(pt[i].name); l > maxWidth {
			maxWidth = l
		}
	}
	fmt.Printf("\n%-[2]*[1]s | Pre-defined | On-the-fly\n", "RType name", maxWidth)
	lines := make([]byte, maxWidth+27)
	for i, n := 0, len(lines); i < n; i++ {
		lines[i] = '-'
	}
	fmt.Printf("%s\n", string(lines))
	for i, n := 0, len(pt); i < n; i++ {
		fmt.Printf("%-[4]*[1]s | %-11[2]d | %[3]d\n", pt[i].name, pt[i].f(), pt[i].z, maxWidth)
	}
	fmt.Printf("\n")
}

func TestRTypeBool(t *testing.T)               { testRType(t, 0) }
func TestRTypeByte(t *testing.T)               { testRType(t, 1) }
func TestRTypeRune(t *testing.T)               { testRType(t, 2) }
func TestRTypeInt(t *testing.T)                { testRType(t, 3) }
func TestRTypeInt8(t *testing.T)               { testRType(t, 4) }
func TestRTypeInt16(t *testing.T)              { testRType(t, 5) }
func TestRTypeInt32(t *testing.T)              { testRType(t, 6) }
func TestRTypeInt64(t *testing.T)              { testRType(t, 7) }
func TestRTypeUint(t *testing.T)               { testRType(t, 8) }
func TestRTypeUint8(t *testing.T)              { testRType(t, 9) }
func TestRTypeUint16(t *testing.T)             { testRType(t, 10) }
func TestRTypeUint32(t *testing.T)             { testRType(t, 11) }
func TestRTypeUint64(t *testing.T)             { testRType(t, 12) }
func TestRTypeFloat32(t *testing.T)            { testRType(t, 13) }
func TestRTypeFloat64(t *testing.T)            { testRType(t, 14) }
func TestRTypeString(t *testing.T)             { testRType(t, 15) }
func TestRTypeStringArray(t *testing.T)        { testRType(t, 16) }
func TestRTypeBytes(t *testing.T)              { testRType(t, 17) }
func TestRTypeBytesArray(t *testing.T)         { testRType(t, 18) }
func TestRTypeMapStringString(t *testing.T)    { testRType(t, 19) }
func TestRTypeMapStringInterface(t *testing.T) { testRType(t, 20) }
func TestRTypeComplex64(t *testing.T)          { testRType(t, 21) }
func TestRTypeComplex128(t *testing.T)         { testRType(t, 22) }
func TestRTypeUintptr(t *testing.T)            { testRType(t, 23) }
func TestRTypeUnsafePointer(t *testing.T)      { testRType(t, 24) }
func TestRTypeTimeTime(t *testing.T)           { testRType(t, 25) }
func TestRTypeTimeDuration(t *testing.T)       { testRType(t, 26) }
func TestRTypeBytesBuffer(t *testing.T)        { testRType(t, 27) }
func TestRTypeBytesBufferPtr(t *testing.T)     { testRType(t, 28) }

func TestRTypeIsAnyNumeric(t *testing.T) { testRTypeIs(t, 0) }
func TestRTypeIsAnyReal(t *testing.T)    { testRTypeIs(t, 1) }
func TestRTypeIsIntAny(t *testing.T)     { testRTypeIs(t, 2) }
func TestRTypeIsIntFixed(t *testing.T)   { testRTypeIs(t, 3) }
func TestRTypeIsUintAny(t *testing.T)    { testRTypeIs(t, 4) }
func TestRTypeIsUintFixed(t *testing.T)  { testRTypeIs(t, 5) }
func TestRTypeIsFloatAny(t *testing.T)   { testRTypeIs(t, 6) }
func TestRTypeIsComplexAny(t *testing.T) { testRTypeIs(t, 7) }
func TestRTypeIsStringLike(t *testing.T) { testRTypeIs(t, 8) }

func BenchmarkRTypeOf(b *testing.B) {
	b.ReportAllocs()

	var j = 0
	for i := 0; i < b.N; i++ {
		_ = ekaunsafe.RTypeOf(tda[j])

		if j == len(tda)-1 {
			j = 0
		} else {
			j++
		}
	}
}
