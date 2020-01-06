// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package dangerous

import "unsafe"
import "github.com/qioalice/gext/internal/gotypes"

// In Golang you can not just take an address of some func
// only if it is not a func literal which is assigned to any var.
//
// But you can assign any kind of func to type-compatible variable
// or interface{}.
// And you can call then that assigned function through used variable.
// But anyway that way still requires a type-compatible variable
// or type-casting interface{} and you can not break that type-compatible
// rules:
//
// - You can not assign func to some variable with wrong func signature,
// - You can not cast interface{} that stores func with one signature
// to func with another.
//
// Unless you have an untyped pointer like void* in C.
//
// reflect.ValueOf(f).Pointer() returns a REAL function address,
// the same as TakeRealAddr(f) returns, but you CAN NOT use it to call.
//
// CAN:
// f := func(){}
// &f
//
// CAN NOT (1):
// func f(){}
// &f
//
// CAN NOT (2):
// type g uint8
// func (g) f()
// var g g
// &g.f
//
// But you can assign any kind of func to type-compatible variable
// or interface{}.

// TODO: Review doc above

// TakeRealAddr takes and returns a real address of function fn or nil if fn is nil.
//
// If fn is not a function, you still get the pointer, but it is unknown
// what that pointer points to.
//
// YOU CAN NOT USE RETURNED POINTER TO CALLING PASSED FUNCTION!
// For that purpose use TakeCallableAddr to take callable pointer directly
// or convert returned pointer to callable pointer using AddrConvert2Callable func.
func TakeRealAddr(fn interface{}) unsafe.Pointer {
	if fn == nil {
		return nil
	}
	return (*gotypes.Interface)(unsafe.Pointer(&fn)).Word
}

// TakeCallableAddr takes and returns an "callable" address of function fn or nil if fn is nil.
//
// If fn is not a function, you still get the pointer, but it is unknown
// what that pointer points to.
//
// You can use that address to call function that address points to, like in C.
// To calling just convert returned untyped pointer to function-typed,
// dereference it and call.
//
// You can AVOID TYPE CHECKS using that way (wrong argument types, wrong return types)
// but there is UB in that way and do it only if you know what you're doing.
func TakeCallableAddr(fn interface{}) unsafe.Pointer {

	// There is no need nil checks,
	// because TakeRealAddr and AddrConvert2Callable already has it
	return AddrConvert2Callable(TakeRealAddr(fn))
}

// AddrConvert2Callable converts a normal function pointer to a pointer using which
// becomes possible to call a function normalPtr points to.
//
// It is assumed that normalPtr has been obtained using TakeRealAddr func.
// PLEASE DO NOT PASS POINTERS OBTAINED BY ANOTHER MEANS.
func AddrConvert2Callable(normalPtr unsafe.Pointer) (callablePtr unsafe.Pointer) {

	type fptr struct {
		ptr unsafe.Pointer
	}

	if normalPtr == nil {
		return nil
	}

	o := new(fptr)
	o.ptr = normalPtr
	return unsafe.Pointer(&o.ptr)
}

// AddrConvert2Normal converts a callable function pointer to a normal, internal func's pointer.
//
// It is assumed that callablePtr has been obtained using TakeCallableAddr func.
// PLEASE DO NOT PASS POINTERS OBTAINED BY ANOTHER MEANS.
func AddrConvert2Normal(callablePtr unsafe.Pointer) (normalPtr unsafe.Pointer) {

	if callablePtr == nil {
		return nil
	}

	return *(*unsafe.Pointer)(callablePtr)
}

// 	type tfptr struct {
// 		ptr unsafe.Pointer
// 	}

// 	if f == nil {
// 		return nil
// 	}

// 	o := new(struct{ ptr unsafe.Pointer })
// 	o.ptr = (*typedInterface)(unsafe.Pointer(&f)).word

// 	if o.ptr == nil {
// 		return nil
// 	}

// 	return unsafe.Pointer(&o.ptr)
// }

// func main() {
// 	// 0x10aff40
// 	// 0x10ea0f0
// 	// (*(*func(*t1))(unsafe.Pointer(uintptr(0x10ea0f0))))(nil)

// 	fObject := f
// 	fObjectPtr := fptr(fObject)
// 	fObjectPtrCasted := (*func(*t2))(fObjectPtr)

// 	fReflectValue := reflect.ValueOf(f)
// 	fReflectPointer := unsafe.Pointer(fReflectValue.Pointer())

// 	fObjectReflectValue := reflect.ValueOf(fObject)
// 	fObjectReflectPointer := unsafe.Pointer(fObjectReflectValue.Pointer())

// 	fmt.Println("new(main.f)", fObjectPtr)
// 	fmt.Println("new(main.f) -> *func(*t2)", fObjectPtrCasted)
// 	fmt.Println("rv(main.f).ptr", fReflectPointer)
// 	fmt.Println("rv(new(main.f)).ptr", fObjectReflectPointer)

// 	if !(fObjectPtr == unsafe.Pointer(fObjectPtrCasted) &&
// 		fObjectPtr == fReflectPointer &&
// 		fObjectPtr == fObjectReflectPointer) {
// 		panic("incompatible addresses")
// 	}

// 	h := reflect.SliceHeader{}
// 	h.Data = uintptr(fObjectPtr)
// 	h.Len = 1
// 	h.Cap = 1

// 	hb := *(*[]byte)(unsafe.Pointer(&h))

// 	if errno := syscall.Mprotect(hb, syscall.PROT_READ|syscall.PROT_EXEC); errno != nil {
// 		panic(errno)
// 	}

// 	f(nil)
// 	(*(*func(*t1))(fObjectPtr))(nil)            // call fObjectPtr as f(*t1)
// 	(*fObjectPtrCasted)(nil)                    // call fObjectPtrCasted as f(*t2)
// 	(*(*func(*t1))(fReflectPointer))(nil)       // call fReflectPointer as f(*t1)
// 	(*(*func(*t1))(fObjectReflectPointer))(nil) // call fObjectReflectPointer as f(*t1)
// }

// func main() {

// 	fObject := f
// 	fObjectPtr := unsafe.Pointer(&fObject)

// 	var fInterface interface{} = fObject

// 	fIfaceHeader := (*typedInterface)((unsafe.Pointer)(&fInterface))

// 	fmt.Println("&fObject", fObjectPtr)
// 	fmt.Println("*fObject", *(*unsafe.Pointer)(fObjectPtr))

// 	fmt.Println("&fObject(I)", unsafe.Pointer(fIfaceHeader.typ), fIfaceHeader.w)

// 	(*(*func(*t2))(fObjectPtr))(nil)

// 	val := fIfaceHeader.w

// 	(*(*func(*t2))(unsafe.Pointer(&val)))(nil)
// }
