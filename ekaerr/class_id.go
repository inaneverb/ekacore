// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"sync/atomic"
)

type (
	// ClassID is just int32 alias and has been introduced to make it easy
	// to replace underlying type to something other in the future.
	ClassID = int32
)

//noinspection GoSnakeCaseUsage
const (
	// _ERR_INVALID_CLASS_ID is a special reserved invalid Class ID
	// that is used to mark that some Class is invalid.
	_ERR_INVALID_CLASS_ID = -1
)

var (
	// classIDPrivateCounter is an internal ClassID counter that is increased
	// at the new error Class creating.
	//
	// DO NOT USE IT DIRECTLY! Use newClassID() instead.
	classIDPrivateCounter ClassID
)

// newClassID returns new ClassID value that can be used as new Class's ID.
// Thread-safety.
func newClassID() ClassID {
	return atomic.AddInt32(&classIDPrivateCounter, 1)
}

// isValidClassID reports whether 'classID' is valid Class's ID.
// Remember, it has been done to avoid cases when Class object has been manually
// instantiated (like 'new(Class)' or 'var _ Class') and thus has not been initialized.
func isValidClassID(classID ClassID) bool {
	// Because in the Class object ID is private and can not be accessed
	// from the outside, it's enough to check whether it's not 0 (manually instantiated),
	// but theoretically classID may be equal _ERR_INVALID_CLASS_ID. Check it too.
	return classID > 0 && classID != _ERR_INVALID_CLASS_ID
}

// classByID returns Class object bases on 'classID' and if Class with provided
// 'classID' stored into the Classes' map, 'lock' indicates whether map's access
// must be protected by its R mutex or not.
//
// WARNING! Make sure you checked whether provided 'classID' is valid using
// isValidClassID() func. UB otherwise (may panic).
func classByID(classID ClassID, lock bool) Class {
	if classID < _ERR_CLASS_ARRAY_CACHE {
		return registeredClassesArr[classID]
	} else {
		if lock {
			registeredClassesMap.RLock()
			defer registeredClassesMap.RUnlock()
		}
		return registeredClassesMap.m[classID]
	}
}
