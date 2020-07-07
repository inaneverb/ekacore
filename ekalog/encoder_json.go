// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"math"
	"time"
	"unsafe"

	"github.com/qioalice/ekago/ekadanger"
	"github.com/qioalice/ekago/ekasys"
	"github.com/qioalice/ekago/internal/field"
	"github.com/qioalice/ekago/internal/letter"

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
	jsonEncoderAddr = ekadanger.TakeRealAddr(jsonEncoder)
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

	allowEmpty := e.LogLetter.Items.Flags.TestAll(FLAG_INTEGRATOR_IGNORE_EMPTY_PARTS)
	continueWriting := false

	s.WriteObjectStart()

	je.encodeBase(s, e, allowEmpty)
	s.WriteMore()

	continueWriting = je.encodeFields(s, e.LogLetter.Items.Fields, allowEmpty)
	if continueWriting {
		s.WriteMore()
	}

	je.encodeStacktrace(s, e, allowEmpty)

	s.WriteObjectEnd()

	buf := s.Buffer()
	copied := make([]byte, len(buf))
	copy(copied, buf)

	return copied
}

// encodeBase encodes e's level, timestamp, message to s.
func (je *JSONEncoder) encodeBase(s *jsoniter.Stream, e *Entry, allowEmpty bool) {

	s.WriteObjectField("level")
	s.WriteString(e.Level.String())
	s.WriteMore()

	s.WriteObjectField("level_value")
	s.WriteUint8(uint8(e.Level))
	s.WriteMore()

	s.WriteObjectField("time")
	s.WriteString(e.Time.Format(time.UnixDate))

	if e.ErrLetter != nil {
		s.WriteMore()
		je.encodeError(s, e.ErrLetter, allowEmpty)
	}

	if len(e.LogLetter.Items.Message) > 0 || allowEmpty {

		s.WriteMore()
		s.WriteObjectField("message")
		s.WriteString(e.LogLetter.Items.Message)
	}
}

//
func (je *JSONEncoder) encodeError(s *jsoniter.Stream, errLetter *letter.Letter, allowEmpty bool) {

	for i, n := 0, len(errLetter.SystemFields); i < n; i++ {
		switch errLetter.SystemFields[i].BaseType() {

		case field.KIND_SYS_TYPE_EKAERR_UUID:
			s.WriteObjectField("error_id")
			s.WriteString(errLetter.SystemFields[i].SValue)

		case field.KIND_SYS_TYPE_EKAERR_CLASS_ID:
			s.WriteObjectField("error_class_id")
			s.WriteInt64(errLetter.SystemFields[i].IValue)

		case field.KIND_SYS_TYPE_EKAERR_CLASS_NAME:
			s.WriteObjectField("error_class_name")
			s.WriteString(errLetter.SystemFields[i].SValue)

		case field.KIND_SYS_TYPE_EKAERR_PUBLIC_MESSAGE:
			if publicMessage := errLetter.SystemFields[i].SValue; len(publicMessage) > 0 || allowEmpty {
				s.WriteObjectField("error_public_message")
				s.WriteString(publicMessage)
				if i < n-1 {
					s.WriteMore()
				}
			}
			continue
		}

		if i < n-1 {
			s.WriteMore()
		}
	}
}

//
func (je *JSONEncoder) encodeFields(

	s *jsoniter.Stream,
	fields []field.Field,
	allowEmpty bool,

) (wasAdded bool) {

	unnamedFieldIdx := 1

	lFields := len(fields)

	if lFields == 0 && !allowEmpty {
		return false
	}

	s.WriteObjectField("fields")

	if lFields > 0 {
		s.WriteArrayStart()

		for _, field_ := range fields[:lFields-1] {
			je.encodeField(s, field_, &unnamedFieldIdx)
			s.WriteMore()
		}
		je.encodeField(s, fields[lFields-1], &unnamedFieldIdx)

		s.WriteArrayEnd()

	} else {
		s.WriteEmptyArray()
	}

	return true
}

//
func (je *JSONEncoder) encodeStacktrace(s *jsoniter.Stream, e *Entry, allowEmpty bool) {

	stacktrace := e.LogLetter.StackTrace
	if len(stacktrace) == 0 && e.ErrLetter != nil {
		stacktrace = e.ErrLetter.StackTrace
	}

	lStacktrace := int16(len(stacktrace))

	if lStacktrace == 0 && !allowEmpty {
		return
	}

	s.WriteObjectField("stacktrace")

	letterItem := (*letter.LetterItem)(nil)
	letterItemIdx := int16(0)
	if e.ErrLetter != nil {
		letterItem = e.ErrLetter.Items
		letterItemIdx = letterItem.StackFrameIdx()
	}

	if lStacktrace > 0 {
		s.WriteArrayStart()

		for i := int16(0); i < lStacktrace; i++ {

			letterItemPassed := (*letter.LetterItem)(nil)
			if letterItem != nil && letterItemIdx == i {
				letterItemPassed = letterItem
				letterItem = letterItem.Next()
				letterItemIdx = letterItem.StackFrameIdx()
			}
			je.encodeStackFrame(s, stacktrace[i], letterItemPassed, allowEmpty)
			if i < lStacktrace-1 {
				s.WriteMore()
			}
		}

		s.WriteArrayEnd()

	} else {
		s.WriteEmptyArray()
	}
}

//
func (je *JSONEncoder) encodeField(s *jsoniter.Stream, f field.Field, unnamedFieldIdx *int) {

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

//
func (je *JSONEncoder) encodeStackFrame(

	s *jsoniter.Stream,
	frame ekasys.StackFrame,
	letterItem *letter.LetterItem,
	allowEmpty bool,
) {
	frame.DoFormat()

	s.WriteObjectStart()

	s.WriteObjectField("func")
	s.WriteString(frame.Format[:frame.FormatFileOffset-1])
	s.WriteMore()

	s.WriteObjectField("file")
	s.WriteString(frame.Format[frame.FormatFileOffset+1 : frame.FormatFullPathOffset-2])
	s.WriteMore()

	s.WriteObjectField("package")
	s.WriteString(frame.Format[frame.FormatFullPathOffset:])

	if letterItem != nil {
		if len(letterItem.Message) > 0 || allowEmpty {
			s.WriteMore()
			s.WriteObjectField("message")
			s.WriteString(letterItem.Message)
		}
		if len(letterItem.Fields) > 0 || allowEmpty {
			s.WriteMore()
			je.encodeFields(s, letterItem.Fields, allowEmpty)
		}
	}

	s.WriteObjectEnd()
}
