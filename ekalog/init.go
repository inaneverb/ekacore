// Copyright © 20202-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"os"
)

func init() {
	entryPool.New = allocEntry

	defaultConsoleEncoder = new(CI_ConsoleEncoder).doBuild()
	defaultJSONEncoder = new(CI_JSONEncoder).doBuild()

	_ = defaultConsoleEncoder
	_ = defaultJSONEncoder

	integrator := new(CommonIntegrator).
		WithEncoder(defaultConsoleEncoder).
		WithMinLevel(LEVEL_DEBUG).
		WithMinLevelForStackTrace(LEVEL_WARNING).
		WriteTo(os.Stdout)
	integrator.build()

	entry := acquireEntry()

	baseLogger = new(Logger).setIntegrator(integrator).setEntry(entry)

	// We know that nopLogger's fields won't be accessed, but we need them
	// to pass Logger.assert() and Logger.IsValid() checks.
	//
	// So, (*CommonIntegrator)(nil) assigned to Integrator's type field
	// won't be nil in Golang terms (it has a type no matter value still is nil),
	// and we don't need a fully initialization of Entry - we won't use it.
	// Logger.setIntegrator(), Logger.setEntry() doesn't have any checks.

	nopLogger = new(Logger).setIntegrator((*CommonIntegrator)(nil)).setEntry(new(Entry))
}
