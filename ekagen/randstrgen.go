package ekagen

import (
	"math/rand"
	"time"
)

const (
	charSetLetters = `abcdefghijklmnopqrstuvwxyz`
	charSetDigits  = `1234567890`
	charSetAll     = charSetLetters + charSetDigits
)

//
var r *rand.Rand

//
func init() {
	r = rand.New(rand.New(rand.NewSource(99)))
	r.Seed(time.Now().UTC().Unix())
}

//
func genWithLenFrom(charSet string, n int) string {

	if n <= 0 {
		return ""
	}

	res := make([]byte, n)

	for i := 0; i < n; i++ {
		res[i] = charSet[r.Intn(len(charSet))]
	}

	return string(res)
}

//
func WithLen(n int) string {
	return genWithLenFrom(charSetAll, n)
}

//
func WithLenOnlyLetters(n int) string {
	return genWithLenFrom(charSetLetters, n)
}

//
func WithLenOnlyNumbers(n int) string {
	return genWithLenFrom(charSetDigits, n)
}
