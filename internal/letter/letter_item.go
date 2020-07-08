// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package letter

import (
	"github.com/qioalice/ekago/v2/internal/field"
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
