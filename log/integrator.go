// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

// -----
// Integrator is the why this package called 'logintegro' earlier.
// The main idea is: "You integrate your log messages with your destination service".
// -----

// Integrator is an interface each type wants to convert log's Entry to some
// real output shall implement.
//
// E.g. If you want to use this package and writes all log's entries to your
// own service declare, define your type and implement Integrator interface
// and reg then it later using 'Logger.Apply' method or 'ApplyThis' func.
// But you also can use any of predefined basic integrators which cover 99% cases.
type Integrator interface {

	// Write writes log entry to some destination (integrator determines
	// what it will be). Thus, Write does the main thing of Integrator:
	// "Integrates your log messages with your log destination service".
	Write(entry *Entry)

	// MinLevelEnabled returns minimum log's Level an integrator will handle
	// log entries with.
	// E.g. if minimum level is 'Warning', 'Debug' logs will be dropped.
	MinLevelEnabled() Level

	// Sync flushes all pending log entries to integrator destination.
	// It useful when integrator does async work and sometimes you need to make sure
	// all pending entries are flushed.
	//
	// Logger type has the same name's method that just calls this method.
	Sync() error
}
