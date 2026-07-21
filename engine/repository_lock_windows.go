//go:build windows

// This file provides blocking repository file locks on Windows.
package engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	lockFileFailImmediately = 0x00000001
	lockFileExclusiveLock   = 0x00000002
	repositoryLockRetryWait = 10 * time.Millisecond
)

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
	return windowsRepositoryLockRelease(file, overlapped), nil
}

// acquireRepositoryLockContext retries the non-blocking Windows lock call so
// context cancellation remains observable instead of being trapped inside a
// synchronous LockFileEx syscall.
func acquireRepositoryLockContext(ctx context.Context, path string) (func() error, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := ensureParent(path); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	overlapped := new(syscall.Overlapped)
	for {
		result, _, callErr := lockFileExProc.Call(
			file.Fd(),
			lockFileExclusiveLock|lockFileFailImmediately,
			0,
			1,
			0,
			uintptr(unsafe.Pointer(overlapped)),
		)
		if result != 0 {
			return windowsRepositoryLockRelease(file, overlapped), nil
		}
		if !windowsLockIsContended(callErr) {
			return nil, errors.Join(windowsLockError("LockFileEx", callErr), file.Close())
		}
		timer := time.NewTimer(repositoryLockRetryWait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, errors.Join(fmt.Errorf("wait for repository lock: %w", ctx.Err()), file.Close())
		case <-timer.C:
		}
	}
}

func windowsRepositoryLockRelease(file *os.File, overlapped *syscall.Overlapped) func() error {
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
	}
}

func windowsLockIsContended(callErr error) bool {
	errno, ok := callErr.(syscall.Errno)
	return ok && (errno == syscall.Errno(33) || errno == syscall.Errno(997))
}

// windowsLockError normalizes a failed kernel32 lock call.
func windowsLockError(operation string, callErr error) error {
	if errno, ok := callErr.(syscall.Errno); ok && errno == 0 {
		return fmt.Errorf("%s failed", operation)
	}
	return fmt.Errorf("%s: %w", operation, callErr)
}
