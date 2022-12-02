// Copyright © 2020-2022. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekamath

import (
	"fmt"

	"github.com/qioalice/ekago/v3/ekaext"
)

type Stat[T ekaext.Numeric] interface {
	fmt.Stringer

	Min() T
	Max() T
	Avg() T
	N() int
	Count(v T)
	Clear()
}

////////////////////////////////////////////////////////////////////////////////

type _StatImpl[T ekaext.Numeric] struct {
	min, max, avg T
	n             int
	isCumulative  bool
}

func (s *_StatImpl[T]) Min() T { return s.min }
func (s *_StatImpl[T]) Max() T { return s.max }
func (s *_StatImpl[T]) N() int { return s.n }

func (s *_StatImpl[T]) Avg() T {

	if s.isCumulative {
		return s.avg / T(s.n)
	} else {
		return s.avg
	}
}

func (s *_StatImpl[T]) Count(v T) {

	if s.n == 0 {
		s.min, s.max, s.avg, s.n = v, v, v, 1
		return
	}

	s.min = Min(s.min, v)
	s.max = Max(s.max, v)

	if !s.isCumulative {
		s.avg = v - s.avg/T(s.n+1) + s.avg
	} else {
		s.avg += v
	}
}

func (s *_StatImpl[T]) Clear() {
	s.min, s.max, s.avg, s.n = 0, 0, 0, 0
}

func (s *_StatImpl[T]) String() string {
	return fmt.Sprintf("<[%d] Min: %v, Max: %v, Avg: %v>",
		s.N(), s.Min(), s.Max(), s.Avg())
}

////////////////////////////////////////////////////////////////////////////////

func NewStatCumulative[T ekaext.Numeric]() Stat[T] {
	return &_StatImpl[T]{isCumulative: true}
}

func NewStatIterative[T ekaext.Numeric]() Stat[T] {
	return &_StatImpl[T]{isCumulative: false}
}
