// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"sync/atomic"
)

type (
	// NamespaceID is just int32 alias and has been introduced to make it easy
	// to replace underlying type to something other in the future.
	NamespaceID = int32
)

//noinspection GoSnakeCaseUsage
const (
	// _ERR_INVALID_NAMESPACE_ID is a special reserved invalid Namespace ID
	// that is used to mark that some Namespace is invalid.
	_ERR_INVALID_NAMESPACE_ID = -1
)

var (
	// namespaceIDPrivateCounter is an internal ClassID counter that is increased
	// at the new error Namespace creating.
	//
	// DO NOT USE IT DIRECTLY! Use newNamespaceID() instead.
	namespacePrivateCounterID ClassID
)

// newNamespaceID returns new NamespaceID value that can be used as new Namespace's ID.
// Thread-safety.
func newNamespaceID() NamespaceID {
	return atomic.AddInt32(&namespacePrivateCounterID, 1)
}

// isValidNamespaceID reports whether 'namespaceID' is valid Namespace's ID.
// Remember, it has been done to avoid cases when Namespace object has been manually
// instantiated (like 'new(Namespace)' or 'var _ Namespace') and thus has not been initialized.
func isValidNamespaceID(namespaceID NamespaceID) bool {
	// Because in the Namespace object ID is private and can not be accessed
	// from the outside, it's enough to check whether it's not 0 (manually instantiated),
	// but theoretically classID may be equal _ERR_INVALID_NAMESPACE_ID. Check it too.
	return namespaceID > 0 && namespaceID != _ERR_INVALID_NAMESPACE_ID
}

// namespaceByID returns Namespace object bases on 'namespaceID'.
// WARNING! Make sure you checked whether provided 'namespaceID' is valid using
// isValidNamespaceID() func. UB otherwise (may panic).
func namespaceByID(namespaceID NamespaceID, lock bool) Namespace {
	if namespaceID < _ERR_NAMESPACE_ARRAY_CACHE {
		return registeredNamespacesArr[namespaceID]
	} else {
		if lock {
			registeredNamespacesMap.RLock()
			defer registeredNamespacesMap.RUnlock()
		}
		return registeredNamespacesMap.m[namespaceID]
	}
}
