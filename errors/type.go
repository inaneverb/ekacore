// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package errors

import (
	"fmt"

	"github.com/modern-go/reflect2"
)

//
type Type struct {
	id       uint64
	parentId uint64

	name     string
	fullName string
}

//
func (t Type) NewSubtype(subtypeName string) Type {

	return newType(getPrivateID(), t.id, subtypeName, t.name)
}

//
func (t Type) New(message string, args ...interface{}) *Error {
	return t.new(message, args)
}

//
func (t Type) new(message string, args []interface{}) *Error {

	e := newError(t.parentId, t.fullName)
	if e == nil {
		return nil
	}

	if len(args) != 0 {
		e.Message = fmt.Sprintf(message, args...)
	} else {
		e.Message = message
	}

	return e
}

//
func (t Type) Wrap(err error, message string, args ...interface{}) *Error {

	// TODO: Maybe call err.Error() and check whether it's empty?
	if reflect2.IsNil(err) {
		return nil
	}

	e := t.new(message, args)
	if e == nil {
		return nil
	}

	e.Hidden = err.Error()
	return e
}

////
//func (t Type) Embed(err error, message string, args ...interface{}) *Error {
//	// TODO: Implement (embed error to public part)
//	return nil
//}
//
////
//func (t Type) Hide(err error, message string, args ...interface{}) *Error {
//	// TODO: Implement (hide error from public)
//	return nil
//}

//
func newType(id, parentId uint64, name, parentName string) Type {

	return Type{
		id:       id,
		parentId: parentId,
		name:     name,
		fullName: parentName + "." + name,
	}
}
