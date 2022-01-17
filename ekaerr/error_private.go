// Copyright © 2020. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/qioalice/ekago/v3/ekasys"
	"github.com/qioalice/ekago/v3/ekatyp"
	"github.com/qioalice/ekago/v3/internal/ekaletter"
)

//noinspection GoSnakeCaseUsage
const (
	// These constants are indexes of some Error's entities stored into
	// en embedded Letter's SystemFields array.
	// See https://github.com/qioalice/ekago/internal/letter/letter.go for more details.

	_ERR_SYS_FIELD_IDX_CLASS_ID   = 0
	_ERR_SYS_FIELD_IDX_CLASS_NAME = 1
	_ERR_SYS_FIELD_IDX_ERROR_ID   = 2
)

// prepare prepares current Error for being used assuming that Error has been
// obtained from the Error's pool. Returns prepared Error.
func (e *Error) prepare() *Error {

	// Because the main reason of Error existence is being logged later,
	// we need to make sure that it will be returned to the pool.
	if e.needSetFinalizer {
		runtime.SetFinalizer(e, releaseErrorForFinalizer)
		e.needSetFinalizer = false
	}

	return e
}

// cleanup frees all allocated resources (RAM in 99% cases) by Error, preparing
// it for being returned to the pool and being reused in the future.
func (e *Error) cleanup() *Error {

	// We don't need to cleanup Namespace's ID or Class's ID,
	// because it will be overwritten and also do not need to update SystemFields
	// they will be overwritten too.

	e.letter.StackTrace = nil

	ekaletter.LReset(e.letter)
	return e
}

// addField checks whether Error is valid and adds an ekaletter.LetterField
// to current Error, if field is addable.
func (e *Error) addField(f ekaletter.LetterField) *Error {
	if e.IsValid() {
		ekaletter.LAddFieldWithCheck(e.letter, f)
	}
	return e
}

// addFields is the same as addField() but works with an array of ekaletter.LetterField.
func (e *Error) addFields(fs []ekaletter.LetterField) *Error {
	if e.IsValid() {
		for i, n := 0, len(fs); i < n; i++ {
			ekaletter.LAddFieldWithCheck(e.letter, fs[i])
		}
	}
	return e
}

// addFieldsParse creates a ekaletter.LetterField objects based on passed values,
// try to treating them as a key-value pairs of that fields.
// Then adds generated ekaletter.LetterField to the Error only if those fields are addable.
func (e *Error) addFieldsParse(fs []interface{}, onlyFields bool) *Error {
	if e.IsValid() && len(fs) > 0 {
		ekaletter.LParseTo(e.letter, fs, onlyFields)
	}
	return e
}

// is reports whether e belongs to at least one of passed cls Class
// or any of Error's parent (base) Class is the same as one of passed (if deep is true).
func (e *Error) is(cls []Class, deep bool) bool {

	if !e.IsValid() || len(cls) == 0 {
		return false
	}

	if deep {
		// lock once, do not lock each time at the classByID() call.
		registeredClassesMap.RLock()
		defer registeredClassesMap.RUnlock()
	}

	n := len(cls)
	for classID := e.classID; isValidClassID(classID); {

		for i := 0; i < n; i++ {
			if cls[i].id == classID {
				return true
			}
		}

		if deep {
			// do not lock, already locked
			classID = classByID(classID, false).parentID
		} else {
			classID = _ERR_INVALID_CLASS_ID
		}
	}

	return false
}

// of reports whether Error belongs to at least one of passed nss Namespace.
func (e *Error) of(nss []Namespace) bool {

	if !e.IsValid() || len(nss) == 0 {
		return false
	}

	for i, n := 0, len(nss); i < n; i++ {
		if isValidNamespaceID(nss[i].id) && e.namespaceID == nss[i].id {
			return true
		}
	}

	return false
}

// init is a part of newError() func (Error's constructor).
// Generates the stacktrace and an unique error's ID (ULID) saving it along with
// classID and namespaceID to the Error and then returns it.
func (e *Error) init(classID ClassID, namespaceID NamespaceID, lightweight bool) *Error {

	skip := 3 // init(), newError(), [Class.New(), Class.Wrap(), Class.LightNew(), Class.LightWrap()]

	if !lightweight {
		e.letter.StackTrace = ekasys.GetStackTrace(skip, -1).ExcludeInternal()
	}

	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_CLASS_ID].IValue = int64(classID)
	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_CLASS_NAME].SValue =
		classByID(classID, true).fullName
	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_ERROR_ID].SValue =
		ekatyp.ULID_New_OrNil().String()

	e.classID = classID
	e.namespaceID = namespaceID

	return e
}

// construct is a part of newError() func (Error's constructor).
// Must be called after init() call. Builds first e's stack frame's message basing on
// passed 'baseMessage' and 'legacyErr'.
func (e *Error) construct(baseMessage string, legacyErr error) *Error {

	baseMessage = strings.TrimSpace(baseMessage)
	legacyErrStr := ""

	if legacyErr != nil {
		legacyErrStr = strings.TrimSpace(legacyErr.Error())
	}

	// isSkipCharByte is for ASCII strings and reports whether 'b' char must be ignored
	// or not while building string based on 'baseMessage' and 'legacyErr'.
	isSkipCharByte := func(b byte) bool {
		return b <= 32 || b == '!' || b == '?' || b == '.' || b == ',' || b == '-'
	}
	// isSkipCharByte is for UTF8 strings and reports whether 'r' rune must be ignored
	// or not while building string based on 'baseMessage' and 'legacyErr'.
	isSkipCharRune := func(r rune) bool {
		return r <= 32 || r == '!' || r == '?' || r == '.' || r == ',' || r == '-'
	}

	switch {
	case legacyErrStr != "" && baseMessage != "":
		l, rl := len(baseMessage), utf8.RuneCountInString(baseMessage)
		if l == rl {
			// 'baseMessage' is ASCII, fast path
			for rl--; rl >= 0 && isSkipCharByte(baseMessage[rl]); rl-- {
			}
			if rl > -1 {
				baseMessage = baseMessage[:rl+1]
			} else {
				baseMessage = ""
			}
		} else {
			// 'baseMessage' is UTF8, slow path.
			// The algorithm is:
			// 1. Reverse 'baseMessage' (because we can't access UTF8 chars at index)
			// 2. Find you where we can start from
			// 3. Cut the string
			// 4. Reverse back

			// 1. Reverse 'baseMessage'
			b := make([]rune, rl) // 'baseMessage' work buffer (reversed bytes atm)
			for _, rune_ := range baseMessage {
				rl--
				b[rl] = rune_
			}
			// 2. Find out where we can start from
			i := 0
			rl = len(b)
			for ; i < rl && isSkipCharRune(b[i]); i++ {
			}

			// 4. Reverse back but only starting from significant chars
			if i < rl {
				rl--
				for j := i; j < rl; {
					b[j], b[rl] = b[rl], b[j]
					j++
					rl--
				}
				// 3. Cut the string, transform types
				baseMessage = string(b[i:])
			} else {
				// the whole string must be ignored (cut)
				baseMessage = ""
			}
		}
		if baseMessage != "" {
			baseMessage += ", cause: " + legacyErrStr + "."
			break
		}
		fallthrough

	case legacyErrStr != "" && baseMessage == "":
		baseMessage = legacyErrStr
	}

	if baseMessage != "" {
		ekaletter.LSetMessage(e.letter, baseMessage, true)
	}

	return e
}

// newError is an Error's constructor.
// There are several steps:
//
//  1. Getting an *Error object (from the pool or allocate at the RAM heap).
//  2. Generate stacktrace, generate unique Error's ID, save it along with
//     provided 'classID' and 'namespaceID' into an Error object.
//  3. Initialize the first message using 'legacyErr' and 'message.
//  4. Parse passed 'args' and also add it as first stack frame's fields.
//  5. Mark first stack frame if generated message (p.3) is not empty.
func newError(

	lightweight bool,
	classID ClassID, namespaceID NamespaceID,
	legacyErr error, message string, args []interface{},

) *Error {

	return acquireError().
		init(classID, namespaceID, lightweight).
		construct(message, legacyErr).
		addFieldsParse(args, false)
}
