// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package process

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procPssCaptureSnapshot = modkernel32.NewProc("PssCaptureSnapshot")
	procPssQuerySnapshot   = modkernel32.NewProc("PssQuerySnapshot")
)

func PssCaptureSnapshot(processHandle syscall.Handle, captureFlags PSSCaptureFlags, threadContextFlags uint32, snapshotHandle *syscall.Handle) (err error) {
	_, _, e1 := syscall.Syscall6(procPssCaptureSnapshot.Addr(), 4, uintptr(processHandle), uintptr(captureFlags), uintptr(threadContextFlags), uintptr(unsafe.Pointer(snapshotHandle)), 0, 0)

	if e1 != windows.ERROR_SUCCESS {
		return e1
	}
	return nil
}

func PssQuerySnapshot(snapshotHandle syscall.Handle, informationClass uint32, buffer *PssThreadInformation, bufferLength uint32) (err error) {
	_, _, e1 := syscall.Syscall6(procPssQuerySnapshot.Addr(), 4, uintptr(snapshotHandle), uintptr(informationClass), uintptr(unsafe.Pointer(buffer)), uintptr(bufferLength), 0, 0)

	if e1 != windows.ERROR_SUCCESS {
		return e1
	}
	return nil
}
