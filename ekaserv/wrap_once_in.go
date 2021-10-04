// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaserv

import (
	"context"
	"sync"

	"github.com/qioalice/ekago/v3/ekatime"
)

type (
	// OnceInWithContextCallback is just a function alias for callback
	// that you may use in WrapOnceInWithContext() function.
	OnceInWithContextCallback func(ctx context.Context, ts ekatime.Timestamp)

	// OnceInWithContextWaitGroupCallback is just a function alias for callback
	// that you may use in WrapOnceInWithContextWaitGroup() function.
	OnceInWithContextWaitGroupCallback func(ctx context.Context, wg *sync.WaitGroup, ts ekatime.Timestamp)
)

// WrapOnceInWithContext allows you to get ekatime.OnceInCallback
// from your OnceInWithContextCallback,
// provided both context.Context `ctx` and sync.WaitGroup `wg`.
//
// Your callback will receive your provided context.Context and ekatime.Timestamp,
// that you may use.
//
// Provided sync.WaitGroup will be captured and used (+1 before calling callback, -1 after)
// only if it's not nil. So it's safe to pass nil `wg` if you don't mind.
//
// It's useful to use this wrapper inside Observe()'s callback.
func WrapOnceInWithContext(ctx context.Context, wg *sync.WaitGroup, cb OnceInWithContextCallback) ekatime.OnceInCallback {
	return func(ts ekatime.Timestamp) {
		if wg != nil {
			wg.Add(1)
			defer wg.Done()
		}
		cb(ctx, ts)
	}
}

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
func WrapOnceInWithContextWaitGroup(ctx context.Context, wg *sync.WaitGroup, cb OnceInWithContextWaitGroupCallback) ekatime.OnceInCallback {
	return func(ts ekatime.Timestamp) {
		cb(ctx, wg, ts)
	}
}
