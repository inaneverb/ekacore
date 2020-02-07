// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

import (
	"fmt"
	"math"
	"time"
	"unsafe"

	"github.com/qioalice/gext/dangerous"
	"github.com/qioalice/gext/sys"

	"github.com/json-iterator/go"
)

//
type JSONEncoder struct {
	jsoniterConfig jsoniter.Config
	jsoniterAPI    jsoniter.API
}

var (
	// Make sure we won't break API.
	_ CommonIntegratorEncoder = (*JSONEncoder)(nil).encode

	// Package's JSON encoder
	jsonEncoder     CommonIntegratorEncoder
	jsonEncoderAddr unsafe.Pointer
)

func init() {
	jsonEncoder = (&JSONEncoder{}).FreezeAndGetEncoder()
	jsonEncoderAddr = dangerous.TakeRealAddr(jsonEncoder)
}

//
func (je *JSONEncoder) FreezeAndGetEncoder() CommonIntegratorEncoder {
	return je.encode
}

//
func (je *JSONEncoder) encode(e *Entry) []byte {

	je.jsoniterAPI = je.jsoniterConfig.Froze()

	cfg := jsoniter.Config{
		IndentionStep:                 4,
		MarshalFloatWith6Digits:       true,
		EscapeHTML:                    false,
		SortMapKeys:                   false,
		UseNumber:                     false,
		DisallowUnknownFields:         false,
		TagKey:                        "",
		OnlyTaggedField:               false,
		ValidateJsonRawMessage:        false,
		ObjectFieldMustBeSimpleString: true,
		CaseSensitive:                 false,
	}.Froze()

	s := cfg.BorrowStream(nil)
	defer cfg.ReturnStream(s)

	s.WriteObjectStart()

	je.encodeBase(e, s)
	s.WriteMore()

	je.encodeFields(e, s)
	s.WriteMore()

	je.encodeStacktrace(e, s)

	s.WriteObjectEnd()

	buf := s.Buffer()
	copied := make([]byte, len(buf))
	copy(copied, buf)

	return copied
}

// encodeBase encodes e's level, timestamp, message to s.
func (je *JSONEncoder) encodeBase(e *Entry, s *jsoniter.Stream) {

	s.WriteObjectField("level")
	s.WriteString(e.Level.String())
	s.WriteMore()

	s.WriteObjectField("level_value")
	s.WriteUint8(e.Level.uint8())
	s.WriteMore()

	s.WriteObjectField("time")
	s.WriteString(e.Time.Format(time.UnixDate))
	s.WriteMore()

	s.WriteObjectField("message")
	s.WriteString(e.Message)

	if len(e.StackTrace) > 0 {
		s.WriteMore()
		s.WriteObjectField("caller")
		s.WriteString(formCaller2(e.StackTrace[0]))
	}
}

//
func (je *JSONEncoder) encodeFields(e *Entry, s *jsoniter.Stream) {

	unnamedFieldIdx := 1

	s.WriteObjectField("fields")

	if l := len(e.Fields); l > 0 {
		s.WriteArrayStart()

		for _, field := range e.Fields[:l-1] {
			je.encodeField(field, s, &unnamedFieldIdx)
			s.WriteMore()
		}
		je.encodeField(e.Fields[l-1], s, &unnamedFieldIdx)

		s.WriteArrayEnd()

	} else {
		s.WriteEmptyArray()
	}
}

//
func (je *JSONEncoder) encodeStacktrace(e *Entry, s *jsoniter.Stream) {

	s.WriteObjectField("stacktrace")

	if l := len(e.StackTrace); l > 0 {
		s.WriteArrayStart()

		for _, stackFrame := range e.StackTrace[:l-1] {
			je.encodeStackFrame(stackFrame, s)
			s.WriteMore()
		}
		je.encodeStackFrame(e.StackTrace[l-1], s)

		s.WriteArrayEnd()

	} else {
		s.WriteEmptyArray()
	}
}

//
func (je *JSONEncoder) encodeField(f Field, s *jsoniter.Stream, unnamedFieldIdx *int) {

	s.WriteObjectStart()

	s.WriteObjectField("key")
	if f.Key != "" {
		s.WriteString(f.Key)
	} else {
		s.WriteString(implicitUnnamedFieldName(*unnamedFieldIdx))
		*unnamedFieldIdx++
	}
	s.WriteMore()

	// TODO: write kind if requested

	s.WriteObjectField("value")
	switch f.Kind {

	case FieldKindBool:
		s.WriteBool(f.IValue != 0)

	case FieldKindInt:
		s.WriteInt(int(f.IValue))

	case FieldKindInt8:
		s.WriteInt8(int8(f.IValue))

	case FieldKindInt16:
		s.WriteInt16(int16(f.IValue))

	case FieldKindInt32:
		s.WriteInt32(int32(f.IValue))

	case FieldKindInt64:
		s.WriteInt64(f.IValue)

	case FieldKindUint:
		s.WriteUint(uint(f.IValue))

	case FieldKindUint8:
		s.WriteUint8(uint8(f.IValue))

	case FieldKindUint16:
		s.WriteUint16(uint16(f.IValue))

	case FieldKindUint32:
		s.WriteUint32(uint32(f.IValue))

	case FieldKindUint64:
		s.WriteUint64(uint64(f.IValue))

	case FieldKindFloat32:
		s.WriteFloat32(math.Float32frombits(uint32(f.IValue)))

	case FieldKindFloat64:
		s.WriteFloat64(math.Float64frombits(uint64(f.IValue)))

	case FieldKindString:
		s.WriteString(f.SValue)
	}

	s.WriteObjectEnd()
}

//
func (je *JSONEncoder) encodeStackFrame(frame sys.StackFrame, s *jsoniter.Stream) {

	s.WriteObjectStart()

	s.WriteObjectField("func")
	s.WriteString(frame.Function)
	s.WriteMore()

	s.WriteObjectField("file")
	s.WriteString(fmt.Sprintf("%s:%d", frame.File, frame.Line))

	s.WriteObjectEnd()
}
