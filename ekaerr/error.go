// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"github.com/qioalice/ekago/v2/internal/ekaletter"
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
	// - NEVER USE ERROR OBJECT AS VALUE, ALWAYS USE AS ITS REFERENCE.
	//
	// -----
	//
	// ERROR OBJECTS CREATED MANUALLY CONSIDERED NOT INITIALIZED AND WILL NOT
	// WORK PROPERLY, WILL NOT CONTAIN ANY YOUR DATA AND WILL NOT WORK AT ALL!
	//
	// Use New() or Wrap() error Class's methods to create an *Error object.
	// See https://github.com/qioalice/ekago/ekaerr/README.md for more details
	// or take a look at the examples.
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
	// AFTER YOU LOG YOUR ERROR THAT ERROR OBJECT WILL BE BROKEN AND CAN NOT BE USED.
	//
	// This is because after calling error's log finisher
	// ( from https://github.com/qioalice/ekago/ekaerr/error_ekalog.go )
	// an internal part of *Error object is returned to the pool for being reused
	// (RAM optimisation).
	//
	// If you do not want to log an error but want to return it manually to the pool
	// (i don't know the case when you needed it, but whatever) you can just call
	// ReleaseError(err).
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

		// stackIdx is an internal counter that is increased by Throw().
		// Allows to specify to which stack frame fields or message will be attached.
		stackIdx int16

		// TODO
		needSetFinalizer bool
	}
)

var (
	// DefaultFlags is the Flag's set that will be used as default flags
	// for each *LetterItem that is created for Error's *Letter.
	DefaultFlags = FLAG_ALLOW_UNNAMED_NIL

	// TODO: Maybe these flags must be protected by the mutex?
)

// IsValid reports whether e is valid Error object or not.
//
// It returns false if e == nil, or e has not been initialized properly
// (instantiated manually instead of Class's constructors calling).
//
// You can also use IsNil() and IsNotNil() to make your code more clean.
func (e *Error) IsValid() bool {
	return e != nil && e.letter != nil
}

// IsNotNil is IsValid's alias. Does absolutely the same thing.
// Introduced to increase your code's cleaning.
func (e *Error) IsNotNil() bool {
	return e.IsValid()
}

// IsNil is !e.IsValid() code alias. Does exactly that thing, nothing more.
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
// See https://github.com/qioalice/ekago/ekaerr/README.md for more details.
func (e *Error) Throw() *Error {
	if e.IsValid() && e.stackIdx+1 < int16(len(e.letter.StackTrace)) {
		e.stackIdx++
	}
	return e
}

// Mark marks your current e's stack frame as important.
// Nil safe. Returns this.
//
// It's very useful at the logging when you need to filter some especially significant
// stack frames. And because Throw() changing current stack frame, you must mark
// the stack frames you need as soon as you working with it.
//
// See https://github.com/qioalice/ekago/ekaerr/README.md for more details.
func (e *Error) Mark() *Error {
	if e.IsValid() {
		e.getCurrentLetterItem().Flags.SetAll(FLAG_MARKED_LETTER_ITEM)
	}
	return e
}

// AddMessage adds a 'message' to your current e's stack frame.
// Nil safe. Returns this.
//
// See https://github.com/qioalice/ekago/ekaerr/README.md for more details.
func (e *Error) AddMessage(message string) *Error {
	if e.IsValid() {
		e.getCurrentLetterItem().Message = message
		if e.stackIdx == 0 {
			e.Mark()
		}
	}
	return e
}

// AddFields extract key-value pairs from 'args' and adds it to your current e's stack frame.
// Nil safe. Returns this.
//
// You can use only strings as keys. If key expected but found not string value,
// it will be added value of unnamed field. You can use whatever you want as values.
//
// See https://github.com/qioalice/ekago/internal/field/field.go ,
// https://github.com/qioalice/ekago/ekaerr/README.md for more details.
//
// Expert mode:
// If you want to minimise reflecting and increase performance, use Field objects
// to generate your fields and pass them as 'args'.
// It's private parts of ekago, but you can access them using ekaexp package.
// See https://github.com/qioalice/ekago/ekaexp/README.md for more details.
func (e *Error) AddFields(args ...interface{}) *Error {
	return e.addFields(args)
}

// ModifyBy calls f callback passing the current Error object into and returning
// the Error object, callback is return what.
// Nil safe.
//
// Does nothing if f == nil.
//
// Why?
// You can write your own fields appender like:
//     func (this *YourType) errAddIdentifiers(err *ekaerr.Error) *ekaerr.Error {
//         return err.AddFields("id", this.ID, "date", this.Date)
//     }
//
// And then use it like:
//     yt := new(YourType)
//     return ekaerr.IllegalState.
//         New("Unexpected state of world").
//         ModifyBy(yt.errAddIdentifiers).
//         Throw()
//
// Brilliant, isn't? And that most important it's so clean.
func (e *Error) ModifyBy(f func(err *Error) *Error) *Error {
	if e.IsValid() && f != nil {
		return f(e)
	}
	return e
}

// Is reports whether e has been instantiated by 'cls' Class's constructors.
// Returns false if either e is not valid Error or 'cls' is invalid.
// Nil safe.
func (e *Error) Is(cls Class) bool {
	return e.IsValid() && isValidClassID(cls.id) && e.classID == cls.id
}

// IsAny reports whether e belongs to at least one of passed 'cls' Classes
// (has been instantiated using one of them).
// Returns false if e is not valid Error or no one class has been passed.
// Nil-safe.
func (e *Error) IsAny(cls ...Class) bool {
	return e.is(cls, false)
}

// Of reports whether e has been instantiated by some Class that belongs to
// 'ns' Namespace. Returns false if either e is not valid or 'ns' is invalid.
// Nil safe.
func (e *Error) Of(ns Namespace) bool {
	return e.IsValid() && isValidNamespaceID(ns.id) && e.namespaceID == ns.id
}

// OfAny reports whether e belongs to at least one of passed 'nss' Namespaces
// (has been instantiated by the Class that belongs to one of 'nss' Namespace).
// Returns false if either e is not valid Error or no one namespace has been passed.
// Nil safe.
func (e *Error) OfAny(nss ...Namespace) bool {
	return e.of(nss)
}

// IsAnuDeep reports whether e belongs to at least one of passed 'cls' Classes
// or any of its parent (base) Classes is the same as one of passed 'cls'
// Returns false if e is not valid Error or no one class has been passed.
// Nil-safe.
//
// IsAnyDeep() has increased algorithmic complexity and MUCH SLOWER than IsAny()
// if you pass subclasses. So, make sure it's what you need.
func (e *Error) IsAnyDeep(cls ...Class) bool {
	return e.is(cls, true)
}

// Class returns e's Class. A special 'invalidClass' is returned if e == nil
// or has been manually instantiated instead of constructor using.
// Nil safe.
func (e *Error) Class() Class {
	if !e.IsValid() {
		return invalidClass
	}
	return classByID(e.classID, true)
}

// PublicMessage returns e's public message that you may set using SetPublicMessage().
// It's a simple abstraction to provide another message that you may show to the user.
// Returns "" if public message is not set or e is not valid Error.
// Nil safe.
func (e *Error) PublicMessage() string {
	if !e.IsValid() {
		return ""
	}
	return e.letter.SystemFields[_ERR_SYS_FIELD_IDX_PUBLIC_MESSAGE].SValue
}

// SetPublicMessage sets e's public message that you may get using PublicMessage().
// Keep in mind that public message IS NOT ATTACHED to the stack frame
// but to the whole Error.
// So, only one public message may be exist at once for each Error object.
// Nil safe.
func (e *Error) SetPublicMessage(newPublicMessage string) *Error {
	if e.IsValid() {
		e.letter.SystemFields[_ERR_SYS_FIELD_IDX_PUBLIC_MESSAGE].SValue =
			newPublicMessage
	}
	return e
}

// ID returns an unique e's ID as UUIDv4. You can tell this ID to the user and
// log this error. Then it will be easy to find an associated error.
// Returns "" if e is not valid Error.
// Nil safe.
func (e *Error) ID() string {
	if !e.IsValid() {
		return ""
	}
	return e.letter.SystemFields[_ERR_SYS_FIELD_IDX_ERROR_ID].SValue
}

// ReleaseError prepares 'err' for being reused in the future and releases
// its internal parts (returning them to the pool).
//
// YOU MUST NOT USE ERROR OBJECT AFTER PASSING THEM INTO THIS FUNCTION.
// YOU DO NOT NEED TO CALL THIS FUNCTION IF YOUR ERROR WILL BE LOGGED USING EKALOG.
func ReleaseError(err **Error) {
	if err != nil && (*err).IsValid() {
		releaseError(*err)
		*err = nil
	}
}
