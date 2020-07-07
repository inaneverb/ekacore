// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekasys

import (
	"path/filepath"
	"runtime"
	"strconv"
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

		fullPackage, fn := filepath.Split(f.Function)
		_, file := filepath.Split(f.File)

		// we need last package from the fullPackage
		lastPackage := filepath.Base(fullPackage)

		// need remove last package from fullPackage
		if len(lastPackage)+2 <= len(fullPackage) && lastPackage != "." {
			fullPackage = fullPackage[:len(fullPackage)-len(lastPackage)-2]
		}

		f.Format += lastPackage + "/" + fn

		f.FormatFileOffset = len(f.Format) + 1
		f.Format += " (" + file + ":" + strconv.Itoa(f.Line) + ")"

		f.FormatFullPathOffset = len(f.Format) + 1
		f.Format += " " + fullPackage
	}

	return f.Format
}
