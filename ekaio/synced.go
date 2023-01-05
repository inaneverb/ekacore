// Copyright © 2021-2023. All rights reserved.
// Author: Ilya Stroy.
// Contacts: iyuryevich@pm.me, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaio

import (
	"io"
	"sync"
)

// rSynced is a wrapper for io.Reader with sync.Mutex protector.
type rSynced struct {
	origin io.Reader
	mu     sync.Mutex
}

// NewSyncedReader creates and returns a new io.Reader,
// read operations of which are protected with sync.Mutex.
func NewSyncedReader(origin io.Reader) io.Reader {
	if origin == nil {
		return NewNopeReadWriteCloser()
	}
	return &rSynced{origin: origin}
}

func (r *rSynced) Read(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.origin.Read(p)
}

////////////////////////////////////////////////////////////////////////////////

// wSynced is a wrapper for io.Writer with sync.Mutex protector.
type wSynced struct {
	origin io.Writer
	mu     sync.Mutex
}

// NewSyncedWriter creates and returns a new io.Writer,
// read operations of which are protected with sync.Mutex.
func NewSyncedWriter(origin io.Writer) io.Writer {
	if origin == nil {
		return NewNopeReadWriteCloser()
	}
	return &wSynced{origin: origin}
}

func (w *wSynced) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.origin.Write(p)
}

////////////////////////////////////////////////////////////////////////////////

// rwSynced is a wrapper for io.ReadWriter with sync.Mutex protector.
type rwSynced struct {
	origin io.ReadWriter
	mu     sync.Mutex
}

// NewSyncedReadWriter creates and returns a new io.ReadWriter,
// read operations of which are protected with sync.Mutex.
func NewSyncedReadWriter(origin io.ReadWriter) io.ReadWriter {
	if origin == nil {
		return NewNopeReadWriteCloser()
	}
	return &rwSynced{origin: origin}
}

func (r *rwSynced) Read(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.origin.Read(p)
}

func (r *rwSynced) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.origin.Write(p)
}

////////////////////////////////////////////////////////////////////////////////

// rwcSynced is a wrapper for io.ReadWriteCloser with sync.Mutex protector.
type rwcSynced struct {
	origin io.ReadWriteCloser
	mu     sync.Mutex
}

// NewSyncedReadWriteCloser creates and returns a new io.ReadWriteCloser,
// read operations of which are protected with sync.Mutex.
func NewSyncedReadWriteCloser(origin io.ReadWriteCloser) io.ReadWriteCloser {
	if origin == nil {
		return NewNopeReadWriteCloser()
	}
	return &rwcSynced{origin: origin}
}

func (r *rwcSynced) Read(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.origin.Read(p)
}

func (r *rwcSynced) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.origin.Write(p)
}

func (r *rwcSynced) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.origin.Close()
}
