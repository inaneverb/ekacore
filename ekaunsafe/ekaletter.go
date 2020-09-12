// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe

import (
	"unsafe"

	"github.com/qioalice/ekago/v2/internal/ekafield"
	"github.com/qioalice/ekago/v2/internal/ekaletter"
)

type (
	// https://github.com/qioalice/ekago/internal/ekaletter/letter.go

	Letter = ekaletter.Letter

	// https://github.com/qioalice/ekago/internal/ekaletter/letter_item.go

	LetterItem = ekaletter.LetterItem
)

// https://github.com/qioalice/ekago/internal/ekaletter/letter.go

func LetterSetLastItem(l *Letter, li *LetterItem) *Letter {
	return ekaletter.SetLastItem(l, li)
}

func LetterGetLastItem(l *Letter) *LetterItem {
	return ekaletter.GetLastItem(l)
}

func LetterSetSomething(l *Letter, ptr unsafe.Pointer) *Letter {
	return ekaletter.SetSomething(l, ptr)
}

func LetterGetSomething(l *Letter) unsafe.Pointer {
	return ekaletter.GetSomething(l)
}

// https://github.com/qioalice/ekago/internal/ekaletter/letter_item.go

func LetterItemParseTo(li *LetterItem, args []interface{}, explicitFields []ekafield.Field, onlyFields bool) {
	ekaletter.ParseTo(li, args, explicitFields, onlyFields)
}

func LetterItemSetNext(current, next *LetterItem) *LetterItem {
	return ekaletter.SetNextItem(current, next)
}

func LetterItemGetNext(current *LetterItem) *LetterItem {
	return ekaletter.GetNextItem(current)
}

func LetterItemSetStackFrameIdx(li *LetterItem, idx int16) *LetterItem {
	return ekaletter.SetStackFrameIdx(li, idx)
}

func LetterItemReset(li *LetterItem) *LetterItem {
	return ekaletter.ResetItem(li)
}
