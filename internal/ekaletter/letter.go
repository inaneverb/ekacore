// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaletter

import (
	"fmt"
	"unsafe"

	"github.com/qioalice/ekago/v3/ekasys"
	"github.com/qioalice/ekago/v3/internal/ekaclike"

	"github.com/modern-go/reflect2"
)

type (
	// Letter is a core for both of ekalog.Entry, ekaerr.Error objects.
	//
	// Both of them may (and should) contain message's body and some attached values.
	// And that's where they are stored.
	Letter struct {

		// StackTrace is just stack trace, nothing more.
		// It fills by ekasys.GetStackTrace() when Letter is under initializing.
		//
		// It can be nil if:
		//
		//  - This Letter belongs to ekalog.Entry, and it's linked
		//    with some ekaerr.Error object to use it as a stacktrace's source.
		//
		//  - This Letter is initialized to be w/o stacktrace.
		//    It's possible for both of ekalog.Entry, ekaerr.Error.
		//
		StackTrace ekasys.StackTrace

		// Messages contains some messages for each stackframe from StackTrace.
		//
		// It's an array, each element of which has an index of stackframe from StackTrace,
		// that element belongs to.
		//
		// It guarantees that each next element's LetterMessage.StackFrameIdx
		// GTE than prev.
		Messages []LetterMessage

		// Fields contains some attached values (fields) for each stackframe from StackTrace.
		//
		// It's an array, each element of which has an index of stackframe from StackTrace,
		// that element belongs to.
		//
		// It guarantees that each next element's LetterField.StackFrameIdx
		// GTE than prev.
		Fields []LetterField

		// SystemFields contains only important system meta information,
		// that could be generated by ekaerr.Error's or ekalog.Logger's methods.
		//
		// It guarantees, that all these field has set KIND_FLAG_SYSTEM bit
		// at the their LetterField.Kind property.
		SystemFields []LetterField

		// ---------------------------- PRIVATE ---------------------------- //

		// stackFrameIdx is a counter that generally uses only for ekaerr.Error object.
		// It's a "pointer" to a current stackframe when stack unwinding is performing.
		//
		// Look:
		// ekaerr.Error has a Throw() method, that does increase this counter
		// (until it reach StackTrace's len -1).
		//
		stackFrameIdx int16
	}
)

var (
	// RTypesBeingIgnoredForParsing is an array of RTypes that will be ignored
	// at the arguments parsing in LParseTo() function.
	RTypesBeingIgnoredForParsing = []uintptr{ RTypeLetterField, RTypeLetterFieldPtr }
)

// LAddField just adds passed LetterField at the end of provided Letter.
func LAddField(l *Letter, f LetterField) {
	f.StackFrameIdx = l.stackFrameIdx
	l.Fields = append(l.Fields, f)
}

// LAddFieldWithCheck calls LAddField() but with some checks before.
// LAddField() won't be called if provided LetterField is invalid or zero-vary.
func LAddFieldWithCheck(l *Letter, f LetterField) {
	if !(f.IsInvalid() || f.RemoveVary() && f.IsZero()) {
		LAddField(l, f)
	}
}

// LSetMessage adds (or overwrites if overwrite is true) message in provided Letter
// that is relevant for the current stack frame idx.
func LSetMessage(l *Letter, msg string, overwrite bool) {
	switch lm := len(l.Messages); {
	case lm == 0 || lm > 0 && l.Messages[lm-1].StackFrameIdx < l.stackFrameIdx:
		l.Messages = append(l.Messages, LetterMessage{
			Body:          msg,
			StackFrameIdx: l.stackFrameIdx,
		})
	case overwrite:
		// Prev cond can be false only if l.Messages[lm-1].StackFrameIdx == l.stackFrameIdx
		// meaning that message for the current stackframe is presented.
		// So, we can overwrite it, if it's allowed.
		l.Messages[lm-1].Body = msg
	}
}

// LGetMessage returns a message that is relevant for the current stack frame idx
// in provided Letter.
func LGetMessage(l *Letter) string {
	if lm := len(l.Messages); lm > 0 && l.Messages[lm-1].StackFrameIdx == l.stackFrameIdx {
		return l.Messages[lm-1].Body
	}
	return ""
}

// LPopLastMessage returns last not empty message from provided Letter,
// replacing it by empty string when found.
func LPopLastMessage(l *Letter) string {
	for i := len(l.Messages)-1; i >= 0; i-- {
		if l.Messages[i].Body != "" {
			ret := l.Messages[i].Body
			l.Messages[i].Body = ""
			return ret
		}
	}
	return ""
}

// LIncStackIdx increments Letter's stackIdx property if it's allowed.
// Returns a new value.
//
// Increment won't happen (and current value is returned) if it's maximum
// of allowed stack idx for the current len of stackframe.
func LIncStackIdx(l *Letter) {
	if l.stackFrameIdx+1 < int16(len(l.StackTrace)) {
		l.stackFrameIdx++
	}
}

// LGetStackIdx returns Letter's stackIdx property.
func LGetStackIdx(l *Letter) int16 {
	return l.stackFrameIdx
}

// LSetStackIdx replaces a Letter's stackIdx to the provided one returning an old.
// IT MUST BE USED CAREFULLY. YOU CAN GET A PANIC IF YOUR STACK INDEX
// IS GREATER THAN A LENGTH OF EMBEDDED STACKFRAME.
func LSetStackIdx(l *Letter, newStackIdx int16) (oldStackIdx int16) {
	oldStackIdx = l.stackFrameIdx
	l.stackFrameIdx = newStackIdx
	return oldStackIdx
}

// LReset resets all internal fields, frees unnecessary RAM, preparing
// it for being returned to the pool and being reused in the future.
func LReset(l *Letter) *Letter {

	for i, n := 0, len(l.Fields); i < n; i++ {
		FieldReset(&l.Fields[i])
	}

	l.stackFrameIdx = 0
	l.Fields = l.Fields[:0]
	l.Messages = l.Messages[:1]

	m := &l.Messages[0]
	m.Body, m.StackFrameIdx = "", 0

	return l
}

// ParseTo is all-in-one function that parses 'args' to extract message
// (if 'onlyFields' is false) and fields to the Letter.
//
//  - If 'onlyField' is true then only fields tries to be extracted from 'args'
//    (explicit or implicit) and then saves to the Letter.
//
//  - If 'onlyField' is false then also a message tried to be extracted (or generated)
//    from 'args' and use it as current stackframe's message
//    and the rest of 'args' will be used to extract fields.
//
// If it's message extraction, then:
//
//   The first item in 'args' that can be used as message
//   (string-like or something that can be stringifed) will be used to do it.
//
//   If it's a string and has a printf verbs, then next N values from 'args'
//   will be used to generate a final string, N is a number of printf verbs in string.
//
//   The rest of 'args' will be parsed as a fields.
//
// Limitations:
// If 'args' contain item that has one of the following type it will be SKIPPED:
//   *Error, Error, *Letter, Letter, *LetterItem, LetterItem.
//   See InitRestrictedTypesBeingParsed() for more details.
//   IT IS NOT POSSIBLE TO USE ANOTHER ERRORS OR THEIR PRIVATE TYPES AS FIELDS.
//   BUILD ONE ERROR THAT WILL CONTAIN ALL YOU WANT INSTEAD OF ERROR SPAWNING.
//
// Requirements:
// 'li' != nil. Otherwise UB (may panic).
//
// Used to:
// - Add fields (explicit/implicit) into *Error,
// - Add fields (explicit/implicit) or/and message to *Logger's *Entry.
func LParseTo(l *Letter, args []interface{}, onlyFields bool) {

	var (
		message = LGetMessage(l)
		messageArgs []interface{}
		messageNeedsArgs int
	)

	if len(args) > 0 && !onlyFields {
		messageNeedsArgs = PrintfVerbsCount(&message)
		if messageNeedsArgs > 0 {
			messageArgs = make([]interface{}, 0, messageNeedsArgs)
		}
	}

	var (
		fieldKey   string // below loop's var
		messageBak string
	)

	if onlyFields {
		messageBak = message
		message = "message must be not empty to avoid its generating"
	}

	// isRestrictedType is an auxiliary function that will be used in the loop above
	// to figure out whether arg item must be ignored.
	isRestrictedType := func(rtype uintptr, basedOn []uintptr) bool {
		for i, n := 0, len(basedOn); i < n; i++ {
			if rtype == basedOn[i] {
				return true
			}
		}
		return false
	}

	for i, n := 0, len(args); i < n; i++ {

		var (
			typeArg = reflect2.TypeOf(args[i])
			rtypeArg = uintptr(0)
		)

		if args[i] != nil {
			rtypeArg = typeArg.RType()
		}

		// let's recognize what kind of arg we have
		switch {

		case rtypeArg == RTypeLetterField:
			var f LetterField
			typeArg.UnsafeSet(unsafe.Pointer(&f), reflect2.PtrOf(args[i]))
			LAddFieldWithCheck(l, f)

		case rtypeArg == RTypeLetterFieldPtr:
			var f *LetterField
			typeArg.UnsafeSet(unsafe.Pointer(&f), reflect2.PtrOf(args[i]))
			if f != nil {
				LAddFieldWithCheck(l, *f)
			}

		case isRestrictedType(rtypeArg, RTypesBeingIgnoredForParsing):
			// DO NOTHING

		case fieldKey != "":
			// it guarantees that if fieldKey is not empty, message's body
			// is already formed (if requested)
			LAddFieldWithCheck(l, FAny(fieldKey, args[i]))
			fieldKey = ""

		case messageNeedsArgs > 0:
			messageNeedsArgs--
			fallthrough

		case message == "" && len(messageArgs) == 0:
			// there is no message's body still and we'll use current arg as message's body.
			messageArgs = append(messageArgs, args[i])

		case rtypeArg == ekaclike.RTypeString:
			// at this code point arg could be only field's key or unnamed arg
			// well, looks like it's a key.
			typeArg.UnsafeSet(unsafe.Pointer(&fieldKey), reflect2.PtrOf(args[i]))

			// it can be "" (empty string) then handle it as unnamed field
			if fieldKey != "" {
				break // break switch, go to next loop's iter
			}
			fallthrough // fallthrough can't be in 'if' statement

		default:
			LAddFieldWithCheck(l, FAny("", args[i]))
		}
	}

	if onlyFields {
		message = messageBak
		return
	}

	// if after loop fieldKey != "", it's the last unnamed field
	if fieldKey != "" {
		LAddFieldWithCheck(l, FString("", fieldKey))
	}

	// TODO: What do we have to do if we had printf verbs < than required ones?
	//  Now it's handled by fmt.Printf, but I guess we shall to handle it manually.

	switch hasPrintArgs := len(messageArgs) > 0; {

	case hasPrintArgs && message != "":
		LSetMessage(l, fmt.Sprintf(message, messageArgs...), true)

	case hasPrintArgs && message == "":
		LSetMessage(l, fmt.Sprint(messageArgs...), true)

	default:
		// already formed
	}

	// TODO: Shall we do something else with empty messages?
}
