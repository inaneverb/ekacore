// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"unsafe"

	"github.com/qioalice/ekago/v4/internal/ekaletter"
)

// -----
// Functions with names started with "bridge..." are not used in this package
// but assigns to the bridge
// ( https://github.com/qioalice/ekago/internal/letter/bridge.go )
// and used at the ekaunsafe package's ekaerr related functions
// ( https://github.com/qioalice/ekago/ekaunsafe/ekaerr.go )
// -----

// bridgeGetLetter return *ekaletter.Letter object from the *ekaerr.Error object
// assuming that 'err' is an untyped pointer to the ekaerr.Error.
// Returns nil if 'err' == nil.
func bridgeGetLetter(err unsafe.Pointer) *ekaletter.Letter {
	if err := (*Error)(err); err.IsValid() {
		return err.letter
	} else {
		return nil
	}
}

// bridgeGetStackIdx returns current stackIdx of the *ekaerr.Error object
// assuming that 'err' is an untyped pointer to the ekaerr.Error.
func bridgeGetStackIdx(err unsafe.Pointer) int16 {
	if err := (*Error)(err); err.IsValid() {
		return ekaletter.LGetStackIdx(err.letter)
	} else {
		return -1
	}
}

// bridgeSetStackIdx sets the new value of the *ekaerr.Error object's 'stackIdx' field
// assuming that 'err' is an untyped pointer to the ekaerr.Error.
// New stack index must be 0 or greater.
//
// IT MUST BE USED CAREFULLY!
// YOU CAN GET PANIC IF YOUR STACK IDX WILL BE GREATER THAN AN EMBEDDED STACKTRACE'S LENGTH.
func bridgeSetStackIdx(err unsafe.Pointer, newStackIdx int16) {
	if err := (*Error)(err); err.IsValid() && newStackIdx >= 0 {
		ekaletter.LSetStackIdx(err.letter, newStackIdx)
	}
}
