// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaletter

import (
	"github.com/qioalice/ekago/v2/ekamath"
)

type (
	// Flag is just ekamath.Flag16 alias and has been introduced to make it easy
	// to replace underlying type to something other in the future.
	Flag = ekamath.Flag16
)

//noinspection GoSnakeCaseUsage
const (

	// These flags are used to determine the behaviour of *Letter, *LetterItem
	// and its linked list.

	// FLAG_MARKED_LETTER_ITEM means that this *LetterItem is marked
	// (maybe it's important one? maybe contain important message or/and field(s)?).
	FLAG_MARKED_LETTER_ITEM Flag = 0x0001

	// FLAG_ONLY_EXPLICIT_FIELDS means that only explicit key-value paired
	// arguments (called fields) were (or may) used for this *LetterItem.
	//
	// See: https://github.com/qioalice/ekago/internal/field/field.go ,
	// https://github.com/qioalice/ekago/ekaexp/README.md for more info.
	FLAG_ONLY_EXPLICIT_FIELDS Flag = 0x0002

	// FLAG_ALLOW_UNNAMED_NIL means that this *LetterItem accept nil unnamed fields
	// (key-value pair in which key is empty and the value is untyped nil).
	FLAG_ALLOW_UNNAMED_NIL Flag = 0x0004

	// FLAG_ALLOW_IMPLICIT_POINTERS means that this *LetterItem will save
	// field that represents a pointer as is w/o attempting to dereference and
	// save a value that pointer points to.
	FLAG_ALLOW_IMPLICIT_POINTERS Flag = 0x0008
)
