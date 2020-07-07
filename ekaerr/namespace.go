// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

type (
	// Namespace is a special type that represents Error's and Class's abstract namespace
	// and provides a mechanism of error and its classes classifying.
	//
	// Using Namespace's object you can register (and instantiate) a new Class.
	// WARNING! If you create a two classes with the same name it will be a two
	// different classes.
	//
	// DO NOT INSTANTIATE Namespace OBJECTS MANUALLY! THEY WILL NOT BE INITIALIZED
	// PROPERLY AND WILL BE CONSIDERED BROKEN. THAT NAMESPACE IS NOT VALID.
	// YOU CAN NOT CREATE A CLASS OBJECTS USING INVALID NAMESPACE.
	//
	// A Namespace is the entry point of the Error's creating chain:
	// Namespace -> Class -> Error.
	// Ekaerr provides you builtin namespace "CommonNamespace" you may use.
	// But if you need your own namespace you can create it using NewNamespace().
	//
	// Namespace is a very lightweight datatype and all Namespace's methods
	// (and functions that takes Namespace as argument) uses Namespace's object by value.
	// It means there is no reason to pass Namespace's object by reference in your code.
	Namespace struct {

		// id is an unique per Namespace its ID.
		// If it's == _ERR_INVALID_NAMESPACE_ID, the Class is considered broken
		// (has been instantiated manually instead of using Namespace constructors).
		id NamespaceID

		// name is this Namespace's name specified by user at the creation.
		// It's just a namespace's name as is, nothing more.
		name string
	}
)

// NewClass is a Class's constructor. Specify the Class's name 'name' and that is!
// A new Class will be created and its copy is returned.
//
// Warnings:
// A two classes with the same names is the two DIFFERENT classes!
//
// Requirements:
// n must be valid Namespace object. Otherwise 'invalidClass' is returned.
func (n Namespace) NewClass(name string) Class {
	if !isValidNamespaceID(n.id) {
		return invalidClass
	}
	fullName := name
	if customNamespaceDefined() {
		fullName = fullClassName(name, n.name, n.id)
	}
	return newClass(_ERR_INVALID_CLASS_ID, n.id, name, fullName)
}

// NewNamespace is a Namespace's constructor. Specify the Namespace's name 'name'
// and that is! A new Namespace will be created and its copy is returned.
//
// Warnings:
// A two namespaces with the same names is the two DIFFERENT namespaces!
func NewNamespace(name string) Namespace {
	return newNamespace(name, true)
}
