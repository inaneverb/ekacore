// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekasys

import (
	"path/filepath"
	"runtime"

	"github.com/qioalice/ekago/v3/ekastr"
)

// StackFrame represents one stack level (frame/item).
// It general purpose is runtime.Frame type extending.
type StackFrame struct {
	runtime.Frame

	// Format is formatted string representation of the current stack frame:
	// "<package>/<func> (<short_file>:<file_line>) <full_package_path>".
	//
	// Note, stack frame is not formatted by default.
	// It means that this field is empty in 99% cases. If you need to generate
	// formatted string use DoFormat method.
	Format string

	// FormatFileOffset is the index of
	// "(<short_file>:<file_line>)..." in Format field.
	FormatFileOffset int

	// FormatFullPathOffset is the index of "<full_package_path>" in Format field.
	FormatFullPathOffset int
}

// DoFormat generates formatted string representation of the current stack frame,
// saves it to Format field and returns it. Output looks like:
// "<package>/<func> (<short_file>:<file_line>) <full_package_path>".
//
// Does not regenerate formatted string if it's already generated (just returns).
func (f *StackFrame) DoFormat() string {

	if f.Format == "" {
		f.Format = f.doFormat()
	}

	return f.Format
}

// doFormat is a private part of DoFormat() function.
func (f *StackFrame) doFormat() string {

	fullPackage, fn := filepath.Split(f.Function)
	_, file := filepath.Split(f.File)

	// we need last package from the fullPackage
	lastPackage := filepath.Base(fullPackage)

	// need remove last package from fullPackage
	if len(lastPackage)+2 <= len(fullPackage) && lastPackage != "." {
		fullPackage = fullPackage[:len(fullPackage)-len(lastPackage)-2]
	}

	requiredBufLen := 9 // 2 spaces, '/', '(', ')', ':' + 3 reserved bytes
	requiredBufLen += len(fullPackage) + len(lastPackage) + len(file) + len(fn)
	requiredBufLen += ekastr.PItoa32(int32(f.Line))

	buf := make([]byte, requiredBufLen)
	offset := 0

	// maybe 'lastPackage' is version of package?
	if len(lastPackage) > 1 && (lastPackage[0] == 'v' || lastPackage[0] == 'V') {
		is := true
		for i, n := 1, len(lastPackage); i < n && is; i++ {
			if lastPackage[i] < '0' || lastPackage[i] > '9' {
				is = false
			}
		}
		if is {
			lastPackage2 := filepath.Base(fullPackage)
			if len(lastPackage2)+2 <= len(fullPackage) && lastPackage2 != "." {
				fullPackage = fullPackage[:len(fullPackage)-len(lastPackage2)-1]
			}
			copy(buf, lastPackage2)
			offset += len(lastPackage2)

			buf[offset] = '.'
			offset++
		}
	}

	copy(buf[offset:], lastPackage)
	offset += len(lastPackage)

	buf[offset] = '.'
	offset++

	copy(buf[offset:], fn)
	offset += len(fn)

	buf[offset] = ' '
	offset++
	buf[offset] = '('
	offset++

	f.FormatFileOffset = offset - 1

	copy(buf[offset:], file)
	offset += len(file)

	buf[offset] = ':'
	offset++

	offset += ekastr.BItoa32(buf[offset:], int32(f.Line))

	buf[offset] = ')'
	offset++
	buf[offset] = ' '
	offset++

	f.FormatFullPathOffset = offset

	copy(buf[offset:], fullPackage)
	offset += len(fullPackage)

	return string(buf[:offset])
}
