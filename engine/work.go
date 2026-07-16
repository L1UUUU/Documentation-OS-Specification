// This file implements deterministic Work workspace generation.
package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// GenerateWork creates a globally unique active Work and refreshes INDEX.md transactionally.
func (e *Engine) GenerateWork(slug string) (WorkResult, error) {
	if err := validateSlug(slug); err != nil {
		return WorkResult{}, err
	}
	if err := e.requireRuntimeRoots(); err != nil {
		return WorkResult{}, err
	}
	activePath := e.path(filepath.ToSlash(filepath.Join(e.Profile.ActiveRoot, slug)))
	completedPath := e.path(filepath.ToSlash(filepath.Join(e.Profile.CompletedRoot, slug)))
	if _, err := os.Stat(activePath); err == nil {
		return WorkResult{}, fmt.Errorf("Work slug %q already exists under active", slug)
	} else if !os.IsNotExist(err) {
		return WorkResult{}, fmt.Errorf("check active Work %q: %w", slug, err)
	}
	if _, err := os.Stat(completedPath); err == nil {
		return WorkResult{}, fmt.Errorf("Work slug %q already exists under completed", slug)
	} else if !os.IsNotExist(err) {
		return WorkResult{}, fmt.Errorf("check completed Work %q: %w", slug, err)
	}
	indexPath := e.path(e.Profile.IndexPath)
	previousIndex, indexExisted, err := readOptional(indexPath)
	if err != nil {
		return WorkResult{}, fmt.Errorf("read existing INDEX: %w", err)
	}
	rollback := func() {
		_ = os.RemoveAll(activePath)
		_ = restoreFile(indexPath, previousIndex, indexExisted)
	}
	if err := os.MkdirAll(filepath.Join(activePath, "issues"), 0o755); err != nil {
		rollback()
		return WorkResult{}, fmt.Errorf("create Work workspace: %w", err)
	}
	prd := []byte("---\nrelationships: []\n---\n\n# " + titleFromSlug(slug) + "\n\n## Objective\n\n<!-- Define the Work objective and constraints. -->\n")
	if err := writeAtomic(filepath.Join(activePath, "PRD.md"), prd, 0o644); err != nil {
		rollback()
		return WorkResult{}, fmt.Errorf("create PRD.md: %w", err)
	}
	if err := writeAtomic(filepath.Join(activePath, "HANDOFF.md"), nil, 0o644); err != nil {
		rollback()
		return WorkResult{}, fmt.Errorf("create HANDOFF.md: %w", err)
	}
	if _, err := e.GenerateIndex(); err != nil {
		rollback()
		return WorkResult{}, fmt.Errorf("generate INDEX after Work creation: %w", err)
	}
	return WorkResult{Slug: slug, Path: activePath, IndexPath: e.relativePath(indexPath)}, nil
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
