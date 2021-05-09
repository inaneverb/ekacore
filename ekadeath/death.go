// Copyright © 2019-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekadeath

import (
	"os"
	"reflect"
)

// ---------------------------------------------------------------------------- //
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
//    Just call Reg(foo), where foo is your func and it's done!
//
// 3. Want to know what exit code has been used to stop service?
//    Pass your destructor func to Reg with "func (code int)" signature.
//    Now in your destructor you have 'code' argument and yes, it's exit code.
//
// 4. Want to stop service?
//    Call Die(). In this case, be default, the exit code is 1.
//    Want another? Pass your desired exit code to Die (e.g. Die(2))
//    and combine it with your own extended destructor "func (code int)".
//
// 5. Do not want to multiply "if exitCode == 1 {} else if exitCode == 2 {} ..."
//    and you hate switch construction? Do not worry!
//    Register your destructors that will be called only for specified exit code.
//    Call Reg(exitCodeToBind int, yourDestructor func()).
//    Want to know in your destructor that you bind it e.g. to exitCode, dunno, 5?
//    No problem, Reg(5, foo), where foo is func(code){...}.
//    Now in your foo, code is 5. It's simple, isn't?
//
// ---------------------------------------------------------------------------- //
//
// What functions you can use as destructors?
// Only 2 signatures of destructors are allowable:
// - func(): no arguments, no returns. Just your callback (or closure).
// - func(code int): one argument, no returns. Argument is exit code Die() called with.
//
// ---------------------------------------------------------------------------- //

type (
	// DestructorSimple is an alias to the function that you may register
	// using Reg() function to be executed after Die() or Exit() is called.
	DestructorSimple = func()
	DestructorWithExitCode = func(code int)
)

// Reg registers destructors to be called when app should be stopped
// when Die() or Exit() is called or SIGTERM/SIGKILL received.
//
// You can use Reg to do:
// 1. Just reg one or many destructor(s): Reg(foo), Reg(foo1, foo2, foo3).
// 2. Reg destructor (one or many) to be called for special exitCode only:
//    Reg(exitCode, foo), Reg(exitCode, foo1, foo2, foo3),
//    where exitCode should be: int, int8, int16, int32, int64 and the same uint's.
//
// Nil destructors will be ignored.
//
// Despite of fact that all int/uint types are available to be used as exit code,
// their values must be in the int32 range. UB otherwise.
//
// Since v3 version of ekago/ekadeath you may register a new destructors
// when a death is requested and other destructors are under executing now.
// In that case, added destructor will be executed just after the destructor,
// that is under executing now.
func Reg(args ...interface{}) {

	if l := len(args); l == 0 {
		return

	} else if v0 := reflect.ValueOf(args[0]); l == 1 {
		reg(false, 0, args[0])

	} else if k := v0.Kind(); k >= reflect.Int && k <= reflect.Int64 {
		reg(true, int(v0.Int()), args[1:]...)

	} else if k >= reflect.Uint && k <= reflect.Uint64 {
		reg(true, int(v0.Uint()), args[1:]...)

	} else {
		reg(false, 0, args...)
	}
}

// Exit is the same as Die(0).
func Exit() {
	Die(0)
}

// Die calls all registered destructors and then shutdowns an app using os.Exit.
// Thus all goroutines will be forcibly stopped.
//
// You can pass one int argument as exit code. In this case only common
// and associated with specified exit code destructors will be called.
// By default, exit code is 1. 2nd and next codes are ignored.
func Die(code ...int) {

	exitCode := 1 // default value, could be overwritten by first arg
	if len(code) > 0 {
		exitCode = code[0]
	}

	for elem, found := destructors.Pop(); found; elem, found = destructors.Pop() {
		destructor := elem.(destructorRegistered)
		if destructor.callAnyway || destructor.bindToExitCode == exitCode {
			invoke(destructor.f, exitCode)
		}
	}

	os.Exit(exitCode)
}

// RegisteredNum reports how much destructors are registered for now.
func RegisteredNum() int {
	return destructors.Len()
}
