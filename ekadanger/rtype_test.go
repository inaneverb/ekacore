// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekadanger_test

import (
	"fmt"
	"testing"

	"github.com/qioalice/ekago/v2/ekadanger"
	"github.com/qioalice/ekago/v2/ekamath"

	"github.com/stretchr/testify/assert"

	"github.com/modern-go/reflect2"
)

//goland:noinspection GoRedundantConversion,GoBoolExpressions
var (
	tda = []interface{}{
		// DO NOT CHANGE AN ORDER OF EXISTED ELEMENTS!
		// DO NOT CHANGE AN EXISTED ELEMENTS AT ALL!
		// ONLY ADD A NEW ONES IF YOU NEED.
		bool(0 == 0),    // 0
		byte(0),         // 1
		rune(0),         // 2
		int(0),          // 3
		int8(0),         // 4
		int16(0),        // 5
		int32(0),        // 6
		int64(0),        // 7
		uint(0),         // 8
		uint8(0),        // 9
		uint16(0),       // 10
		uint32(0),       // 11
		uint64(0),       // 12
		float32(0),      // 13
		float64(0),      // 14
		string(""),      // 15
		[]string(nil),   // 16
		[]byte(nil),     // 17
	}
	td1 = []struct{
		f func() uintptr
		eq uint64
	}{
		{ f: ekadanger.RTypeBool,        eq: 1 << 0          },
		{ f: ekadanger.RTypeByte,        eq: 1 << 1 | 1 << 9 },
		{ f: ekadanger.RTypeRune,        eq: 1 << 2 | 1 << 6 },
		{ f: ekadanger.RTypeInt,         eq: 1 << 3          },
		{ f: ekadanger.RTypeInt8,        eq: 1 << 4          },
		{ f: ekadanger.RTypeInt16,       eq: 1 << 5          },
		{ f: ekadanger.RTypeInt32,       eq: 1 << 6 | 1 << 2 },
		{ f: ekadanger.RTypeInt64,       eq: 1 << 7          },
		{ f: ekadanger.RTypeUint,        eq: 1 << 8          },
		{ f: ekadanger.RTypeUint8,       eq: 1 << 9 | 1 << 1 },
		{ f: ekadanger.RTypeUint16,      eq: 1 << 10         },
		{ f: ekadanger.RTypeUint32,      eq: 1 << 11         },
		{ f: ekadanger.RTypeUint64,      eq: 1 << 12         },
		{ f: ekadanger.RTypeFloat32,     eq: 1 << 13         },
		{ f: ekadanger.RTypeFloat64,     eq: 1 << 14         },
		{ f: ekadanger.RTypeString,      eq: 1 << 15         },
		{ f: ekadanger.RTypeStringArray, eq: 1 << 16         },
		{ f: ekadanger.RTypeByteArray,   eq: 1 << 17         },
	}
	td2 = []struct{
		f func(uintptr) bool
		eq uint64
	}{
		{ f: ekadanger.RTypeIsAnyNumeric, eq: 4095 << 1        }, // [1..12] as idx
		{ f: ekadanger.RTypeIsAnyReal,    eq: 16383 << 1       }, // [1..14] as idx
		{ f: ekadanger.RTypeIsIntAny,     eq: 31 << 3 | 1 << 2 }, // [3..7,2] as idx
		{ f: ekadanger.RTypeIsIntFixed,   eq: 15 << 4 | 1 << 2 }, // [4..7,2] as idx
		{ f: ekadanger.RTypeIsUintAny,    eq: 31 << 8 | 1 << 1 }, // [8..12,1] as idx
		{ f: ekadanger.RTypeIsUintFixed,  eq: 15 << 9 | 1 << 1 }, // [9..12,1] as idx
		{ f: ekadanger.RTypeIsFloatAny,   eq: 3 << 13          }, // [13..14] as idx
	}
	pt = []struct{
		f func() uintptr
		z uintptr
		name string
	}{
		{ f: ekadanger.RTypeBool,        z: reflect2.RTypeOf(tda[0]),  name: "RTypeBool"        },
		{ f: ekadanger.RTypeByte,        z: reflect2.RTypeOf(tda[1]),  name: "RTypeByte"        },
		{ f: ekadanger.RTypeRune,        z: reflect2.RTypeOf(tda[2]),  name: "RTypeRune"        },
		{ f: ekadanger.RTypeInt,         z: reflect2.RTypeOf(tda[3]),  name: "RTypeInt"         },
		{ f: ekadanger.RTypeInt8,        z: reflect2.RTypeOf(tda[4]),  name: "RTypeInt8"        },
		{ f: ekadanger.RTypeInt16,       z: reflect2.RTypeOf(tda[5]),  name: "RTypeInt16"       },
		{ f: ekadanger.RTypeInt32,       z: reflect2.RTypeOf(tda[6]),  name: "RTypeInt32"       },
		{ f: ekadanger.RTypeInt64,       z: reflect2.RTypeOf(tda[7]),  name: "RTypeInt64"       },
		{ f: ekadanger.RTypeUint,        z: reflect2.RTypeOf(tda[8]),  name: "RTypeUint"        },
		{ f: ekadanger.RTypeUint8,       z: reflect2.RTypeOf(tda[9]),  name: "RTypeUint8"       },
		{ f: ekadanger.RTypeUint16,      z: reflect2.RTypeOf(tda[10]), name: "RTypeUint16"      },
		{ f: ekadanger.RTypeUint32,      z: reflect2.RTypeOf(tda[11]), name: "RTypeUint32"      },
		{ f: ekadanger.RTypeUint64,      z: reflect2.RTypeOf(tda[12]), name: "RTypeUint64"      },
		{ f: ekadanger.RTypeFloat32,     z: reflect2.RTypeOf(tda[13]), name: "RTypeFloat32"     },
		{ f: ekadanger.RTypeFloat64,     z: reflect2.RTypeOf(tda[14]), name: "RTypeFloat64"     },
		{ f: ekadanger.RTypeString,      z: reflect2.RTypeOf(tda[15]), name: "RTypeString"      },
		{ f: ekadanger.RTypeStringArray, z: reflect2.RTypeOf(tda[16]), name: "RTypeStringArray" },
		{ f: ekadanger.RTypeByteArray,   z: reflect2.RTypeOf(tda[17]), name: "RTypeByteArray"   },
	}
)

func testRType(t *testing.T, tdIdx uint8) {
	eqIdx := td1[tdIdx].eq
	for i, z, n := 0, uint64(1), ekamath.MinI(64, len(td1)); i < n; i++ {
		if eqIdx & z > 0 {
			assert.Equal(t, reflect2.RTypeOf(tda[i]), td1[tdIdx].f())
		} else {
			assert.NotEqual(t, reflect2.RTypeOf(tda[i]), td1[tdIdx].f())
		}
		z <<= 1
	}
}

func testRTypeIs(t *testing.T, tdIdx uint8) {
	eqIdx := td2[tdIdx].eq
	for i, z, n := 0, uint64(1), ekamath.MinI(64, len(td1)); i < n; i++ {
		if eqIdx & z > 0 {
			assert.True(t, td2[tdIdx].f(reflect2.RTypeOf(tda[i])))
		} else {
			assert.False(t, td2[tdIdx].f(reflect2.RTypeOf(tda[i])))
		}
		z <<= 1
	}
}

func TestRTypePrintTable(t *testing.T) {
	maxWidth := 0
	for i, n := 0, len(pt); i < n; i++ {
		if l := len(pt[i].name); l > maxWidth {
			maxWidth = l
		}
	}
	fmt.Println()
	fmt.Printf("%-[2]*[1]s | Ekadanger RType | reflect2 RType\n", "RType name", maxWidth)
	lines := make([]byte, maxWidth + 36)
	for i, n := 0, len(lines); i < n; i++ {
		lines[i] = '-'
	}
	fmt.Println(string(lines))
	for i, n := 0, len(pt); i < n; i++ {
		fmt.Printf("%-[4]*[1]s | 0x%-13[2]x | 0x%[3]x\n", pt[i].name, pt[i].f(), pt[i].z, maxWidth)
	}
	fmt.Println()
}

func TestRTypeBool        (t *testing.T) { testRType(t, 0)  }
func TestRTypeByte        (t *testing.T) { testRType(t, 1)  }
func TestRTypeRune        (t *testing.T) { testRType(t, 2)  }
func TestRTypeInt         (t *testing.T) { testRType(t, 3)  }
func TestRTypeInt8        (t *testing.T) { testRType(t, 4)  }
func TestRTypeInt16       (t *testing.T) { testRType(t, 5)  }
func TestRTypeInt32       (t *testing.T) { testRType(t, 6)  }
func TestRTypeInt64       (t *testing.T) { testRType(t, 7)  }
func TestRTypeUint        (t *testing.T) { testRType(t, 8)  }
func TestRTypeUint8       (t *testing.T) { testRType(t, 9)  }
func TestRTypeUint16      (t *testing.T) { testRType(t, 10) }
func TestRTypeUint32      (t *testing.T) { testRType(t, 11) }
func TestRTypeUint64      (t *testing.T) { testRType(t, 12) }
func TestRTypeFloat32     (t *testing.T) { testRType(t, 13) }
func TestRTypeFloat64     (t *testing.T) { testRType(t, 14) }
func TestRTypeString      (t *testing.T) { testRType(t, 15) }
func TestRTypeStringArray (t *testing.T) { testRType(t, 16) }
func TestRTypeByteArray   (t *testing.T) { testRType(t, 17) }


func TestRTypeIsAnyNumeric (t *testing.T) { testRTypeIs(t, 0) }
func TestRTypeIsAnyReal    (t *testing.T) { testRTypeIs(t, 1) }
func TestRTypeIsIntAny     (t *testing.T) { testRTypeIs(t, 2) }
func TestRTypeIsIntFixed   (t *testing.T) { testRTypeIs(t, 3) }
func TestRTypeIsUintAny    (t *testing.T) { testRTypeIs(t, 4) }
func TestRTypeIsUintFixed  (t *testing.T) { testRTypeIs(t, 5) }
func TestRTypeIsFloatAny   (t *testing.T) { testRTypeIs(t, 6) }
