// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"sync"
)

//noinspection GoSnakeCaseUsage

const (
	// _ERR_CLASS_ARRAY_CACHE describes how many registered Classes will be
	// stored into internal array instead of map.
	// See registeredClassesArr, registeredClassesMap and classByID() for mor details.
	_ERR_CLASS_ARRAY_CACHE = 128
)

var (
	// invalidClass is a special Class object that will be returned if we can not
	// provide a natural Class object. For example at the: (*Error)(nil).Class().
	invalidClass = Class{id: _ERR_INVALID_CLASS_ID}

	// registeredClassesArr is a registered Classes storage.
	// First N classes (_ERR_CLASS_ARRAY_CACHE) will be saved into this array,
	// others to the registeredClassesMap.
	//
	// Because array accessing faster than map's one it's an array, but
	// user can define more than N classes. It's very rarely case, because
	// I guess N is high enough but if that will happen, new classes will be stored
	// into map. So, if in your app you have more than N error classes, guess
	// performance array > map is not important to you.
	//
	// It's made to provide (*Error).Class() method working.
	registeredClassesArr [_ERR_CLASS_ARRAY_CACHE]Class

	// registeredClassesMap used to save registered Classes when you have their
	// more than N (_ERR_CLASS_ARRAY_CACHE).
	registeredClassesMap = struct {
		sync.RWMutex
		m map[ClassID]Class
	}{
		m: make(map[ClassID]Class),
	}
)

// fullClassName generates and returns full Class's name
// using 'className' as Class's name, and 'namespaceName' as Namespace's
// (or get a Namespace object from the pool using 'namespaceID' if 'namespaceName' is empty).
func fullClassName(className, namespaceName string, namespaceID NamespaceID) string {

	if namespaceName == "" {
		namespaceName = namespaceByID(namespaceID, false).name
	}
	return namespaceName + "::" + className
}

// rebuiltExistedClassNames changes the names of existed (registered) Classes
// (in the pool) use fullClassName() for their full names.
// It guarantees that this function will be called only once when the first custom
// namespace will be created.
func rebuiltExistedClassNames() {
	registeredClassesMap.Lock()
	defer registeredClassesMap.Unlock()

	for i, cls := range registeredClassesArr {
		registeredClassesArr[i].fullName =
			fullClassName(cls.fullName, "", cls.namespaceID)
	}

	for clsID, cls := range registeredClassesMap.m {
		cls.fullName = fullClassName(cls.fullName, "", cls.namespaceID)
		registeredClassesMap.m[clsID] = cls
	}
}

// newClass is a Class's constructor.
// There are several steps:
//
// 1. Getting a new available Class's ID.
//
// 2. Create a new Class object using provided 'parentID', 'namespaceID', 'name'.
//
// 3. Save it into the internal Class's storage basing on the its (Class's) ID.
//
// 4. Done.
func newClass(parentID ClassID, namespaceID NamespaceID, name, fullName string) Class {

	c := Class{
		id:          newClassID(),
		parentID:    parentID,
		namespaceID: namespaceID,
		name:        name,
		fullName:    fullName,
	}

	if c.id >= _ERR_CLASS_ARRAY_CACHE {
		registeredClassesMap.Lock()
		defer registeredClassesMap.Unlock()
		registeredClassesMap.m[c.id] = c
	} else {
		registeredClassesArr[c.id] = c
	}

	return c
}
