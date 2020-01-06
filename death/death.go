// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package death

import "os"
import "os/signal"
import "reflect"
import "syscall"
import "sync"

// ---------------------------------------------------------------------------- //
//
// This package is need to manage graceful shutdown of the whole service.
// Package allows you to call some functions before os.Exit(), SIGKILL, SIGTERM.
//
// It's useful to do some prepares before gonna die, e.g.:
// - close DB connections, close files,
// - flush logs, requests
// - notify about it
//
// All you need is:
//
// ---------------------------------------------------------------------------- //
//
// 1. DO NOT USE os.Exit ANYMORE. HERE AND EVERYWHERE ELSE. NEVER. STOP IT!
//    USE Die() FUNC INSTEAD!
//
// ---------------------------------------------------------------------------- //
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
// What functions you can use as destructors? Context-independent those, I guess.
// Only 2 signatures of destructors are allowable:
// - func(): no arguments, no returns. Just your callback (or closure).
// - func(code int): one argument, no returns. Argument is exit code Die called with.
//
// ---------------------------------------------------------------------------- //
//
// How it works? Very simple!
// Destructors stored in map[exitCode]destructorToBeCalledWithThatExitCode,
// and there is goroutine that calls Die(1) if SIGKILL/SIGTERM is occurred.
//
// ---------------------------------------------------------------------------- //

type (
	DestructorSimple       = func()
	DestructorWithExitCode = func(code int)

	destructorRegistered = struct {
		f              interface{}
		bindToExitCode int
		callAnyway     bool
	}
)

var (
	destructors []destructorRegistered
	mu          sync.Mutex
)

// Package initialization. Spawns goroutine which can handle SIGKILL, SIGTERM
// and call then Die(1).
func init() {

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)

	go func() {
		_ = <-ch // blocks
		Die(1)
	}()
}

// Reg registers destructors to be called when service should be stopped
// (Die is called or SIGTERM/SIGKILL).
//
// You can use func as destructor
// if it's type is either DestructorSimple or DestructorWithExitCode.
//
// You can use Reg to do:
// 1. Just reg one or many destructor(s): Reg(foo), Reg(foo1, foo2, foo3).
// 2. Reg destructor (one or many) to be called for special exitCode only:
//    Reg(exitCode, foo), Reg(exitCode, foo1, foo2, foo3),
//    where exitCode should be: int, int8, int16, int32, int64 and the same uint's.
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

// Die calls all registered destructors (in this moment you can't register new)
// and then shutdowns a service through os.Exit - all goroutines will be forcibly stopped.
//
// You can pass one int argument as exit code. In this case only common
// and associated with specified exit code destructors will be called.
// By default, exit code is 1.
func Die(code ...int) {

	mu.Lock()
	defer mu.Unlock()

	exitCode := 1 // default value, could be overwritten by first arg
	if len(code) > 0 {
		exitCode = code[0]
	}

	// Call destructors in reverse order (LIFO)
	for i := len(destructors) - 1; i >= 0; i-- {
		if destructors[i].callAnyway || destructors[i].bindToExitCode == exitCode {
			call(destructors[i].f, exitCode)
		}
	}

	os.Exit(exitCode)
}

// RegisteredNum reports how much destructors are registered.
func RegisteredNum() int {

	mu.Lock()
	defer mu.Unlock()

	return len(destructors)
}

// reg registers each function from destructorsToBeRegistered as destructor
// that will be called anyway if hasExitCodeBind is false (exitCode is ignored this way)
// or will be called if Die with passed exitCode is called if hasExitCodeBind is true.
func reg(hasExitCodeBind bool, exitCode int, destructorsToBeRegistered ...interface{}) {

	for _, destructor := range destructorsToBeRegistered {
		if !valid(destructor) {
			continue
		}
		destructors = append(destructors, destructorRegistered{
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

// call calls d with no passing arguments if d is DestructorSimple,
// or passing exitCode if d is DestructorWithExitCode.
func call(d interface{}, exitCode int) {
	switch d.(type) {

	case DestructorSimple:
		d.(DestructorSimple)()

	case DestructorWithExitCode:
		d.(DestructorWithExitCode)(exitCode)
	}
}
