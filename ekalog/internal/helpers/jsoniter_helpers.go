// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog_helpers

import (
	"math"

	"github.com/qioalice/ekago/v2/internal/field"
	"github.com/qioalice/ekago/v2/internal/letter"

	"github.com/json-iterator/go"
)

// JsonEncodeFields is JSON encoding helper that encodes 'fields' as JSON array,
// adding it to the JSON document's root with the key "fields".
//
// Puts JSON encoded data into 's' stream,
// doing nothing if 'fields' is empty and 'allowEmpty' is false.
func JsonEncodeFields(

	s *jsoniter.Stream, // jsoniter's Stream object
	allowEmpty bool, // write empty JSON array [] if 'fields' is empty
	fields []field.Field, // fields that must be JSON'ed -> 's'

) (wasAdded bool) { // returns true if at least one byte was written to 's'

	emptySet := len(fields) == 0

	if emptySet && !allowEmpty {
		return false
	}

	s.WriteObjectField("fields")

	if emptySet {
		s.WriteEmptyArray()
		return true
	}

	unnamedFieldIdx := 1

	s.WriteArrayStart()

	for i, n := int16(0), int16(len(fields)); i < n; i++ {
		JsonEncodeField(s, fields[i], &unnamedFieldIdx)
		s.WriteMore()
	}

	b := s.Buffer()
	s.SetBuffer(b[:len(b)-1]) // remove last comma

	s.WriteArrayEnd()

	return true
}

// JsonEncodeField is JSON encoding helper that encodes 'f' field as JSON object
// adding it to the JSON document. Uses 'unnamedFieldIdx' as a number
// that will be transformed into string "unnamed_<number>" and that string will
// be used as field's key if its key is empty.
//
// Puts JSON encoded data into 's' stream.
func JsonEncodeField(

	s *jsoniter.Stream,
	f field.Field,
	unnamedFieldIdx *int,
) {

	s.WriteObjectStart()

	s.WriteObjectField("key")
	if f.Key != "" {
		s.WriteString(f.Key)
	} else {
		s.WriteString(letter.UnnamedAsStr(*unnamedFieldIdx))
		*unnamedFieldIdx++
	}
	s.WriteMore()

	// TODO: write kind if requested

	s.WriteObjectField("value")

	if f.Kind.IsNil() {
		s.WriteNil()
		goto exit
	}

	switch f.Kind.BaseType() {

	case field.KIND_TYPE_BOOL:
		s.WriteBool(f.IValue != 0)

	case field.KIND_TYPE_INT:
		s.WriteInt(int(f.IValue))

	case field.KIND_TYPE_INT_8:
		s.WriteInt8(int8(f.IValue))

	case field.KIND_TYPE_INT_16:
		s.WriteInt16(int16(f.IValue))

	case field.KIND_TYPE_INT_32:
		s.WriteInt32(int32(f.IValue))

	case field.KIND_TYPE_INT_64:
		s.WriteInt64(f.IValue)

	case field.KIND_TYPE_UINT:
		s.WriteUint(uint(f.IValue))

	case field.KIND_TYPE_UINT_8:
		s.WriteUint8(uint8(f.IValue))

	case field.KIND_TYPE_UINT_16:
		s.WriteUint16(uint16(f.IValue))

	case field.KIND_TYPE_UINT_32:
		s.WriteUint32(uint32(f.IValue))

	case field.KIND_TYPE_UINT_64:
		s.WriteUint64(uint64(f.IValue))

	case field.KIND_TYPE_FLOAT_32:
		s.WriteFloat32(math.Float32frombits(uint32(f.IValue)))

	case field.KIND_TYPE_FLOAT_64:
		s.WriteFloat64(math.Float64frombits(uint64(f.IValue)))

	case field.KIND_TYPE_STRING:
		s.WriteString(f.SValue)
	}

exit:
	s.WriteObjectEnd()
}