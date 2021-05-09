// Copyright © 2020-2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekamath

type (
	// Flag8 is a type that allows you to represent a 8 bit set named "flagset".
	// It's often used to optimized storing only true/false values of something.
	Flag8 uint8

	// Flag16 is a type that allows you to represent a 16 bit set named "flagset".
	// It's often used to optimized storing only true/false values of something.
	Flag16 uint16

	// Flag32 is a type that allows you to represent a 32 bit set named "flagset".
	// It's often used to optimized storing only true/false values of something.
	Flag32 uint32

	// Flag64 is a type that allows you to represent a 64 bit set named "flagset".
	// It's often used to optimized storing only true/false values of something.
	Flag64 uint64

	// Flags8 is absolutely the same as Flag8 (8 bit set) but w/o trailing "s".
	// It exists exists to make your code more clear. Using this type you can say:
	// "This variable has exactly many flags, not only one!".
	Flags8 = Flag8

	// Flags16 is absolutely the same as Flag16 (16 bit set) but w/o trailing "s".
	// It exists exists to make your code more clear. Using this type you can say:
	// "This variable has exactly many flags, not only one!".
	Flags16 = Flag16

	// Flags32 is absolutely the same as Flag32 (32 bit set) but w/o trailing "s".
	// It exists exists to make your code more clear. Using this type you can say:
	// "This variable has exactly many flags, not only one!".
	Flags32 = Flag32

	// Flags64 is absolutely the same as Flag64 (64 bit set) but w/o trailing "s".
	// It exists exists to make your code more clear. Using this type you can say:
	// "This variable has exactly many flags, not only one!".
	Flags64 = Flag64
)

// --------------------------------- EXTENDERS -------------------------------- //
// ---------------------------------------------------------------------------- //

// To16 extends and returns current f (8 bit flagset) to 16 bit flagset.
func (f Flag8) To16() Flag16 {
	return Flag16(f)
}

// To32 extends and returns current f (8 bit flagset) to 32 bit flagset.
func (f Flag8) To32() Flag32 {
	return Flag32(f)
}

// To64 extends and returns current f (8 bit flagset) to 64 bit flagset.
func (f Flag8) To64() Flag64 {
	return Flag64(f)
}

// To32 extends and returns current f (16 bit flagset) to 32 bit flagset.
func (f Flag16) To32() Flag32 {
	return Flag32(f)
}

// To64 extends and returns current f (16 bit flagset) to 64 bit flagset.
func (f Flag16) To64() Flag64 {
	return Flag64(f)
}

// To64 extends and returns current f (32 bit flagset) to 64 bit flagset.
func (f Flag32) To64() Flag64 {
	return Flag64(f)
}

// ---------------------------- 8 BIT FLAG METHODS ---------------------------- //
// ---------------------------------------------------------------------------- //

// IsZero reports whether f is the empty flagset. Returns true if f == 0.
func (f Flag8) IsZero() bool {
	return f == 0
}

// TestAny reports whether at least one flag from anotherFlags is set in f.
// Returns false if there is no flags in f from anotherFlags.
func (f Flag8) TestAny(anotherFlags Flag8) bool {
	return f&anotherFlags != 0
}

// TestAll reports whether ALL flags from anotherFlags are set in f.
// Returns false if at least one flag from anotherFlags is NOT set in f.
func (f Flag8) TestAll(anotherFlags Flag8) bool {
	return f&anotherFlags == anotherFlags
}

// SetAll sets all flags from anotherFlags in f. Returns changed f. Keeps flags
// already been set in f. If you want overwrite them use method ReplaceBy() instead.
func (f *Flag8) SetAll(anotherFlags Flag8) *Flag8 {
	*f |= anotherFlags
	return f
}

// ReplaceBy replaces all being set flags in f by newFlags. It's like assign op.
// Returns changed f.
func (f *Flag8) ReplaceBy(newFlags Flag8) *Flag8 {
	*f = newFlags
	return f
}

// Clear downs all flags from cleanedFlags in f. Returns changed f.
func (f *Flag8) Clear(cleanedFlags Flag8) *Flag8 {
	*f &^= cleanedFlags
	return f
}

// Zero downs all set flags in f. After call f.IsZero() == true. Returns clean f.
func (f *Flag8) Zero() *Flag8 {
	*f = 0
	return f
}

// --------------------------- 16 BIT FLAG METHODS ---------------------------- //
// ---------------------------------------------------------------------------- //

// IsZero reports whether f is the empty flagset. Returns true if f == 0.
func (f Flag16) IsZero() bool {
	return f == 0
}

// TestAny reports whether at least one flag from anotherFlags is set in f.
// Returns false if there is no flags in f from anotherFlags.
func (f Flag16) TestAny(anotherFlags Flag16) bool {
	return f&anotherFlags != 0
}

// TestAll reports whether ALL flags from anotherFlags are set in f.
// Returns false if at least one flag from anotherFlags is NOT set in f.
func (f Flag16) TestAll(anotherFlags Flag16) bool {
	return f&anotherFlags == anotherFlags
}

// SetAll sets all flags from anotherFlags in f. Returns changed f. Keeps flags
// already been set in f. If you want overwrite them use method ReplaceBy() instead.
func (f *Flag16) SetAll(anotherFlags Flag16) *Flag16 {
	*f |= anotherFlags
	return f
}

// ReplaceBy replaces all being set flags in f by newFlags. It's like assign op.
// Returns changed f.
func (f *Flag16) ReplaceBy(newFlags Flag16) *Flag16 {
	*f = newFlags
	return f
}

// Clear downs all flags from cleanedFlags in f. Returns changed f.
func (f *Flag16) Clear(cleanedFlags Flag16) *Flag16 {
	*f &^= cleanedFlags
	return f
}

// Zero downs all set flags in f. After call f.IsZero() == true. Returns clean f.
func (f *Flag16) Zero() *Flag16 {
	*f = 0
	return f
}

// --------------------------- 32 BIT FLAG METHODS ---------------------------- //
// ---------------------------------------------------------------------------- //

// IsZero reports whether f is the empty flagset. Returns true if f == 0.
func (f Flag32) IsZero() bool {
	return f == 0
}

// TestAny reports whether at least one flag from anotherFlags is set in f.
// Returns false if there is no flags in f from anotherFlags.
func (f Flag32) TestAny(anotherFlags Flag32) bool {
	return f&anotherFlags != 0
}

// TestAll reports whether ALL flags from anotherFlags are set in f.
// Returns false if at least one flag from anotherFlags is NOT set in f.
func (f Flag32) TestAll(anotherFlags Flag32) bool {
	return f&anotherFlags == anotherFlags
}

// SetAll sets all flags from anotherFlags in f. Returns changed f. Keeps flags
// already been set in f. If you want overwrite them use method ReplaceBy() instead.
func (f *Flag32) SetAll(anotherFlags Flag32) *Flag32 {
	*f |= anotherFlags
	return f
}

// ReplaceBy replaces all being set flags in f by newFlags. It's like assign op.
// Returns changed f.
func (f *Flag32) ReplaceBy(newFlags Flag32) *Flag32 {
	*f = newFlags
	return f
}

// Clear downs all flags from cleanedFlags in f. Returns changed f.
func (f *Flag32) Clear(cleanedFlags Flag32) *Flag32 {
	*f &^= cleanedFlags
	return f
}

// Zero downs all set flags in f. After call f.IsZero() == true. Returns clean f.
func (f *Flag32) Zero() *Flag32 {
	*f = 0
	return f
}

// --------------------------- 64 BIT FLAG METHODS ---------------------------- //
// ---------------------------------------------------------------------------- //

// IsZero reports whether f is the empty flagset. Returns true if f == 0.
func (f Flag64) IsZero() bool {
	return f == 0
}

// TestAny reports whether at least one flag from anotherFlags is set in f.
// Returns false if there is no flags in f from anotherFlags.
func (f Flag64) TestAny(anotherFlags Flag64) bool {
	return f&anotherFlags != 0
}

// TestAll reports whether ALL flags from anotherFlags are set in f.
// Returns false if at least one flag from anotherFlags is NOT set in f.
func (f Flag64) TestAll(anotherFlags Flag64) bool {
	return f&anotherFlags == anotherFlags
}

// SetAll sets all flags from anotherFlags in f. Returns changed f. Keeps flags
// already been set in f. If you want overwrite them use method ReplaceBy() instead.
func (f *Flag64) SetAll(anotherFlags Flag64) *Flag64 {
	*f |= anotherFlags
	return f
}

// ReplaceBy replaces all being set flags in f by newFlags. It's like assign op.
// Returns changed f.
func (f *Flag64) ReplaceBy(newFlags Flag64) *Flag64 {
	*f = newFlags
	return f
}

// Clear downs all flags from cleanedFlags in f. Returns changed f.
func (f *Flag64) Clear(cleanedFlags Flag64) *Flag64 {
	*f &^= cleanedFlags
	return f
}

// Zero downs all set flags in f. After call f.IsZero() == true. Returns clean f.
func (f *Flag64) Zero() *Flag64 {
	*f = 0
	return f
}
