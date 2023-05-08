// Copyright © 2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: inaneverb@pm.me, https://github.com/inaneverb
// License: https://opensource.org/licenses/MIT

package ekatostr

import (
	"bytes"
	"io"
	"strconv"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

func init() {
	gEnc[0] = genEncConstStr("<nil>")

	gEnc[ekaunsafe.RTypeInt()] = genEncItoa[int](strconvAppendInt)
	gEnc[ekaunsafe.RTypeInt8()] = genEncItoa[int8](strconvAppendInt)
	gEnc[ekaunsafe.RTypeInt16()] = genEncItoa[int16](strconvAppendInt)
	gEnc[ekaunsafe.RTypeInt32()] = genEncItoa[int32](strconvAppendInt)
	gEnc[ekaunsafe.RTypeInt64()] = genEncItoa[int64](strconvAppendInt)

	gEnc[ekaunsafe.RTypeUint()] = genEncItoa[uint](strconv.AppendUint)
	gEnc[ekaunsafe.RTypeUint8()] = genEncItoa[uint8](strconv.AppendUint)
	gEnc[ekaunsafe.RTypeUint16()] = genEncItoa[uint16](strconv.AppendUint)
	gEnc[ekaunsafe.RTypeUint32()] = genEncItoa[uint32](strconv.AppendUint)
	gEnc[ekaunsafe.RTypeUint64()] = genEncItoa[uint64](strconv.AppendUint)

	gEnc[ekaunsafe.RTypeUintptr()] = genEncItoa[uintptr](strconv.AppendUint)

	gEnc[ekaunsafe.RTypeFloat32()] = genEncFtoa[float32](32)
	gEnc[ekaunsafe.RTypeFloat64()] = genEncFtoa[float64](64)

	gEnc[ekaunsafe.RTypeComplex64()] = genEncCtoa[complex64](64)
	gEnc[ekaunsafe.RTypeComplex128()] = genEncCtoa[complex128](128)

	gEnc[ekaunsafe.RTypeBool()] = encBool
	gEnc[ekaunsafe.RTypeString()] = encStr

	var ebB = genEncHash(ekaunsafe.RTypeBytesBufferPtr())

	var _ io.Reader = (*bytes.Buffer)(nil)
	gEnc[ekaunsafe.RTypeBytesBufferPtr()] = ebB
	gEnc[ekaunsafe.RTypeBytes()] = wrapEncBytes(ebB)

	gEnc[ekaunsafe.RTypeAny()] = encIface
}
