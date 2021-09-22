// Copyright © 2018-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"fmt"

	"github.com/qioalice/ekago/v3/ekaerr"
	"github.com/qioalice/ekago/v3/internal/ekaletter"
)

// Log writes log message with desired 'level',
// analyzing 'args' in the most powerful and smart way:
//
// - args[0] could be printf-like format string, then next N args will be
//   its printf values (N - num of format's printf verbs),
//   and M-N (where M is total count of args) will be treated as
//   explicit or implicit fields (depends on what kind of fields are allowed);
//
// - If only explicit fields are enabled (Options.OnlyExplicitFields(true)),
//   all args except explicit fields will be form log's message using fmt.Sprint
//   (the same way as LogStrict does), but you still can pass explicit fields also;
//
// - If both of explicit and implicit fields are enabled (by default),
//   only first non-explicit field arg forms log's message using "%+v" printf verb.
//   Others ones are treated as explicit and implicit (including unnamed) fields.
//
// Not bad, huh?
func (l *Logger) Log(level Level, args ...interface{}) (this *Logger) {
	return l.log(level, "", nil, args, nil)
}

// Logf writes log message with desired 'level', generating log message using
// fmt.Sprintf(format, args...) if 'format' != "" or fmt.Sprint(args...) otherwise.
//
// NOTICE!
// You can NOT add explicit/implicit fields using this method. And thus there is
// no reflections (usage of Golang 'reflect' package).
func (l *Logger) Logf(level Level, format string, args ...interface{}) (this *Logger) {
	return l.log(level, fmt.Sprintf(format, args...), nil, nil, nil)
}

func (l *Logger) Logw(level Level, msg string, fields ...ekaletter.LetterField) (this *Logger) {
	return l.log(level, msg, nil, nil, fields)
}
func (l *Logger) Logww(level Level, msg string, fields []ekaletter.LetterField) (this *Logger) {
	return l.log(level, msg, nil, nil, fields)
}

// ---------------------------------------------------------------------------- //

// Debug is the same as Log(LEVEL_DEBUG, args...).
// Read more: Logger.Log().
func (l *Logger) Debug(args ...interface{}) (this *Logger) {
	return l.log(LEVEL_DEBUG, "", nil, args, nil)
}

// Debugf is the same as Logf(LEVEL_DEBUG, format, args...).
// Read more: Logger.Logf().
func (l *Logger) Debugf(format string, args ...interface{}) (this *Logger) {
	return l.log(LEVEL_DEBUG, fmt.Sprintf(format, args...), nil, nil, nil)
}

func (l *Logger) Debugw(msg string, fields ...ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_DEBUG, msg, nil, nil, fields)
}
func (l *Logger) Debugww(msg string, fields []ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_DEBUG, msg, nil, nil, fields)
}

// ---------------------------------------------------------------------------- //

// Info is the same as Log(LEVEL_INFO, args...).
// Read more: Logger.Log().
func (l *Logger) Info(args ...interface{}) (this *Logger) {
	return l.log(LEVEL_INFO, "", nil, args, nil)
}

// Infof is the same as Logf(LEVEL_INFO, format, args...).
// Read more: Logger.Logf().
func (l *Logger) Infof(format string, args ...interface{}) (this *Logger) {
	return l.log(LEVEL_INFO, fmt.Sprintf(format, args...), nil, nil, nil)
}

func (l *Logger) Infow(msg string, fields ...ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_INFO, msg, nil, nil, fields)
}
func (l *Logger) Infoww(msg string, fields []ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_INFO, msg, nil, nil, fields)
}

// ---------------------------------------------------------------------------- //

// Notice is the same as Log(LEVEL_NOTICE, args...).
// Read more: Logger.Log().
func (l *Logger) Notice(args ...interface{}) (this *Logger) {
	return l.log(LEVEL_NOTICE, "", nil, args, nil)
}

// Noticef is the same as Logf(LEVEL_NOTICE, format, args...).
// Read more: Logger.Logf().
func (l *Logger) Noticef(format string, args ...interface{}) (this *Logger) {
	return l.log(LEVEL_NOTICE, fmt.Sprintf(format, args...), nil, nil, nil)
}

func (l *Logger) Noticew(msg string, fields ...ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_NOTICE, msg, nil, nil, fields)
}
func (l *Logger) Noticeww(msg string, fields []ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_NOTICE, msg, nil, nil, fields)
}

// ---------------------------------------------------------------------------- //

// Warn is the same as Log(LEVEL_WARNING, args...).
// Read more: Logger.Log().
func (l *Logger) Warn(args ...interface{}) (this *Logger) {
	return l.log(LEVEL_WARNING, "", nil, args, nil)
}

// Warnf is the same as Logf(LEVEL_WARNING, format, args...).
// Read more: Logger.Logf().
func (l *Logger) Warnf(format string, args ...interface{}) (this *Logger) {
	return l.log(LEVEL_WARNING, fmt.Sprintf(format, args...), nil, nil, nil)
}

func (l *Logger) Warnw(msg string, fields ...ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_WARNING, msg, nil, nil, fields)
}
func (l *Logger) Warnww(msg string, fields []ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_WARNING, msg, nil, nil, fields)
}

func (l *Logger) Warne(msg string, err *ekaerr.Error, kvFields ...interface{}) (this *Logger) {
	return l.log(LEVEL_WARNING, msg, err, kvFields, nil)
}
func (l *Logger) Warnew(msg string, err *ekaerr.Error, fields ...ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_WARNING, msg, err, nil, fields)
}
func (l *Logger) Warneww(msg string, err *ekaerr.Error, fields []ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_WARNING, msg, err, nil, fields)
}

// ---------------------------------------------------------------------------- //

// Error is the same as Log(LEVEL_ERROR, args...).
// Read more: Logger.Log().
func (l *Logger) Error(args ...interface{}) (this *Logger) {
	return l.log(LEVEL_ERROR, "", nil, args, nil)
}

// Errorf is the same as Logf(LEVEL_ERROR, format, args...).
// Read more: Logger.Logf().
func (l *Logger) Errorf(format string, args ...interface{}) (this *Logger) {
	return l.log(LEVEL_ERROR, fmt.Sprintf(format, args...), nil, nil, nil)
}

func (l *Logger) Errorw(msg string, fields ...ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_ERROR, msg, nil, nil, fields)
}
func (l *Logger) Errorww(msg string, fields []ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_ERROR, msg, nil, nil, fields)
}

func (l *Logger) Errore(msg string, err *ekaerr.Error, kvFields ...interface{}) (this *Logger) {
	return l.log(LEVEL_ERROR, msg, err, kvFields, nil)
}
func (l *Logger) Errorew(msg string, err *ekaerr.Error, fields ...ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_ERROR, msg, err, nil, fields)
}
func (l *Logger) Erroreww(msg string, err *ekaerr.Error, fields []ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_ERROR, msg, err, nil, fields)
}

// ---------------------------------------------------------------------------- //

// Crit is the same as Log(LEVEL_CRITICAL, args...).
// Read more: Logger.Log().
func (l *Logger) Crit(args ...interface{}) (this *Logger) {
	return l.log(LEVEL_CRITICAL, "", nil, args, nil)
}

// Critf is the same as Logf(LEVEL_CRITICAL, format, args...).
// Read more: Logger.Logf().
func (l *Logger) Critf(format string, args ...interface{}) (this *Logger) {
	return l.log(LEVEL_CRITICAL, fmt.Sprintf(format, args...), nil, nil, nil)
}

func (l *Logger) Critw(msg string, fields ...ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_CRITICAL, msg, nil, nil, fields)
}
func (l *Logger) Critww(msg string, fields []ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_CRITICAL, msg, nil, nil, fields)
}

func (l *Logger) Crite(msg string, err *ekaerr.Error, kvFields ...interface{}) (this *Logger) {
	return l.log(LEVEL_CRITICAL, msg, err, kvFields, nil)
}
func (l *Logger) Critew(msg string, err *ekaerr.Error, fields ...ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_CRITICAL, msg, err, nil, fields)
}
func (l *Logger) Criteww(msg string, err *ekaerr.Error, fields []ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_CRITICAL, msg, err, nil, fields)
}

// ---------------------------------------------------------------------------- //

// Alert is the same as Log(LEVEL_ALERT, args...).
// Read more: Logger.Log().
func (l *Logger) Alert(args ...interface{}) (this *Logger) {
	return l.log(LEVEL_ALERT, "", nil, args, nil)
}

// Alertf is the same as Logf(LEVEL_ALERT, format, args...).
// Read more: Logger.Logf().
func (l *Logger) Alertf(format string, args ...interface{}) (this *Logger) {
	return l.log(LEVEL_ALERT, fmt.Sprintf(format, args...), nil, nil, nil)
}

func (l *Logger) Alertw(msg string, fields ...ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_ALERT, msg, nil, nil, fields)
}
func (l *Logger) Alertww(msg string, fields []ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_ALERT, msg, nil, nil, fields)
}

func (l *Logger) Alerte(msg string, err *ekaerr.Error, kvFields ...interface{}) (this *Logger) {
	return l.log(LEVEL_ALERT, msg, err, kvFields, nil)
}
func (l *Logger) Alertew(msg string, err *ekaerr.Error, fields ...ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_ALERT, msg, err, nil, fields)
}
func (l *Logger) Alerteww(msg string, err *ekaerr.Error, fields []ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_ALERT, msg, err, nil, fields)
}

// ---------------------------------------------------------------------------- //

// Emerg is the same as Log(LEVEL_EMERGENCY, args...),
// but also then calls ekadeath.Die(1).
// Read more: Logger.Log().
func (l *Logger) Emerg(args ...interface{}) (this *Logger) {
	return l.log(LEVEL_EMERGENCY, "", nil, args, nil)
}

// Emergf is the same as Logf(LEVEL_EMERGENCY, format, args...),
// but also then calls death.Die(1).
// Read more: Logger.Logf().
func (l *Logger) Emergf(format string, args ...interface{}) (this *Logger) {
	return l.log(LEVEL_EMERGENCY, fmt.Sprintf(format, args...), nil, nil, nil)
}

func (l *Logger) Emergw(msg string, fields ...ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_EMERGENCY, msg, nil, nil, fields)
}
func (l *Logger) Emergww(msg string, fields []ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_EMERGENCY, msg, nil, nil, fields)
}

func (l *Logger) Emerge(msg string, err *ekaerr.Error, kvFields ...interface{}) (this *Logger) {
	return l.log(LEVEL_EMERGENCY, msg, err, kvFields, nil)
}
func (l *Logger) Emergew(msg string, err *ekaerr.Error, fields ...ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_EMERGENCY, msg, err, nil, fields)
}
func (l *Logger) Emergeww(msg string, err *ekaerr.Error, fields []ekaletter.LetterField) (this *Logger) {
	return l.log(LEVEL_EMERGENCY, msg, err, nil, fields)
}
