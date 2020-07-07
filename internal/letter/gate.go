// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package letter

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
	// GLogErr is a function that is initialized in the ekalog package and used
	// in the ekaerr package.
	//
	// This function must log an *Error (using its internal private *Letter
	// 'errLetter') with log level 'level' using 'logger' as a untyped pointer
	// to the Logger object.
	GLogErr func(logger unsafe.Pointer, level uint8, errLetter *Letter)

	// GLogErrThroughDefaultLogger is a function that is initialized in the ekalog
	// package and used in the ekaerr package.
	//
	// This function must log an *Error (using its internal private *Letter
	// 'errLetter') with log level 'level' using default package's logger.
	GLogErrThroughDefaultLogger func(level uint8, errLetter *Letter)

	// GErrRelease is a function that is initialized in the ekaerr package
	// and used in the ekalog package.
	//
	// This function must return a 'errLetter' object as Error's *Letter object
	// to its pool for being reused in the future.
	GErrRelease func(errLetter *Letter)
)
