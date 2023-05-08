// Copyright © 2019-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekadeath

import (
	"os"
)

////////////////////////////////////////////////////////////////////////////////
//
// This package is need to manage graceful shutdown of the whole service.
// Package allows you to call some functions before os.Exit(), SIGKILL, SIGTERM.
//
// Registered destructors is called in LIFO order,
// meaning the last registered destructor will be called firstly.
//
// It's useful to do some prepares before gonna die, e.g.:
// - close DB connections, close files,
// - flush logs, requests
// - notify about it
//
// All you need is:
//
// 1. DO NOT USE os.Exit ANYMORE. HERE AND EVERYWHERE ELSE. NEVER. STOP IT!
//    USE Die() FUNC INSTEAD!
//
// 2. You need your foo func to be called when Die() is called or SIGTERM/SIGKILL?
//    Just call Reg(foo), where foo is your func, and it's done!
//
// 3. Want to know what exit code has been used to stop the service?
//    Pass your destructor func to Reg with "func (code int)" signature.
//    Now in your destructor you have 'code' argument and yes, it's an exit code.
//
// 4. Want to stop the service?
//    Call Die(). In this case, be default, the exit code is 1.
//    Want another? Pass your desired exit code to Die (e.g. Die(2))
//    and combine it with your own extended destructor "func (code int)".
//
// 5. Do not want to multiply "if exitCode == 1 {} else if exitCode == 2 {} ..."
//    and you hate switch construction? Do not worry!
//    Register your destructors that will be called only for specified exit code.
//    Call Reg(yourDestructor func(), exitCodes[1,2,3]).
//    Want to know in your destructor that you bind it e.g. to exitCode, dunno, 5?
//    No problem, Reg2(foo, 5), where foo is func(code){...}.
//    Now in your foo, code is 5. It's simple, isn't?
//
////////////////////////////////////////////////////////////////////////////////
//
// What functions you can use as destructors?
// Only 2 signatures of destructors are allowed:
// - func(): no arguments, no returns. Just your callback.
// - func(code int): one argument, no returns.
//   Argument is exit code Die() called with.
//
////////////////////////////////////////////////////////////////////////////////

// Reg registers destructors to be called when app should be stopped
// when Die() or Exit() is called or SIGTERM/SIGKILL received.
//
// You can use Reg to reg destructor to be called:
// 1. For all exit codes (empty 'codes' provided)
// 2. For special exit code only (with specified 'codes').
//
// WARNING!
// Nil destructors will be ignored.
//
// Note.
// You may register a new destructors when death is requested
// and other destructors are under executing now.
// In that case, added destructor will be executed just after the destructor,
// that is under executing now.
func Reg(cb func(), codes ...int) {

	if cb == nil {
		return
	}

	if len(codes) > 0 {
		reg(true, codes[0], wrap(cb))
	} else {
		reg(false, 0, wrap(cb))
	}
}

func Reg2(cb func(int), codes ...int) {

	if cb == nil {
		return
	}

	if len(codes) > 0 {
		reg(true, codes[0], cb)
	} else {
		reg(false, 0, cb)
	}
}

// Exit is the same as Die(0).
func Exit() {
	Die(0)
}

// Die calls all registered destructors and shutdowns an app using os.Exit().
// Thus, all goroutines will be forcibly stopped.
//
// You can pass one int argument as exit code.
// Any numeric variants are allowed (int, int8, ..., uint, uint8, ..., uintptr).
// In this case only common and associated with specified exit code destructors
// will be called. By default, exit code is 1. 2nd and next codes are ignored.
func Die(code ...int) {

	var exitCode = 1
	if len(code) > 0 {
		exitCode = code[0]
	}

	for elem, found := destructors.Pop(); found; elem, found = destructors.Pop() {
		destructor := elem.(dd)
		if destructor.callAnyway || destructor.exitCode == exitCode {
			destructor.f(exitCode)
		}
	}

	os.Exit(exitCode)
}
