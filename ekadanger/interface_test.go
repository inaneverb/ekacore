// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekadanger_test

import (
	"testing"

	"github.com/qioalice/ekago/v2/ekadanger"

	"github.com/stretchr/testify/assert"
)

func TestTypedInterface(t *testing.T) {

	type T1 struct{}
	type T2 struct{}

	t1 := T1{}
	t2 := T2{}

	assert.Equal(t, ekadanger.TypedInterface(&t1).Type, ekadanger.TypedInterface(new(T1)).Type)
	assert.Equal(t, ekadanger.TypedInterface(&t2).Type, ekadanger.TypedInterface(new(T2)).Type)

	assert.NotEqual(t, ekadanger.TypedInterface(&t1).Type, ekadanger.TypedInterface(new(T2)).Type)
	assert.NotEqual(t, ekadanger.TypedInterface(&t2).Type, ekadanger.TypedInterface(new(T1)).Type)

	assert.NotEqual(t, ekadanger.TypedInterface(t1).Type, ekadanger.TypedInterface(t2).Type)
}

type CustomError struct {}
func (_ *CustomError) Error() string { return "<custom error>" }

func TestTakeRealAddrForError(t *testing.T) {

	customNilError := (*CustomError)(nil)
	customNotNilError := new(CustomError)

	var legacyNilError error = customNilError
	var legacyNotNilError error = customNotNilError

	assert.True(t, ekadanger.TakeRealAddr(customNilError) == nil)
	assert.True(t, ekadanger.TakeRealAddr(legacyNilError) == nil)

	assert.True(t, ekadanger.TakeRealAddr(customNotNilError) != nil)
	assert.True(t, ekadanger.TakeRealAddr(legacyNotNilError) != nil)

	// This is why this test exists:
	assert.True(t, legacyNilError != nil)
}
