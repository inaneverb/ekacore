// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package field

// IsZero reports whether f contains zero value of its type (based on kind).
func (f Field) IsZero() bool {
	switch f.Kind.BaseType() {

	case KIND_TYPE_BOOL,
		KIND_TYPE_INT, KIND_TYPE_INT_8, KIND_TYPE_INT_16, KIND_TYPE_INT_32, KIND_TYPE_INT_64,
		KIND_TYPE_UINT, KIND_TYPE_UINT_8, KIND_TYPE_UINT_16, KIND_TYPE_UINT_32, KIND_TYPE_UINT_64,
		KIND_TYPE_UINTPTR, KIND_TYPE_ADDR,
		KIND_TYPE_FLOAT_32, KIND_TYPE_FLOAT_64:
		return f.IValue == 0

	case KIND_TYPE_STRING:
		return f.SValue == ""

	default:
		return f.Value == nil
	}
}
