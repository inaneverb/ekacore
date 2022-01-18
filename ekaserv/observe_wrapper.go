// Copyright © 2022. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaserv

import (
	"github.com/qioalice/ekago/v3/ekatime"
)

type (
	// OnceInWithContextCallback is just a function alias for callback
	// that you may use in WrapOnceInWithContext() function.

	// TODO: Comment
	OnceInForObserveCallback func(oc *ObserveController, ts ekatime.Timestamp)
)

// WrapOnceInWithContextWaitGroup allows you to get ekatime.OnceInCallback
// from your OnceInWithContextWaitGroupCallback,
// provided both context.Context `ctx` and sync.WaitGroup `wg`.
//
// Your callback will receive your provided both of context.Context and sync.WaitGroup,
// along with ekatime.Timestamp that you may use.
//
// The main difference from WrapOnceInWithContext() is sync.WaitGroup's handling:
// This function will pass `wg` to your `cb` callback as is w/o any preparations.
// As an opposite, the WrapOnceInWithContext() won't pass `wg` to your callback
// but "protects" the call of your callback using sync.WaitGroup.
//
// It's useful to use this wrapper inside Observe()'s callback.

// TODO: Comment
func WrapOnceInWithContextWaitGroup(oc *ObserveController, cb OnceInForObserveCallback) ekatime.OnceInCallback {
	return func(ts ekatime.Timestamp) {
		cb(oc, ts)
	}
}
