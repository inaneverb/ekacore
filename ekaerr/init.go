// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"github.com/qioalice/ekago/v2/internal/ekaletter"

	"github.com/modern-go/reflect2"
)

func init() {
	// Create first N *Error objects and fill its pool by them.
	initErrorPool()

	ekaletter.BridgeErrorGetLetter = bridgeGetLetter

	ekaletter.BridgeErrorGetStackIdx = bridgeGetStackIdx
	ekaletter.BridgeErrorSetStackIdx = bridgeSetStackIdx

	// Initialize the gate's functions to link ekalog <-> ekaerr packages.
	ekaletter.GErrRelease = releaseErrorForGate

	// It's prohibited to use some types as Error's fields.
	//
	// See letter.ParseTo() for more details at the
	// https://github.com/qioalice/ekago/internal/letter/letter.go
	ekaletter.TypesBeingIgnoredForParsing = append(
		ekaletter.TypesBeingIgnoredForParsing,

		reflect2.TypeOf(ekaletter.Letter{}),
		reflect2.TypeOf((*ekaletter.Letter)(nil)),

		reflect2.TypeOf(ekaletter.LetterItem{}),
		reflect2.TypeOf((*ekaletter.LetterItem)(nil)),

		reflect2.TypeOf(Error{}),
		reflect2.TypeOf((*Error)(nil)),
	)
}
