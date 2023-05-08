// Copyright © 2020-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaunsafe_test

import (
	"fmt"
	"testing"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

func TestGoRuntimeAddressEncoder(t *testing.T) {
	var s = ekaunsafe.GoRuntimeAddressEncoder().String()
	fmt.Printf("GO RUNTIME BYTE ORDER: %s\n", s)
}
