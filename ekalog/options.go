// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

//import "io"

// TOption is a special alias for func that receives Logger object
// by reference. There are many functions that provides TOptions.
// These objects can change behaviour of logger they're applies to which.
// TOption is a special alias to the function that changes behaviour of
// logger core, it applies to which.

//
type option func(e *Entry)

//
var Options = struct {

	// SetFormat TOption group.
	// Use more than one formatter for one core means that all previous
	// will be overwritten by last.
	//SetFormat tFormatOptioner

	// Set is TOption generator which returns a special
	// TOption that changes the encoding logging behaviour
	// (and log message formatting) to the behaviour argument 'formatter'
	// is represent.
	// formatter argument's type can be:
	// - TFormatter type: Treated as formatter which will be used
	// - string type: TFormatter object will be generated
	// by GenerateFormatter function, and this string will be used
	// as desired format (see GenerateFormatter docs).

	//WriteTo tWriteOptioner
	//
	//Enable tEnableOptioner

	// todo: Probably rename all formatter group applicators to started from "By"

	// EXCEPTION: ToAndBy TOption

	// Destination TOption group.

	// Enable logging TOption group.
}{}

////
//type tFormatOptioner func(format interface{}) TOption
//
////
//func (*tFormatOptioner) AsJSON() TOption {
//
//}
//
////
//func (*tFormatOptioner) AsPlainText() TOption {
//
//}
//
////
//type tWriteOptioner func(writer io.Writer) TOption
//
////
//func (*tWriteOptioner) File(filename string) TOption {
//
//}
//
////
//func (*tWriteOptioner) Stdout() TOption {
//
//}
//
////
//func (*tWriteOptioner) Stderr() TOption {
//
//}
//
////
//type tEnableOptioner func(is ...bool) TOption
//
////
//func (*tEnableOptioner) EmptyMessages(is ...bool) TOption {
//
//}
//
////
//func (*tEnableOptioner) LoggingFrom(level Level) TOption {
//
//}
//
////
//func (*tEnableOptioner) Stacktrace(is ...bool) TOption {
//
//}
//
////
//func (*tEnableOptioner) StacktraceFrom(level Level) TOption {
//
//}
//
////
//func (*tEnableOptioner) AddingCaller(is ...bool) TOption {
//
//}

// parseOptions parses 'anyOptions' and tries to do following things:
//
// 1. Tries to extract Integrator object from anyOptions and if it so,
//    returns it separately (1st return arg).
//
// 2. If there is more than one Integrator in 'anyOptions', all of them will be
//    teed as a new Integrator and it will be returned as 1st return arg.
//
// 3. If there is
func parseOptions(anyOptions []interface{}) (Integrator, []option) {
	return nil, nil
}
