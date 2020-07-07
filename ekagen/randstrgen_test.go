package ekagen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func doTest(t *testing.T, n int) {
	str := WithLen(n)
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
