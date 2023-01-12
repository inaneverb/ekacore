// Copyright © 2020. All rights reserved.
// Author: Eagle Chen. Modifier: Ilya Stroy.
// Original: https://github.com/EagleChen/mapmutex (c133e97)
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatyp

import (
	"math/rand"
	"sync"
	"time"
)

type (
	// Mutex is the mutex with synchronized map
	// it's for reducing unnecessary locks among different keys
	MuMap struct {

		// m is the whole MuMap's mutex. Each lock/unlock operation captures this mutex,
		// but it's not captured when Lock() / TryLock() is waiting for some delay
		// between key's capturing attempts.
		m *sync.Mutex

		// key's "mutexes"
		// map's key it's key; value is a counter.
		locks map[any]int8

		maxRetry  int // how much TryLock() will tries to capture key's "mutex"
		maxRetryR int // how much RTryLock() will tries to capture key's "mutex"

		maxDelay  time.Duration // maximum delay between key's "mutex" capturing attempts
		baseDelay time.Duration // base delay between key's "mutex" capturing attempts

		factor float64 // multiplier of delay between key's "mutex" capturing attempts
		jitter float64 // random for factor
	}
)

func (m *MuMap) RLock(key any) {
	m.assertInitialized()
	m.lock(key, true, true)
}

func (m *MuMap) Lock(key any) {
	m.assertInitialized()
	m.lock(key, true, false)
}

func (m *MuMap) RTryLock(key any) (gotLock bool) {
	m.assertInitialized()
	return m.lock(key, false, true)
}

// TryLock tries to aquire the lock.
func (m *MuMap) TryLock(key any) (gotLock bool) {
	m.assertInitialized()
	return m.lock(key, false, false)
}

func (m *MuMap) RUnlock(key any) {
	m.assertInitialized()
	m.unlock(key, true)
}

// Unlock unlocks for the key
// please call Unlock only after having aquired the lock
func (m *MuMap) Unlock(key any) {
	m.assertInitialized()
	m.unlock(key, false)
}

// assertInitialized checks whether m is initialized and initialized properly.
func (m *MuMap) assertInitialized() {
	if m == nil || m.locks == nil {
		panic("MuMap is not initialized properly")
	}
}

func (m *MuMap) lock(key any, untilSuccess, readOnly bool) (gotLock bool) {

	// First attempt will be done with -1 as attempt's index.
	// -1 does no delay and no sleep before trying to lock
	// 0 does m.baseDelay sleeping before trying to lock
	// any next value does some calculated sleep time before trying to lock

	for i := -1; i < m.maxRetry; i++ {
		if m.lockIter(key, i, readOnly) {
			return true
		}
	}

	if !untilSuccess {
		return false
	}

	for !m.lockIter(key, m.maxRetry, readOnly) {
	}
	return true
}

func (m *MuMap) lockIter(key any, attempt int, readOnly bool) (gotLock bool) {

	var sleepDur time.Duration
	switch {

	case attempt == -1:
		sleepDur = 0

	case attempt == 0:
		sleepDur = m.baseDelay

	case attempt >= m.maxRetry:
		sleepDur = m.maxDelay

	default:
		backoff, max := float64(m.baseDelay), float64(m.maxDelay)
		for ; backoff < max && attempt > 0; attempt-- {
			backoff *= m.factor
		}
		if backoff > max {
			backoff = max
		}
		backoff *= 1 + m.jitter*(rand.Float64()*2-1)
		if backoff < 0 {
			sleepDur = 0
		}
		sleepDur = time.Duration(backoff)
	}

	time.Sleep(sleepDur) // does nothing if sleepDur == 0

	m.m.Lock()

	keyCounter, alreadyLocked := m.locks[key]
	switch {

	case readOnly && keyCounter != -1 && keyCounter != 127:
		// state - RLock or Released, need - RLock
		m.locks[key]++
		gotLock = true

	case !readOnly && !alreadyLocked:
		// state - Released, need - Lock
		m.locks[key] = -1
		gotLock = true
	}

	m.m.Unlock()

	return gotLock
}

func (m *MuMap) unlock(key any, readOnly bool) {

	m.m.Lock()

	keyCounter, alreadyLocked := m.locks[key]
	switch {

	case !alreadyLocked:
		m.m.Unlock()
		panic("MuMap: Can not unlock unlocked mutex")

	case keyCounter > 0 && !readOnly:
		// state - RLock, requested - Unlock
		m.m.Unlock()
		panic("MuMap: Can not unlock read-locked mutex")

	case keyCounter == -1 && readOnly:
		// state - Lock, requested - RUnlock
		m.m.Unlock()
		panic("MuMap: Can not read-unlock locked mutex")

	case keyCounter == 1 || keyCounter == -1:
		// state - RLock, requested - RUnlock, no more other RLock calls or
		// state - Lock, requested - Unlock
		delete(m.locks, key)

	case keyCounter > 1:
		// state - RLock, requested - RUnlock, there is at least one more RLock call
		m.locks[key]--

	default:
		m.m.Unlock()
		panic("MuMap: Unexpected state for unlocking mutex")
	}

	m.m.Unlock()
}

// NewMapMutex returns a mapmutex with default configs
func NewMuMap() *MuMap {
	return NewMuMapCustom(180, 1*time.Second, 10*time.Nanosecond, 1.1, 0.2)
}

// NewCustomizedMapMutex returns a customized mapmutex
func NewMuMapCustom(mRetry int, mDelay, bDelay time.Duration, factor, jitter float64) *MuMap {
	return &MuMap{
		locks:     make(map[any]int8),
		m:         &sync.Mutex{},
		maxRetry:  mRetry,
		maxDelay:  mDelay,
		baseDelay: bDelay,
		factor:    factor,
		jitter:    jitter,
	}
}
