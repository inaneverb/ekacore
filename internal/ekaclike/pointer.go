// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaclike

import (
	"unsafe"
)

/*
TakeRealAddr takes and returns a real address of the variable you passed.
It returns nil if nil is passed.

	var i int
	_ = &i == TakeRealAddr(i) // true

This function exist to extend standard Golang & operator.
Using this function you can take an address of some things,
that "cannot be addressed" in Golang policy. Like functions, for example.

What you CAN do with functions and their addresses:

	f := func(){}
	_ = &f // 0x...

What you CANNOT do with functions and their addresses:

	func foo() {}
	func bar() { _ = &foo } // compilation err

What you CAN do now using TakeRealAddr() with functions and their addresses:

	func foo() {}
	type T struct{}
	func (_ T) bar() {}

	func main() {
	        var t T
	        ptr1 := TakeRealAddr(foo)      // 0x...
	        ptr2 := TakeRealAddr((*T).bar) // 0x...
	        ptr3 := TakeRealAddr(t.bar)    // 0x...
	}

Speaking about functions, if you want to convert it back,
and be able to call a function address of which you got using TakeRealAddr(),
you need to pass it through Addr2Callable() or just use TakeCallableAddr().
*/
func TakeRealAddr(i any) unsafe.Pointer {
	return UnpackInterface(i).Word
}

/*
TakeCallableAddr extends TakeRealAddr functionality,
providing to you a mechanism to get a real address of any function, which you may:

  - Compare with another, obtained the same way, address;
  - Doing other unsafe but interest stuff with function address in C-style;
  - Convert that address back and call the function.

It returns nil if nil is passed.

Tip:
Using this function you may not only compare addresses of function,
convert addresses to functions, vice-versa and call them, but
you can also avoid type checks while any conversion. Yes, C-style, just as we like.

Usage:

	type F = func(f float64) (int32, int32)
	func foo(x1, x2 int32) (int32, int32) { return x1, x2 }
	untypedPtr := TakeCallableAddr(foo)
	x1, x2 := (*(*F)(untypedPtr))(math.PI)

Now, x1 and x2 contains first and last 32 bytes of math.PI constant as Golang int32.
It's a synthetic example. Real world examples will be presented later.

-----

WARNING!
DO NOT PASS SOMETHING BUT EITHER FUNCTIONS, FUNCTORS, CLOSURES, LAMBDAS, ETC.
YOU MUST PASS ONLY THAT THING THAT MAY BE CALLED.
UNDEFINED BEHAVIOUR OTHERWISE!

WARNING!
DESPITE THE FACT, IT IS POSSIBLE TO CONVERT ADDRESS BACK AND CALL THE FUNCTION,
THERE MIGHT BE SOME PROBLEMS WHEN YOU'RE TRYING TO WORK AROUND METHODS,
NOT JUST FUNCTIONS.
DO IT ON YOUR OWN RISK.

WARNING!
USING THIS FUNCTION IS VERY, VERY DANGEROUS.
YOU MAY GET UNDEFINED BEHAVIOUR, PANIC
AND IT IS SO DEPENDED OF THE GOLANG INTERNAL IMPLEMENTATIONS.
THIS FUNCTION NEVER USED IN OTHER PARTS OF THAT LIBRARY.

BE CAREFULLY!
*/
func TakeCallableAddr(fn any) unsafe.Pointer {

	// There is no need nil checks,
	// because TakeRealAddr and AddrConvert2Callable already has it
	return Addr2Callable(TakeRealAddr(fn))
}

/*
Addr2Callable transforms an address of some function obtained by TakeRealAddr()
to that kind of address, you can cast to the typed pointer to the some function,
dereference it and call.

See TakeRealAddr(), TakeCallableAddr() for more info.
*/
func Addr2Callable(realPtr unsafe.Pointer) (callablePtr unsafe.Pointer) {

	type fptr struct {
		ptr unsafe.Pointer
	}

	if realPtr == nil {
		return nil
	}

	o := new(fptr)
	o.ptr = realPtr

	return unsafe.Pointer(&o.ptr)
}

/*
Addr2Real transforms a callable address of some function,
obtained by TakeCallableAddr()
to a simple, real, non-callable address of that function.
*/
func Addr2Real(callablePtr unsafe.Pointer) (realPtr unsafe.Pointer) {

	if callablePtr == nil {
		return nil
	}

	return *(*unsafe.Pointer)(callablePtr)
}
