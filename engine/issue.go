// This file implements deterministic Issue creation within an Active Work.
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

var issueFilePattern = regexp.MustCompile(`^([0-9]{2})-([a-z0-9]+(?:-[a-z0-9]+)*)\.md$`)

// CreateIssue atomically appends one caller-authored Issue to an Active Work.
func (e *Engine) CreateIssue(input CreateIssueInput) (result CreateIssueResult, returnErr error) {
	defer func() {
		returnErr = withLifecycleStage(LifecycleStageIssue, returnErr)
	}()

	if err := validateCreateIssueInput(input); err != nil {
		return CreateIssueResult{}, err
	}
	if err := e.requireRuntimeRoots(); err != nil {
		return CreateIssueResult{}, fmt.Errorf("%w: %v", ErrInvalidRepository, err)
	}

	lockPath := e.path(filepath.ToSlash(filepath.Join(e.Profile.LockRoot, "create-issue.lock")))
	releaseLock, err := acquireRepositoryLock(lockPath)
	if err != nil {
		return CreateIssueResult{}, fmt.Errorf("acquire Create Issue lock: %w", err)
	}
	defer func() {
		if err := releaseLock(); err != nil {
			lockErr := fmt.Errorf("release Create Issue lock: %w", err)
			if returnErr == nil {
				returnErr = lockErr
				return
			}
			returnErr = errors.Join(returnErr, lockErr)
		}
	}()

	issuesPath, err := e.activeIssueWork(input.WorkSlug)
	if err != nil {
		return CreateIssueResult{}, err
	}
	scan, err := e.scanIssues(issuesPath, input)
	if err != nil {
		return CreateIssueResult{}, err
	}
	if scan.existing != nil {
		if !scan.existing.matches(input) {
			return CreateIssueResult{}, fmt.Errorf("%w: Issue slug %q already exists with different title, status, or body", ErrConflict, input.Slug)
		}
		result := CreateIssueResult{Number: scan.existing.number, Name: scan.existing.name, Path: scan.existing.path, Created: false}
		if _, err := e.GenerateIndex(); err != nil {
			return result, fmt.Errorf("reconcile INDEX for Issue %q: %w", scan.existing.name, err)
		}
		return result, nil
	}
	if scan.maxNumber >= 99 {
		return CreateIssueResult{}, fmt.Errorf("%w: Work %q has no Issue number available after 99", ErrConflict, input.WorkSlug)
	}

	number := scan.maxNumber + 1
	name := fmt.Sprintf("%02d-%s.md", number, input.Slug)
	issuePath := filepath.Join(issuesPath, name)
	indexPath := e.path(e.Profile.IndexPath)
	previousIndex, indexExisted, err := readOptional(indexPath)
	if err != nil {
		return CreateIssueResult{}, fmt.Errorf("read existing INDEX: %w", err)
	}
	stagingRoot := e.path(filepath.ToSlash(filepath.Join(e.Profile.RuntimeRoot, ".issue-staging")))
	if err := os.MkdirAll(stagingRoot, 0o755); err != nil {
		return CreateIssueResult{}, fmt.Errorf("create Issue staging root: %w", err)
	}
	stagingPath, err := os.MkdirTemp(stagingRoot, input.WorkSlug+"-")
	if err != nil {
		return CreateIssueResult{}, fmt.Errorf("create Issue staging directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(stagingPath); err != nil {
			cleanupErr := fmt.Errorf("remove Issue staging directory: %w", err)
			if returnErr == nil {
				returnErr = cleanupErr
				return
			}
			returnErr = errors.Join(returnErr, cleanupErr)
		}
	}()
	stagedIssuePath := filepath.Join(stagingPath, name)
	if err := e.writeFileAtomic(stagedIssuePath, renderCreatedIssue(input), 0o644); err != nil {
		return CreateIssueResult{}, fmt.Errorf("publish Issue %q: %w", name, err)
	}
	if err := e.renameFile(stagedIssuePath, issuePath); err != nil {
		return CreateIssueResult{}, fmt.Errorf("publish staged Issue %q: %w", name, err)
	}
	rollback := func(operationErr error) error {
		rollbackErr := errors.Join(
			removeFileIfPresent(issuePath),
			restoreFile(indexPath, previousIndex, indexExisted),
		)
		if rollbackErr != nil {
			return fmt.Errorf("publish Issue and INDEX: %w", errors.Join(operationErr, fmt.Errorf("rollback Issue and INDEX: %w", rollbackErr)))
		}
		return fmt.Errorf("publish Issue and INDEX; Issue and INDEX rolled back: %w", operationErr)
	}
	if _, err := e.GenerateIndex(); err != nil {
		return CreateIssueResult{}, rollback(err)
	}
	return CreateIssueResult{Number: number, Name: name, Path: e.relativePath(issuePath), Created: true}, nil
}

// activeIssueWork validates the target Work state and required Active Runtime structure.
func (e *Engine) activeIssueWork(slug string) (string, error) {
	activePath := e.path(filepath.ToSlash(filepath.Join(e.Profile.ActiveRoot, slug)))
	completedPath := e.path(filepath.ToSlash(filepath.Join(e.Profile.CompletedRoot, slug)))
	activeInfo, activeExists, err := statOptional(activePath)
	if err != nil {
		return "", fmt.Errorf("%w: inspect active Work %q: %v", ErrInvalidRepository, slug, err)
	}
	completedInfo, completedExists, err := statOptional(completedPath)
	if err != nil {
		return "", fmt.Errorf("%w: inspect completed Work %q: %v", ErrInvalidRepository, slug, err)
	}
	if activeExists && completedExists {
		return "", fmt.Errorf("%w: Work %q exists under both active and completed", ErrInvalidRepository, slug)
	}
	if !activeExists && completedExists {
		if !completedInfo.IsDir() {
			return "", fmt.Errorf("%w: completed Work path %q is not a directory", ErrInvalidRepository, slug)
		}
		return "", fmt.Errorf("%w: Work %q exists only under completed", ErrCompletedWork, slug)
	}
	if !activeExists {
		return "", fmt.Errorf("%w: Work %q does not exist", ErrWorkNotFound, slug)
	}
	if !activeInfo.IsDir() {
		return "", fmt.Errorf("%w: active Work path %q is not a directory", ErrInvalidRepository, slug)
	}
	for _, name := range []string{"PRD.md", "HANDOFF.md"} {
		info, exists, statErr := statOptional(filepath.Join(activePath, name))
		if statErr != nil || !exists || !info.Mode().IsRegular() {
			return "", fmt.Errorf("%w: active Work %q requires regular %s", ErrInvalidRepository, slug, name)
		}
	}
	prd, err := os.ReadFile(filepath.Join(activePath, "PRD.md"))
	if err != nil {
		return "", fmt.Errorf("%w: read active Work %q PRD.md: %v", ErrInvalidRepository, slug, err)
	}
	front, err := parseFrontMatter(prd)
	if err != nil || !front.Present || frontMatterFieldPresent(front, "outcome") {
		return "", fmt.Errorf("%w: active Work %q has invalid PRD.md front matter", ErrInvalidRepository, slug)
	}
	issuesPath := filepath.Join(activePath, "issues")
	issuesInfo, issuesExist, err := statOptional(issuesPath)
	if err != nil || !issuesExist || !issuesInfo.IsDir() {
		return "", fmt.Errorf("%w: active Work %q requires an issues directory", ErrInvalidRepository, slug)
	}
	return issuesPath, nil
}

type issueScan struct {
	maxNumber int
	existing  *persistedIssue
}

type persistedIssue struct {
	number int
	name   string
	path   string
	title  string
	status string
	body   string
}

func (issue persistedIssue) matches(input CreateIssueInput) bool {
	return issue.title == input.Title && issue.status == input.Status && issue.body == input.Body
}

// scanIssues validates existing names and finds allocation and idempotency state.
func (e *Engine) scanIssues(issuesPath string, input CreateIssueInput) (issueScan, error) {
	entries, err := os.ReadDir(issuesPath)
	if err != nil {
		return issueScan{}, fmt.Errorf("%w: read issues for Work %q: %v", ErrInvalidRepository, input.WorkSlug, err)
	}
	numbers := map[int]string{}
	slugs := map[string]string{}
	result := issueScan{}
	for _, entry := range entries {
		if entry.IsDir() {
			return issueScan{}, fmt.Errorf("%w: Issue entry %q is a directory", ErrInvalidRepository, entry.Name())
		}
		match := issueFilePattern.FindStringSubmatch(entry.Name())
		if len(match) != 3 {
			return issueScan{}, fmt.Errorf("%w: malformed Issue filename %q", ErrInvalidRepository, entry.Name())
		}
		number, _ := strconv.Atoi(match[1])
		if number < 1 || number > 99 {
			return issueScan{}, fmt.Errorf("%w: Issue filename %q has number outside 01-99", ErrInvalidRepository, entry.Name())
		}
		if previous, duplicate := numbers[number]; duplicate {
			return issueScan{}, fmt.Errorf("%w: Issue number %02d is duplicated by %q and %q", ErrInvalidRepository, number, previous, entry.Name())
		}
		numbers[number] = entry.Name()
		if previous, duplicate := slugs[match[2]]; duplicate {
			return issueScan{}, fmt.Errorf("%w: Issue slug %q is duplicated by %q and %q", ErrInvalidRepository, match[2], previous, entry.Name())
		}
		slugs[match[2]] = entry.Name()
		if number > result.maxNumber {
			result.maxNumber = number
		}
		data, err := os.ReadFile(filepath.Join(issuesPath, entry.Name()))
		if err != nil {
			return issueScan{}, fmt.Errorf("%w: read Issue %q: %v", ErrInvalidRepository, entry.Name(), err)
		}
		front, body, err := parseIssueContent(data)
		if err != nil || !front.Present || !issueStatuses[front.Status] {
			return issueScan{}, fmt.Errorf("%w: Issue %q has invalid front matter: %v", ErrInvalidRepository, entry.Name(), err)
		}
		if match[2] == input.Slug {
			result.existing = &persistedIssue{number: number, name: entry.Name(), path: e.relativePath(filepath.Join(issuesPath, entry.Name())), title: front.Title, status: front.Status, body: body}
		}
	}
	return result, nil
}

func validateCreateIssueInput(input CreateIssueInput) error {
	if err := validateSlug(input.WorkSlug); err != nil {
		return fmt.Errorf("%w: Work %v", ErrInvalidInput, err)
	}
	if err := validateSlug(input.Slug); err != nil {
		return fmt.Errorf("%w: Issue %v", ErrInvalidInput, err)
	}
	if strings.TrimSpace(input.Title) == "" || strings.ContainsAny(input.Title, "\r\n") {
		return fmt.Errorf("%w: Issue title must be a non-empty single line", ErrInvalidInput)
	}
	if !issueStatuses[input.Status] {
		return fmt.Errorf("%w: invalid Issue status %q", ErrInvalidInput, input.Status)
	}
	if strings.TrimSpace(input.Body) == "" {
		return fmt.Errorf("%w: Issue body must contain non-whitespace Markdown", ErrInvalidInput)
	}
	return nil
}

func renderCreatedIssue(input CreateIssueInput) []byte {
	return []byte("---\nstatus: " + input.Status + "\ntitle: " + strconv.Quote(input.Title) + "\n---\n" + input.Body)
}

func parseIssueContent(data []byte) (FrontMatter, string, error) {
	front, err := parseFrontMatter(data)
	if err != nil || !front.Present {
		return front, "", err
	}
	text := string(data)
	firstEnd := strings.IndexByte(text, '\n')
	_, closingEnd, err := frontMatterBounds(text, firstEnd)
	if err != nil {
		return FrontMatter{}, "", err
	}
	bodyStart := closingEnd
	if strings.HasPrefix(text[bodyStart:], "\r\n") {
		bodyStart += 2
	} else if strings.HasPrefix(text[bodyStart:], "\n") {
		bodyStart++
	}
	return front, text[bodyStart:], nil
}

func statOptional(path string) (os.FileInfo, bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return info, true, nil
}

func removeFileIfPresent(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
