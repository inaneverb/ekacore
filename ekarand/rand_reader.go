// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekarand

import (
	crand "crypto/rand"
	mrand "math/rand"
)

type (
	// MathRandReader is an io.Reader interface that provides thread safe Read()
	// method using math/rand random number generator.
	// According with math/rand rules, if you don't instantiate a Rng generator
	// but use package level's function they will be thread safety.
	// You can use both of (*MathRandReader)(nil) or MathRandReader{} variants.
	MathRandReader struct{}

	// CryptoRandReader is an io.Reader interface that provides thread safe Read()
	// method with increased rng complexity comparing with math/rand
	// sacrificing performance.
	// You can use both of (*CryptoRandReader)(nil) or CryptoRandReader{} variants.
	CryptoRandReader struct{}
)

func (_ *MathRandReader) Read(p []byte) (n int, err error) {
	return mrand.Read(p)
}

func (_ *CryptoRandReader) Read(p []byte) (n int, err error) {
	return crand.Read(p)
}
