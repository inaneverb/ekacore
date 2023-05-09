// Copyright © 2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: inaneverb@pm.me, https://github.com/inaneverb
// License: https://opensource.org/licenses/MIT

package ekatostr

import (
	"fmt"
	"io"
	"reflect"
	"sync"
	"unsafe"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

type (
	_Encoder = func(to []byte, v unsafe.Pointer, bh uint8) []byte
)

var (
	gEnc   = make(map[uintptr]_Encoder)
	gEncMu sync.RWMutex
)

var (
	gRtFmtStringer = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	gRtIoReader    = reflect.TypeOf((*io.Reader)(nil)).Elem()
)

func isSkipZero(bh uint8) bool {
	return bh&BH_SKIP_ZERO != 0
}

// getEnc tries to find an _Encoder for given 'rtype', and return it.
// If _Encoder does not exist, it will generate it, save it to the cache,
// and also return.
func getEnc(rtype uintptr) _Encoder {

	var e _Encoder
	if e = getEnc2(rtype, true); e != nil {
		// Try to get encoder from cache with read-lock acquiring.
		return e
	}

	gEncMu.Lock()
	defer gEncMu.Unlock()

	return genEnc(rtype)
}

// genEnc performs a "last call" lookup of _Encoder for provided 'rtype'.
// If _Encoder still is not exist, it will generate, store it into the cache
// and return.
func genEnc(rtype uintptr) _Encoder {

	var e _Encoder
	if e = getEnc2(rtype, false); e != nil {
		// Repeat search in the cache, but now with write-lock acquired.
		// The encoder may be already generated from other goroutine,
		// since there's no guarantee, that write-lock has been acquired
		// right after releasing a read-lock.
		return e
	}

	// Check, whether provided value a non-nil one.
	// It's like an early return and easy case.

	var rt = ekaunsafe.ReflectTypeOfRType(rtype)
	var k = rt.Kind()
	if k == reflect.Invalid {
		e = gEnc[0]
		goto exit
	}

	// One more exception. Maybe, the type implements fmt.Stringer or io.Reader?
	// Handle it, but only if interfaces are enabled.

	switch {
	case rt.Implements(gRtFmtStringer):
		e = genEncStringer(rtype)
		goto exit

	case rt.Implements(gRtIoReader):
		e = genEncHash(rtype)
		goto exit
	}

	switch k {
	case reflect.Pointer, reflect.Uintptr, reflect.UnsafePointer:
		e = genEncPointer(rt)

	case reflect.Array, reflect.Slice:
		e = genEncSlice(rt)

	case reflect.Chan:
		e = genEncChan(rt)

	case reflect.Func:
		e = genEncFunc(rt)

	case reflect.Interface:
		e = genEncInterface(rt)

	case reflect.Map:
		e = genEncMap(rt)

	case reflect.Struct:
		e = genEncStruct(rt)

	default:
		e = genEncConstStr(fmt.Sprintf("<unknown_type_%s>", rt.String()), false)
	}

exit:

	gEnc[rtype] = e
	return e
}

// getEnc2 is a helper for getEnc. It exists just to decrease code complexity.
// Returns an _Encoder from the cache (w/o any attempt to generate if not exist),
// acquiring read-only lock if 'rlock' is true.
func getEnc2(rtype uintptr, rlock bool) _Encoder {
	if rlock {
		gEncMu.RLock()
		defer gEncMu.RUnlock()
	}
	return gEnc[rtype]
}
