// Copyright © 2019-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekadeath

import (
	"github.com/inaneverb/ekacore/ekatyp/v4"
)

// dd is a descriptor of registered destructor.
// Each Reg() call converts provided destructor to this.
type dd struct {
	f          func(code int)
	exitCode   int  // f will be called only if app is down with that exit code
	callAnyway bool // call no matter what exit code is
}

// destructors is a LIFO stack that contains dd objects.
var destructors ekatyp.Stack

// reg registers each function from 'd' as a destructor that:
//   - Will be called anyway if 'bind' is false ('exitCode' is ignored this way);
//   - Will be called if Die() with the same 'exitCode' is called
//     and 'bind' is true.
func reg(bind bool, exitCode int, cb func(code int)) {
	destructors.Push(newDescriptor(cb, exitCode, !bind))
}

// newDescriptor creates a new dd with given args.
func newDescriptor(d func(code int), exitCode int, callAnyway bool) dd {
	return dd{d, exitCode, callAnyway}
}

// wrap wraps destructor w/o accepting exit code to the destructor,
// that accepts.
func wrap(cb func()) func(code int) {
	return func(_ int) { cb() }
}
