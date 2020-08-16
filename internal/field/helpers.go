// Copyright © 2019. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package field

import (
	"io"
	"math"
	"strconv"
)

var (
	fvTrue = []byte("true")
	fvFalse = []byte("false")
	fvQuote = []byte("\"")
	fvNull = []byte("null")
)

// KeyOrUnnamed returns Field's 'Key' field if it's not empty or string
// "unnamed_<unnamedIdx>". Has generated code for 'unnamedIdx' <= 32.
// Increments 'unnamedIdx' before use it (if it have to be used).
//
// Returns an empty string if 'unnamedIdx' == nil and it must be used.
func (f Field) KeyOrUnnamed(unnamedIdx *int) string {

	if f.Key != "" {
		return f.Key
	}

	if unnamedIdx == nil {
		return ""
	}

	*unnamedIdx++

	if *unnamedIdx < 0 || *unnamedIdx > 32 {
		return "unnamed_" + strconv.Itoa(*unnamedIdx)
	}

	switch *unnamedIdx {
	case 0:
		return "unnamed_00"
	case 1:
		return "unnamed_01"
	case 2:
		return "unnamed_02"
	case 3:
		return "unnamed_03"
	case 4:
		return "unnamed_04"
	case 5:
		return "unnamed_05"
	case 6:
		return "unnamed_06"
	case 7:
		return "unnamed_07"
	case 8:
		return "unnamed_08"
	case 9:
		return "unnamed_09"
	case 10:
		return "unnamed_10"
	case 11:
		return "unnamed_11"
	case 12:
		return "unnamed_12"
	case 13:
		return "unnamed_13"
	case 14:
		return "unnamed_14"
	case 15:
		return "unnamed_15"
	case 16:
		return "unnamed_16"
	case 17:
		return "unnamed_17"
	case 18:
		return "unnamed_18"
	case 19:
		return "unnamed_19"
	case 20:
		return "unnamed_20"
	case 21:
		return "unnamed_21"
	case 22:
		return "unnamed_22"
	case 23:
		return "unnamed_23"
	case 24:
		return "unnamed_24"
	case 25:
		return "unnamed_25"
	case 26:
		return "unnamed_26"
	case 27:
		return "unnamed_27"
	case 28:
		return "unnamed_28"
	case 29:
		return "unnamed_29"
	case 30:
		return "unnamed_30"
	case 31:
		return "unnamed_31"
	case 32:
		return "unnamed_32"
	default:
		return ""
	}
}

//
func (f Field) ValueWriteTo(w io.Writer) (n int, err error) {

	if f.IsNil() {
		return w.Write(fvNull)
	}

	switch f.Kind.BaseType() {

	case KIND_TYPE_BOOL:
		if f.IValue != 0 {
			return w.Write(fvTrue)
		} else {
			return w.Write(fvFalse)
		}

	case KIND_TYPE_INT,
		KIND_TYPE_INT_8, KIND_TYPE_INT_16,
		KIND_TYPE_INT_32, KIND_TYPE_INT_64:
		return w.Write([]byte(strconv.FormatInt(f.IValue, 10)))

	case KIND_TYPE_UINT,
		KIND_TYPE_UINT_8, KIND_TYPE_UINT_16,
		KIND_TYPE_UINT_32, KIND_TYPE_UINT_64:
		return w.Write([]byte(strconv.FormatUint(uint64(f.IValue), 10)))

	case KIND_TYPE_FLOAT_32:
		f := float64(math.Float32frombits(uint32(f.IValue)))
		return w.Write([]byte(strconv.FormatFloat(f, 'f', 2, 32)))

	case KIND_TYPE_FLOAT_64:
		f := math.Float64frombits(uint64(f.IValue))
		return w.Write([]byte(strconv.FormatFloat(f, 'f', 2, 64)))

	case KIND_TYPE_STRING:
		if _, err = w.Write(fvQuote); err == nil {
			if n, err = w.Write([]byte(f.SValue)); err == nil {
				if _, err = w.Write(fvQuote); err == nil {
					n += 2
				}
			}
		}
		return n, err

	default:
		return -1, ErrUnsupportedKind
	}
}

// IsZero reports whether f contains zero value of its type (based on kind).
func (f Field) IsZero() bool {

	if f.Kind.IsSystem() {
		switch f.Kind.BaseType() {

		case KIND_SYS_TYPE_EKAERR_UUID, KIND_SYS_TYPE_EKAERR_PUBLIC_MESSAGE,
			KIND_SYS_TYPE_EKAERR_CLASS_NAME:
			return f.SValue == ""

		case KIND_SYS_TYPE_EKAERR_CLASS_ID:
			return f.IValue == 0

		default:
			return true
		}
	}

	switch f.Kind.BaseType() {

	case KIND_TYPE_BOOL,
		KIND_TYPE_INT, KIND_TYPE_INT_8, KIND_TYPE_INT_16, KIND_TYPE_INT_32, KIND_TYPE_INT_64,
		KIND_TYPE_UINT, KIND_TYPE_UINT_8, KIND_TYPE_UINT_16, KIND_TYPE_UINT_32, KIND_TYPE_UINT_64,
		KIND_TYPE_UINTPTR, KIND_TYPE_ADDR,
		KIND_TYPE_FLOAT_32, KIND_TYPE_FLOAT_64:
		return f.IValue == 0

	case KIND_TYPE_STRING:
		return f.SValue == "" || f.SValue == "00000000-0000-0000-0000-000000000000"

	default:
		return f.Value == nil
	}
}
