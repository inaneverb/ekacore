// Copyright © 2019-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekadeath

import (
	"github.com/inaneverb/ekacore/ekatyp/v4"
	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

// destructorDescriptor is a descriptor of registered destructor.
// Each Reg() call converts provided destructor to this.
type destructorDescriptor struct {
	f          DestructorWithExitCode
	exitCode   int  // f will be called only if app is down with that exit code
	callAnyway bool // call no matter what exit code is
}

// destructors is a LIFO stack that contains destructorDescriptor objects.
var destructors ekatyp.Stack

// reg registers each function from 'd' as a destructor that:
//   - Will be called anyway if 'bind' is false ('exitCode' is ignored this way);
//   - Will be called if Die() with the same 'exitCode' is called
//     and 'bind' is true.
func reg(bind bool, exitCode int, d []any) {
	for i, n := 0, len(d); i < n; i++ {
		if f := parse(d[i]); f != nil {
			destructors.Push(newDescriptor(f, exitCode, !bind))
		}
	}
}

// parse tries to treat provided argument as a destructor.
// Returns non-nil one if it so; nil otherwise.
func parse(d any) DestructorWithExitCode {
	if ekaunsafe.TakeRealAddr(d) == nil {
		return nil
	}
	switch d := d.(type) {
	case DestructorSimple:
		return func(_ int) { d() } // "cast" to DestructorWithExitCode
	case DestructorWithExitCode:
		return d
	default:
		return nil
	}
}

// newDescriptor creates a new destructorDescriptor with given args.
func newDescriptor(
	d DestructorWithExitCode,
	exitCode int, callAnyway bool) destructorDescriptor {

	return destructorDescriptor{d, exitCode, callAnyway}
}
