//go:build windows

package engine

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"
)

const repositoryLockHelperPathEnv = "DOS_ENGINE_LOCK_HELPER_PATH"

func TestCreateIssueReturnsWithinBoundWhenAnotherProcessHoldsLock(t *testing.T) {
	if lockPath := os.Getenv(repositoryLockHelperPathEnv); lockPath != "" {
		release, err := acquireRepositoryLock(lockPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "acquire helper lock: %v\n", err)
			os.Exit(2)
		}
		fmt.Println("locked")
		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		if err := release(); err != nil {
			fmt.Fprintf(os.Stderr, "release helper lock: %v\n", err)
			os.Exit(3)
		}
		return
	}

	instance := newTestEngine(t)
	beginIssueWork(t, instance, "contended-lock")
	releaseHelper := startRepositoryLockHelper(t, instance.path(".scratch/.locks/create-issue.lock"))

	result := make(chan error, 1)
	started := time.Now()
	go func() {
		_, err := instance.CreateIssue(validCreateIssueInput("contended-lock"))
		result <- err
	}()

	select {
	case err := <-result:
		if err == nil {
			t.Fatal("CreateIssue() error = nil while another process held the repository lock")
		}
		if elapsed := time.Since(started); elapsed < 4*time.Second || elapsed > 7*time.Second {
			t.Fatalf("CreateIssue() lock wait = %s, want bounded retry between 4s and 7s", elapsed)
		}
	case <-time.After(7 * time.Second):
		releaseHelper()
		select {
		case <-result:
		case <-time.After(2 * time.Second):
			t.Fatal("CreateIssue() remained blocked after the helper released the lock")
		}
		t.Fatal("CreateIssue() remained blocked for more than 7 seconds under cross-process lock contention")
	}
}

func TestCreateIssueContextCancellationReleasesResourcesAndAllowsRetry(t *testing.T) {
	instance := newTestEngine(t)
	beginIssueWork(t, instance, "cancelled-lock")
	releaseHelper := startRepositoryLockHelper(t, instance.path(".scratch/.locks/create-issue.lock"))

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	started := time.Now()
	_, err := instance.CreateIssueContext(ctx, validCreateIssueInput("cancelled-lock"))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("CreateIssueContext() error = %v, want context deadline exceeded", err)
	}
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("CreateIssueContext() cancellation took %s, want under 1s", elapsed)
	}

	releaseHelper()
	result, err := instance.CreateIssue(validCreateIssueInput("cancelled-lock"))
	if err != nil || result.Number != 1 || !result.Created {
		t.Fatalf("CreateIssue(retry) = %+v, %v, want newly allocated Issue 01", result, err)
	}
}

func TestCreateIssueWaitsForCrossProcessLockThenAllocatesOnce(t *testing.T) {
	instance := newTestEngine(t)
	beginIssueWork(t, instance, "released-lock")
	releaseHelper := startRepositoryLockHelper(t, instance.path(".scratch/.locks/create-issue.lock"))
	released := make(chan struct{})
	go func() {
		time.Sleep(150 * time.Millisecond)
		releaseHelper()
		close(released)
	}()

	started := time.Now()
	result, err := instance.CreateIssue(validCreateIssueInput("released-lock"))
	if err != nil || result.Number != 1 || !result.Created {
		t.Fatalf("CreateIssue() = %+v, %v, want newly allocated Issue 01", result, err)
	}
	if elapsed := time.Since(started); elapsed < 100*time.Millisecond || elapsed > time.Second {
		t.Fatalf("CreateIssue() lock wait = %s, want retry after helper release", elapsed)
	}
	<-released

	retried, err := instance.CreateIssue(validCreateIssueInput("released-lock"))
	if err != nil || retried.Number != 1 || retried.Created {
		t.Fatalf("CreateIssue(retry) = %+v, %v, want idempotent Issue 01", retried, err)
	}
}

func startRepositoryLockHelper(t *testing.T, lockPath string) func() {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestCreateIssueReturnsWithinBoundWhenAnotherProcessHoldsLock$")
	cmd.Env = append(os.Environ(), repositoryLockHelperPathEnv+"="+lockPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("helper stdout: %v", err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("helper stdin: %v", err)
	}
	var stderr strings.Builder
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start lock helper: %v", err)
	}
	if line, err := bufio.NewReader(stdout).ReadString('\n'); err != nil || strings.TrimSpace(line) != "locked" {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		t.Fatalf("helper readiness = %q, %v; stderr=%s", line, err, stderr.String())
	}

	var releaseOnce sync.Once
	releaseHelper := func() {
		releaseOnce.Do(func() {
			_, _ = stdin.Write([]byte("\n"))
			_ = stdin.Close()
			if err := cmd.Wait(); err != nil {
				t.Errorf("lock helper: %v; stderr=%s", err, stderr.String())
			}
		})
	}
	t.Cleanup(releaseHelper)
	return releaseHelper
}
