// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"strings"
)

type (
	// Level is the log message's severity level.
	// There are 7 log levels the same as used in syslog.
	//
	// You can use log constants to determine
	// which log entry will you receive in your Integrator.
	//
	// Keep in mind, logging using LEVEL_EMERGENCY cause calling ekadeath.Die(),
	// after writing a log message.
	//
	// Read more:
	// https://en.wikipedia.org/wiki/Syslog
	//
	Level uint8
)

//noinspection GoSnakeCaseUsage
const (
	LEVEL_EMERGENCY Level = iota
	LEVEL_ALERT
	LEVEL_CRITICAL
	LEVEL_ERROR
	LEVEL_WARNING
	LEVEL_NOTICE
	LEVEL_INFO
	LEVEL_DEBUG
)

// String returns a capitalized string of the current log level.
// Returns an empty string if it's unexpected log level.
func (l Level) String() string {
	switch l {
	case LEVEL_EMERGENCY: return "Emergency"
	case LEVEL_ALERT:     return "Alert"
	case LEVEL_CRITICAL:  return "Critical"
	case LEVEL_ERROR:     return "Error"
	case LEVEL_WARNING:   return "Warning"
	case LEVEL_NOTICE:    return "Notice"
	case LEVEL_INFO:      return "Info"
	case LEVEL_DEBUG:     return "Debug"
	default:              return ""
	}
}

// String3 returns a capitalized short-hand string of the current log level.
// It has a length of 3 chars for all but LEVEL_EMERGENCY. Takes 5 for that.
// Returns an empty string if it's unexpected log level.
func (l Level) String3() string {
	switch l {
	case LEVEL_EMERGENCY: return "Emerg"
	case LEVEL_ALERT:     return "Ale"
	case LEVEL_CRITICAL:  return "Cri"
	case LEVEL_ERROR:     return "Err"
	case LEVEL_WARNING:   return "War"
	case LEVEL_NOTICE:    return "Noe"
	case LEVEL_INFO:      return "Inf"
	case LEVEL_DEBUG:     return "Deb"
	default:              return ""
	}
}

// ToUpper returns an uppercase variant of String() call.
func (l Level) ToUpper() string { return strings.ToUpper(l.String()) }

// ToLower returns a lowercase variant of String() call.
func (l Level) ToLower() string { return strings.ToLower(l.String()) }

// ToUpper3 returns an uppercase variant of String3() call.
func (l Level) ToUpper3() string { return strings.ToUpper(l.String3()) }

// ToLower3 returns an uppercase variant of String3() call.
func (l Level) ToLower3() string { return strings.ToLower(l.String3()) }
