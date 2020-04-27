// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package loge

import (
	"strings"
)

type (
	// Level represents the log message's level.
	// All constants this type is package private (see below),
	// but there is 'Levels' package-level struct, each field of that
	// represents one log level.
	Level uint8
)

//noinspection GoSnakeCaseUsage (Intellij IDEA suppress snake case warning).
const (
	// Levels represents the log message's levels.
	// You can use these constants to determine which log entry will you receive
	// in your Integrator.
	//
	// The int values represents each log level. Thus you can disable or enable all
	// lt/lte/gt/gte log levels than your desired in Integrator. These values:
	//   - Debug: 50,
	//   - Info: 60,
	//   - Warn: 70,
	//   - Error: 80,
	//   - Fatal: 90.
	// It has step == 10 to allow you to determine up to 9 your own custom,
	// intermediate log levels.

	/* 50 */ LEVEL_DEBUG Level = 50 + 10*iota
	/* 60 */ LEVEL_INFO
	/* 70 */ LEVEL_WARNING
	/* 80 */ LEVEL_ERROR
	/* 90 */ LEVEL_FATAL
)

var (
	// TODO: Thread-safety and controlling it if it's not necessary
	names       = make(map[Level]string)
	fatalLevels = make([]Level, 0, 10)
)

func init() {
	RegisterLevelName(LEVEL_DEBUG, "Debug")
	RegisterLevelName(LEVEL_INFO, "Info")
	RegisterLevelName(LEVEL_WARNING, "Warning")
	RegisterLevelName(LEVEL_ERROR, "Error")
	RegisterLevelName(LEVEL_FATAL, "Fatal")
}

// RegisterLevelName registers new log level's name that will be returned
// by Level.String() method. Allows you to overwrite standard log level names
// and name your own custom log levels. There is no-op if name is empty.
func RegisterLevelName(level Level, name string) {
	if name == "" {
		return
	}
	names[level] = name
}

// MarkLevelAsFatal marks passed level as level when you write log with,
// causes death.Die() (the same behaviour as standard Fatal handlers).
//
// You can undone it passing 'false' as second arg.
// All next arg or 'true' as 1st are ignored.
//
// WARNING!
// You can not mark standard levels as fatal.
func MarkLevelAsFatal(level Level, enable ...bool) {

	switch enable := !(len(enable) > 0 && !enable[0]); {
	case enable && (level == LEVEL_DEBUG || level == LEVEL_INFO || level == LEVEL_WARNING ||
		level == LEVEL_ERROR):
		// can not mark standard levels

	case level == LEVEL_FATAL:
		// already marked or can not be unmarked

	case enable && !level.mustDie():
		// not marked yet
		fatalLevels = append(fatalLevels, level)

	case !enable:
		for i, fatalLevel := range fatalLevels {
			if level == fatalLevel {
				fatalLevels = append(fatalLevels[:i], fatalLevels[i+1:]...)
				break
			}
		}
	}
}

// String returns a capitalized string representing the log level 'l'.
func (l Level) String() string {
	if name := names[l]; name != "" {
		return name
	} else {
		return "Unknown"
	}
}

// ToUpper returns an uppercase string representing the log level 'l'.
func (l Level) ToUpper() string { return strings.ToUpper(l.String()) }

// ToLower returns an lowercase string representing the log level 'l'.
func (l Level) ToLower() string { return strings.ToLower(l.String()) }

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
