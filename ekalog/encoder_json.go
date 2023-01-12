// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"time"

	"github.com/json-iterator/go"

	"github.com/qioalice/ekago/v4/internal/ekaletter"
)

// noinspection GoSnakeCaseUsage
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

		oneDepthLevel bool

		// api is jsoniter's API object.
		// Created at the first doBuild() call for object.
		api jsoniter.API

		preEncodedFieldsStreamIndentX1 *jsoniter.Stream
		preEncodedFieldsStreamIndentX2 *jsoniter.Stream

		fieldNames map[CI_JSONEncoder_Field]string

		timeFormatter func(t time.Time) string
	}

	// CI_JSONEncoder_Field is a special type that represents a type of CI_JSONEncoder
	// field names.
	// This type exist to declare corresponding constants and be able to change
	// default field's names to their user-defined alternatives.
	CI_JSONEncoder_Field uint8
)

// noinspection GoSnakeCaseUsage
const (
	CI_JSON_ENCODER_FIELD_LEVEL CI_JSONEncoder_Field = 1 + iota
	CI_JSON_ENCODER_FIELD_LEVEL_VALUE
	CI_JSON_ENCODER_FIELD_TIME
	CI_JSON_ENCODER_FIELD_MESSAGE
	CI_JSON_ENCODER_FIELD_ERROR_ID
	CI_JSON_ENCODER_FIELD_ERROR_CLASS_ID
	CI_JSON_ENCODER_FIELD_ERROR_CLASS_NAME
	CI_JSON_ENCODER_FIELD_STACKTRACE
	CI_JSON_ENCODER_FIELD_1DL_STACKTRACE_MESSAGES
	CI_JSON_ENCODER_FIELD_FIELDS
	CI_JSON_ENCODER_FIELD_1DL_LOG_FIELDS_PREFIX
	CI_JSON_ENCODER_FIELD_1DL_STACKTRACE_FIELDS_PREFIX
)

// noinspection GoSnakeCaseUsage
const (
	CI_JSON_ENCODER_FIELD_DEFAULT_LEVEL                        = "level"
	CI_JSON_ENCODER_FIELD_DEFAULT_LEVEL_VALUE                  = "level_value"
	CI_JSON_ENCODER_FIELD_DEFAULT_TIME                         = "time"
	CI_JSON_ENCODER_FIELD_DEFAULT_MESSAGE                      = "message"
	CI_JSON_ENCODER_FIELD_DEFAULT_ERROR_ID                     = "error_id"
	CI_JSON_ENCODER_FIELD_DEFAULT_ERROR_CLASS_ID               = "error_class_id"
	CI_JSON_ENCODER_FIELD_DEFAULT_ERROR_CLASS_NAME             = "error_class_name"
	CI_JSON_ENCODER_FIELD_DEFAULT_STACKTRACE                   = "stacktrace"
	CI_JSON_ENCODER_FIELD_DEFAULT_1DL_STACKTRACE_MESSAGES      = "stacktrace_messages"
	CI_JSON_ENCODER_FIELD_DEFAULT_FIELDS                       = "fields"
	CI_JSON_ENCODER_FIELD_DEFAULT_1DL_LOG_FIELDS_PREFIX        = "field_"
	CI_JSON_ENCODER_FIELD_DEFAULT_1DL_STACKTRACE_FIELDS_PREFIX = "field_stacktrace_{{num}}_"
)

var (
	// Make sure we won't break API.
	_ CI_Encoder = (*CI_JSONEncoder)(nil)
)

// SetOneDepthLevel sets a depth level of an output of CI_JSONEncoder.
//
// By default, the depth level is > 1, meaning an output may look like:
//
//	{
//	    "message": "Error message",
//	    "fields": {
//	        "key1": "value1" // <-- here 2 depth level
//	    }
//	    "stacktrace": [{
//	        "func": "ekago.v3.ekalog_test.foo", // <-- here 2 depth level
//	        "file": "logger_test.go:22",
//	        "package": "github.com/qioalice",
//	        "fields": {
//	            "test": 42 // <-- here 3 depth level
//	        }
//	    }]
//	}
//
// But enabling 1 depth level you will get:
//
//	{
//	    "message": "Error message",
//	    "field_key1": "value1",
//	    "stacktrace": [
//	        "github.com/qioalice/ekago.v3.ekalog_test.foo (logger_test.go:22)"
//	    ],
//	    "field_stacktrace_0_test": 42
//	}
//
// Calling this method many times will overwrite previous value.
//
// This method MUST NOT be called after CI_JSONEncoder is registered
// with CommonIntegrator using CommonIntegrator.WithEncoder() method.
func (je *CI_JSONEncoder) SetOneDepthLevel(enable bool) *CI_JSONEncoder {

	je.oneDepthLevel = enable
	return je
}

// SetIndent sets an indentation of JSON encoding format.
// Any value <= 0 meaning "no indentation".
//
// Calling this method many times will overwrite previous value of format string.
//
// This method MUST NOT be called after CI_JSONEncoder is registered
// with CommonIntegrator using CommonIntegrator.WithEncoder() method.
func (je *CI_JSONEncoder) SetIndent(num int) *CI_JSONEncoder {

	if num < 0 {
		num = 0
	}
	je.indent = num
	return je
}

// SetNameForField allows you to rename default name for some fields.
//
// Keep in mind, using this method you can overwrite SYSTEM field's names
// (the names of those fields that are ekalog.Entry contains).
// You CANNOT change the name of user-added fields to ekalog.Entry.
//
// Calling this method many times with the same `fieldType`
// will overwrite previous value.
//
// This method MUST NOT be called after CI_JSONEncoder is registered
// with CommonIntegrator using CommonIntegrator.WithEncoder() method.
func (je *CI_JSONEncoder) SetNameForField(fieldType CI_JSONEncoder_Field, name string) *CI_JSONEncoder {

	if je.fieldNames == nil {
		je.fieldNames = make(map[CI_JSONEncoder_Field]string)
	}
	je.fieldNames[fieldType] = name
	return je
}

// SetTimeFormatter allows you to set formatter that will be encode `time` field
// of the ekalog.Entry. Presented `formatter` MUST BE not nil, ignored otherwise.
//
// Calling this method many times will overwrite previous value of formatter.
//
// This method MUST NOT be called after CI_JSONEncoder is registered
// with CommonIntegrator using CommonIntegrator.WithEncoder() method.
func (je *CI_JSONEncoder) SetTimeFormatter(formatter func(t time.Time) string) *CI_JSONEncoder {

	if formatter != nil {
		je.timeFormatter = formatter
	}
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
	if f.Kind&ekaletter.KIND_FLAG_USER_DEFINED != 0 {
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
	copied := make([]byte, len(b)+1)
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
