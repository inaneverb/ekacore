// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekarand_test

import (
	"fmt"
	"testing"

	"github.com/qioalice/ekago/ekarand/v4"
)

func TestHaikunate(t *testing.T) {
	fmt.Println(ekarand.HaikunateWithRange(100, 200))
	fmt.Println(ekarand.HaikunateWithRange(200, 500))
	fmt.Println(ekarand.HaikunateWithRange(30, 10))
	fmt.Println(ekarand.Haikunate())
	fmt.Println(ekarand.Haikunate())
	fmt.Println(ekarand.Haikunate())
	fmt.Println(ekarand.Haikunate())
}

func BenchmarkHaikunate(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ekarand.Haikunate()
	}
}
