// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

import (
	"fmt"
	"reflect"
	"sync"
	"time"
	"unsafe"

	"github.com/qioalice/gext/sys"

	"github.com/modern-go/reflect2"
)

// General rules of arguments' parsing:
//
// 1. First string (or something that looks like string (1) (ducktypes))
//    is a message's body or it's printf-like format string.
//
// 2. If it's format string, N next args (except "explicit fields" (2))
//    are printf args, where N is printf verbs' number in format string.
//
// 3. All "explicit fields" (2) handled out of turn even if they
//    between key and value of some "implicit field" (3).
//
// 4. All next M-N args (M - total count of args, N described in p.2)
//    treated as key-value pairs of "implicit fields" (3).
//
// 5. If next arg in the processed queue should be field's key
//    it should be only string!
//    If it's not, it will be treated as "unnamed" (4) argument's value,
//    and processing will go to the next arg.
//
// --------
//
// (1). "Something that looks like string".
//      We expect not exactly string as message's body or it's printf-like
//      format string, but also []byte, fmt.Stringer values.
//      More flexible, more powerful.
//
// (2). "Explicit fields".
//      Explicit fields are the fields that specified explicitly.
//      It's just a values of Field, Fields or map[string]interface{}
//      types.
//
// (3). "Implicit fields".
//      All those fields that are not explicit.
//      If argument is not log message's body, printf format, printf value,
//      then it's either implicit field's key or implicit field's value.
//      Depended by its type and processing order.
//      E.g.
//      If prev arg was an implicit field's key, the next one is its value.
//      If prev arg was an implicit field's value, the next one is
//      a key of next implicit field (if it's string) or ...
//
// (4). "Unnamed implicit field".
//      ... or a value of next unnamed implicit field.
//      Thus unnamed implicit field is the field that has no key (name).
//      Yes it could be when, e.g:
//
//      	"arg1", 1, 2 - are in processing order, then
//
//      There are two implicit field: one named, one unnamed:
//      First has key "arg1" and value 1: {"arg1",1} pair. Named field.
//      Second has only value: 2: {2}. Unnamed field that will have
//      autogenerated key (name).

//
type Entry struct {
	l *Logger

	// todo: remove flag mask
	flagMask entryFlags

	Package string
	Func    string
	Class   string
	Method  string

	Time time.Time

	Level Level

	Message string

	Fields []Field

	StackTrace sys.StackTrace
	ssf        int // skip stack frames
	ssfp       int // skip stack frames private ("3" by default)

	beforeWrite BeforeWriteCallback
	logArgs     []interface{}
}

// todo: lazy generation field's array with the same key (issue, low prioirity, disableble)

// entryFlags represents
type entryFlags uint8

// BeforeWriteCallback is type of user's callback that fires exactly before
// log message's entry will be written to all registered destinations.
//
// Callback must return passed Entry if user wants Entry to be written.
// Otherwise return nil and passed Entry won't be written (and processed at all!).
type BeforeWriteCallback func(entry *Entry) *Entry

const (
	bEntryFlagDisableStacktrace entryFlags = 0x01 << iota
	bEntryFlagOverwriteStacktrace
	bEntryFlagAutoGenerateCaller

	bEntryFlagDontSkipEmptyMessages

	bEntryFlagOnlyExplicitFields
	bEntryFlagAllowImplicitPointers
	bEntryFlagAllowUnnamedNil
)

var _PoolEntry = sync.Pool{
	New: func() interface{} {
		return new(Entry)
	},
}

//
var EntryFlags = struct {
	AutoGenerateCaller entryFlags
}{
	AutoGenerateCaller: bEntryFlagAutoGenerateCaller,
}

//
func (e *Entry) testFlag(flag entryFlags) bool {
	return e.flagMask&flag == flag
}

//
func (e *Entry) setFlag(flag entryFlags) *Entry {
	e.flagMask |= flag
	return e
}

//
func (e *Entry) resetFlag(flag entryFlags) *Entry {
	e.flagMask &^= flag
	return e
}

//
func (e *Entry) apply(options []option) (this *Entry) {

	for _, option := range options {
		option(e)
	}

	return e
}

//
func (e *Entry) reset() (this *Entry) {

	e.l = nil

	e.flagMask = 0 |
		bEntryFlagAutoGenerateCaller |
		bEntryFlagAllowUnnamedNil

	e.Package = ""
	e.Func = ""
	e.Class = ""
	e.Func = ""

	for _, field := range e.Fields {
		field.reset()
	}

	e.Fields = e.Fields[:0]

	e.StackTrace = nil

	e.ssf = 0
	e.ssfp = 3

	for i, n := 0, len(e.logArgs); i < n; i++ {
		e.logArgs[i] = nil
	}

	e.logArgs = e.logArgs[:0]

	return e
}

// clone clones the current Entry object 'e' only if it's not thread-safety
// entry or if 'force' is true.
// Returns the copied clone, if it was clone, otherwise returns 'e'.
// todo: remove comment of force arg
// TODO: PUBLIC BECAUSE I WANT
func (e *Entry) clone() *Entry {

	clonedEntry := getEntry()

	clonedEntry.flagMask = e.flagMask

	clonedEntry.Package = e.Package
	clonedEntry.Func = e.Func
	clonedEntry.Class = e.Class

	clonedEntry.ssf = e.ssf
	clonedEntry.ssfp = e.ssfp

	clonedEntry.beforeWrite = e.beforeWrite

	// There is no necessary to zero Time, Level, Message fields
	// because they used only in one place and will be overwritten anyway.

	return clonedEntry
}

func (e *Entry) setPackageName(packageName string) (this *Entry) {

	e.resetFlag(bEntryFlagAutoGenerateCaller)
	e.Func, e.Class = "", ""
	e.Package = packageName
	return e
}

func (e *Entry) setFuncName(funcName string) (this *Entry) {

	e.resetFlag(bEntryFlagAutoGenerateCaller)
	e.Class = ""
	e.Func = funcName
	return e
}

func (e *Entry) setClassName(className string) (this *Entry) {

	e.resetFlag(bEntryFlagAutoGenerateCaller)
	e.Func = ""
	e.Class = className
	return e
}

func (e *Entry) setMethodName(methodName string) (this *Entry) {

	e.resetFlag(bEntryFlagAutoGenerateCaller)
	e.Method = methodName
	return e
}

func (e *Entry) forceStacktrace(ignoreFrames int) (this *Entry) {

	e.setFlag(bEntryFlagAutoGenerateCaller)
	e.Package, e.Func, e.Class = "", "", ""
	e.ssf = ignoreFrames
	return e
}

// addStacktrace adds caller (if it's not added yet manually) and stacktrace
// (if it's not added before from error or if forced overwrite is enabled).
func (e *Entry) addStacktrace() (this *Entry) {

	if !e.testFlag(bEntryFlagDisableStacktrace) &&
		(e.StackTrace == nil || e.testFlag(bEntryFlagOverwriteStacktrace)) {
		e.StackTrace = sys.GetStackTrace(e.ssf+e.ssfp, -1).ExcludeInternal()
	}

	if !e.testFlag(bEntryFlagAutoGenerateCaller) {
		return e
	}

	stacktrace := e.StackTrace
	if len(stacktrace) == 0 {
		stacktrace = sys.GetStackTrace(e.ssf+e.ssfp, 1)
		if len(stacktrace) == 0 {
			// TODO: Internal error, can't get a stacktrace => can't get a caller info
			return e
		}
	}

	e.Func = formCaller2(stacktrace[0])
	return e
}

// parseLogArgs tries to parse 'args' or 'explicitFields' and maybe combine it with
// the passed 'format' to fill Message and Field e's fields. Returns this.
//
// WARNING!
// Do not pass 'args' and 'explicitFields' at the same time! It will break everything!
func (e *Entry) parseLogArgs(format string, args []interface{}, explicitFields []Field) (this *Entry) {

	// ----
	// IT IS STRONGLY GUARANTEES THAT EITHER 'args' OR 'explicitFields'
	// IS NOT EMPTY AT THE SAME TIME!
	// ----

	// assume e.logArgs has enough space to store all args
	// that could be used to generate printf-like message or just a message.
	e.initLogArgs(len(args))

	var (
		requiredPrintfValues int    // doesn't calculated if 'explicitFields' isn't empty
		fieldKey             string // below loop's var
	)

	if len(explicitFields) > 0 {
		// easy case, args == nil
		e.Fields = explicitFields

	} else if len(args) > 0 {
		requiredPrintfValues = printfHowMuchVerbs(&format)
	}

	for _, arg := range args {

		// let's recognize what kind of arg we have
		switch argType := reflect2.TypeOf(arg); {

		case argType == reflectTypeField || argType == reflectTypeFieldPtr:
			e.addExplicitField(arg, argType)

		case fieldKey != "":
			// it guarantees that if 'fieldKey' is not empty, message's body
			// is already formed
			e.addImplicitField(fieldKey, arg, argType)
			fieldKey = ""

		case requiredPrintfValues > 0:
			// printf wants value to its associated verb,
			// e.logArgs is used to store all printf values
			requiredPrintfValues--
			fallthrough

		case e.testFlag(bEntryFlagOnlyExplicitFields):
			// even if it's overwrites required printf values, use it
			// we must not drop any passed arg!
			fallthrough

		case format == "" && len(e.logArgs) == 0:
			// there is no message's body still and we'll use current arg
			// as message's body but only if there is no another one the same
			// assuming that implicit fields are enabled
			// (and all other args will be treated as explicit/implicit fields)
			e.logArgs = append(e.logArgs, arg)

		case argType.Kind() == reflect.String:
			// at this code point arg could be only field's key or unnamed arg
			// well, looks like it's a key of implicit field.
			//
			// the same as fieldKey = field.(string)
			argType.UnsafeSet(unsafe.Pointer(&fieldKey), reflect2.PtrOf(arg))

			// it can be "" (empty string) then handle it as unnamed field
			if fieldKey != "" {
				break
			}
			fallthrough // fallthrough can't be in 'if' statement

		case arg != nil || e.testFlag(bEntryFlagAllowUnnamedNil):
			// unnamed field
			e.addImplicitField("", arg, argType)
		}
	}

	// if after loop 'fieldKey' != "", it's unnamed field
	if fieldKey != "" {
		e.Fields = append(e.Fields, String("", fieldKey))
	}

	// TODO: What do we have to do if we had printf verbs < than required ones?
	//  Now it's handled by fmt.Printf, but I guess we shall to handle it manually.

	switch hasPrintArgs := len(e.logArgs) > 0; {

	case hasPrintArgs && format != "":
		format = fmt.Sprintf(format, e.logArgs...)

	case hasPrintArgs && format == "":
		format = fmt.Sprint(e.logArgs...)
	}

	if format == "" {
		// okay, it's empty message
		// TODO: Shall we do something else with empty log messages?
	}

	e.Message = format
	return e
}

// with adds the both of explicit/implicit fields (anyFields) or only explicit ones
// to the current Entry and returns it.
func (e *Entry) with(anyFields []interface{}, explicitFields []Field) (this *Entry) {

	// It guarantees that at this code point at least one of anyFields, explicitFields
	// are not empty.

	if len(explicitFields) > 0 {
		e.Fields = explicitFields
	}

	var fieldKey string
	for _, field := range anyFields {

		// let's recognize what kind of arg we have
		switch fieldType := reflect2.TypeOf(field); {

		case fieldType == reflectTypeField || fieldType == reflectTypeFieldPtr:
			e.addExplicitField(field, fieldType)

		case fieldKey != "":
			// it guarantees that if 'fieldKey' is not empty, message's body
			// is already formed
			e.addImplicitField(fieldKey, field, fieldType)
			fieldKey = ""

		case fieldType.Kind() == reflect.String:
			// at this code point arg could be only field's key or unnamed arg
			// well, looks like it's a key of implicit field
			//
			// the same as fieldKey = field.(string)
			fieldType.UnsafeSet(unsafe.Pointer(&fieldKey), reflect2.PtrOf(field))

			// it can be "" (empty string) then handle it as unnamed field
			if fieldKey != "" {
				break
			}
			fallthrough // fallthrough can't be in 'if' statement

		case field != nil || e.testFlag(bEntryFlagAllowUnnamedNil):
			// unnamed field
			e.addImplicitField("", field, fieldType)
		}
	}

	// if after loop 'fieldKey' != "", it's unnamed field
	if fieldKey != "" {
		e.Fields = append(e.Fields, String("", fieldKey))
	}

	return e
}

func (e *Entry) addImplicitField(name string, value interface{}, typ reflect2.Type) {

	// 'reflect' does not have 'fmt.Stringer' kind, check it earlier,
	// cause switch has default
	if typ.Implements(reflectTypeFmtStringer) {

		var stringer fmt.Stringer
		typ.Set(&stringer, &value)
		e.Fields = append(e.Fields, Stringer(name, stringer))

		return
	}

	switch typ.Kind() {

	case reflect.Ptr:
		if e.testFlag(bEntryFlagAllowImplicitPointers) {
			e.Fields = append(e.Fields, Addr(name, value))

		} else {
			value = typ.Indirect(value)
			e.addImplicitField(name, value, reflect2.TypeOf(value))
		}

	case reflect.Bool:
		var boolVal bool
		typ.UnsafeSet(unsafe.Pointer(&boolVal), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Bool(name, boolVal))

	case reflect.Int:
		var intVal int
		typ.UnsafeSet(unsafe.Pointer(&intVal), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Int(name, intVal))

	case reflect.Int8:
		var int8Val int8
		typ.UnsafeSet(unsafe.Pointer(&int8Val), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Int8(name, int8Val))

	case reflect.Int16:
		var int16Val int16
		typ.UnsafeSet(unsafe.Pointer(&int16Val), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Int16(name, int16Val))

	case reflect.Int32:
		var int32Val int32
		typ.UnsafeSet(unsafe.Pointer(&int32Val), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Int32(name, int32Val))

	case reflect.Int64:
		var int64Val int64
		typ.UnsafeSet(unsafe.Pointer(&int64Val), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Int64(name, int64Val))

	case reflect.Uint:
		var uintVal uint64
		typ.UnsafeSet(unsafe.Pointer(&uintVal), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Uint64(name, uintVal))

	case reflect.Uint8:
		var uint8Val uint8
		typ.UnsafeSet(unsafe.Pointer(&uint8Val), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Uint8(name, uint8Val))

	case reflect.Uint16:
		var uint16Val uint16
		typ.UnsafeSet(unsafe.Pointer(&uint16Val), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Uint16(name, uint16Val))

	case reflect.Uint32:
		var uint32Val uint32
		typ.UnsafeSet(unsafe.Pointer(&uint32Val), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Uint32(name, uint32Val))

	case reflect.Uint64:
		var uint64Val uint64
		typ.UnsafeSet(unsafe.Pointer(&uint64Val), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Uint64(name, uint64Val))

	case reflect.Float32:
		var float32Val float32
		typ.UnsafeSet(unsafe.Pointer(&float32Val), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Float32(name, float32Val))

	case reflect.Float64:
		var float64Val float64
		typ.UnsafeSet(unsafe.Pointer(&float64Val), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Float64(name, float64Val))

	case reflect.Complex64:
		var complex64Val complex64
		typ.UnsafeSet(unsafe.Pointer(&complex64Val), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Complex64(name, complex64Val))

	case reflect.Complex128:
		var complex128Val complex128
		typ.UnsafeSet(unsafe.Pointer(&complex128Val), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, Complex128(name, complex128Val))

	case reflect.String:
		var stringVal string
		typ.UnsafeSet(unsafe.Pointer(&stringVal), reflect2.PtrOf(value))
		e.Fields = append(e.Fields, String(name, stringVal))

	case reflect.Uintptr:
		var uintptrVal uintptr
		typ.Set(&uintptrVal, &value)
		e.Fields = append(e.Fields, Addr(name, uintptrVal))

	case reflect.UnsafePointer:
		var unsafePtrVal unsafe.Pointer
		typ.Set(&unsafePtrVal, &value)
		e.Fields = append(e.Fields, Addr(name, unsafePtrVal))

	}
}

func (e *Entry) addExplicitField(field interface{}, fieldType reflect2.Type) {

	var explicitFieldPtr *Field

	if fieldType == reflectTypeFieldPtr {
		fieldType.Set(&explicitFieldPtr, &field)

	} else {
		explicitFieldPtr = new(Field)
		fieldType.Set(explicitFieldPtr, &field)
	}

	if explicitFieldPtr != nil {
		e.Fields = append(e.Fields, *explicitFieldPtr)
	}
}

//
func (e *Entry) initLogArgs(requiredElems int) {
	// TODO: implement?
}

//
func getEntry() *Entry {
	return _PoolEntry.Get().(*Entry)
}

//
func reuseEntry(e *Entry) {
	_PoolEntry.Put(e.reset())
}