// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package letter

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/qioalice/ekago/ekadanger"
	"github.com/qioalice/ekago/internal/field"

	"github.com/modern-go/reflect2"
)

// addImplicitField adds new field to l.Fields treating 'name' as field's name,
// 'value' as field's value and using 'typ' (assuming it's value's type)
// to recognize how to convert Golang's interface{} to the 'Field' object.
func (li *LetterItem) addImplicitField(name string, value interface{}, typ reflect2.Type) {

	varyField := name != "" && name[len(name)-1] == '?'
	if varyField {
		name = name[:len(name)-1]
	}

	var f field.Field

	switch {
	case value == nil && varyField:
		// do nothing
		return

	case value == nil:
		li.Fields = append(li.Fields, field.NilValue(name, field.KIND_TYPE_INVALID))
		return

	case typ.Implements(field.ReflectedTypeFmtStringer):
		var stringer fmt.Stringer
		typ.Set(&stringer, &value)
		f = field.Stringer(name, stringer)
		goto recognizer
	}

	switch typ.Kind() {

	case reflect.Ptr:
		// Maybe it's typed nil pointer? We can't dereference nil pointer
		// but I guess it's important to log nil pointer as is even if
		// FLAG_ALLOW_IMPLICIT_POINTERS is not set (because what can we do otherwise?)
		logPtrAsIs :=
			li.Flags.TestAll(FLAG_ALLOW_IMPLICIT_POINTERS) ||
				ekadanger.TakeRealAddr(value) == nil

		if logPtrAsIs {
			f = field.Addr(name, value)
		} else {
			value = typ.Indirect(value)
			li.addImplicitField(name, value, reflect2.TypeOf(value))
		}

	case reflect.Bool:
		var boolVal bool
		typ.UnsafeSet(unsafe.Pointer(&boolVal), reflect2.PtrOf(value))
		f = field.Bool(name, boolVal)

	case reflect.Int:
		var intVal int
		typ.UnsafeSet(unsafe.Pointer(&intVal), reflect2.PtrOf(value))
		f = field.Int(name, intVal)

	case reflect.Int8:
		var int8Val int8
		typ.UnsafeSet(unsafe.Pointer(&int8Val), reflect2.PtrOf(value))
		f = field.Int8(name, int8Val)

	case reflect.Int16:
		var int16Val int16
		typ.UnsafeSet(unsafe.Pointer(&int16Val), reflect2.PtrOf(value))
		f = field.Int16(name, int16Val)

	case reflect.Int32:
		var int32Val int32
		typ.UnsafeSet(unsafe.Pointer(&int32Val), reflect2.PtrOf(value))
		f = field.Int32(name, int32Val)

	case reflect.Int64:
		var int64Val int64
		typ.UnsafeSet(unsafe.Pointer(&int64Val), reflect2.PtrOf(value))
		f = field.Int64(name, int64Val)

	case reflect.Uint:
		var uintVal uint64
		typ.UnsafeSet(unsafe.Pointer(&uintVal), reflect2.PtrOf(value))
		f = field.Uint64(name, uintVal)

	case reflect.Uint8:
		var uint8Val uint8
		typ.UnsafeSet(unsafe.Pointer(&uint8Val), reflect2.PtrOf(value))
		f = field.Uint8(name, uint8Val)

	case reflect.Uint16:
		var uint16Val uint16
		typ.UnsafeSet(unsafe.Pointer(&uint16Val), reflect2.PtrOf(value))
		f = field.Uint16(name, uint16Val)

	case reflect.Uint32:
		var uint32Val uint32
		typ.UnsafeSet(unsafe.Pointer(&uint32Val), reflect2.PtrOf(value))
		f = field.Uint32(name, uint32Val)

	case reflect.Uint64:
		var uint64Val uint64
		typ.UnsafeSet(unsafe.Pointer(&uint64Val), reflect2.PtrOf(value))
		f = field.Uint64(name, uint64Val)

	case reflect.Float32:
		var float32Val float32
		typ.UnsafeSet(unsafe.Pointer(&float32Val), reflect2.PtrOf(value))
		f = field.Float32(name, float32Val)

	case reflect.Float64:
		var float64Val float64
		typ.UnsafeSet(unsafe.Pointer(&float64Val), reflect2.PtrOf(value))
		f = field.Float64(name, float64Val)

	case reflect.Complex64:
		var complex64Val complex64
		typ.UnsafeSet(unsafe.Pointer(&complex64Val), reflect2.PtrOf(value))
		f = field.Complex64(name, complex64Val)

	case reflect.Complex128:
		var complex128Val complex128
		typ.UnsafeSet(unsafe.Pointer(&complex128Val), reflect2.PtrOf(value))
		f = field.Complex128(name, complex128Val)

	case reflect.String:
		var stringVal string
		typ.UnsafeSet(unsafe.Pointer(&stringVal), reflect2.PtrOf(value))
		f = field.String(name, stringVal)

	case reflect.Uintptr:
		var uintptrVal uintptr
		typ.Set(&uintptrVal, &value)
		f = field.Addr(name, uintptrVal)

	case reflect.UnsafePointer:
		var unsafePtrVal unsafe.Pointer
		typ.Set(&unsafePtrVal, &value)
		f = field.Addr(name, unsafePtrVal)
	}

recognizer:
	if !(varyField && f.IsZero()) {
		li.Fields = append(li.Fields, f)
	}
}

// addExplicitField assumes that 'field' either 'Field' or '*Field' type
// (checks it by 'fieldType') and adds it to the l.Fields.
// If field == (*Field)(nil) there is no-op.
func (li *LetterItem) addExplicitField(fieldValue interface{}, fieldType reflect2.Type) {

	var fPtr *field.Field

	if fieldType == field.ReflectedTypePtr {
		fPtr = fieldValue.(*field.Field)

	} else {
		fPtr = new(field.Field)
		*fPtr = fieldValue.(field.Field)
	}

	if fPtr != nil {
		varyField := fPtr.Key != "" && fPtr.Key[len(fPtr.Key)-1] == '?'
		if varyField {
			fPtr.Key = fPtr.Key[:len(fPtr.Key)-1]
		}
		if !(varyField && fPtr.IsZero()) {
			li.Fields = append(li.Fields, *fPtr)
		}
	}
}
