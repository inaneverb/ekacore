package dangerous

import "testing"
import "github.com/stretchr/testify/assert"

func TestTypedInterface(t *testing.T) {

	type T1 struct{}
	type T2 struct{}

	t1 := T1{}
	t2 := T2{}

	assert.Equal(t, TypedInterface(&t1).Type, TypedInterface(new(T1)).Type)
	assert.Equal(t, TypedInterface(&t2).Type, TypedInterface(new(T2)).Type)

	assert.NotEqual(t, TypedInterface(&t1).Type, TypedInterface(new(T2)).Type)
	assert.NotEqual(t, TypedInterface(&t2).Type, TypedInterface(new(T1)).Type)

	assert.NotEqual(t, TypedInterface(t1).Type, TypedInterface(t2).Type)
}
