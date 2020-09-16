// Copyright © 2020. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/qioalice/ekago/v2/ekasys"
	"github.com/qioalice/ekago/v2/ekatyp"
	"github.com/qioalice/ekago/v2/internal/ekaletter"
)

//noinspection GoSnakeCaseUsage
const (
	// These constants are indexes of some Error's entities stored into
	// en embedded Letter's SystemFields array.
	// See https://github.com/qioalice/ekago/internal/letter/letter.go for more details.

	_ERR_SYS_FIELD_IDX_CLASS_ID       = 0
	_ERR_SYS_FIELD_IDX_CLASS_NAME     = 1
	_ERR_SYS_FIELD_IDX_PUBLIC_MESSAGE = 2
	_ERR_SYS_FIELD_IDX_ERROR_ID       = 3
)

// prepare prepares current Error for being used assuming that Error has been
// obtained from the Error's pool. Returns prepared Error.
func (e *Error) prepare() *Error {

	// We need to set flags to the defaults anyway, cause *LetterItem may be marked.
	for item := e.letter.Items; item != nil; item = ekaletter.GetNextItem(item) {
		item.Flags = DefaultFlags
	}

	// Because the main reason of Error existence is being logged later,
	// we need to make sure that it will be returned to the pool.
	if e.needSetFinalizer {
		runtime.SetFinalizer(e.letter, releaseErrorForFinalizer)
		e.needSetFinalizer = false
	}

	return e
}

// cleanup prepares internal *Letter object and its internal *LetterItem items
// for returning to the pool and being reused in the future.
func (e *Error) cleanup() *Error {

	// It guarantees that cleanup() will be called only at the head of the
	// *LetterItem linked list, and that *LetterItem linked list has been allocated
	// chunk by chunk and its capacity % _LETTER_ITEM_ALLOC_CHUNK_SIZE == 0.

	// First of all we need to release all unnecessary chunks
	// (if they were allocated and added).

	// It guarantees that Items != nil, and lastItem != nil because *Letter
	// always have preallocated *LetterItem chunks.

	item := e.letter.Items
	for i := int16(0); i < _LETTER_REUSE_MAX_LETTER_ITEMS && item != nil; i++ {

		for j := int16(0); j < _LETTER_ITEM_ALLOC_CHUNK_SIZE-1; j++ {
			ekaletter.ResetItem(item)
			item = ekaletter.GetNextItem(ekaletter.ResetItem(item))
		}

		next := ekaletter.GetNextItem(item)

		if i == _LETTER_REUSE_MAX_LETTER_ITEMS-1 {
			if next != nil {
				pruneLetterItemsChunk(next)
				ekaletter.SetNextItem(item, nil)
			}
		} else {
			item = next
		}
	}

	// We don't need to cleanup Namespace's ID or Class's ID,
	// because it will be overwritten and also do not need to update SystemFields
	// they will be overwritten too.

	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_PUBLIC_MESSAGE].SValue = ""

	e.letter.StackTrace = nil
	e.stackIdx = 0
	ekaletter.SetLastItem(e.letter, e.letter.Items)

	return e
}

// is reports whether e belongs to at least one of passed 'cls' Classes
// or to any of their parents (base) Classes (if 'deep' is true).
func (e *Error) is(cls []Class, deep bool) bool {

	if !e.IsValid() || len(cls) == 0 {
		return false
	}

	if deep {
		// lock once, do not lock each time at the classByID() call.
		registeredClassesMap.RLock()
		defer registeredClassesMap.RUnlock()
	}

	for i, n := 0, len(cls); i < n; i++ {
		if isValidClassID(cls[i].id) && e.classID == cls[i].id {
			return true
		}
		if deep {
			classID := cls[i].parentID
			for isValidClassID(classID) {
				if e.classID == classID {
					return true
				}
				// do not lock, already locked
				classID = classByID(classID, false).parentID
			}
		}
	}

	return false
}

// of reports whether e belongs to at least one of passed 'nss' Namespaces.
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

// getCurrentLetterItem returns a *LetterItem from e's *Letter for e's.stackIdx.
func (e *Error) getCurrentLetterItem() *ekaletter.LetterItem {

	lastItem := ekaletter.GetLastItem(e.letter)
	if e.stackIdx > -1 && lastItem.StackFrameIdx() == -1 {
		ekaletter.SetStackFrameIdx(lastItem, e.stackIdx)
	}

	if e.stackIdx > lastItem.StackFrameIdx() {
		// We need another *LetterItem,because now 'stackIdx' > last *LetterItem's idx.
		// Are preallocated *LetterItem s over?
		nextItem := ekaletter.GetNextItem(lastItem)

		if nextItem == nil {
			nextItem, _ = allocLetterItemsChunk()
			ekaletter.SetNextItem(lastItem, nextItem)
		}

		ekaletter.SetLastItem(e.letter, nextItem)
		ekaletter.SetStackFrameIdx(nextItem, e.stackIdx)
		lastItem = nextItem
	}

	return lastItem
}

// init is a part of newError() func (Error's constructor).
// Generates the stacktrace and an unique error's ID (UUID) saving it along with
// 'classID' and 'namespaceID' to the e and then returns it.
//
// Requirements:
// e != nil. Otherwise UB (may panic).
func (e *Error) init(classID ClassID, namespaceID NamespaceID) *Error {

	skip := 3 // init(), newError(), [Class.New(), Class.Wrap()]
	e.letter.StackTrace = ekasys.GetStackTrace(skip, -1).ExcludeInternal()

	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_CLASS_ID].IValue = int64(classID)
	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_CLASS_NAME].SValue =
		classByID(classID, true).fullName
	e.letter.SystemFields[_ERR_SYS_FIELD_IDX_ERROR_ID].SValue =
		ekatyp.UUID_NewV4_OrNil().String()

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
		e.getCurrentLetterItem().Message = baseMessage
	}

	return e
}

// newError is an Error's constructor.
// There are several steps:
//
// 1. Getting an *Error object (from the pool or allocate at the RAM heap).
//
// 2. Generate stacktrace, generate unique Error's ID, save it along with
//    provided 'classID' and 'namespaceID' into an Error object.
//
// 3. Initialize the first message using 'legacyErr' and 'message.
//
// 4. Parse passed 'args' and also add it as first stack frame's fields.
//
// 5. Mark first stack frame if generated message (p.3) is not empty.
//
// 6. Done.
func newError(classID ClassID, namespaceID NamespaceID, legacyErr error, message string, args []interface{}) *Error {

	err := acquireError().init(classID, namespaceID).construct(message, legacyErr).AddFields(args)
	if err.letter.Items.Message != "" || len(err.letter.Items.Fields) > 0 {
		err.Mark()
	}

	return err
}
