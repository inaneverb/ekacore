// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr

import (
	"sync"
	"sync/atomic"
)

//noinspection GoSnakeCaseUsage
const (
	// _ERR_NAMESPACE_ARRAY_CACHE describes how many registered Namespaces will be
	// stored into internal array instead of map.
	// See registeredClassesArr, registeredClassesMap and classByID() for mor details.
	_ERR_NAMESPACE_ARRAY_CACHE = 16
)

var (
	// isCustomNamespaceDefined reports whether at least one custom namespace
	// has been defined.
	//
	// DO NOT USE IT DIRECTLY! Use customNamespaceDefined() instead.
	isCustomNamespaceDefined int32

	// registeredNamespacesArr is a registered Namespaces storage.
	// First N namespaces (_ERR_NAMESPACE_ARRAY_CACHE) will be saved into this array,
	// others to the registeredNamespacesMap.
	//
	// Because array accessing faster than map's one it's an array, but
	// user can define more than N namespaces. It's very rarely case, because
	// I guess N is high enough but if that will happen, new namespaces will be stored
	// into map. So, if in your app you have more than N error namespaces, guess
	// performance array > map is not important to you.
	//
	// It's made to provide (*Error).Namespace() method working.
	registeredNamespacesArr [_ERR_NAMESPACE_ARRAY_CACHE]Namespace

	// registeredNamespacesMap used to save registered Namespaces when you have their
	// more than N (_ERR_NAMESPACE_ARRAY_CACHE).
	registeredNamespacesMap = struct {
		sync.RWMutex
		m map[NamespaceID]Namespace
	}{
		m: make(map[NamespaceID]Namespace),
	}
)

// customNamespaceDefined reports whether at least one custom namespace has been
// created by the user. Thread-safety.
func customNamespaceDefined() bool {
	return atomic.LoadInt32(&isCustomNamespaceDefined) != 0
}

// newNamespace is a Namespace's constructor.
// There are several steps:
//
// 1. Getting a new available Namespace's ID.
//
// 2. Create a new Namespace object using provided 'name'.
//
// 3. Save it into internal Namespace's storage basing on its (Namespace's) ID.
//
// 4. If it was first call newNamespace() with 'custom' == true,
//    regenerate all existed Classes' names adding namespace's names to the
//    their full names.
//
// 5. Done.
func newNamespace(name string, custom bool) Namespace {

	if custom {
		firstCustomNamespace :=
			atomic.CompareAndSwapInt32(&isCustomNamespaceDefined, 0, 1)
		if firstCustomNamespace {
			rebuiltExistedClassNames()
		}
	}

	n := Namespace{
		id:   newNamespaceID(),
		name: name,
	}

	if n.id >= _ERR_NAMESPACE_ARRAY_CACHE {
		registeredNamespacesMap.Lock()
		defer registeredNamespacesMap.Unlock()
		registeredNamespacesMap.m[n.id] = n
	} else {
		registeredNamespacesMap.m[n.id] = n
	}

	return n
}
