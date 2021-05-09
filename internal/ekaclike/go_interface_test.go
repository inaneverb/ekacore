// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaclike_test

import (
	"testing"

	"github.com/qioalice/ekago/v3/internal/ekaclike"

	"github.com/stretchr/testify/assert"
)

func TestUnpackInterface(t *testing.T) {

	type T1 struct{}
	type T2 struct{}

	t1 := T1{}
	t2 := T2{}

	assert.Equal(t, ekaclike.UnpackInterface(&t1).Type, ekaclike.UnpackInterface(new(T1)).Type)
	assert.Equal(t, ekaclike.UnpackInterface(&t2).Type, ekaclike.UnpackInterface(new(T2)).Type)

	assert.NotEqual(t, ekaclike.UnpackInterface(&t1).Type, ekaclike.UnpackInterface(new(T2)).Type)
	assert.NotEqual(t, ekaclike.UnpackInterface(&t2).Type, ekaclike.UnpackInterface(new(T1)).Type)

	assert.NotEqual(t, ekaclike.UnpackInterface(t1).Type, ekaclike.UnpackInterface(t2).Type)
}
