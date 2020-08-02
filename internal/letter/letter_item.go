// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package letter

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/qioalice/ekago/v2/internal/field"

	"github.com/modern-go/reflect2"
)

//noinspection GoNameStartsWithPackageName
type (
	// Letter is a core of Logger's Entry or an Error's object.
	//
	// Both of Logger's Entry and Error's object may (and should) contain
	// message's body and some arguments. These fields are public and they
	// have absolutely the same logic.
	//
	// When you can access 'Message' or 'Field' outside it's considered generated.
	// Until then these fields may contain some temporary data that is useful
	// to finish transforming input parse arguments to 'Message' and 'Fields'.
	LetterItem struct {

		// TODO
		stackFrameIdx int16

		// TODO
		next *LetterItem

		// Message contains Logger Entry's message body or Error's one.
		//
		// If you were used printf-like style to generate this message, when
		// you can access this field outside this package - it's formed.
		Message string

		// Fields contains all data you attaching to the Logger's Entry or Error.
		// Your fields may be named or unnamed and may contain various data.
		//
		// Even if you did use implicit arguments as fields, they were converted
		// to explicit and saved here.
		Fields []field.Field

		// Flags describes how should the parsing process being proceed.
		// Inherited from parent *Letter object, overwrites in the prepare() call.
		Flags Flag
	}
)

// ParseTo is all-in-one function that actually does one of two things:
// - Saves 'explicitFields' -> 'li' as is (add) if len('explicitFields') > 0 or
// - Parses 'args' to extract message and fields -> 'li' if len('args') > 0.
//
// If it's saving there is easy case. Just assign or append() call and nothing more.
// If it's parsing, then:
//
// - If 'onlyField' is true then only fields tries to be extracted from 'args'
//   (explicit or implicit) and then saves -> 'li'.
// - If 'onlyField' is false then also a message tried to be extracted (or generated)
//   from 'args' and use it as 'li's message and the rest of 'args' will be used
//   to extract fields.
//
// If it's message extraction, then:
//   The first item in 'args' that can be used as message (or its generation)
//   will be used to do it. If it's not string-like something - just use it as message.
//   But if it's something that looks like string (duck types), it will be tried to
//   used as printf format string if it has printf verbs. If it so, the N values
//   from 'args' will be used (and then the rest after N args will be used as fields),
//   where N is the count of printf verbs that has been detected in the printf format.
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
func ParseTo(li *LetterItem, args []interface{}, explicitFields []field.Field, onlyFields bool) {

	// REMINDER!
	// IT IS STRONGLY GUARANTEES THAT BOTH OF 'args' AND 'explicitFields'
	// CAN NOT BE AN EMPTY (OR SET) AT THE SAME TIME!

	var (
		messageNeedsArgs int
		messageArgs      []interface{}
	)

	switch {
	case len(explicitFields) > 0 && len(li.Fields) > 0:
		li.Fields = append(li.Fields, explicitFields...) // easy case
		return

	case len(explicitFields) > 0:
		li.Fields = explicitFields // easy case
		return

	case len(args) > 0 && !onlyFields:
		messageNeedsArgs = PrintfVerbsCount(&li.Message)
		messageArgs = make([]interface{}, 0, messageNeedsArgs)
	}

	var (
		fieldKey   string // below loop's var
		messageBak string
	)

	if onlyFields {
		messageBak = li.Message
		li.Message = "message must be not empty to avoid its generating"
	}

	// isRestrictedType is an auxiliary function that will be used in the loop above
	// to figure out whether 'args's item must be ignored.
	isRestrictedType := func(typ reflect2.Type, basedOn []reflect2.Type) bool {
		for i, n := 0, len(basedOn); i < n; i++ {
			if typ == basedOn[i] {
				return true
			}
		}
		return false
	}

	// -- MAIN LOOP --

	for i, n := 0, len(args); i < n; i++ {

		// let's recognize what kind of arg we have
		switch argType := reflect2.TypeOf(args[i]); {

		case argType == field.ReflectedType:
			li.addExplicitField2(args[i].(field.Field))

		case argType == field.ReflectedTypePtr:
			li.addExplicitFieldByPtr(args[i].(*field.Field))

		case isRestrictedType(argType, TypesBeingIgnoredForParsing):
			// DO NOTHING

		case fieldKey != "":
			// it guarantees that if 'fieldKey' is not empty, message's body
			// is already formed
			li.addImplicitField(fieldKey, args[i], argType)
			fieldKey = ""

		case messageNeedsArgs > 0:
			// printf wants value to its associated verb;
			// l.messagePrintfArgs is used to store all printf values
			messageNeedsArgs--
			fallthrough

		case li.Flags.TestAll(FLAG_ONLY_EXPLICIT_FIELDS) && !onlyFields:
			// even if it's overwrites required printf values, use it
			// we must not drop any passed arg!
			fallthrough

		case li.Message == "" && len(messageArgs) == 0:
			// there is no message's body still and we'll use current arg
			// as message's body but only if there is no another one the same
			// assuming that implicit fields are enabled
			// (and all other args will be treated as explicit/implicit fields)
			messageArgs = append(messageArgs, args[i])

		case argType.Kind() == reflect.String:
			// at this code point arg could be only field's key or unnamed arg
			// well, looks like it's a key of implicit field.
			//
			// the same as fieldKey = field.(string)
			argType.UnsafeSet(unsafe.Pointer(&fieldKey), reflect2.PtrOf(args[i]))

			// it can be "" (empty string) then handle it as unnamed field
			if fieldKey != "" {
				break // break switch, go to next loop's iter
			}
			fallthrough // fallthrough can't be in 'if' statement

		case args[i] != nil || li.Flags.TestAll(FLAG_ALLOW_UNNAMED_NIL):
			// unnamed field
			li.addImplicitField("", args[i], argType)
		}
	} // end loop

	if onlyFields {
		li.Message = messageBak
		return
	}

	// if after loop 'fieldKey' != "", it's the last unnamed field
	if fieldKey != "" {
		li.Fields = append(li.Fields, field.String("", fieldKey))
	}

	// TODO: What do we have to do if we had printf verbs < than required ones?
	//  Now it's handled by fmt.Printf, but I guess we shall to handle it manually.

	switch hasPrintArgs := len(messageArgs) > 0; {

	case hasPrintArgs && li.Message != "":
		li.Message = fmt.Sprintf(li.Message, messageArgs...)

	case hasPrintArgs && li.Message == "":
		li.Message = fmt.Sprint(messageArgs...)

	default:
		// li.Message = li.Message (already formed)
	}

	// TODO: Shall we do something else with empty log messages?
	//if li.Message == "" { ??? }
}

// StackFrameIdx returns current's *LetterItem stack frame ID or -1 if 'li' is invalid.
// Nil safe.
func (li *LetterItem) StackFrameIdx() int16 {

	if li == nil {
		return -1
	}
	return li.stackFrameIdx
}

// Next returns a next *LetterItem item from the linked list followed to the current's one.
// Nil safe.
func (li *LetterItem) Next() *LetterItem {
	if li != nil && li.next != nil && li.next.stackFrameIdx != -1 {
		return li.next
	}
	return nil
}

// LI_SetNext is just 'current'.next = 'next'. Returns modified 'current'.
//
// It's a function, not a method, because it's a part of internal package and
// I want to use this inside other ekago's packages (can't make it private method),
// but don't want user to use this method (can't make it public method).
//
// Requirements:
// 'current' != nil. Otherwise UB (may panic).
//
//noinspection GoSnakeCaseUsage
func LI_SetNextItem(current, next *LetterItem) *LetterItem {
	current.next = next
	return current
}

// LI_GetNextItem just returns 'current'.next. It's the same as 'current'.Next()
// but it will return a valid *LetterItem even is returned item is broken or invalid.
//
// It's a function, not a method, because it's a part of internal package and
// I want to use this inside other ekago's packages (can't make it private method),
// but don't want user to use this method (can't make it public method).
//
// Requirements:
// 'current' != nil. Otherwise UB (may panic).
//
//noinspection GoSnakeCaseUsage
func LI_GetNextItem(current *LetterItem) *LetterItem {
	return current.next
}

// LI_SetStackFrameIdx is just 'current'.stackFrameIdx = 'stackFrameIdx'.
// Returns modified 'current'.
//
// It's a function, not a method, because it's a part of internal package and
// I want to use this inside other ekago's packages (can't make it private method),
// but don't want user to use this method (can't make it public method).
//
// Requirements:
// 'current' != nil. Otherwise UB (may panic).
//
//noinspection GoSnakeCaseUsage
func LI_SetStackFrameIdx(current *LetterItem, stackFrameIdx int16) *LetterItem {
	current.stackFrameIdx = stackFrameIdx
	return current
}

// LI_ResetItem resets all internal fields, frees unnecessary RAM, preparing
// it for being returned to the pool (as a part of *Letter)
// and being reused in the future.
// DOES NOT CHANGE LINKED LIST LINKING! It means that 'current'.next won't be changed!
//
// It's a function, not a method, because it's a part of internal package and
// I want to use this inside other ekago's packages (can't make it private method),
// but don't want user to use this method (can't make it public method).
//
// Requirements:
// 'current' != nil. Otherwise UB (may panic).
//
//noinspection GoSnakeCaseUsage
func LI_ResetItem(current *LetterItem) *LetterItem {

	// Flags will be restored manually

	for i, n := 0, len(current.Fields); i < n; i++ {
		field.Reset(&current.Fields[i])
	}

	current.stackFrameIdx = -1
	current.Fields = current.Fields[:0]
	current.Message = ""

	return current
}
