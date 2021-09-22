// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

type (
	// Class is a special type that represents Error's abstract class
	// and provides a mechanism of error classifying.
	//
	// Using Class's object you can instantiate an *Error object
	// using New() or Wrap() methods. See them for more info.
	//
	// DO NOT INSTANTIATE Class OBJECTS MANUALLY! THEY WILL NOT BE INITIALIZED
	// PROPERLY AND WILL BE CONSIDERED BROKEN. THAT CLASS IS NOT VALID.
	// YOU CAN NOT CREATE ERROR OBJECTS USING INVALID CLASS.
	// YOU CAN NOT CREATE ANOTHER DERIVED CLASSES USING INVALID CLASS AS BASE.
	//
	// A Class can not exist w/o a Namespace. Ekaerr provides you builtin classes
	// you may use and, typically, that covers ~90% all error classifying models.
	// But if you need your own classes or even subclasses (yes, you can derive them)
	// you may:
	//
	// - Create your own Class.
	//   Use <Namespace>.NewClass(...) method.
	//   You can use builtin "Common" namespace (allowed by CommonNamespace variable)
	//   or create your own Namespace at first and then create a Class that you want.
	//
	// - Create a new class but derive them from another.
	//   Use <Class>.NewSubClass(...) method.
	//   You can use builtin classes (see them at the "class_builtin.go" file)
	//   or crate your own Class (see above how) and derive from them.
	//
	// You also need to know that there is a special Class - an 'invalidClass'.
	// This is a special class that you may get if you instantiate Class object
	// manually like "new(Class)" or "var c Class". That class IS NOT VALID.
	// YOU MUST NOT USE INVALID CLASS!
	// But you may get an invalid class also if you will call "(*Error)(nil).Class()".
	// So, you can check whether your class is valid using IsValid() method.
	//
	// Class is a very lightweight datatype and all Class's methods
	// (and function that takes Class as argument) uses Class's objects by value.
	// It means there is no reason to use Class's object by reference in your code.
	Class struct {

		// id is an unique per Class its ID.
		// If it's == _ERR_INVALID_CLASS_ID, the Class is considered broken
		// (has been instantiated manually instead of using Class constructors).
		id ClassID

		// parentID is an ID of some Class that is base Class for this Class.
		// (it means that this Class has been created by Class.NewSubClass()).
		parentID ClassID

		// namespaceID is an ID of some Namespace this Class belongs to.
		// (it means that this Class has been created by Namespace.NewClass() or
		// derived from the Class that was created that way).
		namespaceID NamespaceID

		// name is this Class's name specified by user at the creation.
		// It's just a class's name as is, nothing more.
		name string

		// fullName is this Class's name but with base Classes and Namespaces
		// chaining. E.g: "<Namespace>.<BaseClass1>.<BaseClass2>. ... .<ThisClass>".
		//
		// WARNING!
		// READ THIS FIELD ONLY OF OBJECTS YOU OBTAIN FROM THE CLASS'S POOL!
		fullName string
	}
)

// New is an Error's constructor. Specify what happen by 'message' and key-value
// paired arguments 'args' and that is! A new *Error object is returned.
// Arguments adding works the same as (*Error).AddFields() do.
//
// Requirements:
// c must be valid Class object. Otherwise nil Error is returned.
func (c Class) New(message string, args ...interface{}) *Error {
	if !isValidClassID(c.id) {
		return nil
	}
	return newError(false, c.id, c.namespaceID, nil, message, args)
}

// LightNew is the same as just New() but creates a lightweight Error instead.
// Read more what lightweight error is in Error's doc.
func (c Class) LightNew(message string, args ...interface{}) *Error {
	if !isValidClassID(c.id) {
		return nil
	}
	return newError(true, c.id, c.namespaceID, nil, message, args)
}

// Wrap is an Error's constructor. Specify what legacy Golang error you need
// to wrap using 'err', what happen by 'message' and key-value paired arguments 'args'
// and that is! A new *Error object is returned.
// Arguments adding works the same as (*Error).AddFields() do.
//
// Requirements:
// c must be valid Class object. Otherwise nil Error is returned.
// 'err' != nil. Otherwise nil Error is returned.
func (c Class) Wrap(err error, message string, args ...interface{}) *Error {
	if !isValidClassID(c.id) || err == nil {
		return nil
	}
	return newError(false, c.id, c.namespaceID, err, message, args)
}

// LightWrap is the same as just Wrap() but creates a lightweight Error instead.
// Read more what lightweight error is in Error's doc.
func (c Class) LightWrap(err error, message string, args ...interface{}) *Error {
	if !isValidClassID(c.id) || err == nil {
		return nil
	}
	return newError(true, c.id, c.namespaceID, err, message, args)
}

// IsValid reports whether c is valid Class object or not.
//
// It returns false if c has not been initialized properly (instantiated manually
// instead of Class's or Namespace's constructors calling or obtaining from the
// Error's Class() method).
func (c Class) IsValid() bool {
	return isValidClassID(c.id)
}

// ParentClass returns a copy of Class object that has been used as a base class
// for the current's one or a special 'invalidClass' if this Class has been
// created directly from the namespace not as a subclass.
func (c Class) ParentClass() Class {
	if !c.IsValid() || !isValidClassID(c.parentID) {
		return invalidClass
	}
	return classByID(c.parentID, true)
}

// Name returns a current Class's name that was used at the Class creation.
// If you want to get a full name (without Namespace's and Base classes' use
// c.FullName() instead).
func (c Class) Name() string {
	if !c.IsValid() {
		return ""
	}
	return c.name
}

// FullName returns a current Class's full name:
// a Class's name but with base Classes and Namespaces chaining.
// E.g: "<Namespace>.<BaseClass1>.<BaseClass2>. ... .<ThisClass>".
func (c Class) FullName() string {
	if !c.IsValid() {
		return ""
	}
	return classByID(c.id, true).fullName
}

// NewSubClass is a Class's constructor. Specify the derived Class's name 'name'
// and that is! A new Class will be created and its copy is returned
// (created class will belongs to the same Namespace as based one).
//
// Warnings:
// A two classes with the same names is the two DIFFERENT classes!
//
// Requirements:
// c must be valid Class object. Otherwise 'invalidClass' is returned.
func (c Class) NewSubClass(subClassName string) Class {
	if !isValidClassID(c.id) {
		return invalidClass
	}
	fullName := classByID(c.id, true).fullName + "." + subClassName
	return newClass(c.id, c.namespaceID, subClassName, fullName)
}
