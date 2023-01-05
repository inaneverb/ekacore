// Copyright © 2021-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

/*
Original package: https://github.com/ef-ds/stack
Design draft about queue + thread safe:
https://github.com/golang/proposal/blob/master/design/27935-unbounded-queue-package.md
*/

package ekatyp

import (
	"sync"

	"github.com/ef-ds/stack"
)

// Stack is LIFO data structure.
// Thanks to https://github.com/ef-ds/stack it's blazing fast.
// Stack must not be used by value, only by reference.
// Stack is thread UNSAFE. Use StackSafe if you need thread safety version.
type Stack = stack.Stack

// StackSafe is the same as Stack but provides thread-safety operations
// protecting them by sync.Mutex.
// StackSafe must not be used by value, only by reference.
type StackSafe struct {
	s Stack
	m sync.Mutex
}

// NewStack returns a new initialized thread UNSAFE stack.
func NewStack() *Stack {
	return stack.New()
}

// NewStackSafe returns a new initialized thread safe stack.
func NewStackSafe() *StackSafe {
	return new(StackSafe)
}

// Init initializes or clears thread safe stack.
func (s *StackSafe) Init() *StackSafe {
	if s == nil {
		return NewStackSafe()
	}
	s.m.Lock()
	s.s.Init()
	s.m.Unlock()
	return s
}

// Len returns the number of elements of thread safe stack.
// The complexity is O(1).
// StackSafe must be not nil. Panic otherwise.
func (s *StackSafe) Len() int {
	s.m.Lock()
	ret := s.s.Len()
	s.m.Unlock()
	return ret
}

// Back returns the last element of thread safe stack or nil if the StackSafe is empty.
// The second, bool result indicates whether a valid value was returned;
// if the stack is empty, false will be returned.
// The complexity is O(1).
// StackSafe must be not nil. Panic otherwise.
func (s *StackSafe) Back() (any, bool) {
	s.m.Lock()
	elem, found := s.s.Back()
	s.m.Unlock()
	return elem, found
}

// Pop retrieves and removes the current element from the back
// of the thread safe stack.
// The second, bool result indicates whether a valid value was returned;
// if the stack is empty, false will be returned.
// The complexity is O(1).
// StackSafe must be not nil. Panic otherwise.
func (s *StackSafe) Pop() (any, bool) {
	s.m.Lock()
	elem, found := s.s.Pop()
	s.m.Unlock()
	return elem, found
}

// Push adds value v to the the back of the thread safe stack.
// The complexity is O(1).
// StackSafe must be not nil. Panic otherwise.
func (s *StackSafe) Push(v any) {
	s.m.Lock()
	s.s.Push(v)
	s.m.Unlock()
}
