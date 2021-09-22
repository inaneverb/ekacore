// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaletter

import (
	"unsafe"
)

// It's a special file that contains gate functions.
//
// So, we need to import 'ekalog' -> 'ekaerr' and vice-versa, but cross imports
// are prohibited by Go rules. So, it's a simple hack. We initializing these
// functions at the packages' init() calls and then use them and use unsafe.Pointer
// as arguments to describe arguments but do not specify types.

var (

	// BridgeErrorGetLetter is a function that is initialized
	// in the ekaerr package and used in the ekaunsafe package.
	//
	// This function must return an underlying *Letter object
	// from the 'err' - *ekaerr.Error object.
	BridgeErrorGetLetter func(err unsafe.Pointer) *Letter

	// BridgeErrorGetStackIdx, BridgeErrorSetStackIdx are a functions that are initialized
	// in the ekaerr package and used in the ekaunsafe package.
	//
	// These functions are *ekaerr.Error object's 'stackIdx' field getter/setter.

	BridgeErrorGetStackIdx func(err unsafe.Pointer) int16
	BridgeErrorSetStackIdx func(err unsafe.Pointer, newStackIdx int16)
)
