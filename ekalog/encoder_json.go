// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"github.com/qioalice/ekago/v2/internal/ekaletter"

	"github.com/json-iterator/go"
)

//noinspection GoSnakeCaseUsage
type (
	// CI_JSONEncoder is a type that built to be used as a part of CommonIntegrator
	// as an log Entry encoder to the some output as JSON.
	// Custom indentation supported.
	//
	// If you want to use CI_JSONEncoder, you need to instantiate object,
	// set indentation (if you need, default is 0: no indentation) and that is.
	// The last thing you need to do is to register CI_ConsoleEncoder with
	// CommonIntegrator using CommonIntegrator.WithEncoder().
	//
	// You MUST NOT to call EncodeEntry() method manually.
	// It is used by associated CommonIntegrator and it WILL lead to UB
	// if you will try to use it manually. May even panic.
	CI_JSONEncoder struct {

		// You know what "JSON indent" is (pretty JSON, etc), right?
		// How much spaces will be added to the beginning of line for each
		// nested JSON entity for each nesting level.
		//
		// So, for indent == 4, you will get:
		//
		// 		{
		// 		    "key1": "value",
		// 		    "nested": {
		// 		        "nested_key": "value"
		// 		    }
		// 		}
		//
		// So, keys for 1st nesting level JSON entities has 4 spaces before
		// ("key1", "nested") and for 2nd nesting level - 8 spaces before
		// ("nested_key").
		//
		// You may set this value using SetIndent() method.
		indent int

		// api is jsoniter's API object.
		// Created at the first doBuild() call for object.
		api jsoniter.API

		preEncodedFieldsStreamIndentX1 *jsoniter.Stream
		preEncodedFieldsStreamIndentX2 *jsoniter.Stream
	}
)

// SetIndent sets an indentation of JSON encoding format.
// Any value <= 0 meaning "no indentation".
//
// Calling this method many times will overwrite previous value of format string.
//
// An indentation MUST NOT be changed after CI_JSONEncoder is registered
// with CommonIntegrator using CommonIntegrator.WithEncoder() method.
func (je *CI_JSONEncoder) SetIndent(num int) *CI_JSONEncoder {

	if num < 0 {
		num = 0
	}
	je.indent = num
	return je
}

// PreEncodeField allows you to pre-encode some ekaletter.LetterField,
// that is must be used with EACH Entry that will be encoded using this CI_JSONEncoder.
//
// It's useful when you want some unchanged runtime data for each log message,
// like git hash commit, version, etc. Or if you want to create many loggers
// attach some different fields to them and log different logs using them.
//
// Unnamed fields are not allowed.
//
// By default, encoded field will be added to the "fields" root section.
// If you want to place field directly to the root section,
// use ekaletter.KIND_FLAG_USER_DEFINED for ekaletter.LetterField's Kind property
// (set it).
//
// WARNING!
// PreEncodeField() MUST BE USED ONLY IF CI_JSONEncoder HAS BEEN REGISTERED
// WITH SOME CommonIntegrator ALREADY. UB OTHERWISE, MAY PANIC!
func (je *CI_JSONEncoder) PreEncodeField(f ekaletter.LetterField) {

	// Avoid calls of PreEncodeField() when CI_ConsoleEncoder has not built yet.
	if f.Key == "" || je.api == nil || f.IsInvalid() || f.RemoveVary() && f.IsZero() {
		return
	}

	stream := je.preEncodedFieldsStreamIndentX2
	if f.Kind & ekaletter.KIND_FLAG_USER_DEFINED != 0 {
		stream = je.preEncodedFieldsStreamIndentX1
	}

	if wasAdded := je.encodeField(stream, f); wasAdded {
		stream.WriteMore()
	}
}

// EncodeEntry encodes passed Entry in JSON format using provided indentation.
//
// EncodeEntry is for internal purposes only and MUST NOT be called directly.
// UB otherwise, may panic.
func (je *CI_JSONEncoder) EncodeEntry(e *Entry) []byte {

	s := je.api.BorrowStream(nil)
	defer je.api.ReturnStream(s)

	// Use last ekaerr.Error's message as Entry's one if it's empty.
	if e.ErrLetter != nil {
		if l := len(e.ErrLetter.Messages); l > 0 && e.LogLetter.Messages[0].Body == "" {
			e.LogLetter.Messages[0].Body = e.ErrLetter.Messages[l-1].Body
			e.ErrLetter.Messages[l-1].Body = ""
		}
	}

	s.WriteObjectStart()

	je.encodeBase(s, e)
	s.WriteMore()

	// Write pre-encoded fields in root section
	if b := je.preEncodedFieldsStreamIndentX1.Buffer(); len(b) > 0 {
		s.SetBuffer(bufw2(s.Buffer(), b))
		//s.WriteMore() // unnecessary, WriteMore() already called for field stream
	}

	// Handle special case when ekaerr.Error's ekaletter.Letter has a fields
	// but has no stacktrace. It means that lightweight error has been created.
	lightweightErrorFields := []ekaletter.LetterField(nil)
	if e.ErrLetter != nil && len(e.ErrLetter.StackTrace) == 0 && len(e.ErrLetter.Fields) > 0 {
		lightweightErrorFields = e.ErrLetter.Fields
	}

	if wasAdded := je.encodeFields(s, e.LogLetter.Fields, lightweightErrorFields, true); wasAdded {
		s.WriteMore()
	}

	if wasAdded := je.encodeStacktrace(s, e); wasAdded {
		s.WriteMore()
	}

	// ------------ Add new sections here ------------ //

	// We writing the JSON's comma at the each section, expecting that the next
	// section will be written too. But it might be an empty.
	// So, we need to remove the last comma. There is no more sections to be written.
	b := s.Buffer()
	s.SetBuffer(b[:len(b)-1])

	s.WriteObjectEnd()

	b = s.Buffer()
	copied := make([]byte, len(b) +1)
	copy(copied, b)

	copied[len(copied)-1] = '\n'

	// Restore ekaerr.Error's last message that was used as Entry's message.
	if e.ErrLetter != nil {
		if l := len(e.ErrLetter.Messages); l > 0 && e.ErrLetter.Messages[l-1].Body == "" {
			e.ErrLetter.Messages[l-1].Body = e.LogLetter.Messages[0].Body
			e.LogLetter.Messages[0].Body = ""
		}
	}

	return copied
}
