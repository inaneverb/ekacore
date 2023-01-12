// Copyright © 2019-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekadeath

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {

	var d1 DestructorSimple = nil
	var d2 DestructorWithExitCode = nil
	var d3 = DestructorSimple(func() {})
	var d4 = DestructorWithExitCode(func(_ int) {})

	require.Nil(t, parse(d1))
	require.Nil(t, parse(d2))
	require.NotNil(t, parse(d3))
	require.NotNil(t, parse(d4))
}
