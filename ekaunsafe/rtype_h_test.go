// Copyright © 2020-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

// --------------------------------------------------------------------------
//  IF ANY OF THOSE TESTS ARE FAILED, AND YOU'RE SURE THAT CODE IS CORRECT,
//  JUST CHECK THAT TEST DATA VARIABLES IN THE FILE rtype_test.go,
//  NAMED tda, td1, td2, pt CONTAINS CONSISTENT AND CORRECT DATA!
// --------------------------------------------------------------------------

func TestRTypeOfReflectType(t *testing.T) {
	// THESE TESTS ARE WRITTEN CORRECTLY AND THERE ARE NO MISTAKES OR BUGS!
	for i, n := 0, len(td1); i < n; i++ {
		var y = ekaunsafe.RTypeOfReflectType(reflect.TypeOf(tda[i]))
		assert.Equal(t, td1[i].f(), y, "%s", pt[i].name)
	}
}

func TestRTypeOfDeref(t *testing.T) {
	// THESE TESTS ARE WRITTEN CORRECTLY AND THERE ARE NO MISTAKES OR BUGS!
	for i, n := 0, len(td1); i < n; i++ {
		var y = reflect.New(reflect.TypeOf(tda[i])).Interface()
		assert.Equal(t, td1[i].f(), ekaunsafe.RTypeOfDeref(y), "%s", pt[i].name)
	}
}

func TestReflectTypeOfRType(t *testing.T) {
	// THESE TESTS ARE WRITTEN CORRECTLY AND THERE ARE NO MISTAKES OR BUGS!
	for i, n := 0, len(td1); i < n; i++ {
		var y = ekaunsafe.ReflectTypeOfRType(td1[i].f())
		assert.True(t, reflect.TypeOf(tda[i]) == y, "%s", pt[i].name)
	}
}

func BenchmarkReflectTypeOfRType(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ekaunsafe.ReflectTypeOfRType(ekaunsafe.RTypeMapStringAny())
	}
}
