// This file provides repository construction, path safety, and atomic file helpers.
package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var workstreamSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// New creates an engine rooted at an existing directory.
func New(root string) (*Engine, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve repository root: %w", err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, fmt.Errorf("%w: stat repository root %q: %w", ErrInvalidRepository, abs, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%w: repository root %q is not a directory", ErrInvalidRepository, abs)
	}
	return &Engine{
		Root:                  filepath.Clean(abs),
		Profile:               DefaultProfile(),
		renameFile:            os.Rename,
		writeFileAtomic:       writeAtomic,
		afterOutcomePersisted: func() error { return nil },
	}, nil
}

// path resolves a profile-relative path beneath the repository root.
func (e *Engine) path(relative string) string {
	return filepath.Join(e.Root, filepath.FromSlash(relative))
}

// relativePath returns a stable slash-separated path for user-facing output.
func (e *Engine) relativePath(absolute string) string {
	relative, err := filepath.Rel(e.Root, absolute)
	if err != nil {
		return filepath.ToSlash(absolute)
	}
	return filepath.ToSlash(relative)
}

// pathInside reports whether a path is contained by the repository root.
func (e *Engine) pathInside(path string) bool {
	relative, err := filepath.Rel(e.Root, path)
	if err != nil {
		return false
	}
	return relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) && !filepath.IsAbs(relative)
}

// validateSlug checks the profile's globally addressable Work slug syntax.
func validateSlug(slug string) error {
	if !workstreamSlugPattern.MatchString(slug) {
		return fmt.Errorf("workstream slug %q must be lowercase kebab-case", slug)
	}
	return nil
}

// ensureParent creates a file's parent directory.
func ensureParent(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create parent directory for %s: %w", path, err)
	}
	return nil
}

// writeAtomic writes a file through a same-directory temporary file.
func writeAtomic(path string, data []byte, mode os.FileMode) error {
	if err := ensureParent(path); err != nil {
		return err
	}
	if existing, err := os.ReadFile(path); err == nil && string(existing) == string(data) {
		return nil
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".dos-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary file for %s: %w", path, err)
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	if err := temporary.Chmod(mode); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("set mode for temporary file %s: %w", path, err)
	}
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("write temporary file %s: %w", path, err)
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("sync temporary file %s: %w", path, err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary file %s: %w", path, err)
	}
	if err := replaceFile(temporaryName, path); err != nil {
		return fmt.Errorf("replace %s: %w", path, err)
	}
	return nil
}

// replaceFile replaces a destination on platforms where rename-over-existing is unavailable.
func replaceFile(source, destination string) error {
	if err := os.Rename(source, destination); err == nil {
		return nil
	}
	if err := os.Remove(destination); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.Rename(source, destination)
}

// restoreFile restores an optional previous file state for transactional operations.
func restoreFile(path string, previous []byte, existed bool) error {
	if existed {
		return writeAtomic(path, previous, 0o644)
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// readOptional reads a file while distinguishing absence from I/O failure.
func readOptional(path string) ([]byte, bool, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return data, true, nil
}
