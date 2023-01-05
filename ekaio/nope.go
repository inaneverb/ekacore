// Copyright © 2021-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaio

import (
	"io"
)

// rwcNope implements io.Reader, io.Writer, io.Closer with no-op of them.
type rwcNope struct{}

func (r *rwcNope) Read(_ []byte) (int, error)  { return 0, nil }
func (r *rwcNope) Write(_ []byte) (int, error) { return 0, nil }
func (r *rwcNope) Close() error                { return nil }

func NewNopeReadWriteCloser() io.ReadWriteCloser {
	return (*rwcNope)(nil)
}
