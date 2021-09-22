// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"time"

	"github.com/qioalice/ekago/v3/internal/ekaletter"
)

type (
	// Entry is a Logger's object that contains all data about one log event.
	//
	// Entry contains time, message, fields, attached error (if any), level
	// and a Logger using which this Entry must be logged.
	//
	// You MUST NOT instantiate this object manually.
	// It's exposed only for custom encoders.
	Entry struct {
		l *Logger // The logger this entry created by or belongs to.

		// LogLetter contains message, fields, stacktrace of Entry.
		// Public because of providing access from the user's Integrator.
		LogLetter *ekaletter.Letter

		// ErrLetter is attached error's letter. Contains its stacktrace, message, fields.
		// Public because of providing access from the user's Integrator.
		ErrLetter *ekaletter.Letter

		// Level indicates how important occurred event (this log entry represents) is.
		Level Level

		// Time contains the time when an event occurred this log entry represents.
		// Generated automatically by time.Now() call in log finisher.
		Time time.Time

		needSetFinalizer bool
	}
)
