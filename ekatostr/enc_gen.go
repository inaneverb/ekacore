// Copyright © 2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: inaneverb@pm.me, https://github.com/inaneverb
// License: https://opensource.org/licenses/MIT

package ekatostr

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/inaneverb/ekacore/ekaext/v4"
	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

// strconvAppendInt is a strconv.AppendInt() func, but with the same
// func type as strconv.AppendUint() has.
//
//go:linkname strconvAppendInt strconv.AppendInt
func strconvAppendInt([]byte, uint64, int) []byte

func genEncConstStr(s string, skipZero bool) _Encoder {
	return func(to []byte, _ unsafe.Pointer, bh uint8) []byte {
		if !isSkipZero(bh) || !skipZero {
			to = append(to, s...)
		}
		return to
	}
}

func genEncItoa[T ekaext.Integer](cb func([]byte, uint64, int) []byte) _Encoder {
	return func(to []byte, v unsafe.Pointer, bh uint8) []byte {
		if !isSkipZero(bh) || uint64(*(*T)(v)) != 0 {
			return cb(to, uint64(*(*T)(v)), 10)
		}
		return to
	}
}

func genEncFtoa[T ekaext.Float](bs int) _Encoder {
	return func(to []byte, v unsafe.Pointer, bh uint8) []byte {
		if !isSkipZero(bh) || float64(*(*T)(v)) != 0 {
			return strconv.AppendFloat(to, float64(*(*T)(v)), 'f', -1, bs)
		}
		return to
	}
}

func genEncCtoa[T ekaext.Complex](bs int) _Encoder {
	return func(to []byte, v unsafe.Pointer, bh uint8) []byte {
		if !isSkipZero(bh) || complex128(*(*T)(v)) != 0 {
			var s = strconv.FormatComplex(complex128(*(*T)(v)), 'f', -1, bs)
			return append(to, s...)
		}
		return to
	}
}

// genEncFunc generates and returns an _Encoder that encodes 't',
// assuming that it's a Golang's func type.
func genEncFunc(t reflect.Type) _Encoder {
	var ts = t.String()

	return func(to []byte, v unsafe.Pointer, bh uint8) []byte {
		var skipZero = isSkipZero(bh)
		var isZero = v == nil

		switch {
		case skipZero && isZero:
			return to

		case bh&TOSTR_BH_HIGH_ON == 0:
			return to

		case !skipZero && isZero:
			to = append(to, "[<nil>] "...)

		default:
			to = append(to, "[0x"...)
			to = strconv.AppendUint(to, uint64(uintptr(v)), 10)
			to = append(to, "] "...)
		}

		return append(to, ts...)
	}
}

// genEncChan generates and returns an _Encoder that encodes 't',
// assuming that it's a Golang's channel type.
func genEncChan(t reflect.Type) _Encoder {
	return genEncFunc(t)
}

// genEncInterface generates and returns an _Encoder that encodes 't',
// assuming that it's a Golang's interface (with methods) type.
func genEncInterface(t reflect.Type) _Encoder {
	return genEncFunc(t)
}

// genEncPointer generates and returns an _Encoder that encodes 't',
// assuming that it's a Golang's pointer type.
func genEncPointer(t reflect.Type) _Encoder {

	//goland:noinspection GoSnakeCaseUsage
	const GO_PTR_DEEP = 4

	var k = t.Kind()
	var ts = t.String()

	return func(to []byte, v unsafe.Pointer, bh uint8) []byte {
		var v2 = uintptr(v)
		var isZero = v2 == 0

		if !isZero && k != reflect.UnsafePointer {
			v2 = *(*uintptr)(v)
			isZero = v2 == 0
		}

		if isSkipZero(bh) && isZero {
			return to
		}

		if k == reflect.Pointer && bh&TOSTR_BH_GO_PTR != 0 {
			var rtype = ekaunsafe.RTypeOfReflectType(t)

			for i := 0; i < GO_PTR_DEEP && k == reflect.Pointer && v != nil; i++ {
				rtype = ekaunsafe.RTypeOfDerefRType(rtype) // Deref rtype
				t = ekaunsafe.ReflectTypeOfRType(rtype)    // Update reflect.Type
				if k = t.Kind(); k != reflect.Struct {     // Update kind
					v = *(*unsafe.Pointer)(v) // Deref value
				}
			}

			if k != reflect.Pointer {
				return getEnc(rtype)(to, v, bh)
			}
		}

		if k == reflect.Pointer {
			to = append(to, '(')
			to = append(to, ts...)
			to = append(to, ")("...)
		}

		if isZero {
			to = append(to, "<nil>"...)
		} else {
			to = append(to, "0x"...)
			to = strconv.AppendUint(to, uint64(v2), 10)
		}

		if k == reflect.Pointer {
			to = append(to, ')')
		}

		return to
	}
}

// genEncHash generates and returns an _Encoder that reconstructs io.Reader
// from given 'rtype' and runtime 'v' (as a data), and hashing this data.
// It goes hex encoding, if data length <= 64. Sha256 otherwise.
func genEncHash(rtype uintptr) _Encoder {
	return func(to []byte, v unsafe.Pointer, bh uint8) []byte {

		switch {
		case v == nil && isSkipZero(bh):
			return to
		case v == nil:
			return append(to, "<nil>"...)
		}

		//goland:noinspection GoSnakeCaseUsage
		const MAX_LEN_GO_HEX = 64

		var r = ekaunsafe.PackInterface(rtype, v).(io.Reader)
		var head [MAX_LEN_GO_HEX + 1]byte // First bytes from 'r'
		var n, _ = r.Read(head[:])

		switch {
		case n == 0 && isSkipZero(bh):
			return to
		case n == 0:
			return append(to, "<nil>"...)
		}

		// Go hex, if we have read bytes <= MAX_LEN_GO_HEX.
		// Moreover, we already read them. Otherwise, go sha256.

		if n > MAX_LEN_GO_HEX {
			to = append(to, "sha256:"...)
		}

		var b = bytes.NewBuffer(to)
		var hexEnc = hex.NewEncoder(b)

		if n > MAX_LEN_GO_HEX {
			var h = sha256.New()           // Create new SHA-256 hasher
			_, _ = h.Write(head[:])        // Hash first 65 bytes
			_, _ = io.Copy(h, r)           // Hash the rest ones
			_ = h.Sum(head[:0])            // Write hash to `head`
			_, _ = hexEnc.Write(head[:32]) // Hex the SHA-256 hash
		} else {
			_, _ = hexEnc.Write(head[:n]) // Hex the real data
		}

		return b.Bytes()
	}
}

func genEncStringer(rtype uintptr) _Encoder {
	return func(to []byte, v unsafe.Pointer, bh uint8) []byte {

		switch {
		case v == nil && isSkipZero(bh):
			return to
		case v == nil:
			return append(to, "<nil>"...)
		}

		var sr = ekaunsafe.PackInterface(rtype, v).(fmt.Stringer)
		var s = sr.String()
		return encStr(to, unsafe.Pointer(&s), bh)
	}
}

func genEncSlice(t reflect.Type) _Encoder {

	var tElem = t.Elem()
	var e = genEnc(ekaunsafe.RTypeOfReflectType(tElem))
	var elemSize = tElem.Size()

	var fixedSize = -1
	var mark byte = '~'

	if t.Kind() == reflect.Array {
		fixedSize = t.Len()
		mark = '&'
	}

	return func(to []byte, v unsafe.Pointer, bh uint8) []byte {
		switch {
		case v == nil && isSkipZero(bh):
			return to

		case v == nil:
			return append(to, "<nil>"...)

		case bh&TOSTR_BH_LOW_JSON != 0:
			return encJson(to, t, v)

		case bh&TOSTR_BH_LOW_ON == 0:
			return to
		}

		to = append(to, '[')
		var sh = (*reflect.SliceHeader)(v)

		if fixedSize == 0 || (fixedSize == -1 && sh.Len == 0) {
			return append(to, ']')
		}

		if fixedSize > -1 {
			var rtype = ekaunsafe.RTypeOfReflectType(t)
			var rv = reflect.ValueOf(ekaunsafe.PackInterface(rtype, v))
			rv = rv.Slice(0, fixedSize)
			v = ekaunsafe.UnpackInterface(rv.Interface()).Word
			sh = (*reflect.SliceHeader)(v)
		}

		to = append(to, '<', mark)
		to = strconvAppendInt(to, uint64(sh.Len), 10)
		to = append(to, '>', ' ')

		v = unsafe.Add(nil, sh.Data)
		for i := 0; i < sh.Len; i++ {
			to = e(to, v, bh)
			to = append(to, ',')
			v = unsafe.Add(v, elemSize)
		}

		if to[len(to)-1] == ',' {
			to = to[:len(to)-1]
		}

		return append(to, ']')
	}
}

func genEncMap(t reflect.Type) _Encoder {

	var kRtype = ekaunsafe.RTypeOfReflectType(t.Key())
	var vRtype = ekaunsafe.RTypeOfReflectType(t.Elem())

	var kIsDynamic = kRtype == ekaunsafe.RTypeAny()
	var vIsDynamic = vRtype == ekaunsafe.RTypeAny()

	var ek = genEnc(kRtype)
	var ev = genEnc(vRtype)

	return func(to []byte, v unsafe.Pointer, bh uint8) []byte {
		switch {
		case v == nil && isSkipZero(bh):
			return to

		case v == nil:
			return append(to, "<nil>"...)

		case bh&TOSTR_BH_LOW_JSON != 0:
			return encJson(to, t, v)

		case bh&TOSTR_BH_LOW_ON == 0:
			return to
		}

		var i = ekaunsafe.PackInterface(ekaunsafe.RTypeOfReflectType(t), v)
		var rv = reflect.ValueOf(i)
		var iter = rv.MapRange()

		to = append(to, '{')
		for iter.Next() {

			var iK = ekaunsafe.UnpackInterface(iter.Key().Interface())
			var iV = ekaunsafe.UnpackInterface(iter.Value().Interface())

			var n1 = len(to)
			if kIsDynamic {
				to = getEnc(iK.Type)(to, iK.Word, bh)
			} else {
				to = ek(to, iK.Word, bh)
			}

			var n2 = len(to)
			if n1 == n2 {
				continue
			}

			to = append(to, ':', ' ')

			if vIsDynamic {
				to = getEnc(iV.Type)(to, iV.Word, bh)
			} else {
				to = ev(to, iV.Word, bh)
			}

			var n3 = len(to)
			if n2 == n3 {
				to = to[:n1] // If key was added, but not value, "remove" key.
				continue
			}

			to = append(to, ',', ' ')
		}

		if to[len(to)-2] == ',' && to[len(to)-1] == ' ' {
			to = to[:len(to)-2]
		}

		return append(to, '}')
	}
}

// genEncStruct generates and returns an _Encoder that encodes given 't',
// assuming that it's some Golang's struct.
// It has a complicated and tricky way to perform encoding. During generation
// of encoder, it parses given type, extracting meta info about fields:
// their offsets, their names. Then a special _Encoder is prepared for each
// that field, wrapping with field's name.
// At the moment of executing _Encoder, the prepared set of fields' _Encoder s
// will be used, iterating over struct using address algebra.
func genEncStruct(t reflect.Type) _Encoder {

	type _FieldEncoder struct {
		Offset  uintptr
		Encoder _Encoder
	}

	// We will collect all exported and sized fields to this array.
	// Allocate some capacity (of maximum possible) for fields meta-info.

	var fiEnc = make([]_FieldEncoder, 0, t.NumField())

	// Iterate over fields, but we will check field only if it's public
	// (exported) and sized (the value has size >= 1 byte in RAM).
	// For each that field we determine, what kind of _Encoder shall be used.

	for i, n := 0, t.NumField(); i < n; i++ {
		var sf = t.Field(i)

		if sf.Type.Size() == 0 || !sf.IsExported() {
			continue
		}

		// genEnc here is kinda recursive call, cause this func (genEncStruct)
		// is a part of genEnc function.
		// Since we're already in genEnc, it means that write-lock is acquired.
		// Moreover, genEnc may be used as getter, since it has "final check"
		// for getEnc that suits our needs as an encoder's getter.

		var e = genEnc(ekaunsafe.RTypeOfReflectType(sf.Type)) // <- recursive
		e = wrapEncWithPrefix(sf.Name+"=", e)
		fiEnc = append(fiEnc, _FieldEncoder{sf.Offset, e})
	}

	// TODO: Meticulous optimization? Here we re-allocate final slice
	//  to get an array of actual used fields encoders, instead of holding
	//  some unused memory in slice capacity. Maybe it's too much?

	if len(fiEnc) != cap(fiEnc) {
		var fiEncNew = make([]_FieldEncoder, len(fiEnc))
		copy(fiEncNew, fiEnc)
		fiEnc, fiEncNew = fiEncNew, nil
	}

	// We're ready to present our final encoder for that struct.
	// Just iterate over slice of fields' encoders.

	return func(to []byte, v unsafe.Pointer, bh uint8) []byte {
		switch {
		case v == nil && isSkipZero(bh):
			return to

		case v == nil:
			return append(to, "<nil>"...)

		case bh&TOSTR_BH_LOW_JSON != 0:
			return encJson(to, t, v)

		case bh&TOSTR_BH_LOW_ON == 0:
			return to

		case len(fiEnc) == 0:
			return append(to, "{}"...)
		}

		to = append(to, '{')
		var n = len(fiEnc)
		for i := 0; i < n-1; i++ {
			to = fiEnc[i].Encoder(to, unsafe.Add(v, fiEnc[i].Offset), bh)
			to = append(to, ", "...)
		}

		to = fiEnc[n-1].Encoder(to, unsafe.Add(v, fiEnc[n-1].Offset), bh)
		return append(to, '}')
	}
}
