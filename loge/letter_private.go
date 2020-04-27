// Copyright © 2020. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package loge

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/qioalice/gext/mathe"

	"github.com/modern-go/reflect2"
)

type (
	// letter is a core of Logger's Entry or an Error's object.
	// TBH I didn't know what name I can use for this entity, so let's call it as is.
	//
	// Both of Logger's Entry and Error's object may (and should) contain
	// message's body and some arguments. These fields are public and they
	// have absolutely the same logic.
	//
	// When you can access 'Message' or 'Field' outside it's considered generated.
	// Until then these fields may contain some temporary data that is useful
	// to finish transforming input parse arguments to 'Message' and 'Fields'.
	letter struct {

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
		Fields []Field

		// flags describes how should the parsing process proceed.
		// See constants for more info.
		flags mathe.Flags8
	}
)

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
const (
	_LF_ONLY_EXPLICIT_FIELDS    mathe.Flag8 = 0b_0000_0010
	_LF_ALLOW_UNNAMED_NIL       mathe.Flag8 = 0b_0000_0100
	_LF_ALLOW_IMPLICIT_POINTERS mathe.Flag8 = 0b_0000_1000
)

// message is a l.Message setter. Saves 'message' to l.Message and returns l.
func (l *letter) message(message string) *letter {
	l.Message = message
	return l
}

//
// printf-like format string or just message; can be an empty
// Both of args and explicitArgs can be empty at the same time but
// they can't be non-empty both at the same time.
// nil if len(explicitArgs) > 0
// nil if len(args) > 0

// parse just adds 'explicitArgs' (if it's not empty) to l.Fields or
// parses 'args' for being used in both of Message and Fields fields
// (or only in Fields if 'onlyFields' is true).
//
// This method is used when you just want to add some fields or start the whole
// process of generating finally message (printf-like style supported) and fields
// if it's possible.
//
// So:
// 1. If l.Message == "",
func (l *letter) parse(args []interface{}, explicitArgs []Field, onlyFields bool) *letter {

	// ----
	// REMINDER!
	// IT IS STRONGLY GUARANTEES THAT BOTH OF 'args' AND 'explicitArgs'
	// CAN NOT BE AN EMPTY AT THE SAME TIME!
	// ----

	var (
		messageNeedsArgs int
		messageArgs      []interface{}
	)

	switch {
	case len(explicitArgs) > 0 && len(l.Fields) > 0:
		l.Fields = append(l.Fields, explicitArgs...) // easy case
		return l

	case len(explicitArgs) > 0:
		l.Fields = explicitArgs // easy case
		return l

	case len(args) > 0 && !onlyFields:
		messageNeedsArgs = printfHowMuchVerbs(&l.Message)
		messageArgs = make([]interface{}, 0, messageNeedsArgs)
	}

	var (
		bakFlags mathe.Flags8 // see 'onlyFields' if statement
		fieldKey string       // below loop's var
	)

	if onlyFields {
		// Avoid l.Message related cases in l.parseArg()'s switch statements:
		// make sure they never happen
		bakFlags, l.flags = l.flags, l.flags&^_LF_ONLY_EXPLICIT_FIELDS
		// We don't need to backup l.Message because it can not be reused.
		// l.Message regenerates every time when finisher is called.
		l.Message = "."
	}

	// -- MAIN LOOP --

	for i, n := 0, len(args); i < n; i++ {

		// let's recognize what kind of arg we have
		switch argType := reflect2.TypeOf(args[i]); {

		case argType == reflectTypeField || argType == reflectTypeFieldPtr:
			l.addExplicitField(args[i], argType)

		case fieldKey != "":
			// it guarantees that if 'fieldKey' is not empty, message's body
			// is already formed
			l.addImplicitField(fieldKey, args[i], argType)
			fieldKey = ""

		case messageNeedsArgs > 0:
			// printf wants value to its associated verb;
			// l.messagePrintfArgs is used to store all printf values
			messageNeedsArgs--
			fallthrough

		case l.flags.TestAll(_LF_ONLY_EXPLICIT_FIELDS):
			// even if it's overwrites required printf values, use it
			// we must not drop any passed arg!
			fallthrough

		case l.Message == "" && len(messageArgs) == 0:
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

		case args[i] != nil || l.flags.TestAll(_LF_ALLOW_UNNAMED_NIL):
			// unnamed field
			l.addImplicitField("", args[i], argType)
		}
	} // end loop

	if onlyFields {
		l.flags = bakFlags
		//l.Message = "" // tbh, this statement is redundant
		return l
	}

	// if after loop 'fieldKey' != "", it's unnamed field
	if fieldKey != "" {
		l.Fields = append(l.Fields, String("", fieldKey))
	}

	// TODO: What do we have to do if we had printf verbs < than required ones?
	//  Now it's handled by fmt.Printf, but I guess we shall to handle it manually.

	switch hasPrintArgs := len(messageArgs) > 0; {

	case hasPrintArgs && l.Message != "":
		l.Message = fmt.Sprintf(l.Message, messageArgs...)

	case hasPrintArgs && l.Message == "":
		l.Message = fmt.Sprint(messageArgs...)

	default:
		// l.Message = l.Message (already formed)
	}

	if l.Message == "" {
		// okay, it's empty message
		// TODO: Shall we do something else with empty log messages?
	}

	return l
}

// addImplicitField adds new field to l.Fields treating 'name' as field's name,
// 'value' as field's value and using 'typ' (assuming it's value's type)
// to recognize how to convert Golang's interface{} to the 'Field' object.
func (l *letter) addImplicitField(name string, value interface{}, typ reflect2.Type) {

	switch {
	case value == nil:
		l.Fields = append(l.Fields, fieldNilValue(name, FieldKindInvalid))
		return

	case typ.Implements(reflectTypeFmtStringer):
		var stringer fmt.Stringer
		typ.Set(&stringer, &value)
		l.Fields = append(l.Fields, Stringer(name, stringer))
		return
	}

	switch typ.Kind() {

	case reflect.Ptr:
		if l.flags.TestAll(_LF_ALLOW_IMPLICIT_POINTERS) {
			l.Fields = append(l.Fields, Addr(name, value))

		} else {
			value = typ.Indirect(value)
			l.addImplicitField(name, value, reflect2.TypeOf(value))
		}

	case reflect.Bool:
		var boolVal bool
		typ.UnsafeSet(unsafe.Pointer(&boolVal), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Bool(name, boolVal))

	case reflect.Int:
		var intVal int
		typ.UnsafeSet(unsafe.Pointer(&intVal), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Int(name, intVal))

	case reflect.Int8:
		var int8Val int8
		typ.UnsafeSet(unsafe.Pointer(&int8Val), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Int8(name, int8Val))

	case reflect.Int16:
		var int16Val int16
		typ.UnsafeSet(unsafe.Pointer(&int16Val), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Int16(name, int16Val))

	case reflect.Int32:
		var int32Val int32
		typ.UnsafeSet(unsafe.Pointer(&int32Val), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Int32(name, int32Val))

	case reflect.Int64:
		var int64Val int64
		typ.UnsafeSet(unsafe.Pointer(&int64Val), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Int64(name, int64Val))

	case reflect.Uint:
		var uintVal uint64
		typ.UnsafeSet(unsafe.Pointer(&uintVal), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Uint64(name, uintVal))

	case reflect.Uint8:
		var uint8Val uint8
		typ.UnsafeSet(unsafe.Pointer(&uint8Val), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Uint8(name, uint8Val))

	case reflect.Uint16:
		var uint16Val uint16
		typ.UnsafeSet(unsafe.Pointer(&uint16Val), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Uint16(name, uint16Val))

	case reflect.Uint32:
		var uint32Val uint32
		typ.UnsafeSet(unsafe.Pointer(&uint32Val), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Uint32(name, uint32Val))

	case reflect.Uint64:
		var uint64Val uint64
		typ.UnsafeSet(unsafe.Pointer(&uint64Val), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Uint64(name, uint64Val))

	case reflect.Float32:
		var float32Val float32
		typ.UnsafeSet(unsafe.Pointer(&float32Val), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Float32(name, float32Val))

	case reflect.Float64:
		var float64Val float64
		typ.UnsafeSet(unsafe.Pointer(&float64Val), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Float64(name, float64Val))

	case reflect.Complex64:
		var complex64Val complex64
		typ.UnsafeSet(unsafe.Pointer(&complex64Val), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Complex64(name, complex64Val))

	case reflect.Complex128:
		var complex128Val complex128
		typ.UnsafeSet(unsafe.Pointer(&complex128Val), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, Complex128(name, complex128Val))

	case reflect.String:
		var stringVal string
		typ.UnsafeSet(unsafe.Pointer(&stringVal), reflect2.PtrOf(value))
		l.Fields = append(l.Fields, String(name, stringVal))

	case reflect.Uintptr:
		var uintptrVal uintptr
		typ.Set(&uintptrVal, &value)
		l.Fields = append(l.Fields, Addr(name, uintptrVal))

	case reflect.UnsafePointer:
		var unsafePtrVal unsafe.Pointer
		typ.Set(&unsafePtrVal, &value)
		l.Fields = append(l.Fields, Addr(name, unsafePtrVal))
	}
}

// addExplicitField assumes that 'field' either 'Field' or '*Field' type
// (checks it by 'fieldType') and adds it to the l.Fields.
// If field == (*Field)(nil) there is no-op.
func (l *letter) addExplicitField(field interface{}, fieldType reflect2.Type) {

	var explicitFieldPtr *Field

	if fieldType == reflectTypeFieldPtr {
		explicitFieldPtr = field.(*Field)

	} else {
		explicitFieldPtr = new(Field)
		*explicitFieldPtr = field.(Field)
	}

	if explicitFieldPtr != nil {
		l.Fields = append(l.Fields, *explicitFieldPtr)
	}
}
