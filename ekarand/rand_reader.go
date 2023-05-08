// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekarand

import (
	"io"

	crand "crypto/rand"
	mrand "math/rand"

	"github.com/awnumar/fastrand"
	"lukechampine.com/frand"
)

// MathRandReader is an io.Reader interface that provides thread safe Read()
// method using math/rand random number generator.
// According to math/rand rules, if you don't instantiate a Rng generator
// but use package level's function they will be thread safety.
// You can use both of (*MathRandReader)(nil) or MathRandReader{} variants.
type MathRandReader struct{}

// CryptoRandReader is an io.Reader interface that provides thread safe Read()
// method with increased rng complexity comparing with math/rand
// sacrificing performance.
// You can use both of (*CryptoRandReader)(nil) or CryptoRandReader{} variants.
type CryptoRandReader struct{}

type FastRandReader struct{}

type CryptoFastRandReader struct{}

func (_ *MathRandReader) Read(p []byte) (n int, err error) {
	return mrand.Read(p)
}

func (_ *CryptoRandReader) Read(p []byte) (n int, err error) {
	return crand.Read(p)
}

func (_ *FastRandReader) Read(p []byte) (n int, err error) {
	return frand.Read(p)
}

func (_ *CryptoFastRandReader) Read(p []byte) (n int, err error) {
	return fastrand.Reader.Read(p)
}

func GetMathRandReader() io.Reader       { return (*MathRandReader)(nil) }
func GetCryptoRandReader() io.Reader     { return (*CryptoRandReader)(nil) }
func GetFastRandReader() io.Reader       { return (*FastRandReader)(nil) }
func GetCryptoFastRandReader() io.Reader { return (*CryptoFastRandReader)(nil) }
