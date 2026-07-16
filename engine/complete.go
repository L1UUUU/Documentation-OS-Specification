// This file implements the atomic Complete stage and idempotent Cleanup stage.
package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

var issueFilenamePattern = regexp.MustCompile(`^[0-9]{2}-[a-z0-9]+(?:-[a-z0-9]+)*\.md$`)

var nonTerminalIssueStatuses = map[string]bool{
	"open": true, "in-progress": true, "blocked": true,
}

// ephemeralBackup captures one temporary file for Complete rollback.
type ephemeralBackup struct {
	Path      string
	Data      []byte
	Mode      os.FileMode
	Directory bool
}

// Complete executes Complete and Cleanup, or resumes Cleanup for a completed Work.
func (e *Engine) Complete(slug, outcome string) (CompleteResult, error) {
	if err := validateSlug(slug); err != nil {
		return CompleteResult{}, fmt.Errorf("%w: %v", ErrPreflight, err)
	}
	if !validOutcome(outcome) {
		return CompleteResult{}, fmt.Errorf("%w: outcome %q is not one of succeeded, cancelled, superseded, failed", ErrPreflight, outcome)
	}
	activePath := e.path(filepath.ToSlash(filepath.Join(e.Profile.ActiveRoot, slug)))
	completedPath := e.path(filepath.ToSlash(filepath.Join(e.Profile.CompletedRoot, slug)))
	activeExists, activeErr := pathExists(activePath)
	if activeErr != nil {
		return CompleteResult{}, activeErr
	}
	completedExists, completedErr := pathExists(completedPath)
	if completedErr != nil {
		return CompleteResult{}, completedErr
	}
	if activeExists && completedExists {
		return CompleteResult{}, fmt.Errorf("%w: Work %q exists under both active and completed", ErrPreflight, slug)
	}
	if activeExists {
		if info, err := os.Stat(activePath); err != nil || !info.IsDir() {
			return CompleteResult{}, fmt.Errorf("%w: active Work path %q is not a directory", ErrPreflight, slug)
		}
	}
	if completedExists {
		if info, err := os.Stat(completedPath); err != nil || !info.IsDir() {
			return CompleteResult{}, fmt.Errorf("%w: completed Work path %q is occupied by a non-directory", ErrPreflight, slug)
		}
		index, err := e.GenerateIndex()
		result := CompleteResult{Slug: slug, Completed: true, CleanupCompleted: err == nil, RetriedCleanup: true, Outcome: outcome}
		if err != nil {
			return result, fmt.Errorf("cleanup completed Work %q: %w", slug, err)
		}
		_ = index
		return result, nil
	}
	if !activeExists {
		return CompleteResult{}, fmt.Errorf("%w: active Work %q does not exist", ErrPreflight, slug)
	}
	if err := e.completePreflight(activePath, outcome); err != nil {
		return CompleteResult{}, err
	}
	prdPath := filepath.Join(activePath, "PRD.md")
	originalPRD, err := os.ReadFile(prdPath)
	if err != nil {
		return CompleteResult{}, fmt.Errorf("%w: read PRD.md: %v", ErrPreflight, err)
	}
	backups, err := captureEphemeralFiles(activePath)
	if err != nil {
		return CompleteResult{}, fmt.Errorf("prepare Complete rollback: %w", err)
	}
	updatedPRD, err := setFrontMatterField(originalPRD, "outcome", outcome)
	if err != nil {
		return CompleteResult{}, fmt.Errorf("%w: record outcome in PRD.md: %v", ErrPreflight, err)
	}
	if err := writeAtomic(prdPath, updatedPRD, 0o644); err != nil {
		return CompleteResult{}, fmt.Errorf("record Work outcome: %w", err)
	}
	if err := removeEphemeralFiles(activePath); err != nil {
		_ = restoreCompleteTransaction(activePath, prdPath, originalPRD, backups)
		return CompleteResult{}, fmt.Errorf("Complete cleanup before move: %w", err)
	}
	if err := e.renameFile(activePath, completedPath); err != nil {
		_ = restoreCompleteTransaction(activePath, prdPath, originalPRD, backups)
		return CompleteResult{}, fmt.Errorf("Complete stage move %s to %s: %w", e.relativePath(activePath), e.relativePath(completedPath), err)
	}
	result := CompleteResult{Slug: slug, Completed: true, Outcome: outcome}
	if _, err := e.GenerateIndex(); err != nil {
		return result, fmt.Errorf("Complete stage succeeded; Cleanup failed for %q, retry Complete to regenerate INDEX: %w", slug, err)
	}
	result.CleanupCompleted = true
	return result, nil
}

// completePreflight validates all conditions that must hold before mutation.
func (e *Engine) completePreflight(workPath, outcome string) error {
	for _, relative := range []string{"PRD.md", "issues", "HANDOFF.md"} {
		info, err := os.Stat(filepath.Join(workPath, relative))
		if err != nil {
			return fmt.Errorf("%w: missing Core Runtime Asset %s: %v", ErrPreflight, relative, err)
		}
		if relative == "issues" && !info.IsDir() {
			return fmt.Errorf("%w: issues is not a directory", ErrPreflight)
		}
		if relative != "issues" && info.IsDir() {
			return fmt.Errorf("%w: %s is a directory", ErrPreflight, relative)
		}
	}
	issueEntries, err := os.ReadDir(filepath.Join(workPath, "issues"))
	if err != nil {
		return fmt.Errorf("%w: read issues directory: %v", ErrPreflight, err)
	}
	var issueNames []string
	for _, entry := range issueEntries {
		if !entry.IsDir() {
			issueNames = append(issueNames, entry.Name())
		}
	}
	sort.Strings(issueNames)
	if len(issueNames) == 0 {
		return fmt.Errorf("%w: issues directory must contain at least one Issue before Complete", ErrPreflight)
	}
	prd, err := os.ReadFile(filepath.Join(workPath, "PRD.md"))
	if err != nil {
		return fmt.Errorf("%w: read PRD.md: %v", ErrPreflight, err)
	}
	front, err := parseFrontMatter(prd)
	if err != nil || !front.Present {
		return fmt.Errorf("%w: PRD.md must contain valid YAML front matter", ErrPreflight)
	}
	if frontMatterFieldPresent(front, "outcome") {
		return fmt.Errorf("%w: active Work already declares an outcome", ErrPreflight)
	}
	for _, name := range issueNames {
		if !issueFilenamePattern.MatchString(name) {
			return fmt.Errorf("%w: Issue filename %q must match NN-<slug>.md", ErrPreflight, name)
		}
		data, readErr := os.ReadFile(filepath.Join(workPath, "issues", name))
		if readErr != nil {
			return fmt.Errorf("%w: read Issue %s: %v", ErrPreflight, name, readErr)
		}
		issue, parseErr := parseFrontMatter(data)
		if parseErr != nil || !issue.Present || issue.Status == "" {
			return fmt.Errorf("%w: Issue %s must declare a status", ErrPreflight, name)
		}
		if !issueStatuses[issue.Status] {
			return fmt.Errorf("%w: Issue %s has invalid status %q", ErrPreflight, name, issue.Status)
		}
		if outcome == OutcomeSucceeded && nonTerminalIssueStatuses[issue.Status] {
			return fmt.Errorf("%w: succeeded Work has non-terminal Issue %s with status %s", ErrPreflight, name, issue.Status)
		}
	}
	return nil
}

// pathExists distinguishes an absent path from filesystem errors.
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// captureEphemeralFiles snapshots non-core files before Complete cleanup.
func captureEphemeralFiles(workPath string) ([]ephemeralBackup, error) {
	entries, err := os.ReadDir(workPath)
	if err != nil {
		return nil, err
	}
	var backups []ephemeralBackup
	for _, entry := range entries {
		if entry.Name() == "PRD.md" || entry.Name() == "issues" || entry.Name() == "HANDOFF.md" {
			continue
		}
		root := filepath.Join(workPath, entry.Name())
		walkErr := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			relative, relErr := filepath.Rel(workPath, path)
			if relErr != nil {
				return relErr
			}
			if info.IsDir() {
				backups = append(backups, ephemeralBackup{Path: relative, Mode: info.Mode(), Directory: true})
				return nil
			}
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			backups = append(backups, ephemeralBackup{Path: relative, Data: data, Mode: info.Mode()})
			return nil
		})
		if walkErr != nil {
			return nil, walkErr
		}
	}
	return backups, nil
}

// removeEphemeralFiles removes all entries except Core Runtime Assets.
func removeEphemeralFiles(workPath string) error {
	entries, err := os.ReadDir(workPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Name() == "PRD.md" || entry.Name() == "issues" || entry.Name() == "HANDOFF.md" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(workPath, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

// restoreCompleteTransaction restores the active Work after a failed directory move.
func restoreCompleteTransaction(workPath, prdPath string, originalPRD []byte, backups []ephemeralBackup) error {
	if err := writeAtomic(prdPath, originalPRD, 0o644); err != nil {
		return err
	}
	for _, backup := range backups {
		if !backup.Directory {
			continue
		}
		if err := os.MkdirAll(filepath.Join(workPath, backup.Path), backup.Mode.Perm()); err != nil {
			return err
		}
	}
	for _, backup := range backups {
		if backup.Directory {
			continue
		}
		if err := writeAtomic(filepath.Join(workPath, backup.Path), backup.Data, backup.Mode); err != nil {
			return err
		}
	}
	return nil
}
