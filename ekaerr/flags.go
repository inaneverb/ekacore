// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"github.com/qioalice/ekago/v2/internal/ekaletter"
)

//noinspection GoSnakeCaseUsage
const (

	// These flags are used to determine the behaviour of *Letter, *LetterItem
	// and its linked list.

	// FLAG_MARKED_LETTER_ITEM means that this *LetterItem is marked
	// (maybe it's important one? maybe contain important message or/and field(s)?).
	FLAG_MARKED_LETTER_ITEM ekaletter.Flag = ekaletter.FLAG_MARKED_LETTER_ITEM

	// FLAG_ONLY_EXPLICIT_FIELDS means that only explicit key-value paired
	// arguments (called fields) were (or may) used for this *LetterItem.
	//
	// See: https://github.com/qioalice/ekago/internal/field/field.go ,
	// https://github.com/qioalice/ekago/ekaexp/README.md for more info.
	FLAG_ONLY_EXPLICIT_FIELDS = ekaletter.FLAG_ONLY_EXPLICIT_FIELDS

	// FLAG_ALLOW_UNNAMED_NIL means that this *LetterItem accept nil unnamed fields
	// (key-value pair in which key is empty and the value is untyped nil).
	FLAG_ALLOW_UNNAMED_NIL = ekaletter.FLAG_ALLOW_UNNAMED_NIL

	// FLAG_ALLOW_IMPLICIT_POINTERS means that this *LetterItem will save
	// field that represents a pointer as is w/o attempting to dereference and
	// save a value that pointer points to.
	FLAG_ALLOW_IMPLICIT_POINTERS = ekaletter.FLAG_ALLOW_IMPLICIT_POINTERS
)
