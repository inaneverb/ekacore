// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package letter

import (
	"unsafe"

	"github.com/qioalice/ekago/v2/internal/field"
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

	// BridgeLogErr2 and BridgeLogwErr2 are a functions that are initialized
	// in the ekalog package and used in the ekaerr package.
	//
	// These functions must log an *ekaerr.Error
	// (using its internal private *Letter 'errLetter') with log level 'level',
	// using 'logger' as untyped pointer to the *ekalog.Logger object
	// (or keep it nil if standard package's level logger must be used).

	BridgeLogErr2 func(logger unsafe.Pointer, level uint8, errLetter *Letter, errArgs []interface{})
	BridgeLogwErr2 func(logger unsafe.Pointer, level uint8, errLetter *Letter, errMessage string, errFields []field.Field)

	// GErrRelease is a function that is initialized in the ekaerr package
	// and used in the ekalog package.
	//
	// This function must return a 'errLetter' object as Error's *Letter object
	// to its pool for being reused in the future.
	GErrRelease func(errLetter *Letter)
)
