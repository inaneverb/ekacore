// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/qioalice/ekago/v3/ekasys"
	"github.com/qioalice/ekago/v3/internal/ekaletter"

	"github.com/json-iterator/go"
)

var (
	// Make sure we won't break API by declaring package's console encoder
	defaultJSONEncoder CI_Encoder
)

// doBuild builds the current CI_JSONEncoder only if it has not built yet.
// There is no-op if encoder already built.
func (je *CI_JSONEncoder) doBuild() *CI_JSONEncoder {

	switch {
	case je == nil:
		return nil

	case je.api != nil:
		// do not build if it's so already
		return je
	}

	je.api = jsoniter.Config{
		IndentionStep:                 je.indent,
		MarshalFloatWith6Digits:       true,
		ObjectFieldMustBeSimpleString: true,
	}.Froze()

	preEncodedFieldsApi := jsoniter.Config{
		IndentionStep:                 je.indent * 2,
		MarshalFloatWith6Digits:       true,
		ObjectFieldMustBeSimpleString: true,
	}.Froze()

	je.preEncodedFieldsStreamIndentX1 = je.api.BorrowStream(nil)
	je.preEncodedFieldsStreamIndentX2 = preEncodedFieldsApi.BorrowStream(nil)

	return je
}

// encodeBase encodes Entry's level, timestamp, message to s.
func (je *CI_JSONEncoder) encodeBase(s *jsoniter.Stream, e *Entry) {

	s.WriteObjectField("level")
	s.WriteString(e.Level.String())
	s.WriteMore()

	s.WriteObjectField("level_value")
	s.WriteUint8(uint8(e.Level))
	s.WriteMore()

	s.WriteObjectField("time")
	s.WriteString(e.Time.Format(time.RFC3339))

	s.WriteMore()
	s.WriteObjectField("message")
	s.WriteString(e.LogLetter.Messages[0].Body)

	if e.ErrLetter != nil {
		s.WriteMore()
		je.encodeErrorHeader(s, e.ErrLetter)
	}
}

// encodeErrorHeader writes ekaerr.Error's header object treating provided
// ekaletter.Letter as ekaerr.Error's one.
//
// It won't encode stacktrace, neither its messages nor fields.
// encodeStackFrame() does that.
func (je *CI_JSONEncoder) encodeErrorHeader(s *jsoniter.Stream, errLetter *ekaletter.Letter) {

	for i, n := 0, len(errLetter.SystemFields); i < n; i++ {
		switch errLetter.SystemFields[i].BaseType() {

		case ekaletter.KIND_SYS_TYPE_EKAERR_UUID:
			s.WriteObjectField("error_id")
			s.WriteString(errLetter.SystemFields[i].SValue)

		case ekaletter.KIND_SYS_TYPE_EKAERR_CLASS_ID:
			s.WriteObjectField("error_class_id")
			s.WriteInt64(errLetter.SystemFields[i].IValue)

		case ekaletter.KIND_SYS_TYPE_EKAERR_CLASS_NAME:
			s.WriteObjectField("error_class_name")
			s.WriteString(errLetter.SystemFields[i].SValue)

		default:
			continue
		}

		if i < n-1 {
			s.WriteMore()
		}
	}

	to := s.Buffer()
	if l := len(to); to[l-1] == ',' {
		s.SetBuffer(to[:l-1])
	}
}

func (je *CI_JSONEncoder) encodeStacktrace(s *jsoniter.Stream, e *Entry) (wasAdded bool) {

	stacktrace := e.LogLetter.StackTrace
	if len(stacktrace) == 0 && e.ErrLetter != nil {
		stacktrace = e.ErrLetter.StackTrace
	}

	n := int16(len(stacktrace))
	if n == 0 {
		return false
	}

	var (
		fi       = 0 // fi for fields' index
		mi       = 0 // mi for messages' index
		fields   []ekaletter.LetterField
		messages []ekaletter.LetterMessage
	)

	if e.ErrLetter != nil {
		fields = e.ErrLetter.Fields
		messages = e.ErrLetter.Messages
	}

	s.WriteObjectField("stacktrace")
	s.WriteArrayStart()

	for i := int16(0); i < n; i++ {
		messageForStackFrame := ekaletter.LetterMessage{}
		fieldsForStackFrame := []ekaletter.LetterField(nil)
		fiEnd := 0

		//goland:noinspection GoNilness
		if mi < len(messages) && messages[mi].StackFrameIdx == i {
			messageForStackFrame = messages[mi]
			mi++
		}

		if fi < len(fields) && fields[fi].StackFrameIdx == i {
			fiEnd = fi + 1
			for fiEnd < len(fields) && fields[fiEnd].StackFrameIdx == i {
				fiEnd++
			}
		}

		if fiEnd != 0 {
			fieldsForStackFrame = fields[fi:fiEnd]
		}

		je.encodeStackFrame(s, stacktrace[i], fieldsForStackFrame, messageForStackFrame)

		if i < n-1 {
			s.WriteMore()
		}
	}

	s.WriteArrayEnd()
	return true
}

func (je *CI_JSONEncoder) encodeStackFrame(

	s *jsoniter.Stream,
	frame ekasys.StackFrame,
	fields []ekaletter.LetterField,
	message ekaletter.LetterMessage,

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

	if message.Body != "" {
		s.WriteMore()
		s.WriteObjectField("message")
		s.WriteString(message.Body)
	}

	if len(fields) > 0 {
		s.WriteMore()
		if wasAdded := je.encodeFields(s, fields, nil, false); !wasAdded {
			b := s.Buffer()
			s.SetBuffer(b[:len(b)-1])
		}
	}

	s.WriteObjectEnd()
}

func (je *CI_JSONEncoder) encodeFields(s *jsoniter.Stream, fs, addFs []ekaletter.LetterField, addPreEncoded bool) (wasAdded bool) {

	if len(fs) == 0 && len(addFs) == 0 {
		return false
	}

	var (
		unnamedFieldIdx, writtenFields int16
	)

	s.WriteObjectField("fields")
	s.WriteObjectStart()

	addField := func(s *jsoniter.Stream, f *ekaletter.LetterField, unnamedFieldIdx, writtenFields *int16) {
		if f.IsZero() || f.Kind.IsInvalid() {
			return
		}
		keyBak := f.Key
		if f.Key == "" && !f.IsSystem() {
			f.Key = f.KeyOrUnnamed(unnamedFieldIdx)
		}
		if wasAdded = je.encodeField(s, *f); wasAdded {
			s.WriteMore()
			*writtenFields++
		}
		f.Key = keyBak
	}

	for i, n := int16(0), int16(len(fs)); i < n; i++ {
		addField(s, &fs[i], &unnamedFieldIdx, &writtenFields)
	}
	for i, n := int16(0), int16(len(addFs)); i < n; i++ {
		addField(s, &addFs[i], &unnamedFieldIdx, &writtenFields)
	}

	to := s.Buffer()

	// Write pre-encoded fields in "fields" section
	if addPreEncoded {
		to = bufw2(to, je.preEncodedFieldsStreamIndentX2.Buffer())
	}

	i := len(to) - 1

	// Remove last comma.
	if to[i] == ',' {
		i--
	}

	// Maybe no fields were added?
	if writtenFields == 0 {
		for i >= 0 && to[i] != 'f' { // start of "fields"
			i--
		}
	}

	s.SetBuffer(to[:i+1])

	s.WriteObjectEnd()
	return true
}

func (je *CI_JSONEncoder) encodeField(s *jsoniter.Stream, f ekaletter.LetterField) (wasAdded bool) {

	// Maybe field must be skipped? Field should be skipped if it's vary field
	// (the name has '?' at the end and the value is zero).
	// Also fields that is started from "sys." is skipped.
	if f.Kind.IsInvalid() || strings.HasPrefix(f.Key, "sys.") {
		return false
	}

	s.WriteObjectField(f.Key)
	je.encodeFieldValue(s, f)

	return true
}

func (je *CI_JSONEncoder) encodeFieldValue(s *jsoniter.Stream, f ekaletter.LetterField) {

	if f.Kind.IsSystem() {
		switch f.Kind.BaseType() {

		case ekaletter.KIND_SYS_TYPE_EKAERR_UUID, ekaletter.KIND_SYS_TYPE_EKAERR_CLASS_NAME:
			s.WriteString(f.SValue)

		case ekaletter.KIND_SYS_TYPE_EKAERR_CLASS_ID:
			b := s.Buffer()
			b = strconv.AppendInt(b, f.IValue, 10)
			s.SetBuffer(b)

		default:
			s.WriteString("<unsupported system field>")
		}

	} else if f.Kind.IsNil() {
		s.WriteNil()

	} else if f.Kind.IsInvalid() {
		s.WriteString("<invalid_field>")

	} else {
		switch f.Kind.BaseType() {

		case ekaletter.KIND_TYPE_BOOL:
			b := s.Buffer()
			b = strconv.AppendBool(b, f.IValue != 0)
			s.SetBuffer(b)

		case ekaletter.KIND_TYPE_INT,
			ekaletter.KIND_TYPE_INT_8, ekaletter.KIND_TYPE_INT_16,
			ekaletter.KIND_TYPE_INT_32, ekaletter.KIND_TYPE_INT_64:
			b := s.Buffer()
			b = strconv.AppendInt(b, f.IValue, 10)
			s.SetBuffer(b)

		case ekaletter.KIND_TYPE_UINT,
			ekaletter.KIND_TYPE_UINT_8, ekaletter.KIND_TYPE_UINT_16,
			ekaletter.KIND_TYPE_UINT_32, ekaletter.KIND_TYPE_UINT_64:
			b := s.Buffer()
			b = strconv.AppendUint(b, uint64(f.IValue), 10)
			s.SetBuffer(b)

		case ekaletter.KIND_TYPE_FLOAT_32:
			b := s.Buffer()
			f := float64(math.Float32frombits(uint32(f.IValue)))
			b = strconv.AppendFloat(b, f, 'f', 2, 32)
			s.SetBuffer(b)

		case ekaletter.KIND_TYPE_FLOAT_64:
			b := s.Buffer()
			f := math.Float64frombits(uint64(f.IValue))
			b = strconv.AppendFloat(b, f, 'f', 2, 64)
			s.SetBuffer(b)

		case ekaletter.KIND_TYPE_UINTPTR, ekaletter.KIND_TYPE_ADDR:
			b := s.Buffer()
			b = bufw(b, "0x")
			b = strconv.AppendInt(b, f.IValue, 16)
			s.SetBuffer(b)

		case ekaletter.KIND_TYPE_STRING:
			s.WriteString(f.SValue)

		case ekaletter.KIND_TYPE_COMPLEX_64:
			b := s.Buffer()
			r := math.Float32frombits(uint32(f.IValue >> 32))
			i := math.Float32frombits(uint32(f.IValue))
			c := complex128(complex(r, i))
			// TODO: Use strconv.AppendComplex() when it will be released.
			b = bufw(b, strconv.FormatComplex(c, 'f', 2, 32))
			s.SetBuffer(b)

		case ekaletter.KIND_TYPE_COMPLEX_128:
			b := s.Buffer()
			c := f.Value.(complex128)
			// TODO: Use strconv.AppendComplex() when it will be released.
			b = bufw(b, strconv.FormatComplex(c, 'f', 2, 64))
			s.SetBuffer(b)

		case ekaletter.KIND_TYPE_UNIX:
			s.WriteString(time.Unix(f.IValue, 0).Format("Jan 2 15:04:05"))

		case ekaletter.KIND_TYPE_UNIX_NANO:
			s.WriteString(time.Unix(0, f.IValue).Format("Jan 2 15:04:05.000000000"))

		case ekaletter.KIND_TYPE_DURATION:
			s.WriteString(time.Duration(f.IValue).String())

		case ekaletter.KIND_TYPE_MAP, ekaletter.KIND_TYPE_EXTMAP,
			ekaletter.KIND_TYPE_STRUCT, ekaletter.KIND_TYPE_ARRAY:
			// TODO: Add support of extracted maps.
			s.WriteVal(f.Value)

		default:
			s.WriteString("<unsupported_field>")
		}
	}
}
