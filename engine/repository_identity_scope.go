// This file resolves Git common-directory scope for repository-wide Identity allocation.
package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// identityAllocationScope returns the shared lock path and visible worktree roots.
func (e *Engine) identityAllocationScope(category KnowledgeCategory) (string, []string, error) {
	commonGitDir, managed, err := e.gitCommonDir()
	if err != nil {
		return "", nil, err
	}
	lockName := strings.ToLower(category.Prefix) + "-identity.lock"
	if !managed {
		return e.path(filepath.ToSlash(filepath.Join(e.Profile.LockRoot, lockName))), []string{e.Root}, nil
	}
	roots, err := e.commonGitWorktreeRoots(commonGitDir)
	if err != nil {
		return "", nil, err
	}
	lockPath := filepath.Join(commonGitDir, filepath.FromSlash(e.Profile.LockRoot), lockName)
	return lockPath, roots, nil
}

// gitCommonDir resolves the repository's shared Git directory without invoking Git.
func (e *Engine) gitCommonDir() (string, bool, error) {
	gitEntry := filepath.Join(e.Root, ".git")
	info, err := os.Stat(gitEntry)
	if os.IsNotExist(err) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("inspect Git metadata: %w", err)
	}
	gitDir := gitEntry
	if !info.IsDir() {
		content, err := os.ReadFile(gitEntry)
		if err != nil {
			return "", false, fmt.Errorf("read Git worktree metadata: %w", err)
		}
		const prefix = "gitdir:"
		value := strings.TrimSpace(string(content))
		if !strings.HasPrefix(strings.ToLower(value), prefix) {
			return "", false, fmt.Errorf("Git worktree metadata does not declare gitdir")
		}
		gitDir = strings.TrimSpace(value[len(prefix):])
		if !filepath.IsAbs(gitDir) {
			gitDir = filepath.Join(e.Root, gitDir)
		}
	}
	commonDirPath := filepath.Join(gitDir, "commondir")
	content, err := os.ReadFile(commonDirPath)
	if os.IsNotExist(err) {
		return filepath.Clean(gitDir), true, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("read Git common directory metadata: %w", err)
	}
	commonGitDir := strings.TrimSpace(string(content))
	if !filepath.IsAbs(commonGitDir) {
		commonGitDir = filepath.Join(gitDir, commonGitDir)
	}
	return filepath.Clean(commonGitDir), true, nil
}

// commonGitWorktreeRoots lists the primary checkout and registered linked worktrees.
func (e *Engine) commonGitWorktreeRoots(commonGitDir string) ([]string, error) {
	roots := map[string]bool{filepath.Clean(e.Root): true}
	if filepath.Base(commonGitDir) == ".git" {
		roots[filepath.Clean(filepath.Dir(commonGitDir))] = true
	}
	entries, err := os.ReadDir(filepath.Join(commonGitDir, "worktrees"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("list Git worktrees: %w", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		content, err := os.ReadFile(filepath.Join(commonGitDir, "worktrees", entry.Name(), "gitdir"))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read Git worktree %s metadata: %w", entry.Name(), err)
		}
		gitEntry := strings.TrimSpace(string(content))
		if gitEntry != "" {
			roots[filepath.Clean(filepath.Dir(gitEntry))] = true
		}
	}
	result := make([]string, 0, len(roots))
	for root := range roots {
		if info, err := os.Stat(root); err == nil && info.IsDir() {
			result = append(result, root)
		}
	}
	sort.Strings(result)
	return result, nil
}
