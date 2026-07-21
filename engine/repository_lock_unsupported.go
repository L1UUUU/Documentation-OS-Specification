//go:build !windows && !darwin && !dragonfly && !freebsd && !linux && !netbsd && !openbsd

// This file reports unsupported repository locking platforms explicitly.
package engine

import (
	"context"
	"fmt"
	"runtime"
)

// acquireRepositoryLock rejects identity allocation where file locking is unavailable.
func acquireRepositoryLock(_ string) (func() error, error) {
	return nil, fmt.Errorf("repository file locking is not supported on %s", runtime.GOOS)
}

func acquireRepositoryLockContext(ctx context.Context, path string) (func() error, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return acquireRepositoryLock(path)
}
