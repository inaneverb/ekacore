// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekarand_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qioalice/ekago/v2/ekarand"
)

func doTest(t *testing.T, n int) {
	str := ekarand.WithLen(n)
	assert.Len(t, str, n)
	fmt.Println(str)
}

func TestWithLen(t *testing.T) {
	doTest(t, 10)
	doTest(t, 10)
	doTest(t, 10)
	doTest(t, 10)
	doTest(t, 10)
	doTest(t, 10)
}
