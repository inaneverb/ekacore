// Copyright © 2022. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekasys

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/inaneverb/ekacore/ekadeath/v4"
)

type ObserveController struct {
	ctx  context.Context
	wg   sync.WaitGroup
	sema sync.Mutex
}

// ObserveCallback is just a function alias for callback
// that you will use in Observe() function.
type ObserveCallback func(oc *ObserveController)

const (
	mutexLocked = 1 << iota
)

var (
	// observeCalled is a special package-level variable,
	// that checks whether Observe() is called only once using atomic CAS operation.
	observeCalled uint32
)

func (oc *ObserveController) Join() bool {
	ptr := unsafe.Pointer(&oc.sema)
	if !atomic.CompareAndSwapInt32((*int32)(ptr), 0, mutexLocked) {
		return false
	}

	oc.wg.Add(1)
	oc.sema.Unlock()
	return true
}

func (oc *ObserveController) Leave() {
	oc.wg.Done()
}

func (oc *ObserveController) Context() context.Context {
	return oc.ctx
}

func (oc *ObserveController) wait() {
	oc.sema.Lock()
	oc.wg.Wait()
}

func newObserveController(ctx context.Context) *ObserveController {
	return &ObserveController{ctx: ctx}
}

// Observe is your service initializer.
// It provides a mechanism of graceful shutdown, supporting OS signals (like SIGINT)
// and requested shutdowns from ekadeath or ekalog.
//
// First of all. You need a special "main" context.
// If you pass nil as context.Context, the context.Background() will be used.
// The difference is:
//   - Overwritten background context can be canceled ONLY by OS signals;
//   - Your context can be cancelled by your hands manually AND by OS signals.
//
// This context will be used to be able to say your jobs that it's time to stop
// their work and they need to be completed completely.
//
// WARNING!
// You can call Observe() ONLY ONE per whole lifetime of your application.
// The next call will lead to panic.
func Observe(ctx context.Context, cb ObserveCallback) {

	if !atomic.CompareAndSwapUint32(&observeCalled, 0, 1) {
		panic("Observe must be called only once at all!")
	}

	if cb == nil {
		panic("Observe initialization callback must not be nil")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	var cancelFunc func()
	ctx, cancelFunc = context.WithCancel(ctx)

	var oc = newObserveController(ctx)

	ekadeath.Reg(func() { cancelFunc(); oc.wait() })
	cb(oc)
}
