//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd

// This file provides blocking repository file locks on Unix systems.
package engine

import (
	"context"
	"errors"
	"os"
	"syscall"
	"time"
)

const repositoryLockRetryWait = 10 * time.Millisecond

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

// acquireRepositoryLockContext preserves flock exclusion while polling the
// non-blocking form so callers can cancel a contended wait.
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
	for {
		err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			return func() error {
				return errors.Join(syscall.Flock(int(file.Fd()), syscall.LOCK_UN), file.Close())
			}, nil
		}
		if !errors.Is(err, syscall.EWOULDBLOCK) && !errors.Is(err, syscall.EAGAIN) {
			return nil, errors.Join(err, file.Close())
		}
		timer := time.NewTimer(repositoryLockRetryWait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, errors.Join(ctx.Err(), file.Close())
		case <-timer.C:
		}
	}
}
