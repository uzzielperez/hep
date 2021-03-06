// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrdfs

import (
	"context"
	"io"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
)

// File implements access to a content and meta information of file over XRootD.
type File interface {
	io.ReaderAt
	io.WriterAt

	// Compression returns the compression info.
	Compression() *FileCompression

	// Info returns the cached stat info.
	// Note that it may return nil if info was not yet fetched and info may be not up-to-date.
	Info() *EntryStat

	// Handle returns the file handle.
	Handle() FileHandle

	// Close closes the file.
	Close(ctx context.Context) error

	// CloseVerify closes the file and checks whether the file has the provided size.
	// A zero size suppresses the verification.
	CloseVerify(ctx context.Context, size int64) error

	// Sync commits all pending writes to an open file.
	Sync(ctx context.Context) error

	// ReadAtContext reads len(p) bytes into p starting at offset off.
	ReadAtContext(ctx context.Context, p []byte, off int64) (n int, err error)

	// WriteAtContext writes len(p) bytes from p to the file at offset off.
	WriteAtContext(ctx context.Context, p []byte, off int64) error

	// Truncate changes the size of the file.
	Truncate(ctx context.Context, size int64) error

	// Stat fetches the stat info of this file from the XRootD server.
	// Note that Stat re-fetches value returned by the Info, so after the call to Stat
	// calls to Info may return different value than before.
	Stat(ctx context.Context) (EntryStat, error)

	// StatVirtualFS fetches the virtual stat info of this file from the XRootD server.
	StatVirtualFS(ctx context.Context) (VirtualFSStat, error)

	// VerifyWriteAt writes len(p) bytes from p to the file at offset off using crc32 verification.
	//
	// TODO: note that verifyw is not supported by the XRootD server.
	// See https://github.com/xrootd/xrootd/issues/738 for the details.
	VerifyWriteAt(ctx context.Context, p []byte, off int64) error
}

// FileHandle is the file handle, which should be treated as opaque data.
type FileHandle [4]byte

// FileCompression holds the compression parameters such as the page size and the type of compression.
type FileCompression struct {
	PageSize int32
	Type     [4]byte
}

// MarshalXrd implements xrdproto.Marshaler
func (o FileCompression) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteI32(o.PageSize)
	wBuffer.WriteBytes(o.Type[:])
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler
func (o *FileCompression) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	o.PageSize = rBuffer.ReadI32()
	rBuffer.ReadBytes(o.Type[:])
	return nil
}
