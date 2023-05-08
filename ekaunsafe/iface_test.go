// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

func TestUnpackInterface(t *testing.T) {

	type T1 struct{}
	type T2 struct{}

	var t1 = T1{}
	var t2 = T2{}

	assert.Equal(t, ekaunsafe.UnpackInterface(&t1).Type, ekaunsafe.UnpackInterface(new(T1)).Type)
	assert.Equal(t, ekaunsafe.UnpackInterface(&t2).Type, ekaunsafe.UnpackInterface(new(T2)).Type)

	assert.NotEqual(t, ekaunsafe.UnpackInterface(&t1).Type, ekaunsafe.UnpackInterface(new(T2)).Type)
	assert.NotEqual(t, ekaunsafe.UnpackInterface(&t2).Type, ekaunsafe.UnpackInterface(new(T1)).Type)

	assert.NotEqual(t, ekaunsafe.UnpackInterface(t1).Type, ekaunsafe.UnpackInterface(t2).Type)
}
