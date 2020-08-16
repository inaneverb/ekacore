// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog_encoder_datadog

import (
	"path/filepath"
	"strings"

	"github.com/qioalice/ekago/v2/ekaexp"
	"github.com/qioalice/ekago/v2/ekalog"
	"github.com/qioalice/ekago/v2/internal/field"
	"github.com/qioalice/ekago/v2/internal/letter"

	"github.com/json-iterator/go"
)

// doBuild builds the current CI_DatadogEncoder only if it has not built yet.
// There is no-op if encoder already built.
func (de *CI_DatadogEncoder) doBuild() *CI_DatadogEncoder {

	switch {
	case de == nil:
		return nil

	case de.api != nil:
		// do not build if it's so already
		return de
	}

	de.api = jsoniter.Config{
		IndentionStep:                 0,
		MarshalFloatWith6Digits:       true,
		ObjectFieldMustBeSimpleString: true,
	}.Froze()

	return de
}

// encode is an entry point of encoding log's entry to the JSON for DataDog service.
// It's the function that is returned by FreezeAndGetEncoder() method.
func (de *CI_DatadogEncoder) encode(e *ekalog.Entry) []byte {

	s := de.api.BorrowStream(nil)
	defer de.api.ReturnStream(s)

	var (
		allowEmpty = e.LogLetter.Items.Flags.TestAll(ekalog.FLAG_INTEGRATOR_IGNORE_EMPTY_PARTS)
		wasAdded bool
	)

	s.WriteObjectStart()

	de.encodeBase(s, e, allowEmpty)
	s.WriteMore()

	wasAdded =
		de.encodeSysPrefixedFields(s, e.LogLetter.Items.Fields)
	if wasAdded {
		s.WriteMore()
	}

	wasAdded =
		de.encodeStacktrace(s, e, allowEmpty)
	if wasAdded {
		s.WriteMore()
	}

	wasAdded =
		de.encodeFields(s, e.LogLetter.Items.Fields, e.ErrLetter, allowEmpty)
	if wasAdded {
		s.WriteMore()
	}

	// ------------ Add new sections here ------------ //

	b := s.Buffer()
	s.SetBuffer(b[:len(b)-1])

	s.WriteObjectEnd()

	b = s.Buffer()
	copied := make([]byte, len(b))
	copy(copied, b)

	return copied
}

// encodeBase encodes as JSON:
//   - e's level (if no custom level registered, writes only string level representation,
//     otherwise - an integer level value too),
//   - e's timestamp (as ISO8601),
//   - e's message (skipping if empty and 'allowEmpty' is false),
//   - attached error's data (error's ID, all error's fields).
//
// Puts JSON encoded data into 's' stream, skipping empty log's message
// and error's public message if 'allowEmpty' is false.
func (de *CI_DatadogEncoder) encodeBase(

	s *jsoniter.Stream,
	e *ekalog.Entry,
	allowEmpty bool,
) {

	const (
		// Datadog accepts ISO8601, RFC3164 or UNIX (the milliseconds EPOCH format).
		// So, here ISO8601 format, cause std Golang time package doesn't have that.
		ISO8601 = "2006-01-02T15:04:05.000-0700"
	)

	s.WriteObjectField("level")
	s.WriteString(e.Level.String())
	s.WriteMore()

	if ekalog.RegisteredCustomLevels() > 0 {
		s.WriteObjectField("level_value")
		s.WriteUint8(uint8(e.Level))
		s.WriteMore()
	}

	// Datadog supports following time names:
	// "@timestamp", "timestamp", "_timestamp", "Timestamp", "eventTime", "date",
	// "published_date", "syslog.timestamp".
	//
	// But:
	// Datadog rejects a log entry if its official date is older than 6 hours in the past.
	//
	// So, we don't want to handle log entry's time as a log publishing time
	// because there may be a situation when log writer will log the log entry
	// after 6h when it was occurred, and it will be rejected from Datadog.
	//
	// Thus, we have to use a different key.
	s.WriteObjectField("timestamp_real")
	s.WriteString(e.Time.Format(ISO8601))

	if e.ErrLetter != nil {
		s.WriteMore()
		de.encodeError(s, e.ErrLetter)
	}

	message := e.LogLetter.Items.Message
	if message == "" && e.ErrLetter != nil {
		for letterItem := e.ErrLetter.Items; letterItem != nil && message == ""; {
			if letterItem.Message != "" {
				message = letterItem.Message
			}
			letterItem = letterItem.Next()
		}
	}

	if message != "" || allowEmpty {
		s.WriteMore()
		s.WriteObjectField("message")
		s.WriteString(message)
	}
}

// encodeErrorHeader encodes as JSON e's error's ID (if attached error is presented),
// and all error's fields.
//
// Puts JSON encoded data into 's' stream.
func (de *CI_DatadogEncoder) encodeError(

	s *jsoniter.Stream,
	errLetter *letter.Letter,
) {
	for i, n := 0, len(errLetter.SystemFields); i < n; i++ {
		if errLetter.SystemFields[i].BaseType() == field.KIND_SYS_TYPE_EKAERR_UUID {
			s.WriteObjectField("error_id")
			s.WriteString(errLetter.SystemFields[i].SValue)
			return
		}
	}
}

// encodeStacktrace encodes as JSON a presented stacktrace:
// either log entry's or attached error's ones. Also encodes attached error's
// stack frame messages (if presented).
//
// Puts JSON encoded data into 's' stream, skipping adding stacktrace
// if it's an empty and 'allowEmpty' is false.
func (de *CI_DatadogEncoder) encodeStacktrace(

	s *jsoniter.Stream,
	e *ekalog.Entry,
	allowEmpty bool,

) (wasAdded bool) {

	stacktrace := e.LogLetter.StackTrace
	if len(stacktrace) == 0 && e.ErrLetter != nil {
		stacktrace = e.ErrLetter.StackTrace
	}

	lStacktrace := int16(len(stacktrace))
	if lStacktrace == 0 && !allowEmpty {
		return false
	}

	// write stacktrace as JSON array of stack frames formatted
	// in the following format: "(<stackIdx>) <stackFunc>(<file>:<line>)"

	s.WriteObjectField("stack_trace")

	if lStacktrace == 0 {
		s.WriteEmptyArray()
		s.WriteMore()
		s.WriteObjectField("stack_messages")
		s.WriteEmptyArray()
		return true
	}

	s.WriteArrayStart()

	for i := int16(0); i < lStacktrace; i++ {
		s.WriteRaw("\"(")
		de.encodeNumWithZeroes(s, i, lStacktrace)
		s.WriteRaw("): ")
		s.WriteRaw(stacktrace[i].Function)
		s.WriteRaw("(")
		s.WriteRaw(filepath.Base(stacktrace[i].File))
		s.WriteRaw(":")
		s.WriteInt(stacktrace[i].Line)
		s.WriteRaw(")\"")
		s.WriteMore()
	}

	b := s.Buffer()
	s.SetBuffer(b[:len(b)-1]) // remove last comma

	s.WriteArrayEnd()
	s.WriteMore()

	// write stacktrace's messages as JSON array in the following format:
	// "(<stackIdx>) <message>".

	s.WriteObjectField("stack_messages")
	s.WriteArrayStart()

	letterItem := (*letter.LetterItem)(nil)
	if e.ErrLetter != nil {
		letterItem = e.ErrLetter.Items
	}

	nonEmptyStackMessages := false

	for ; letterItem != nil; letterItem = letterItem.Next() {
		if letterItem.Message != "" {
			nonEmptyStackMessages = true

			s.WriteRaw("\"(")
			de.encodeNumWithZeroes(s, letterItem.StackFrameIdx(), lStacktrace)
			s.WriteRaw("): ")
			s.WriteRaw(letterItem.Message)
			s.WriteRaw("\"")
			s.WriteMore()
		}
	}

	if nonEmptyStackMessages {
		b := s.Buffer()
		s.SetBuffer(b[:len(b)-1]) // remove last comma
	}

	s.WriteArrayEnd()
	return true
}

// encodeFields encodes as JSON 'logFields' (assuming that this is log entry's
// fields) and attached error's ones (assuming that 'errLetter' is error's
// internal parts).
// Skips those fields from 'logFields' which names started from "sys.".
//
// Puts JSON encoded data into 's' stream, skipping adding fields
// if there is no fields (log's and error's) and 'allowEmpty' is false.
func (_ *CI_DatadogEncoder) encodeFields(

	s *jsoniter.Stream,
	logFields []ekaexp.Field,
	errLetter *letter.Letter,
	allowEmpty bool,

) (wasAdded bool) {

	// TODO: Refactor for O(2N) -> O(N).

	emptySet := true
	emptySet = emptySet && len(logFields) == 0
	if errLetter != nil {
		for letterItem := errLetter.Items; letterItem != nil && emptySet; {
			emptySet = emptySet && len(letterItem.Fields) == 0
			letterItem = letterItem.Next()
		}
	}

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

	for i, n := int16(0), int16(len(logFields)); i < n; i++ {
		// all "sys." prefixed fields will be processed at the
		// encodeSysPrefixedFields() method.
		if !strings.HasPrefix(logFields[i].Key, "sys.") {
			s.WriteObjectField(logFields[i].KeyOrUnnamed(&unnamedFieldIdx))
			_, _ = logFields[i].WriteTo(s)
			s.WriteMore()
		}
	}

	if errLetter != nil {
		for letterItem := errLetter.Items; letterItem != nil; {
			for i, n := int16(0), int16(len(letterItem.Fields)); i < n; i++ {
				s.WriteObjectField(logFields[i].KeyOrUnnamed(&unnamedFieldIdx))
				_, _ = logFields[i].WriteTo(s)
				s.WriteMore()
			}
			letterItem = letterItem.Next()
		}
	}

	b := s.Buffer()
	if b[len(b)-1] == ',' {
		// emptySet might be false, but no one fields are written
		// (all of them are "sys." prefixed and skipped
		s.SetBuffer(b[:len(b)-1])
	}

	s.WriteArrayEnd()
	return true
}

// encodeSysPrefixedFields encodes as JSON only those fields from
// 'fieldsMayContainSysPrefixed' that has prefix "sys.dd.", removing it
// and adding as is to the JSON root.
//
// Puts JSON encoded data into 's' stream, skipping adding fields
// if there is no fields (log's and error's) and 'allowEmpty' is false.
func (de *CI_DatadogEncoder) encodeSysPrefixedFields(

	s *jsoniter.Stream,
	fieldsMayContainSysPrefixed []field.Field,

) (wasAdded bool) {
	fields := fieldsMayContainSysPrefixed // just an alias, nothing more

	if len(fields) == 0 {
		return false
	}

	for i, n := int16(0), int16(len(fields)); i < n; i++ {

		needToWriteField := len(fields[i].Key) > 7 &&
			strings.HasPrefix(fields[i].Key, "sys.dd.") &&
			fields[i].SValue != ""

		if needToWriteField {
			wasAdded = true

			s.WriteObjectField(fields[i].Key[7:])
			s.WriteString(fields[i].SValue)

			s.WriteMore()
		}
	}

	if wasAdded {
		b := s.Buffer()
		s.SetBuffer(b[:len(b)-1]) // remove last comma
	}

	return wasAdded
}

// encodeNumWithZeroes encodes as JSON 'num' with a few leading zeroes,
// depending of 'maxNum'.
//
// Puts JSON encoded data into 's' stream.
//
// So, let's assume we have a "buffer" with the width that exactly enough
// to place a 'maxNum',
// and we need to place 'num at right and fills the free space by zeroes.
//
// E.g: 'num' == 26, 'maxNum' = 12345.
// Then this method will add a string "00026" to the 's'.
func (_ *CI_DatadogEncoder) encodeNumWithZeroes(s *jsoniter.Stream, num, maxNum int16) {

	var (
		needSections = int16(0)
		maxSections = int16(0)
	)

	switch {
	case num < 10:    needSections = 1
	case num < 100:   needSections = 2
	case num < 1000:  needSections = 3
	case num < 10000: needSections = 4
	default:          needSections = 5
	}

	switch {
	case maxNum < 10:    maxSections = 1
	case maxNum < 100:   maxSections = 2
	case maxNum < 1000:  maxSections = 3
	case maxNum < 10000: maxSections = 4
	default:             maxSections = 5
	}

	switch maxSections - needSections {
	case 1: s.WriteRaw("0")
	case 2: s.WriteRaw("00")
	case 3: s.WriteRaw("000")
	case 4: s.WriteRaw("0000")
	}

	s.WriteInt16(num)
}