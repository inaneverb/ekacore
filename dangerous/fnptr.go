// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package dangerous

import (
	"unsafe"
)

type typedInterface struct {
	typ  uintptr
	word unsafe.Pointer
}

// FnPtr returns a real address of function f or nil if f is nil.
//
// If f is not a function, you still get the pointer, but it is unknown
// what that pointer points to.
//
// YOU CAN NOT USE RETURNED POINTER TO CALLING PASSED FUNCTION!
// For that purpose use PtrCallable.
func FnPtr(f interface{}) unsafe.Pointer {

	// goiface represents what "interface{}" means in internal Golang parts.
	type goiface struct {
		typ uintptr        // pointer to the type definition struct
		val unsafe.Pointer // pointer to the value
	}

	if f == nil {
		return nil
	}

	return (*goiface)(unsafe.Pointer(&f)).val
}

// FnPtrCallable returns an address of function f using which you can call
// a function that was passed or nil if f is nil.
// To calling just convert returned untyped pointer to function-typed,
// dereference it and call.
//
// You can AVOID TYPE CHECKS using that way (wrong argument types, wrong
// return types) but in that way the BEHAVIOUR is UNDEFINED and do it only
// if you know what you're doing.
//
// If f is not a function, you still get the pointer, but it is unknown
// what that pointer points to.
func FnPtrCallable(f interface{}) unsafe.Pointer {

	// PtrReal and ptrReal2Callable has nil checks.
	return fnPtrReal2Callable(FnPtr(f))
}

// ptrReal2Callable converts a real function pointer to a pointer using which
// becomes possible to call a function ptr points to.
func fnPtrReal2Callable(ptr unsafe.Pointer) unsafe.Pointer {

	type fptr struct {
		ptr unsafe.Pointer
	}

	if ptr == nil {
		return nil
	}

	o := new(fptr)
	o.ptr = ptr
	return unsafe.Pointer(&o.ptr)
}

// // Ptr returns an address of function f or nil if f is nil.
// // If f is not a function, you still get the pointer, but it is unknown
// // what that pointer points to.
// func Ptr(f interface{}) unsafe.Pointer {

// 	// In Golang you can not just take an address of some func
// 	// only if it is not a func literal which is assigned to any var.
// 	//
// 	// But you can assign any kind of func to type-compatible variable
// 	// or interface{}.
// 	// And you can call then that assigned function through used variable.
// 	// But anyway that way still requires a type-compatible variable
// 	// or type-casting interface{} and you can not break that type-compatible
// 	// rules:
// 	// You can not assign func to some variable with wrong func signature,
// 	// You can not cast interface{} that stores func with one signature
// 	// to func with another.
// 	//
// 	// Unless you have an untyped pointer like void*.
// 	//
// 	// reflect.ValueOf(f).Pointer() returns a REAL function address,
// 	// but you CAN NOT use it to call.
// 	// I do not know how it is implemented in an internal Golang parts, but it is.
// 	//
// 	//
// 	// CAN:
// 	// f := func(){}
// 	// &f
// 	//
// 	// CAN NOT:
// 	// func f(){}
// 	// &f

// 	// But you can assign any kind of func to type-compatible variable
// 	// or interface{}.

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
