//go:build windows

// This file provides blocking repository file locks on Windows.
package engine

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const lockFileExclusiveLock = 0x00000002

var (
	kernel32DLL      = syscall.NewLazyDLL("kernel32.dll")
	lockFileExProc   = kernel32DLL.NewProc("LockFileEx")
	unlockFileExProc = kernel32DLL.NewProc("UnlockFileEx")
)

// acquireRepositoryLock exclusively locks the first byte of a repository lock file.
func acquireRepositoryLock(path string) (func() error, error) {
	if err := ensureParent(path); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	overlapped := new(syscall.Overlapped)
	result, _, callErr := lockFileExProc.Call(
		file.Fd(),
		lockFileExclusiveLock,
		0,
		1,
		0,
		uintptr(unsafe.Pointer(overlapped)),
	)
	if result == 0 {
		_ = file.Close()
		return nil, windowsLockError("LockFileEx", callErr)
	}
	return func() error {
		result, _, callErr := unlockFileExProc.Call(
			file.Fd(),
			0,
			1,
			0,
			uintptr(unsafe.Pointer(overlapped)),
		)
		var unlockErr error
		if result == 0 {
			unlockErr = windowsLockError("UnlockFileEx", callErr)
		}
		return errors.Join(unlockErr, file.Close())
	}, nil
}

// windowsLockError normalizes a failed kernel32 lock call.
func windowsLockError(operation string, callErr error) error {
	if errno, ok := callErr.(syscall.Errno); ok && errno == 0 {
		return fmt.Errorf("%s failed", operation)
	}
	return fmt.Errorf("%s: %w", operation, callErr)
}
