// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekatyp_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/inaneverb/ekacore/ekatyp/v4"
)

func TestUuidV8(t *testing.T) {

	const X uint32 = 0x1234ABCD

	var u ekatyp.Uuid
	var err = ekatyp.NewUuidWithCustomTo(&u, X, 32)

	require.NoError(t, err)
	require.EqualValues(t, X, ekatyp.UuidExtractCustomData(&u, 32))

	fmt.Printf("%v, %s\n", &u, ekatyp.UuidExtractTimestamp(&u))
}
