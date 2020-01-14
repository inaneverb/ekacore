package haikunator

import (
	"fmt"
	"testing"
)

func TestHaikunate(t *testing.T) {
	fmt.Println(HaikunateWithRange(100, 200))
	fmt.Println(HaikunateWithRange(200, 500))
	fmt.Println(HaikunateWithRange(30, 10))
	fmt.Println(Haikunate())
	fmt.Println(Haikunate())
	fmt.Println(Haikunate())
	fmt.Println(Haikunate())
}
