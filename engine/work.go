// This file implements deterministic Work workspace generation.
package engine

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

// GenerateWork creates a globally unique active Work and refreshes INDEX.md transactionally.
func (e *Engine) GenerateWork(slug string) (WorkResult, error) {
	return e.beginWork(BeginInput{Slug: slug, Title: titleFromSlug(slug)}, false)
}

// BeginWork atomically creates a caller-defined active Work and refreshes INDEX.md.
func (e *Engine) BeginWork(input BeginInput) (WorkResult, error) {
	return e.beginWork(input, true)
}

// beginWork shares transactional creation while preserving GenerateWork's legacy duplicate behavior.
func (e *Engine) beginWork(input BeginInput, allowIdempotentRetry bool) (result WorkResult, returnErr error) {
	if err := validateBeginInput(input); err != nil {
		return WorkResult{}, err
	}
	if err := e.requireRuntimeRoots(); err != nil {
		return WorkResult{}, err
	}
	activePath := e.path(filepath.ToSlash(filepath.Join(e.Profile.ActiveRoot, input.Slug)))
	completedPath := e.path(filepath.ToSlash(filepath.Join(e.Profile.CompletedRoot, input.Slug)))
	lockPath := e.path(filepath.ToSlash(filepath.Join(e.Profile.LockRoot, "begin-work.lock")))
	releaseLock, err := acquireRepositoryLock(lockPath)
	if err != nil {
		return WorkResult{}, fmt.Errorf("acquire Begin Work lock: %w", err)
	}
	defer func() {
		if err := releaseLock(); err != nil {
			lockErr := fmt.Errorf("release Begin Work lock: %w", err)
			if returnErr == nil {
				returnErr = lockErr
				return
			}
			returnErr = errors.Join(returnErr, lockErr)
		}
	}()
	if existing, found, err := e.existingBeginResult(input, allowIdempotentRetry, activePath, completedPath); found || err != nil {
		return existing, err
	}

	stagingRoot := e.path(filepath.ToSlash(filepath.Join(e.Profile.RuntimeRoot, ".begin-staging")))
	if err := os.MkdirAll(stagingRoot, 0o755); err != nil {
		return WorkResult{}, fmt.Errorf("create Begin Work staging root: %w", err)
	}
	stagingPath, err := os.MkdirTemp(stagingRoot, input.Slug+"-")
	if err != nil {
		return WorkResult{}, fmt.Errorf("create Begin Work staging directory: %w", err)
	}
	defer func() {
		if stagingPath == "" {
			return
		}
		if err := os.RemoveAll(stagingPath); err != nil {
			cleanupErr := fmt.Errorf("remove Begin Work staging directory: %w", err)
			if returnErr == nil {
				returnErr = cleanupErr
				return
			}
			returnErr = errors.Join(returnErr, cleanupErr)
		}
	}()
	rollbackStaging := func(stage string, operationErr error) (WorkResult, error) {
		cleanupErr := os.RemoveAll(stagingPath)
		stagingPath = ""
		if cleanupErr != nil {
			return WorkResult{}, fmt.Errorf("%s: %w", stage, errors.Join(operationErr, fmt.Errorf("rollback staging: %w", cleanupErr)))
		}
		return WorkResult{}, fmt.Errorf("%s; staging rolled back: %w", stage, operationErr)
	}
	if err := os.MkdirAll(filepath.Join(stagingPath, "issues"), 0o755); err != nil {
		return rollbackStaging("create Work workspace", err)
	}
	prdFrontMatter := "---\nrelationships: []\n"
	if allowIdempotentRetry {
		prdFrontMatter += "begin_input_hash: " + beginInputFingerprint(input) + "\n"
	}
	prd := []byte(prdFrontMatter + "---\n\n# " + input.Title + "\n\n## Objective\n\n<!-- Define the Work objective and constraints. -->\n")
	if err := e.writeFileAtomic(filepath.Join(stagingPath, "PRD.md"), prd, 0o644); err != nil {
		return rollbackStaging("create PRD.md", err)
	}
	if err := e.writeFileAtomic(filepath.Join(stagingPath, "HANDOFF.md"), nil, 0o644); err != nil {
		return rollbackStaging("create HANDOFF.md", err)
	}
	for index, issue := range input.Issues {
		name := fmt.Sprintf("%02d-%s.md", index+1, issue.Slug)
		content := []byte("---\nstatus: " + issue.Status + "\ntitle: " + strconv.Quote(issue.Title) + "\n---\n")
		if err := e.writeFileAtomic(filepath.Join(stagingPath, "issues", name), content, 0o644); err != nil {
			return rollbackStaging("create initial Issue "+name, err)
		}
	}

	indexPath := e.path(e.Profile.IndexPath)
	previousIndex, indexExisted, err := readOptional(indexPath)
	if err != nil {
		return WorkResult{}, fmt.Errorf("read existing INDEX: %w", err)
	}
	if err := e.renameFile(stagingPath, activePath); err != nil {
		if existing, found, existingErr := e.existingBeginResult(input, allowIdempotentRetry, activePath, completedPath); found || existingErr != nil {
			return existing, existingErr
		}
		return rollbackStaging("publish staged Work", err)
	}
	stagingPath = ""
	rollbackPublished := func(stage string, operationErr error) (WorkResult, error) {
		rollbackErr := errors.Join(
			os.RemoveAll(activePath),
			restoreFile(indexPath, previousIndex, indexExisted),
		)
		if rollbackErr != nil {
			return WorkResult{}, fmt.Errorf("%s: %w", stage, errors.Join(operationErr, fmt.Errorf("rollback workspace and INDEX: %w", rollbackErr)))
		}
		return WorkResult{}, fmt.Errorf("%s; workspace and INDEX rolled back: %w", stage, operationErr)
	}
	if _, err := e.GenerateIndex(); err != nil {
		return rollbackPublished("generate INDEX after Work creation", err)
	}
	return WorkResult{Slug: input.Slug, Path: activePath, IndexPath: e.relativePath(indexPath)}, nil
}

// existingBeginResult resolves a final Work observed while holding the Begin Work lock.
func (e *Engine) existingBeginResult(input BeginInput, allowIdempotentRetry bool, activePath, completedPath string) (WorkResult, bool, error) {
	if _, err := os.Stat(activePath); err == nil {
		if !allowIdempotentRetry {
			return WorkResult{}, true, fmt.Errorf("Work slug %q already exists under active", input.Slug)
		}
		persistedFingerprint, fingerprintErr := beginInputFingerprintFromPRD(filepath.Join(activePath, "PRD.md"))
		if fingerprintErr != nil {
			return WorkResult{}, true, fingerprintErr
		}
		if persistedFingerprint != beginInputFingerprint(input) {
			return WorkResult{}, true, fmt.Errorf("%w: active Work %q was begun with different input", ErrConflict, input.Slug)
		}
		result := WorkResult{Slug: input.Slug, Path: activePath, IndexPath: e.Profile.IndexPath}
		if _, err := e.GenerateIndex(); err != nil {
			return result, true, fmt.Errorf("reconcile INDEX for active Work %q: %w", input.Slug, err)
		}
		return result, true, nil
	} else if !os.IsNotExist(err) {
		return WorkResult{}, true, fmt.Errorf("check active Work %q: %w", input.Slug, err)
	}
	if _, err := os.Stat(completedPath); err == nil {
		if allowIdempotentRetry {
			return WorkResult{}, true, fmt.Errorf("%w: Work %q is already completed", ErrConflict, input.Slug)
		}
		return WorkResult{}, true, fmt.Errorf("Work slug %q already exists under completed", input.Slug)
	} else if !os.IsNotExist(err) {
		return WorkResult{}, true, fmt.Errorf("check completed Work %q: %w", input.Slug, err)
	}
	return WorkResult{}, false, nil
}

// beginInputFingerprint records only the public inputs needed to distinguish retries.
func beginInputFingerprint(input BeginInput) string {
	normalized := BeginInput{Slug: input.Slug, Title: input.Title, Issues: append([]BeginIssue{}, input.Issues...)}
	data, _ := json.Marshal(normalized)
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

// beginInputFingerprintFromPRD reads the immutable retry identity persisted at creation.
func beginInputFingerprintFromPRD(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("%w: read active PRD.md: %v", ErrInvalidRepository, err)
	}
	front, err := parseFrontMatter(data)
	if err != nil || !front.Present {
		return "", fmt.Errorf("%w: active PRD.md must contain valid front matter", ErrInvalidRepository)
	}
	fingerprint, ok := front.Values["begin_input_hash"].(string)
	if !ok || fingerprint == "" {
		return "", fmt.Errorf("%w: active PRD.md does not record Begin input", ErrInvalidRepository)
	}
	return fingerprint, nil
}

// validateBeginInput rejects values that cannot be rendered as stable Work assets.
func validateBeginInput(input BeginInput) error {
	if err := validateSlug(input.Slug); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	if strings.TrimSpace(input.Title) == "" || strings.ContainsAny(input.Title, "\r\n") {
		return fmt.Errorf("%w: Work title must be a non-empty single line", ErrInvalidInput)
	}
	if len(input.Issues) > 99 {
		return fmt.Errorf("%w: Begin Work supports at most 99 initial Issues", ErrInvalidInput)
	}
	seen := map[string]bool{}
	for _, issue := range input.Issues {
		if err := validateSlug(issue.Slug); err != nil {
			return fmt.Errorf("%w: Issue %v", ErrInvalidInput, err)
		}
		if seen[issue.Slug] {
			return fmt.Errorf("%w: duplicate initial Issue slug %q", ErrInvalidInput, issue.Slug)
		}
		seen[issue.Slug] = true
		if strings.TrimSpace(issue.Title) == "" || strings.ContainsAny(issue.Title, "\r\n") {
			return fmt.Errorf("%w: Issue %q title must be a non-empty single line", ErrInvalidInput, issue.Slug)
		}
		if !issueStatuses[issue.Status] {
			return fmt.Errorf("%w: Issue %q has invalid status %q", ErrInvalidInput, issue.Slug, issue.Status)
		}
	}
	return nil
}

// requireRuntimeRoots verifies that Generate Work operates on an initialized profile.
func (e *Engine) requireRuntimeRoots() error {
	for _, relative := range []string{e.Profile.RuntimeRoot, e.Profile.ActiveRoot, e.Profile.CompletedRoot} {
		info, err := os.Stat(e.path(relative))
		if err != nil {
			return fmt.Errorf("repository is not initialized: %s: %w", relative, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("repository profile path %s is not a directory", relative)
		}
	}
	return nil
}

// titleFromSlug derives a readable deterministic Work title.
func titleFromSlug(slug string) string {
	words := strings.Split(slug, "-")
	for i, word := range words {
		runes := []rune(word)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}
