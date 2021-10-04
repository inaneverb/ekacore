// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaserv

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/qioalice/ekago/v3/ekadeath"
)

type (
	// ObserveCallback is just a function alias for callback
	// that you will use in Observe() function.
	ObserveCallback func(ctx context.Context, wg *sync.WaitGroup)
)

var (
	// observeCalled is a special package-level variable,
	// that checks whether Observe() is called only once using atomic CAS operation.
	observeCalled uint32
)

// Observe is your service initializer.
// It provides a mechanism of graceful shutdown, supporting OS signals (like SIGINT)
// and requested shutdowns from ekadeath or ekalog.
//
// First of all. You need a special "main" context.
// If you pass nil as context.Context, the context.Background() will be used.
// The difference is:
//  - Overwritten background context can be canceled ONLY by OS signals;
//  - Your context can be cancelled by your hands manually AND by OS signals.
// This context will be used to be able to say your jobs that it's time to stop
// their work and they need to be completed completely.
//
// Next thing is sync.WaitGroup.
// This is the same logic. If you pass nil, the new sync.WaitGroup will be created
// and used. The limitations and difference is also the same.
// You can do some job with your sync.WaitGroup expanding their area of responsibility,
// created sync.WaitGroup will live only inside Observe() and its internal parts.
//
// The last thing is your callback.
// Think about that callback like it's your "main()" function but with
// sync.WaitGroup's based protector and context.Context's based shutdown notifier.
// So, you can run your jobs, pass them shutdown context (or modified)
// and wait group to be able to wait until their jobs are completed.
// This is the main purpose of that function.
//
// WARNING!
// You can call Observe() ONLY ONE per whole lifetime of your binary.
// The next call will lead to panic.
//
// NOTE.
// You may think that it's always better to use your own context and wait group,
// but actually it depends on your needs. Moreover it's normal and better practice
// if you don't need such context and wait group outside the Observe()
// according with design pattern. Thus an ability to use your context and/or wait group
// is an "extension" to main functionality and purpose of Observe().
func Observe(parentCtx context.Context, parentWg *sync.WaitGroup, cb ObserveCallback) {

	if !atomic.CompareAndSwapUint32(&observeCalled, 0, 1) {
		panic("Observe must be called only once at all!")
	}

	if cb == nil {
		panic("Observe initialization callback must not be nil")
	}

	if parentWg == nil {
		parentWg = new(sync.WaitGroup)
	}

	if parentCtx == nil {
		parentCtx = context.Background()
	}

	parentCtx, cancelFunc := context.WithCancel(parentCtx)

	parentWg.Add(1)
	ekadeath.Reg(func() {
		cancelFunc()
		parentWg.Wait()
	})

	cb(parentCtx, parentWg)
	parentWg.Done()
}
