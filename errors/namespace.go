// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package errors

//
type Namespace struct {
	id uint64

	name string
}

//
func (n Namespace) NewType(typeName string) Type {

	return newType(getPrivateID(), n.id, typeName, n.name)
}

//
func NewNamespace(name string) Namespace {

	return Namespace{
		id:   getPrivateID(),
		name: name,
	}
}
