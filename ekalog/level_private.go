// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
)

var (
	// TODO: Thread-safety and controlling it if it's not necessary
	names       = make(map[Level]string)
	fatalLevels = make([]Level, 0, 10)

	registeredNewLevels uint32
)


// mustDie reports whether l shall cause death.Die() call.
func (l Level) mustDie() bool {

	switch l {
	case LEVEL_DEBUG, LEVEL_INFO, LEVEL_WARNING, LEVEL_ERROR:
		return false

	case LEVEL_FATAL:
		return true

	default:
		for _, fatalLevel := range fatalLevels {
			if l == fatalLevel {
				return true
			}
		}
		return false
	}
}

// initLevels initializes standard log levels, their names.
func initLevels() {
	RegisterLevelName(LEVEL_DEBUG, "Debug")
	RegisterLevelName(LEVEL_INFO, "Info")
	RegisterLevelName(LEVEL_WARNING, "Warning")
	RegisterLevelName(LEVEL_ERROR, "Error")
	RegisterLevelName(LEVEL_FATAL, "Fatal")

	// RegisterLevelName() increments registeredNewLevels, but there were
	// a default log levels. Overwrite to 0.
	registeredNewLevels = 0
}