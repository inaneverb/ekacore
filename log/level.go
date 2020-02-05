// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

import "strings"

// Level represents the log message's level.
// All constants this type is package private (see below),
// but there is 'Levels' package-level struct, each field of that
// represents one log level.
type Level uint8

// Log level constants. Used in all internal functions.
// User should use 'Levels' fields.
const (
	/* 50 */ lvlDebug Level = 50 + 10*iota
	/* 60 */ lvlInfo
	/* 70 */ lvlWarning
	/* 80 */ lvlError
	/* 90 */ lvlFatal
)

// Levels represents log message's levels.
// You can use these constants to determine what log entry you will receive
// in your Integrator.
//
// There are int values that represents each log level.
// Thus you can disable or enable all lt/lte/gt/gte log levels than your desired
// in Integrator.
// These values:
// Debug: 50, Info: 60, Warn: 70, Error: 80, Fatal: 90.
// It has step == 10 to allow you to determine up to 9 your own custom log levels.
//
// WARNING!
// Since golang allows to overwrite struct fields, you shouldn't do this.
// Otherwise you will not have access to overwritten log levels anymore.
var (
	Levels = struct {
		Debug Level
		Info  Level
		Warn  Level
		Error Level
		Fatal Level
	}{
		Debug: lvlDebug,
		Info:  lvlInfo,
		Warn:  lvlWarning,
		Error: lvlError,
		Fatal: lvlFatal,
	}

	// TODO: Thread-safety and controlling it if it's not necessary
	names       = make(map[Level]string)
	fatalLevels = make([]Level, 0, 10)
)

func init() {
	RegisterLevelName(lvlDebug, "Debug")
	RegisterLevelName(lvlInfo, "Info")
	RegisterLevelName(lvlWarning, "Warning")
	RegisterLevelName(lvlError, "Error")
	RegisterLevelName(lvlFatal, "Fatal")
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
	case enable && (level == lvlDebug || level == lvlInfo || level == lvlWarning ||
		level == lvlError):
		// can not mark standard levels

	case level == lvlFatal:
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
	case lvlDebug, lvlInfo, lvlWarning, lvlError:
		return false

	case lvlFatal:
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

func (l Level) uint8() uint8 {
	return uint8(l)
}
