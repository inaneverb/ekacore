// Copyright © 2022. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaext

func If[T any](cond bool, vThen, vElse T) T {
	if cond {
		return vThen
	}
	return vElse
}

func ZeroIf[T comparable](v T, cond bool) T {
	return If(cond, *(new(T)), v)
}
