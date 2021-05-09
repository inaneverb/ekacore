/*
MIT License

Copyright (c) 2018 ef-ds

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

/*
Original package: https://github.com/ef-ds/deque
Design draft about deque + thread safe:
https://github.com/golang/proposal/blob/master/design/27935-unbounded-queue-package.md
*/

package ekatyp

import (
	"sync"

	"github.com/ef-ds/deque"
)

type (
	// Deque is double-ended queue providing both of FIFO, LIFO design.
	// Thanks to https://github.com/ef-ds/deque it's blazing fast.
	// Deque must not be used by value, only by reference.
	// Deque is thread UNSAFE. Use DequeSafe if you need thread safety version.
	Deque = deque.Deque

	// DequeSafe is the same as Deque but provides thread-safety operations
	// protecting them by sync.Mutex.
	// DequeSafe must not be used by value, only by reference.
	DequeSafe struct {
		q Deque
		m sync.Mutex
	}
)

// Init initializes or clears thread safe double ended queue.
func (dq *DequeSafe) Init() *DequeSafe {
	if dq == nil {
		return NewDequeSafe()
	}
	dq.m.Lock()
	dq.q.Init()
	dq.m.Unlock()
	return dq
}

// Len returns the number of elements of thread safe double ended queue.
// The complexity is O(1).
// DequeSafe must be not nil. Panic otherwise.
func (dq *DequeSafe) Len() int {
	dq.m.Lock()
	ret := dq.q.Len()
	dq.m.Unlock()
	return ret
}

// Back returns the last element of thread safe double ended queue
// or nil if the DequeSafe is empty.
// The second, bool result indicates whether a valid value was returned;
// if the deque is empty, false will be returned.
// The complexity is O(1).
// DequeSafe must be not nil. Panic otherwise.
func (dq *DequeSafe) Back() (interface{}, bool) {
	dq.m.Lock()
	elem, found := dq.q.Back()
	dq.m.Unlock()
	return elem, found
}

// Front returns the first element of thread safe double ended queue
// or nil if the DequeSafe is empty.
// The second, bool result indicates whether a valid value was returned;
// if the deque is empty, false will be returned.
// The complexity is O(1).
// DequeSafe must be not nil. Panic otherwise.
func (dq *DequeSafe) Front() (interface{}, bool) {
	dq.m.Lock()
	elem, found := dq.q.Front()
	dq.m.Unlock()
	return elem, found
}

// PopBack retrieves and removes the current element from the back
// of the thread safe double ended queue.
// The second, bool result indicates whether a valid value was returned;
// if the DequeSafe is empty, false will be returned.
// The complexity is O(1).
// DequeSafe must be not nil. Panic otherwise.
func (dq *DequeSafe) PopBack() (interface{}, bool) {
	dq.m.Lock()
	elem, found := dq.q.PopBack()
	dq.m.Unlock()
	return elem, found
}

// PopFront retrieves and removes the current element from the front
// of the thread safe double ended queue.
// The second, bool result indicates whether a valid value was returned;
// if the DequeSafe is empty, false will be returned.
// The complexity is O(1).
// DequeSafe must be not nil. Panic otherwise.
func (dq *DequeSafe) PopFront() (interface{}, bool) {
	dq.m.Lock()
	elem, found := dq.q.PopFront()
	dq.m.Unlock()
	return elem, found
}

// PushBack adds value v to the the back of the thread safe double ended queue.
// The complexity is O(1).
// DequeSafe must be not nil. Panic otherwise.
func (dq *DequeSafe) PushBack(v interface{}) {
	dq.m.Lock()
	dq.q.PushBack(v)
	dq.m.Unlock()
}

// PushFront adds value v to the the front of the thread safe double ended queue.
// The complexity is O(1).
// DequeSafe must be not nil. Panic otherwise.
func (dq *DequeSafe) PushFront(v interface{}) {
	dq.m.Lock()
	dq.q.PushFront(v)
	dq.m.Unlock()
}

// NewDeque returns a new initialized thread UNSAFE double ended queue.
func NewDeque() *Deque {
	return deque.New()
}

// NewDequeSafe returns a new initialized thread safe double ended queue.
func NewDequeSafe() *DequeSafe {
	return new(DequeSafe)
}
