// This file implements the read-only structural, identity, relationship, and lifecycle checks.
package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var knowledgeIDPattern = regexp.MustCompile(`^(ARCH|ADR|STD)-([0-9]{4})$`)
var knowledgeFilenamePattern = regexp.MustCompile(`^([0-9]{4})-.+\.md$`)
var draftFilenamePattern = regexp.MustCompile(`^(ARCH|ADR|STD)-DRAFT(?:-([a-z0-9-]+))?\.md$`)

var knowledgeStatuses = map[string]bool{
	"created": true, "active": true, "updated": true, "validated": true, "archived": true, "retired": true,
}

var issueStatuses = map[string]bool{
	"open": true, "in-progress": true, "done": true, "blocked": true, "cancelled": true, "superseded": true,
}

var relationshipTypes = map[string]bool{
	"references": true, "depends-on": true, "implements": true, "produces": true, "affects": true, "supersedes": true,
}

type knowledgeArtifact struct {
	ID            string
	Status        string
	Path          string
	Relationships []Relationship
}

type runtimeRelationship struct {
	Source string
	Path   string
	Links  []Relationship
}

// Validate verifies repository consistency without changing any file.
func (e *Engine) Validate() (ValidationReport, error) {
	return e.ValidateContext(context.Background())
}

// ValidateContext verifies repository consistency while honoring cancellation.
// It is behaviorally identical to Validate in normal builds. A conformance-tag
// build may install the internal validation seam through NewConformance.
func (e *Engine) ValidateContext(ctx context.Context) (ValidationReport, error) {
	if err := ctx.Err(); err != nil {
		return ValidationReport{}, err
	}
	if e.conformance != nil {
		if err := e.conformance.BeforeValidate(ctx); err != nil {
			return ValidationReport{}, withLifecycleStage(LifecycleStageValidate, err)
		}
	}
	report := ValidationReport{Status: "passed"}
	e.validateRepositoryStructure(&report)
	knowledge := e.validateKnowledge(&report)
	runtime := e.validateRuntime(&report)
	e.validateRelationships(&report, knowledge, runtime)
	e.validateGeneratedIndex(&report)
	report.sortIssues()
	return report, nil
}

// validateRepositoryStructure checks required profile locations and guidance mirrors.
func (e *Engine) validateRepositoryStructure(report *ValidationReport) {
	for _, relative := range []string{e.Profile.KnowledgeRoot, e.Profile.RuntimeRoot, e.Profile.ActiveRoot, e.Profile.CompletedRoot, e.Profile.ArchitectureDir, e.Profile.ADRDir, e.Profile.StandardsDir, e.Profile.InboxDir} {
		if err := e.requireDirectory(relative); err != nil {
			report.add(ValidationIssue{Category: "Repository Structure", Artifact: relative, Rule: "required-directory", Reason: err.Error(), Recovery: "create the required Single Repository Profile directory"})
		}
	}
	for _, relative := range []string{e.Profile.RootAgents, e.Profile.RootClaude, e.Profile.ScratchAgents, e.Profile.ScratchClaude, e.Profile.IndexPath} {
		if err := e.requireRegularFile(relative); err != nil {
			report.add(ValidationIssue{Category: "Repository Structure", Artifact: relative, Rule: "required-file", Reason: err.Error(), Recovery: "run `dos init` or restore the required file"})
		}
	}
	e.validateGuidanceMirrors(report)
	archive := e.path(filepath.Join(e.Profile.RuntimeRoot, "archive"))
	if info, err := os.Stat(archive); err == nil && info.IsDir() {
		report.add(ValidationIssue{Category: "Repository Structure", Artifact: e.relativePath(archive), Rule: "no-archive-directory", Reason: "the Single Repository Profile has no .scratch/archive directory", Recovery: "move Work directories to active/ or completed/"})
	}
}

// validateGuidanceMirrors verifies every directory that declares Agent guidance.
func (e *Engine) validateGuidanceMirrors(report *ValidationReport) {
	seen := map[string]bool{}
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
		if entry.Name() != "AGENTS.md" && entry.Name() != "CLAUDE.md" {
			return nil
		}
		directory := filepath.Dir(path)
		if seen[directory] {
			return nil
		}
		seen[directory] = true
		agentsPath := filepath.Join(directory, "AGENTS.md")
		claudePath := filepath.Join(directory, "CLAUDE.md")
		agents, agentsErr := os.ReadFile(agentsPath)
		claude, claudeErr := os.ReadFile(claudePath)
		if os.IsNotExist(agentsErr) || os.IsNotExist(claudeErr) {
			report.add(ValidationIssue{Category: "Repository Structure", Artifact: e.relativePath(directory), Rule: "agent-entry-mirror", Reason: "AGENTS.md and CLAUDE.md must both exist in a guidance scope", Recovery: "restore the missing mirror"})
			return nil
		}
		if agentsErr != nil || claudeErr != nil {
			report.add(ValidationIssue{Category: "Repository Structure", Artifact: e.relativePath(directory), Rule: "agent-entry-readable", Reason: "could not read both Agent Entry files", Recovery: "check file permissions and encoding"})
			return nil
		}
		if string(agents) != string(claude) {
			report.add(ValidationIssue{Category: "Repository Structure", Artifact: e.relativePath(directory), Rule: "agent-entry-equivalence", Reason: "AGENTS.md and CLAUDE.md are not content-equivalent", Recovery: "synchronize the mirror from the canonical AGENTS.md"})
		}
		return nil
	})
	if err != nil {
		report.add(ValidationIssue{Category: "Repository Structure", Artifact: ".", Rule: "guidance-discovery", Reason: err.Error(), Recovery: "check repository traversal permissions"})
	}
}

// validateKnowledge validates identity-managed Knowledge artifacts and returns resolvable IDs.
func (e *Engine) validateKnowledge(report *ValidationReport) map[string]knowledgeArtifact {
	artifacts := map[string]knowledgeArtifact{}
	for _, category := range e.Profile.KnowledgeEntries {
		root := e.path(category.Directory)
		if _, err := os.Stat(root); err != nil {
			continue
		}
		walkErr := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				report.add(ValidationIssue{Category: "Identity", Artifact: e.relativePath(path), Rule: "read-artifact", Reason: err.Error(), Recovery: "restore the artifact or fix permissions"})
				return nil
			}
			if entry.IsDir() {
				return nil
			}
			if filepath.Ext(entry.Name()) != ".md" || entry.Name() == "AGENTS.md" || entry.Name() == "CLAUDE.md" {
				return nil
			}
			if draftFilenamePattern.MatchString(entry.Name()) {
				return nil
			}
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				report.add(ValidationIssue{Category: "Identity", Artifact: e.relativePath(path), Rule: "read-artifact", Reason: readErr.Error(), Recovery: "restore the Knowledge artifact"})
				return nil
			}
			front, parseErr := parseFrontMatter(data)
			if parseErr != nil {
				report.add(ValidationIssue{Category: "Identity", Artifact: e.relativePath(path), Rule: "front-matter", Reason: parseErr.Error(), Recovery: "repair the YAML front matter"})
				return nil
			}
			if !front.Present {
				report.add(ValidationIssue{Category: "Identity", Artifact: e.relativePath(path), Rule: "identity-front-matter", Reason: "Knowledge artifacts must declare YAML front matter", Recovery: "add id, status, and optional relationships"})
				return nil
			}
			match := knowledgeIDPattern.FindStringSubmatch(front.ID)
			if len(match) != 3 {
				report.add(ValidationIssue{Category: "Identity", Artifact: e.relativePath(path), Rule: "identifier-format", Reason: fmt.Sprintf("identifier %q must match PREFIX-NNNN", front.ID), Recovery: "allocate a final ARCH-NNNN, ADR-NNNN, or STD-NNNN identifier"})
				return nil
			}
			if match[1] != category.Prefix {
				report.add(ValidationIssue{Category: "Identity", Artifact: e.relativePath(path), Rule: "category-prefix", Reason: fmt.Sprintf("identifier prefix %s does not belong in %s", match[1], category.Name), Recovery: "move the artifact or correct its identifier"})
			}
			filenameMatch := knowledgeFilenamePattern.FindStringSubmatch(entry.Name())
			if len(filenameMatch) != 2 || filenameMatch[1] != match[2] {
				report.add(ValidationIssue{Category: "Identity", Artifact: e.relativePath(path), Rule: "filename-number", Reason: "the filename's leading number must match the identifier", Recovery: "rename the file to NNNN-<slug>.md"})
			}
			if !knowledgeStatuses[front.Status] {
				report.add(ValidationIssue{Category: "Lifecycle", Artifact: e.relativePath(path), Rule: "knowledge-status", Reason: fmt.Sprintf("invalid Knowledge status %q", front.Status), Recovery: "use created, active, updated, validated, archived, or retired"})
			}
			if _, exists := artifacts[front.ID]; exists {
				report.add(ValidationIssue{Category: "Identity", Artifact: e.relativePath(path), Rule: "identifier-unique", Reason: fmt.Sprintf("identifier %s is duplicated", front.ID), Recovery: "allocate a distinct monotonic identifier"})
			} else {
				artifacts[front.ID] = knowledgeArtifact{ID: front.ID, Status: front.Status, Path: path, Relationships: front.Relationships}
			}
			return nil
		})
		if walkErr != nil {
			report.add(ValidationIssue{Category: "Identity", Artifact: category.Directory, Rule: "knowledge-discovery", Reason: walkErr.Error(), Recovery: "check directory traversal permissions"})
		}
	}
	return artifacts
}

// validateRuntime validates Work locations, Core Runtime Assets, Issues, and outcomes.
func (e *Engine) validateRuntime(report *ValidationReport) map[string]string {
	states := map[string]string{}
	for _, state := range []string{"active", "completed"} {
		rootRelative := e.Profile.ActiveRoot
		if state == "completed" {
			rootRelative = e.Profile.CompletedRoot
		}
		root := e.path(rootRelative)
		entries, err := os.ReadDir(root)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			report.add(ValidationIssue{Category: "Lifecycle", Artifact: rootRelative, Rule: "work-discovery", Reason: err.Error(), Recovery: "check Runtime directory permissions"})
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			slug := entry.Name()
			if err := validateSlug(slug); err != nil {
				report.add(ValidationIssue{Category: "Lifecycle", Artifact: filepath.ToSlash(filepath.Join(rootRelative, slug)), Rule: "workstream-slug", Reason: err.Error(), Recovery: "rename the Work to lowercase kebab-case"})
			}
			if previous, exists := states[slug]; exists {
				report.add(ValidationIssue{Category: "Lifecycle", Artifact: filepath.ToSlash(filepath.Join(rootRelative, slug)), Rule: "work-global-unique", Reason: fmt.Sprintf("Work slug also exists under %s", previous), Recovery: "keep one observable Work location"})
			} else {
				states[slug] = state
			}
			e.validateWork(report, rootRelative, slug, state)
		}
	}
	return states
}

// validateWork validates one Work's structure and metadata.
func (e *Engine) validateWork(report *ValidationReport, rootRelative, slug, state string) {
	workRelative := filepath.ToSlash(filepath.Join(rootRelative, slug))
	workPath := e.path(workRelative)
	if _, err := os.Stat(filepath.Join(workPath, "work.yaml")); err == nil {
		report.add(ValidationIssue{Category: "Repository Structure", Artifact: filepath.ToSlash(filepath.Join(workRelative, "work.yaml")), Rule: "no-work-yaml", Reason: "Runtime Work metadata belongs in PRD.md front matter; work.yaml is not part of the profile", Recovery: "move Work metadata into PRD.md front matter"})
	}
	for _, name := range []string{"PRD.md", "issues", "HANDOFF.md"} {
		info, err := os.Stat(filepath.Join(workPath, name))
		if err != nil {
			report.add(ValidationIssue{Category: "Repository Structure", Artifact: filepath.ToSlash(filepath.Join(workRelative, name)), Rule: "core-runtime-asset", Reason: fmt.Sprintf("missing Core Runtime Asset: %v", err), Recovery: "restore PRD.md, issues/, or HANDOFF.md"})
			continue
		}
		if name == "issues" && !info.IsDir() {
			report.add(ValidationIssue{Category: "Repository Structure", Artifact: filepath.ToSlash(filepath.Join(workRelative, name)), Rule: "issues-directory", Reason: "issues must be a directory", Recovery: "create an issues/ directory"})
		}
		if name != "issues" && info.IsDir() {
			report.add(ValidationIssue{Category: "Repository Structure", Artifact: filepath.ToSlash(filepath.Join(workRelative, name)), Rule: "core-runtime-file", Reason: name + " must be a file", Recovery: "restore the Core Runtime Asset as a file"})
		}
	}
	prdPath := filepath.Join(workPath, "PRD.md")
	var prd FrontMatter
	if data, err := os.ReadFile(prdPath); err == nil {
		front, parseErr := parseFrontMatter(data)
		if parseErr != nil {
			report.add(ValidationIssue{Category: "Lifecycle", Artifact: e.relativePath(prdPath), Rule: "prd-front-matter", Reason: parseErr.Error(), Recovery: "repair PRD front matter"})
		} else {
			prd = front
			if state == "active" && frontMatterFieldPresent(front, "outcome") && !validOutcome(front.Outcome) {
				report.add(ValidationIssue{Category: "Lifecycle", Artifact: e.relativePath(prdPath), Rule: "active-outcome", Reason: "active Work recovery marker must declare a legal outcome", Recovery: "record succeeded, cancelled, superseded, or failed, or remove the invalid outcome"})
			}
			if state == "completed" {
				if !front.Present || !frontMatterFieldPresent(front, "outcome") || !validOutcome(front.Outcome) {
					report.add(ValidationIssue{Category: "Lifecycle", Artifact: e.relativePath(prdPath), Rule: "completed-outcome", Reason: "completed Work must declare a legal outcome", Recovery: "record succeeded, cancelled, superseded, or failed in PRD front matter"})
				}
			}
		}
	} else if !os.IsNotExist(err) {
		report.add(ValidationIssue{Category: "Lifecycle", Artifact: e.relativePath(prdPath), Rule: "prd-readable", Reason: err.Error(), Recovery: "restore PRD.md"})
	}
	if state == "completed" && prd.Outcome == "" {
		return
	}
	issuesPath := filepath.Join(workPath, "issues")
	entries, err := os.ReadDir(issuesPath)
	if err != nil {
		return
	}
	numbers := map[string]bool{}
	issueCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			report.add(ValidationIssue{Category: "Repository Structure", Artifact: e.relativePath(filepath.Join(issuesPath, entry.Name())), Rule: "issue-file", Reason: "issues/ may contain Issue files, not directories", Recovery: "flatten the Issue file into issues/"})
			continue
		}
		issueCount++
		if !issueFilenamePattern.MatchString(entry.Name()) {
			report.add(ValidationIssue{Category: "Repository Structure", Artifact: e.relativePath(filepath.Join(issuesPath, entry.Name())), Rule: "issue-filename", Reason: "Issue filename must match NN-<slug>.md", Recovery: "rename the Issue to a two-digit kebab-case filename"})
			continue
		}
		number := entry.Name()[:2]
		if numbers[number] {
			report.add(ValidationIssue{Category: "Repository Structure", Artifact: e.relativePath(filepath.Join(issuesPath, entry.Name())), Rule: "issue-number-unique", Reason: fmt.Sprintf("Issue number %s is duplicated", number), Recovery: "assign a unique two-digit Issue number"})
		}
		numbers[number] = true
		data, readErr := os.ReadFile(filepath.Join(issuesPath, entry.Name()))
		if readErr != nil {
			report.add(ValidationIssue{Category: "Repository Structure", Artifact: e.relativePath(filepath.Join(issuesPath, entry.Name())), Rule: "issue-readable", Reason: readErr.Error(), Recovery: "restore the Issue file"})
			continue
		}
		front, parseErr := parseFrontMatter(data)
		if parseErr != nil || !front.Present || !issueStatuses[front.Status] {
			reason := "Issue front matter must declare a legal status"
			if parseErr != nil {
				reason = parseErr.Error()
			}
			report.add(ValidationIssue{Category: "Lifecycle", Artifact: e.relativePath(filepath.Join(issuesPath, entry.Name())), Rule: "issue-status", Reason: reason, Recovery: "use open, in-progress, done, blocked, cancelled, or superseded"})
			continue
		}
		if state == "completed" && prd.Outcome == OutcomeSucceeded && nonTerminalIssueStatuses[front.Status] {
			report.add(ValidationIssue{Category: "Lifecycle", Artifact: e.relativePath(filepath.Join(issuesPath, entry.Name())), Rule: "succeeded-issues-terminal", Reason: fmt.Sprintf("succeeded Work retains Issue status %s", front.Status), Recovery: "resolve the Issue to done, cancelled, or superseded"})
		}
	}
	if state == "completed" && issueCount == 0 {
		report.add(ValidationIssue{Category: "Repository Structure", Artifact: e.relativePath(issuesPath), Rule: "completed-issue-required", Reason: "completed Work must retain at least one Issue", Recovery: "restore the Work's Issue definition"})
	}
}

// validateRelationships verifies explicit relationship types and resolvable targets.
func (e *Engine) validateRelationships(report *ValidationReport, knowledge map[string]knowledgeArtifact, states map[string]string) {
	for _, artifact := range knowledge {
		e.checkRelationshipSet(report, artifact.Path, artifact.Relationships, knowledge, states)
	}
	for _, state := range []string{"active", "completed"} {
		root := e.Profile.ActiveRoot
		if state == "completed" {
			root = e.Profile.CompletedRoot
		}
		for slug := range states {
			if states[slug] != state {
				continue
			}
			prdPath := e.path(filepath.ToSlash(filepath.Join(root, slug, "PRD.md")))
			data, err := os.ReadFile(prdPath)
			if err != nil {
				continue
			}
			front, parseErr := parseFrontMatter(data)
			if parseErr == nil && front.Present {
				e.checkRelationshipSet(report, prdPath, front.Relationships, knowledge, states)
			}
		}
	}
}

// checkRelationshipSet validates one artifact's relationship list.
func (e *Engine) checkRelationshipSet(report *ValidationReport, source string, links []Relationship, knowledge map[string]knowledgeArtifact, states map[string]string) {
	for _, link := range links {
		if !relationshipTypes[link.Type] {
			report.add(ValidationIssue{Category: "Relationships", Artifact: e.relativePath(source), Rule: "relationship-type", Reason: fmt.Sprintf("invalid relationship type %q", link.Type), Recovery: "use the normative relationship vocabulary"})
		}
		if link.Target == "" {
			report.add(ValidationIssue{Category: "Relationships", Artifact: e.relativePath(source), Rule: "relationship-target", Reason: "relationship target must not be empty", Recovery: "declare an identity or Work-scoped target"})
			continue
		}
		if match := knowledgeIDPattern.FindStringSubmatch(link.Target); len(match) == 3 {
			if _, ok := knowledge[link.Target]; !ok {
				report.add(ValidationIssue{Category: "Relationships", Artifact: e.relativePath(source), Rule: "target-resolvable", Reason: fmt.Sprintf("Knowledge target %s does not resolve", link.Target), Recovery: "create the target artifact or correct the identifier"})
			}
			continue
		}
		if strings.Contains(link.Target, "/") {
			parts := strings.SplitN(link.Target, "/", 2)
			if parts[0] == "active" || parts[0] == "completed" {
				report.add(ValidationIssue{Category: "Relationships", Artifact: e.relativePath(source), Rule: "work-scoped-target", Reason: "Runtime relationship targets must not include active/ or completed/", Recovery: "use <workstream-slug>/<file>"})
				continue
			}
			if _, ok := states[parts[0]]; !ok || strings.Contains(parts[1], "..") || filepath.IsAbs(filepath.FromSlash(parts[1])) {
				report.add(ValidationIssue{Category: "Relationships", Artifact: e.relativePath(source), Rule: "target-resolvable", Reason: fmt.Sprintf("Runtime target %s does not resolve", link.Target), Recovery: "use a valid Work-scoped path"})
				continue
			}
			found := false
			for _, root := range []string{e.Profile.ActiveRoot, e.Profile.CompletedRoot} {
				candidate := e.path(filepath.ToSlash(filepath.Join(root, parts[0], filepath.FromSlash(parts[1]))))
				if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
					found = true
				}
			}
			if !found {
				report.add(ValidationIssue{Category: "Relationships", Artifact: e.relativePath(source), Rule: "target-resolvable", Reason: fmt.Sprintf("Runtime target %s does not resolve", link.Target), Recovery: "restore the Work-scoped asset"})
			}
			continue
		}
		if _, ok := states[link.Target]; !ok {
			report.add(ValidationIssue{Category: "Relationships", Artifact: e.relativePath(source), Rule: "target-resolvable", Reason: fmt.Sprintf("target %s is neither a known identity nor Work slug", link.Target), Recovery: "correct the relationship target"})
		}
	}
}

// validateGeneratedIndex verifies INDEX.md is the exact reproducible derivative.
func (e *Engine) validateGeneratedIndex(report *ValidationReport) {
	path := e.path(e.Profile.IndexPath)
	actual, err := os.ReadFile(path)
	if err != nil {
		return
	}
	expected, err := e.renderIndex()
	if err != nil {
		report.add(ValidationIssue{Category: "Generated Content", Artifact: e.Profile.IndexPath, Rule: "index-reproducible", Reason: err.Error(), Recovery: "repair Work metadata and regenerate INDEX.md"})
		return
	}
	if string(actual) != string(expected) {
		report.add(ValidationIssue{Category: "Generated Content", Artifact: e.Profile.IndexPath, Rule: "index-current", Reason: "INDEX.md does not match repository-derived content", Recovery: "run `dos sync` to regenerate INDEX.md"})
	}
}

// validateIdentityAndRelationships performs the post-allocation validation gate.
func (e *Engine) validateIdentityAndRelationships() ValidationReport {
	report := ValidationReport{Status: "passed"}
	knowledge := e.validateKnowledge(&report)
	runtime := e.validateRuntime(&report)
	e.validateRelationships(&report, knowledge, runtime)
	report.sortIssues()
	return report
}

// requireDirectory reports a missing or non-directory profile path.
func (e *Engine) requireDirectory(relative string) error {
	info, err := os.Stat(e.path(relative))
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", relative)
	}
	return nil
}

// requireRegularFile reports a missing or directory profile file.
func (e *Engine) requireRegularFile(relative string) error {
	info, err := os.Stat(e.path(relative))
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("%s is a directory, expected a file", relative)
	}
	return nil
}
