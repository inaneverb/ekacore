// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"fmt"

	"github.com/qioalice/ekago/internal/field"
)

// -----
// This file contains only Logger's finishers:
// Those methods that really generates log message's body and starts
// a process of writing'em.
//
// They moved to a separately file because of bleeding my eyes,
// while I trying to work with another Logger's methods.
//
// Separated by 3 spaces into log level categories.
//
// Generally, there are common, regular methods, that are really documented
// and could be used to write log message with user defined log level.
//
// All the rest methods are just "aliases" to most used (and useful) log levels.
//
// TODO: Supporting Group by finishers
// TODO: Supporting Options by finishers
// TODO: Logger.Logec, Logger.LogecStrict support printf string (not only error)
// TODO: IfErr func implement (the same as If but with err and attach it)
// -----

// Log writes log message with desired 'level',
// analyzing 'args' in the most powerful and smart way:
//
// - args[0] could be printf-like format string, then next N args will be
//   its printf values (N - num of format's printf verbs),
//   and M-N (where M is total count of args) will be treated as
//   explicit or implicit fields (depends on what kind of fields are allowed);
//
// - args[0] could be an error => error's message will be used as log's one,
//   and wherein, error's message could be printf-like format string,
//   then next N args will be handled as in the case above,
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
	return l.log(level, format, nil, args, nil)
}

// Logw writes log message 'msg' with desired 'level', and passed implicit fields.
func (l *Logger) Logw(level Level, msg string, fields ...field.Field) (this *Logger) {

	return l.log(level, msg, nil, nil, fields)
}

// Debug is the same as Log(Level.Debug, args...).
// Read more: Entry.Log.
func (l *Logger) Debug(args ...interface{}) (this *Logger) {

	return l.log(LEVEL_DEBUG, "", nil, args, nil)
}

// DebugStrict is the same as LogStrict(Level.Debug, args...).
// Read more: Entry.LogStrict.
func (l *Logger) DebugStrict(args ...interface{}) (this *Logger) {

	return l.log(LEVEL_DEBUG, fmt.Sprint(args...), nil, nil, nil)
}

// Debugf is the same as Logf(Level.Debug, format, args...).
// Read more: Entry.Logf.
func (l *Logger) Debugf(format string, args ...interface{}) (this *Logger) {

	return l.log(LEVEL_DEBUG, format, nil, args, nil)
}

// DebugfStrict is the same as LogfStrict(Level.Debug, format, args...).
// Read more: Entry.LogfStrict.
func (l *Logger) DebugfStrict(format string, args ...interface{}) (this *Logger) {

	return l.log(LEVEL_DEBUG, fmt.Sprintf(format, args...), nil, nil, nil)
}

// Debugw is the same as Logw(Level.Debug, msg, fields...).
// Read more: Entry.Logw.
func (l *Logger) Debugw(msg string, fields ...field.Field) (this *Logger) {

	return l.log(LEVEL_DEBUG, msg, nil, nil, fields)
}

// Info is the same as Log(Level.Info, args...).
// Read more: Entry.Log.
func (l *Logger) Info(args ...interface{}) (this *Logger) {

	return l.log(LEVEL_INFO, "", nil, args, nil)
}

// InfoStrict is the same as LogStrict(Level.Info, args...).
// Read more: Entry.LogStrict.
func (l *Logger) InfoStrict(args ...interface{}) (this *Logger) {

	return l.log(LEVEL_INFO, fmt.Sprint(args...), nil, nil, nil)
}

// Infof is the same as Logf(Level.Info, format, args...).
// Read more: Entry.Logf.
func (l *Logger) Infof(format string, args ...interface{}) (this *Logger) {

	return l.log(LEVEL_INFO, format, nil, args, nil)
}

// InfofStrict is the same as LogfStrict(Level.Info, format, args...).
// Read more: Entry.LogfStrict.
func (l *Logger) InfofStrict(format string, args ...interface{}) (this *Logger) {

	return l.log(LEVEL_INFO, fmt.Sprintf(format, args...), nil, nil, nil)
}

// Infow is the same as Logw(Level.Info, msg, fields...).
// Read more: Entry.Logw.
func (l *Logger) Infow(msg string, fields ...field.Field) (this *Logger) {

	return l.log(LEVEL_INFO, msg, nil, nil, fields)
}

// Warn is the same as Log(Level.Warn, args...).
// Read more: Entry.Log.
func (l *Logger) Warn(args ...interface{}) (this *Logger) {

	return l.log(LEVEL_WARNING, "", nil, args, nil)
}

// WarnStrict is the same as LogStrict(Level.Warn, args...).
// Read more: Entry.LogStrict.
func (l *Logger) WarnStrict(args ...interface{}) (this *Logger) {

	return l.log(LEVEL_WARNING, fmt.Sprint(args...), nil, nil, nil)
}

// Warnf is the same as Logf(Level.Warn, format, args...).
// Read more: Entry.Logf.
func (l *Logger) Warnf(format string, args ...interface{}) (this *Logger) {

	return l.log(LEVEL_WARNING, format, nil, args, nil)
}

// WarnfStrict is the same as LogfStrict(Level.Warn, format, args...).
// Read more: Entry.LogfStrict.
func (l *Logger) WarnfStrict(format string, args ...interface{}) (this *Logger) {

	return l.log(LEVEL_WARNING, fmt.Sprintf(format, args...), nil, nil, nil)
}

// Warnw is the same as Logw(Level.Warn, msg, fields...).
// Read more: Entry.Logw.
func (l *Logger) Warnw(msg string, fields ...field.Field) (this *Logger) {

	return l.log(LEVEL_WARNING, msg, nil, nil, fields)
}

// Error is the same as Log(Level.Error, args...).
// Read more: Entry.Log.
func (l *Logger) Error(args ...interface{}) (this *Logger) {

	return l.log(LEVEL_ERROR, "", nil, args, nil)
}

// ErrorStrict is the same as LogStrict(Level.Error, args...).
// Read more: Entry.LogStrict.
func (l *Logger) ErrorStrict(args ...interface{}) (this *Logger) {

	return l.log(LEVEL_ERROR, fmt.Sprint(args...), nil, nil, nil)
}

// Errorf is the same as Logf(Level.Error, format, args...).
// Read more: Entry.Logf.
func (l *Logger) Errorf(format string, args ...interface{}) (this *Logger) {

	return l.log(LEVEL_ERROR, format, nil, args, nil)
}

// ErrorfStrict is the same as LogfStrict(Level.Error, format, args...).
// Read more: Entry.LogfStrict.
func (l *Logger) ErrorfStrict(format string, args ...interface{}) (this *Logger) {

	return l.log(LEVEL_ERROR, fmt.Sprintf(format, args...), nil, nil, nil)
}

// Errorw is the same as Logw(Level.Error, msg, fields...).
// Read more: Entry.Logw.
func (l *Logger) Errorw(msg string, fields ...field.Field) (this *Logger) {

	return l.log(LEVEL_ERROR, msg, nil, nil, fields)
}

// Fatal is the same as Log(Level.Fatal, args...),
// but also then calls death.Die(1).
// Read more: Entry.Log.
func (l *Logger) Fatal(args ...interface{}) (this *Logger) {

	return l.log(LEVEL_FATAL, "", nil, args, nil)
}

// FatalStrict is the same as LogStrict(Level.Fatal, args...),
// but also then calls death.Die(1).
// Read more: Entry.LogStrict.
func (l *Logger) FatalStrict(args ...interface{}) (this *Logger) {

	return l.log(LEVEL_FATAL, fmt.Sprint(args...), nil, nil, nil)
}

// Fatalf is the same as Logf(Level.Fatal, format, args...),
// but also then calls death.Die(1).
// Read more: Entry.Logf.
func (l *Logger) Fatalf(format string, args ...interface{}) (this *Logger) {

	return l.log(LEVEL_FATAL, format, nil, args, nil)
}

// FatalfStrict is the same as LogfStrict(Level.Fatal, format, args...),
// but also then calls death.Die(1).
// Read more: Entry.LogfStrict.
func (l *Logger) FatalfStrict(format string, args ...interface{}) (this *Logger) {

	return l.log(LEVEL_FATAL, fmt.Sprintf(format, args...), nil, nil, nil)
}

// Fatalw is the same as Logw(Level.Fatal, msg, fields...),
// but also then calls death.Die(1).
// Read more: Entry.Logw.
func (l *Logger) Fatalw(msg string, fields ...field.Field) (this *Logger) {

	return l.log(LEVEL_FATAL, msg, nil, nil, fields)
}
