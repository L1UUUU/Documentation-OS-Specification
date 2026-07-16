//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd

// This file provides blocking repository file locks on Unix systems.
package engine

import (
	"errors"
	"os"
	"syscall"
)

// acquireRepositoryLock exclusively locks a repository lock file.
func acquireRepositoryLock(path string) (func() error, error) {
	if err := ensureParent(path); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		_ = file.Close()
		return nil, err
	}
	return func() error {
		return errors.Join(syscall.Flock(int(file.Fd()), syscall.LOCK_UN), file.Close())
	}, nil
}
