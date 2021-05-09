// Copyright © 2019-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekadeath

import (
	"github.com/qioalice/ekago/v3/ekatyp"
)

type (
	// destructorRegistered is a destructor descriptor.
	// Each Reg() call converts passed destructor to that descriptor.
	destructorRegistered = struct {
		f              interface{} // can be DestructorSimple or DestructorWithExitCode
		bindToExitCode int         // f will be called only if app is down with that exit code
		callAnyway     bool        // call no matter what exit code is
	}
)

var (
	// destructors is a LIFO stack that contains destructorRegistered objects.
	destructors ekatyp.Stack
)

// reg registers each function from destructorsToBeRegistered as destructor
// that will be called anyway if hasExitCodeBind is false (exitCode is ignored this way)
// or will be called if Die with passed exitCode is called if hasExitCodeBind is true.
func reg(hasExitCodeBind bool, exitCode int, destructorsToBeRegistered ...interface{}) {
	for _, destructor := range destructorsToBeRegistered {
		if !valid(destructor) {
			continue
		}
		destructors.Push(destructorRegistered{
			f:              destructor,
			bindToExitCode: exitCode,
			callAnyway:     !hasExitCodeBind,
		})
	}
}

// valid reports whether d is valid destructor:
// - it's type either DestructorSimple or DestructorWithExitCode,
// - it's value is not nil.
func valid(d interface{}) bool {
	switch d.(type) {
	case DestructorSimple:
		return d.(DestructorSimple) != nil
	case DestructorWithExitCode:
		return d.(DestructorWithExitCode) != nil
	default:
		return false
	}
}

// invoke calls d with no passing arguments if d is DestructorSimple,
// or passing exitCode if d is DestructorWithExitCode.
func invoke(d interface{}, exitCode int) {
	switch d.(type) {
	case DestructorSimple:
		d.(DestructorSimple)()
	case DestructorWithExitCode:
		d.(DestructorWithExitCode)(exitCode)
	}
}
