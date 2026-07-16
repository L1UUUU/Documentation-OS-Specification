// This file allocates final monotonic Knowledge identifiers from Work-local drafts.
package engine

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var finalIdentifierPattern = regexp.MustCompile(`^(ARCH|ADR|STD)-([0-9]{4})$`)

// AllocateKnowledgeIdentifier integrates one Work-local draft into a final Knowledge artifact.
func (e *Engine) AllocateKnowledgeIdentifier(categoryName, draftRelativePath string) (result AllocationResult, returnErr error) {
	category, err := e.resolveCategory(categoryName)
	if err != nil {
		return AllocationResult{}, err
	}
	draftPath := draftRelativePath
	if !filepath.IsAbs(draftPath) {
		draftPath = e.path(filepath.FromSlash(draftRelativePath))
	}
	draftPath, err = filepath.Abs(draftPath)
	if err != nil || !e.pathInside(draftPath) {
		return AllocationResult{}, fmt.Errorf("draft path must be inside the repository")
	}
	draftPath = filepath.Clean(draftPath)
	base := filepath.Base(draftPath)
	match := draftFilenamePattern.FindStringSubmatch(base)
	if len(match) != 3 || match[1] != category.Prefix {
		return AllocationResult{}, fmt.Errorf("draft %q must use %s-DRAFT-<slug>.md", e.relativePath(draftPath), category.Prefix)
	}
	categoryPath := e.path(category.Directory)
	relativeToCategory, err := filepath.Rel(categoryPath, filepath.Dir(draftPath))
	if err != nil || relativeToCategory == ".." || strings.HasPrefix(relativeToCategory, ".."+string(filepath.Separator)) {
		return AllocationResult{}, fmt.Errorf("draft must reside beneath %s", category.Directory)
	}
	if _, err := os.Stat(draftPath); err != nil {
		return AllocationResult{}, fmt.Errorf("read draft %s: %w", e.relativePath(draftPath), err)
	}
	lockPath := e.path(filepath.ToSlash(filepath.Join(e.Profile.LockRoot, strings.ToLower(category.Prefix)+"-identity.lock")))
	releaseLock, err := acquireRepositoryLock(lockPath)
	if err != nil {
		return AllocationResult{}, fmt.Errorf("acquire %s identity allocation lock: %w", category.Prefix, err)
	}
	defer func() {
		if err := releaseLock(); err != nil {
			lockErr := fmt.Errorf("release %s identity allocation lock: %w", category.Prefix, err)
			if returnErr == nil {
				returnErr = lockErr
				return
			}
			returnErr = errors.Join(returnErr, lockErr)
		}
	}()
	next, err := e.nextKnowledgeNumber(category)
	if err != nil {
		return AllocationResult{}, err
	}
	slug := match[2]
	if slug == "" {
		slug = strings.ToLower(category.Prefix) + "-draft"
	}
	identifier := fmt.Sprintf("%s-%04d", category.Prefix, next)
	newBase := fmt.Sprintf("%04d-%s.md", next, slug)
	newPath := filepath.Join(filepath.Dir(draftPath), newBase)
	if _, err := os.Stat(newPath); err == nil {
		return AllocationResult{}, fmt.Errorf("allocated Knowledge path already exists: %s", e.relativePath(newPath))
	} else if !os.IsNotExist(err) {
		return AllocationResult{}, err
	}
	originalDraft, err := os.ReadFile(draftPath)
	if err != nil {
		return AllocationResult{}, err
	}
	updatedDraft, err := setFrontMatterField(originalDraft, "id", identifier)
	if err != nil {
		return AllocationResult{}, fmt.Errorf("set final identifier: %w", err)
	}
	oldToken := strings.TrimSuffix(base, ".md")
	changedFiles, err := e.findManagedReferenceChanges(oldToken, base, newBase, identifier)
	if err != nil {
		return AllocationResult{}, err
	}
	delete(changedFiles, draftPath)
	rollback := func() {
		_ = os.Remove(newPath)
		_ = writeAtomic(draftPath, originalDraft, 0o644)
		for path, data := range changedFiles {
			_ = writeAtomic(path, data, 0o644)
		}
	}
	if err := os.Rename(draftPath, newPath); err != nil {
		return AllocationResult{}, fmt.Errorf("rename draft to final artifact: %w", err)
	}
	if err := writeAtomic(newPath, updatedDraft, 0o644); err != nil {
		rollback()
		return AllocationResult{}, err
	}
	updatedReferences := 0
	for path, original := range changedFiles {
		updated := replaceDraftReferences(original, oldToken, base, newBase, identifier)
		if string(updated) == string(original) {
			continue
		}
		if err := writeAtomic(path, updated, 0o644); err != nil {
			rollback()
			return AllocationResult{}, fmt.Errorf("update managed reference %s: %w", e.relativePath(path), err)
		}
		updatedReferences++
	}
	gate := e.validateIdentityAndRelationships()
	if !gate.Passed() {
		rollback()
		return AllocationResult{}, fmt.Errorf("identifier allocation validation gate failed: %s", gate.String())
	}
	return AllocationResult{Identifier: identifier, OldPath: e.relativePath(draftPath), NewPath: e.relativePath(newPath), References: updatedReferences}, nil
}

// resolveCategory resolves a profile category name or prefix.
func (e *Engine) resolveCategory(name string) (KnowledgeCategory, error) {
	normalized := strings.ToLower(strings.TrimSpace(name))
	for _, category := range e.Profile.KnowledgeEntries {
		if normalized == strings.ToLower(category.Name) || normalized == strings.ToLower(category.Prefix) || normalized == strings.TrimSuffix(strings.ToLower(category.Directory), "s") {
			return category, nil
		}
	}
	return KnowledgeCategory{}, fmt.Errorf("unknown Knowledge category %q", name)
}

// nextKnowledgeNumber finds the next monotonic category number from repository state.
func (e *Engine) nextKnowledgeNumber(category KnowledgeCategory) (int, error) {
	max := 0
	err := filepath.WalkDir(e.path(category.Directory), func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" || draftFilenamePattern.MatchString(entry.Name()) {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		front, parseErr := parseFrontMatter(data)
		if parseErr == nil {
			if match := finalIdentifierPattern.FindStringSubmatch(front.ID); len(match) == 3 && match[1] == category.Prefix {
				value, _ := strconv.Atoi(match[2])
				if value > max {
					max = value
				}
			}
		}
		if match := knowledgeFilenamePattern.FindStringSubmatch(entry.Name()); len(match) == 2 {
			value, _ := strconv.Atoi(match[1])
			if value > max {
				max = value
			}
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("scan %s identifiers: %w", category.Directory, err)
	}
	return max + 1, nil
}

// findManagedReferenceChanges snapshots Markdown files containing the draft marker.
func (e *Engine) findManagedReferenceChanges(oldToken, oldBase, newBase, identifier string) (map[string][]byte, error) {
	changed := map[string][]byte{}
	err := filepath.WalkDir(e.Root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(entry.Name()) != ".md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		updated := replaceDraftReferences(data, oldToken, oldBase, newBase, identifier)
		if string(updated) != string(data) {
			changed[path] = data
		}
		return nil
	})
	return changed, err
}

// replaceDraftReferences rewrites exact managed draft tokens without semantic inference.
func replaceDraftReferences(data []byte, oldToken, oldBase, newBase, identifier string) []byte {
	text := strings.ReplaceAll(string(data), oldBase, newBase)
	text = strings.ReplaceAll(text, oldToken, identifier)
	return []byte(text)
}
