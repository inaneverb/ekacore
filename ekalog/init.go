// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"github.com/qioalice/ekago/v2/internal/ekaletter"
)

func init() {

	// Init log levels.
	initLevels()

	// Create first N *Error objects and fill its pool by them.
	initEntryPool()

	// Initialize package's Logger (first base logger).
	initBaseLogger()

	// Initialize the gate's functions to link ekalog <-> ekaerr packages.
	ekaletter.BridgeLogErr2 = logErr
	ekaletter.BridgeLogwErr2 = logErrw
}
