// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"fmt"
	"strings"
	"time"

	"github.com/qioalice/ekago/v3/internal/ekaletter"
)

type (
	// Error is an object representation of your abstract error.
	// Error accumulates and stores some your data that will help you to log it later.
	//
	// In your code you must now use *Error as error indicating. Like
	//     func foo() *Error
	// You still may return nil if no error has occurred, or create an Error
	// object and return it (by reference).
	//
	// TLDR:
	// - DO NOT CREATE ERROR OBJECTS MANUALLY, USE Class's CONSTRUCTORS INSTEAD.
	// - DO NOT FORGOT TO USE Throw().
	// - IF YOU WANT TO DO SOMETHING WITH YOUR ERROR, DO IT BEFORE LOGGING.
	// - ALL ERROR OBJECTS ARE THREAD-UNSAFE. AVOID POTENTIAL DATA RACES!
	// - NEVER USE ERROR OBJECT AS VALUE, ALWAYS USE BY REFERENCE.
	//
	// -----
	//
	// ERROR OBJECTS CREATED MANUALLY CONSIDERED NOT INITIALIZED AND WILL NOT
	// WORK PROPERLY, WILL NOT CONTAIN ANY YOUR DATA AND WILL NOT WORK AT ALL!
	//
	// Use Class.New(), Class.Wrap(), Class.LightNew(), Class.LightWrap() methods
	// to create an *Error object.
	//
	// BECAUSE OF THE MAIN IDEA OF THIS PACKAGE AND EKALOG PACKAGE,
	// ERROR AND ITS LOGGING ARE STRONGLY LINKED AND ERROR EXISTENCE WITHOUT LOGGING
	// IS MEANINGLESS.
	//
	// In accordance with the above, and in pursuit of better performance
	// (decreasing unnecessary allocations / deallocations => RAM reusing)
	// once allocated *Error objects may be reused in the future.
	// Because of that, follow these rules:
	//
	// IF YOU LOG AN ERROR, YOU MUST NOT USE AN ERROR OBJECT THEN!
	// DO WHATEVER YOU WANT WITH AN ERROR OBJECT BEFORE YOU WILL LOG IT USING EKALOG.
	// AFTER YOU LOG YOUR ERROR THAT ERROR OBJECT WILL BE BROKEN AND CANNOT BE USED.
	//
	// This is because after calling error's log finisher
	// ( from https://github.com/qioalice/ekago/ekaerr/error_ekalog.go )
	// an internal part of *Error object is returned to the pool for being reused
	// (RAM optimisation).
	//
	// If you do not want to log an error but want to return it manually to the pool
	// (I don't know the case when you needed it, but whatever) you can just call
	// ReleaseError(err).
	// If you won't do it, an allocated Error will returned to the pool automatically
	// when it goes out from the scope.
	//
	// -----
	//
	// Lightweight errors.
	//
	// Lightweight errors is just the same Error but w/o stacktrace generating.
	// You can also add fields or messages but they won't be linked to some stacktrace.
	// Thus, Throw() call in that case is meaningless and will do nothing.
	//
	// Often it's useful when you don't want to log your error but do something instead.
	//
	// If you will log lightweight Error, make sure your encoder supports lightweight errors.
	// Both of ekalog.CI_ConsoleEncoder, ekalog.CI_JSONEncoder provides that.
	//
	// -----
	//
	// Error's ID.
	//
	// Each Error object (even lightweight) has its own ID.
	// You can use that ID to find and determine Error by its ID.
	//
	// Earlier and UUIDv4 was used to generate ID.
	// Since 3.0 ver, an ULID is used instead.
	//
	// Read more about ULID here: https://github.com/ulid/spec
	// and here https://github.com/oklog/ulid .
	Error struct {

		// letter is the main internal part of Error object.
		// It stores a stacktrace, fields, messages, public message, unique error's ID.
		//
		// If it's nil, the Error object is considered broken and not valid.
		// Breaks after Error's logging. Under pool reusing (RAM optimisation).
		//
		// See https://github.com/qioalice/ekago/internal/letter/letter.go .
		letter *ekaletter.Letter

		// classID is just an ID of Class what has been used to create this object.
		classID ClassID

		// namespaceID is just an ID of Namespace
		// the Class that has been used to create this object, belongs to.
		namespaceID NamespaceID

		needSetFinalizer bool
	}
)

// IsValid reports whether Error is valid Error object or not.
//
// It returns false if either Error is nil, or Error has not been initialized properly
// (instantiated manually instead of Class's constructors calling).
//
// You can also use IsNil() and IsNotNil() to make your code more clean.
func (e *Error) IsValid() bool {
	return e != nil && e.letter != nil
}

// IsNotNil is IsValid() alias. Does absolutely the same thing.
// Introduced to increase your code's cleaning.
func (e *Error) IsNotNil() bool {
	return e.IsValid()
}

// IsNil is `!Error.IsValid()` code alias. Does exactly that thing, nothing more.
// Introduced to increase your code's cleaning (and easy chaining typing).
func (e *Error) IsNil() bool {
	return !e.IsValid()
}

// Throw is an OOP style of raising an error up.
// YOU MUST CALL THIS METHOD EACH TIME YOU RETURNING AN ERROR OBJECT FROM THE FUNC.
// THIS CALL MUST BE THE LAST ONE AT THE YOUR RETURN STATEMENT.
// Nil safe. Returns this.
//
// Typically, you have two cases you must follow.
// 1. An Error instantiating:
//        func foo() *Error {
//            return ekaerr.IllegalState.New("something happen").Throw()
//        }
// 2. Already existed error raising up:
//    (an example with adding custom stack frame's message and field)
//        if err := foo(); err != nil {
//            return err.S("foo failed").W("id", 42).Throw()
//        }
//
// Does nothing for lightweight errors.
// Nil safe.
func (e *Error) Throw() *Error {
	if e.IsValid() {
		ekaletter.LIncStackIdx(e.letter)
	}
	return e
}

// AddMessage adds a message to your current Error's stack frame.
// Nil safe. Returns this.
func (e *Error) AddMessage(message string) *Error {
	if e.IsValid() {
		if message = strings.TrimSpace(message); message != "" {
			ekaletter.LSetMessage(e.letter, message, true)
		}
	}
	return e
}

// With adds provided ekaletter.LetterField to your current Error stack frame.
// Nil safe. Returns this.
func (e *Error) With(f ekaletter.LetterField) *Error { return e.addField(f) }

// Methods below are code-generated.

func (e *Error) WithBool(key string, value bool) *Error {
	return e.addField(ekaletter.FBool(key, value))
}
func (e *Error) WithInt(key string, value int) *Error {
	return e.addField(ekaletter.FInt(key, value))
}
func (e *Error) WithInt8(key string, value int8) *Error {
	return e.addField(ekaletter.FInt8(key, value))
}
func (e *Error) WithInt16(key string, value int16) *Error {
	return e.addField(ekaletter.FInt16(key, value))
}
func (e *Error) WithInt32(key string, value int32) *Error {
	return e.addField(ekaletter.FInt32(key, value))
}
func (e *Error) WithInt64(key string, value int64) *Error {
	return e.addField(ekaletter.FInt64(key, value))
}
func (e *Error) WithUint(key string, value uint) *Error {
	return e.addField(ekaletter.FUint(key, value))
}
func (e *Error) WithUint8(key string, value uint8) *Error {
	return e.addField(ekaletter.FUint8(key, value))
}
func (e *Error) WithUint16(key string, value uint16) *Error {
	return e.addField(ekaletter.FUint16(key, value))
}
func (e *Error) WithUint32(key string, value uint32) *Error {
	return e.addField(ekaletter.FUint32(key, value))
}
func (e *Error) WithUint64(key string, value uint64) *Error {
	return e.addField(ekaletter.FUint64(key, value))
}
func (e *Error) WithUintptr(key string, value uintptr) *Error {
	return e.addField(ekaletter.FUintptr(key, value))
}
func (e *Error) WithFloat32(key string, value float32) *Error {
	return e.addField(ekaletter.FFloat32(key, value))
}
func (e *Error) WithFloat64(key string, value float64) *Error {
	return e.addField(ekaletter.FFloat64(key, value))
}
func (e *Error) WithComplex64(key string, value complex64) *Error {
	return e.addField(ekaletter.FComplex64(key, value))
}
func (e *Error) WithComplex128(key string, value complex128) *Error {
	return e.addField(ekaletter.FComplex128(key, value))
}
func (e *Error) WithString(key string, value string) *Error {
	return e.addField(ekaletter.FString(key, value))
}
func (e *Error) WithStringFromBytes(key string, value []byte) *Error {
	return e.addField(ekaletter.FStringFromBytes(key, value))
}
func (e *Error) WithBoolp(key string, value *bool) *Error {
	return e.addField(ekaletter.FBoolp(key, value))
}
func (e *Error) WithIntp(key string, value *int) *Error {
	return e.addField(ekaletter.FIntp(key, value))
}
func (e *Error) WithInt8p(key string, value *int8) *Error {
	return e.addField(ekaletter.FInt8p(key, value))
}
func (e *Error) WithInt16p(key string, value *int16) *Error {
	return e.addField(ekaletter.FInt16p(key, value))
}
func (e *Error) WithInt32p(key string, value *int32) *Error {
	return e.addField(ekaletter.FInt32p(key, value))
}
func (e *Error) WithInt64p(key string, value *int64) *Error {
	return e.addField(ekaletter.FInt64p(key, value))
}
func (e *Error) WithUintp(key string, value *uint) *Error {
	return e.addField(ekaletter.FUintp(key, value))
}
func (e *Error) WithUint8p(key string, value *uint8) *Error {
	return e.addField(ekaletter.FUint8p(key, value))
}
func (e *Error) WithUint16p(key string, value *uint16) *Error {
	return e.addField(ekaletter.FUint16p(key, value))
}
func (e *Error) WithUint32p(key string, value *uint32) *Error {
	return e.addField(ekaletter.FUint32p(key, value))
}
func (e *Error) WithUint64p(key string, value *uint64) *Error {
	return e.addField(ekaletter.FUint64p(key, value))
}
func (e *Error) WithFloat32p(key string, value *float32) *Error {
	return e.addField(ekaletter.FFloat32p(key, value))
}
func (e *Error) WithFloat64p(key string, value *float64) *Error {
	return e.addField(ekaletter.FFloat64p(key, value))
}
func (e *Error) WithType(key string, value interface{}) *Error {
	return e.addField(ekaletter.FType(key, value))
}
func (e *Error) WithStringer(key string, value fmt.Stringer) *Error {
	return e.addField(ekaletter.FStringer(key, value))
}
func (e *Error) WithAddr(key string, value interface{}) *Error {
	return e.addField(ekaletter.FAddr(key, value))
}
func (e *Error) WithUnixFromStd(key string, value time.Time) *Error {
	return e.addField(ekaletter.FUnixFromStd(key, value))
}
func (e *Error) WithUnixNanoFromStd(key string, value time.Time) *Error {
	return e.addField(ekaletter.FUnixNanoFromStd(key, value))
}
func (e *Error) WithUnix(key string, value int64) *Error {
	return e.addField(ekaletter.FUnix(key, value))
}
func (e *Error) WithUnixNano(key string, value int64) *Error {
	return e.addField(ekaletter.FUnixNano(key, value))
}
func (e *Error) WithDuration(key string, value time.Duration) *Error {
	return e.addField(ekaletter.FDuration(key, value))
}
func (e *Error) WithArray(key string, value interface{}) *Error {
	return e.addField(ekaletter.FArray(key, value))
}
func (e *Error) WithObject(key string, value interface{}) *Error {
	return e.addField(ekaletter.FObject(key, value))
}
func (e *Error) WithMap(key string, value interface{}) *Error {
	return e.addField(ekaletter.FMap(key, value))
}
func (e *Error) WithExtractedMap(key string, value map[string]interface{}) *Error {
	return e.addField(ekaletter.FExtractedMap(key, value))
}
func (e *Error) WithAny(key string, value interface{}) *Error {
	return e.addField(ekaletter.FAny(key, value))
}

func (e *Error) WithMany(fields ...ekaletter.LetterField) *Error {
	return e.addFields(fields)
}

func (e *Error) WithManyAny(fields ...interface{}) *Error {
	return e.addFieldsParse(fields, true)
}

func (e *Error) WithDescription(description string) *Error {
	return e.WithString("description", description)
}

// Apply calls f callback passing the current Error object into and returning
// the Error object, callback is return what.
// Nil safe.
//
// Does nothing if f is nil.
//
// Why?
// You can write your own fields appender like:
//     func (this *YourType) errAddIdentifiers(err *ekaerr.Error) *ekaerr.Error {
//         return err.WithManyAny("id", this.ID, "date", this.Date)
//     }
//
// And then use it like:
//     yt := new(YourType)
//     return ekaerr.IllegalState.New("Unexpected state of world").
//         Apply(yt.errAddIdentifiers).
//         Throw()
//
// Brilliant, isn't? And most importantly it's so clean.
func (e *Error) Apply(f func(err *Error) *Error) *Error {
	if e.IsValid() && f != nil {
		return f(e)
	}
	return e
}

// Is reports whether Error has been instantiated by cls Class's constructors.
// Returns false if either Error is not valid or Class is invalid.
// Nil safe.
func (e *Error) Is(cls Class) bool {
	return e.IsValid() && isValidClassID(cls.id) && e.classID == cls.id
}

// IsAny reports whether Error belongs to at least one of passed cls Class
// (has been instantiated using one of them).
// Returns false if Error is not valid or no one class has been passed.
// Nil-safe.
func (e *Error) IsAny(cls ...Class) bool {
	return e.is(cls, false)
}

// Of reports whether Error has been instantiated by some Class that belongs to
// ns Namespace. Returns false if either Error is not valid or Namespace is invalid.
// Nil safe.
func (e *Error) Of(ns Namespace) bool {
	return e.IsValid() && isValidNamespaceID(ns.id) && e.namespaceID == ns.id
}

// OfAny reports whether Error belongs to at least one of passed nss Namespace
// (has been instantiated by the Class that belongs to one of nss Namespace).
// Returns false if either Error is not valid or no one namespace has been passed.
// Nil safe.
func (e *Error) OfAny(nss ...Namespace) bool {
	return e.of(nss)
}

// IsAnyDeep reports whether Error belongs to at least one of passed cls Class
// or any of its parent (base) Class is the same as one of passed cls.
// Returns false if Error is not valid or no one class has been passed.
// Nil-safe.
//
// IsAnyDeep has increased algorithmic complexity (uses recursive algorithm)
// and MUCH SLOWER than IsAny() if you pass subclasses. So, make sure it's what you need.
func (e *Error) IsAnyDeep(cls ...Class) bool {
	return e.is(cls, true)
}

// Class returns Error's Class. A special invalidClass is returned if Error is nil
// or has been manually instantiated instead of constructor using.
// Nil safe.
func (e *Error) Class() Class {
	if !e.IsValid() {
		return invalidClass
	}
	return classByID(e.classID, true)
}

// ReplaceClass replaces both of Error's Class and Namespace to the provided Class
// and its Namespace.
// Does nothing if any of Error or new Class is invalid.
// Nil safe.
func (e *Error) ReplaceClass(newClass Class) *Error {
	if e.IsValid() && newClass.IsValid() {
		e.classID = newClass.id
		e.namespaceID = newClass.namespaceID
	}
	return e
}

// ID returns an unique Error's ID as ULID. You can tell this ID to the user and
// log this error. Then it will be easy to find an associated error.
// Returns "" if Error is not valid.
// Nil safe.
func (e *Error) ID() string {
	if !e.IsValid() {
		return ""
	}
	return e.letter.SystemFields[_ERR_SYS_FIELD_IDX_ERROR_ID].SValue
}

// ReleaseError prepares Error for being reused in the future and releases
// its internal parts (returning them to the pool).
//
// YOU MUST NOT USE ERROR OBJECT AFTER PASSING THEM INTO THIS FUNCTION.
// YOU DO NOT NEED TO CALL THIS FUNCTION IF YOUR ERROR WILL BE LOGGED USING EKALOG.
func ReleaseError(err *Error) {
	if err.IsValid() {
		releaseError(err)
	}
}
