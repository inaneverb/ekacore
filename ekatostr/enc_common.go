// Copyright © 2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: inaneverb@pm.me, https://github.com/inaneverb
// License: https://opensource.org/licenses/MIT

package ekatostr

import (
	"bytes"
	"reflect"
	"unsafe"

	"github.com/goccy/go-json"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

func encBool(to []byte, v unsafe.Pointer, bh uint8) []byte {
	switch b := *(*bool)(v); {
	case b:
		return append(to, "true"...)
	case !b && !isSkipZero(bh):
		return append(to, "false"...)
	}
	return to
}

func encStr(to []byte, v unsafe.Pointer, bh uint8) []byte {
	var s = *(*string)(v)
	if !isSkipZero(bh) || s != "" {
		to = append(to, '"')
		to = append(to, s...)
		to = append(to, '"')
	}
	return to
}

func encIface(to []byte, v unsafe.Pointer, bh uint8) []byte {
	var i = (*ekaunsafe.Interface)(v)
	return getEnc(i.Type)(to, i.Word, bh)
}

// encJson encodes 't', that may represent any Golang values (except nil)
// using 'v' as a source of data.
func encJson(to []byte, t reflect.Type, v unsafe.Pointer) []byte {

	var b = bytes.NewBuffer(to)
	var n = len(to)
	var rtype = ekaunsafe.RTypeOfReflectType(t)

	var err = json.NewEncoder(b).Encode(ekaunsafe.PackInterface(rtype, v))
	to = b.Bytes() // Maybe encoder replaces our buffer?

	// Technically, an error may be returned only if unsupported JSON type
	// is provided, but we have that check above. So, only array, map or struct
	// should be here.
	// Anyway, just for case, we'll overwrite any output that may be already
	// written to the 'to' with an error string. That's why we need a backup
	// of the writing offset in original buffer ('n').

	if err != nil {
		to = (to)[:n]
		to = append(to, "JsonEncodeErr: "...)
		to = append(to, err.Error()...)
	}

	if to[len(to)-1] == '\n' {
		return to[:len(to)-1]
	}

	return to
}
