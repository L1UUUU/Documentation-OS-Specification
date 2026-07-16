// This file implements the supported deterministic Documentation OS migration path.
package engine

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var versionPattern = regexp.MustCompile(`(?m)(?:Documentation OS Specification|\*\*Version:\*\*)\s*:?\s*([0-9]+\.[0-9]+)`)

// Version reports specification, profile, repository, engine, and CLI versions.
func (e *Engine) Version() VersionInfo {
	return VersionInfo{
		SpecificationVersion: SpecificationVersion,
		RepositoryVersion:    e.repositoryVersion(),
		RepositoryProfile:    e.Profile.Name,
		EngineVersion:        EngineVersion,
		CLIVersion:           EngineVersion,
	}
}

// repositoryVersion derives the repository's declared version from repository guidance.
func (e *Engine) repositoryVersion() string {
	for _, relative := range []string{e.Profile.RootAgents, "README.md"} {
		data, err := os.ReadFile(e.path(relative))
		if err != nil {
			continue
		}
		match := versionPattern.FindStringSubmatch(string(data))
		if len(match) == 2 {
			return match[1]
		}
	}
	if _, err := os.Stat(e.path(e.Profile.RuntimeRoot)); err == nil {
		return SpecificationVersion
	}
	return "unknown"
}

// Migrate validates and refreshes a repository for the supported target version.
func (e *Engine) Migrate(target string) (MigrationReport, error) {
	target = normalizeVersion(target)
	if target == "" {
		target = SpecificationVersion
	}
	if target != SpecificationVersion {
		return MigrationReport{Target: target, Success: false}, fmt.Errorf("unsupported Documentation OS target version %q; supported version is %s", target, SpecificationVersion)
	}
	before := e.repositoryVersion()
	result := MigrationReport{Before: before, Target: target, After: target, Success: false}
	if err := e.requireRuntimeRoots(); err != nil {
		return result, fmt.Errorf("migration prerequisites not satisfied: %w", err)
	}
	previousIndex, indexExisted, err := readOptional(e.path(e.Profile.IndexPath))
	if err != nil {
		return result, err
	}
	index, err := e.GenerateIndex()
	if err != nil {
		return result, fmt.Errorf("migration transform failed: %w", err)
	}
	if index.Changed {
		result.Operations = append(result.Operations, "regenerated .scratch/INDEX.md")
	}
	report, err := e.Validate()
	if err != nil {
		_ = restoreFile(e.path(e.Profile.IndexPath), previousIndex, indexExisted)
		return result, err
	}
	if !report.Passed() {
		_ = restoreFile(e.path(e.Profile.IndexPath), previousIndex, indexExisted)
		return result, fmt.Errorf("migration validation failed: %s", report.String())
	}
	if before == "unknown" {
		result.Warnings = append(result.Warnings, "repository version was not explicitly declared; validated as Single Repository Profile 1.0")
	}
	result.Operations = append(result.Operations, "validated repository semantics")
	result.Success = true
	return result, nil
}

// normalizeVersion trims a user-provided migration version.
func normalizeVersion(version string) string {
	return strings.TrimSpace(version)
}
