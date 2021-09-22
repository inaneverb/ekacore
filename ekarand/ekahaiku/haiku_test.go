// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekahaiku_test

import (
	"fmt"
	"testing"

	"github.com/qioalice/ekago/v3/ekarand/ekahaiku"
)

func TestHaikunate(t *testing.T) {
	fmt.Println(ekahaiku.HaikunateWithRange(100, 200))
	fmt.Println(ekahaiku.HaikunateWithRange(200, 500))
	fmt.Println(ekahaiku.HaikunateWithRange(30, 10))
	fmt.Println(ekahaiku.Haikunate())
	fmt.Println(ekahaiku.Haikunate())
	fmt.Println(ekahaiku.Haikunate())
	fmt.Println(ekahaiku.Haikunate())
}

func BenchmarkHaikunate(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ekahaiku.Haikunate()
	}
}
