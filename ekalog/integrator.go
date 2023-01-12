// Copyright © 2018-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"github.com/qioalice/ekago/v4/internal/ekaletter"
)

// Integrator is the way, an Entry must be encoded how and written to.
//
// It also allows to drop an Entry at the earlier steps if it's Level less than
// level is returned by MinLevelEnabled().
//
// Main two methods are PreEncodeField() that allows to encode some fields once
// and then use them for each Entry.
type Integrator interface {

	// PreEncodeField allows to encode ekaletter.LetterField once
	// and then use them for each Entry.
	// Using this method you will spent less resources at the ekaletter.LetterField's
	// encoding, doing it only once for each that field.
	PreEncodeField(f ekaletter.LetterField)

	// EncodeAndWrite encodes Entry and then writes to some destination
	// (integrator determines what and how it will be).
	//
	// Thus, EncodeAndWrite does the main thing of Integrator:
	// "Integrates your log messages with your log destination service".
	//
	// WARNING.
	// Integrator MUST NOT to hold an Entry or its parts after this method is done.
	// The Logger will return Entry to its pool after this method is complete,
	// so accessing to it or its parts is unsafe and WILL lead to UB.
	EncodeAndWrite(entry *Entry)

	// MinLevelEnabled returns minimum Level an integrator will handle Entry with.
	// E.g. if minimum level is LEVEL_WARNING, LEVEL_DEBUG logs will be dropped
	// at the most earlier moment - even before parsing.
	MinLevelEnabled() Level

	// MinLevelForStackTrace must return a minimum level starting with a stacktrace
	// must be generated and added to the Logger's Entry only if it's not presented
	// yet by attached ekaerr.Error object.
	MinLevelForStackTrace() Level

	// Sync flushes all pending log entries to integrator destination.
	// It useful when integrator does async work and sometimes you need to make sure
	// all pending entries are flushed.
	//
	// Logger type has the same name's method that just calls this method.
	Sync() error
}
