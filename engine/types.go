package engine

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	// SpecificationVersion is the Documentation OS specification version implemented here.
	SpecificationVersion = "1.0"
	// SpecificationStatus is the maturity of the supported specification revision.
	SpecificationStatus = "draft"
	// SpecificationRevision identifies the exact Draft semantics implemented here.
	SpecificationRevision = "12"
	// RepositoryProfileVersion is the supported Single Repository Profile version.
	RepositoryProfileVersion = "1.0"
	// EngineVersion is the implementation version exposed by the engine and CLI.
	EngineVersion = "0.1.0-rc.3"
	// CLIVersion is the command-line contract version.
	CLIVersion = "0.1.0-rc.3"
	// ProfileName is the repository profile implemented by this engine.
	ProfileName = "Single Repository"

	// OutcomeSucceeded records a successfully completed Work.
	OutcomeSucceeded = "succeeded"
	// OutcomeCancelled records an abandoned Work.
	OutcomeCancelled = "cancelled"
	// OutcomeSuperseded records a Work replaced by another Work.
	OutcomeSuperseded = "superseded"
	// OutcomeFailed records a Work that could not be completed.
	OutcomeFailed = "failed"
)

var (
	// ErrPreflight identifies a Complete operation rejected before mutation.
	ErrPreflight = errors.New("complete preflight failed")
	// ErrInvalidRepository identifies a repository that cannot be operated on.
	ErrInvalidRepository = errors.New("invalid Documentation OS repository")
	// ErrIdentityConflict identifies an Identity assigned to distinct artifacts across worktrees.
	ErrIdentityConflict = errors.New("identity conflict")
	// ErrInvalidInput identifies a public operation input outside its contract.
	ErrInvalidInput = errors.New("invalid input")
	// ErrConflict identifies a retry that contradicts persisted repository state.
	ErrConflict = errors.New("operation conflicts with persisted state")
)

// KnowledgeImpact is the required result of Knowledge Impact Analysis.
type KnowledgeImpact string

const (
	// KnowledgeImpactChanged declares that the caller changed persistent Knowledge.
	KnowledgeImpactChanged KnowledgeImpact = "changed"
	// KnowledgeImpactNoChange declares that no persistent Knowledge edits were required.
	KnowledgeImpactNoChange KnowledgeImpact = "no-change"
)

// KnowledgeCategory describes one identity-managed Knowledge directory.
type KnowledgeCategory struct {
	Name      string
	Prefix    string
	Directory string
}

// Profile defines all repository-relative paths used by the Single Repository Profile.
type Profile struct {
	Name             string
	Version          string
	KnowledgeRoot    string
	ArchitectureDir  string
	ADRDir           string
	StandardsDir     string
	InboxDir         string
	RuntimeRoot      string
	ActiveRoot       string
	CompletedRoot    string
	LockRoot         string
	IndexPath        string
	RootAgents       string
	RootClaude       string
	ScratchAgents    string
	ScratchClaude    string
	KnowledgeEntries []KnowledgeCategory
}

// DefaultProfile returns the normative Single Repository Profile paths.
func DefaultProfile() Profile {
	return Profile{
		Name:            ProfileName,
		Version:         RepositoryProfileVersion,
		KnowledgeRoot:   "docs",
		ArchitectureDir: "docs/architecture",
		ADRDir:          "docs/adr",
		StandardsDir:    "docs/standards",
		InboxDir:        "docs/inbox",
		RuntimeRoot:     ".scratch",
		ActiveRoot:      ".scratch/active",
		CompletedRoot:   ".scratch/completed",
		LockRoot:        ".scratch/.locks",
		IndexPath:       ".scratch/INDEX.md",
		RootAgents:      "AGENTS.md",
		RootClaude:      "CLAUDE.md",
		ScratchAgents:   ".scratch/AGENTS.md",
		ScratchClaude:   ".scratch/CLAUDE.md",
		KnowledgeEntries: []KnowledgeCategory{
			{Name: "Architecture", Prefix: "ARCH", Directory: "docs/architecture"},
			{Name: "ADR", Prefix: "ADR", Directory: "docs/adr"},
			{Name: "Standards", Prefix: "STD", Directory: "docs/standards"},
		},
	}
}

// Engine executes deterministic operations against one repository root.
type Engine struct {
	Root            string
	Profile         Profile
	renameFile      func(string, string) error
	writeFileAtomic func(string, []byte, os.FileMode) error
}

// ValidationIssue describes one deterministic repository validation observation.
type ValidationIssue struct {
	Category string `json:"category"`
	Artifact string `json:"artifact"`
	Rule     string `json:"rule"`
	Severity string `json:"severity"`
	Reason   string `json:"reason"`
	Recovery string `json:"recovery,omitempty"`
}

// ValidationReport contains a read-only repository validation result.
type ValidationReport struct {
	Status string            `json:"status"`
	Issues []ValidationIssue `json:"issues"`
}

// Passed reports whether no validation errors were found.
func (r ValidationReport) Passed() bool {
	for _, issue := range r.Issues {
		if issue.Severity == "error" {
			return false
		}
	}
	return true
}

// String renders a deterministic human-readable validation report.
func (r ValidationReport) String() string {
	var builder strings.Builder
	builder.WriteString(r.Status)
	for _, issue := range r.Issues {
		fmt.Fprintf(&builder, "\n[%s] %s: %s (%s) - %s", issue.Severity, issue.Artifact, issue.Reason, issue.Rule, issue.Recovery)
	}
	return builder.String()
}

// add appends a validation issue and keeps the report status synchronized.
func (r *ValidationReport) add(issue ValidationIssue) {
	if issue.Severity == "" {
		issue.Severity = "error"
	}
	r.Issues = append(r.Issues, issue)
	if issue.Severity == "error" {
		r.Status = "failed"
	} else if r.Status == "passed" {
		r.Status = "passed-with-warnings"
	}
}

// sortIssues makes validation output reproducible across filesystem traversal orders.
func (r *ValidationReport) sortIssues() {
	sort.Slice(r.Issues, func(i, j int) bool {
		left, right := r.Issues[i], r.Issues[j]
		for _, pair := range [][2]string{{left.Category, right.Category}, {left.Artifact, right.Artifact}, {left.Rule, right.Rule}, {left.Reason, right.Reason}} {
			if pair[0] != pair[1] {
				return pair[0] < pair[1]
			}
		}
		return left.Severity < right.Severity
	})
	if r.Status == "" {
		r.Status = "passed"
	}
}

// IssueSummary describes one Issue returned by repository inspection.
type IssueSummary struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// WorkSummary describes one Work for inspection and health reports.
type WorkSummary struct {
	Slug        string         `json:"slug"`
	State       string         `json:"state"`
	PRDPath     string         `json:"prd_path"`
	HANDOFFPath string         `json:"handoff_path"`
	IssueCount  int            `json:"issue_count"`
	Issues      []IssueSummary `json:"issues"`
	Outcome     string         `json:"outcome,omitempty"`
}

// InspectReport summarizes repository state without modifying it.
type InspectReport struct {
	SpecificationVersion string        `json:"specification_version"`
	EngineVersion        string        `json:"engine_version"`
	RepositoryProfile    string        `json:"repository_profile"`
	Root                 string        `json:"root"`
	KnowledgeDocuments   int           `json:"knowledge_documents"`
	ActiveWorks          int           `json:"active_works"`
	CompletedWorks       int           `json:"completed_works"`
	Works                []WorkSummary `json:"works"`
}

// VersionInfo exposes the independent Documentation OS version dimensions.
type VersionInfo struct {
	SpecificationVersion     string `json:"specification_version"`
	SpecificationStatus      string `json:"specification_status"`
	SpecificationRevision    string `json:"specification_revision"`
	RepositoryVersion        string `json:"repository_version"`
	RepositoryProfile        string `json:"repository_profile"`
	RepositoryProfileVersion string `json:"repository_profile_version"`
	EngineVersion            string `json:"engine_version"`
	CLIVersion               string `json:"cli_version"`
}

// WorkResult describes a generated Work workspace.
type WorkResult struct {
	Slug      string `json:"slug"`
	Path      string `json:"path"`
	IndexPath string `json:"index_path"`
}

// BeginInput describes the caller-owned identity and initial plan for a new Work.
type BeginInput struct {
	Slug   string       `json:"slug"`
	Title  string       `json:"title"`
	Issues []BeginIssue `json:"issues"`
}

// BeginIssue describes one initial Issue created with a Work.
type BeginIssue struct {
	Slug   string `json:"slug"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// IndexResult describes generated derived Runtime metadata.
type IndexResult struct {
	Path    string `json:"path"`
	Changed bool   `json:"changed"`
}

// CompleteResult describes the two stages of a Complete operation.
type CompleteResult struct {
	Slug             string `json:"slug"`
	Completed        bool   `json:"completed"`
	CleanupCompleted bool   `json:"cleanup_completed"`
	RetriedCleanup   bool   `json:"retried_cleanup"`
	Outcome          string `json:"outcome"`
}

// AllocationResult describes a final Knowledge identifier allocation.
type AllocationResult struct {
	Identifier string `json:"identifier"`
	OldPath    string `json:"old_path"`
	NewPath    string `json:"new_path"`
	References int    `json:"references_updated"`
}

// SyncInput carries the caller's Knowledge Impact Analysis declaration.
type SyncInput struct {
	KnowledgeImpact KnowledgeImpact `json:"knowledge_impact"`
}

// SyncResult makes the caller's Knowledge synchronization declaration observable.
type SyncResult struct {
	KnowledgeImpact   KnowledgeImpact `json:"knowledge_impact"`
	NoKnowledgeChange bool            `json:"no_knowledge_change"`
	Index             IndexResult     `json:"index"`
}

// HealthCategory is one advisory Health assessment category.
type HealthCategory struct {
	Name            string   `json:"name"`
	Level           string   `json:"level"`
	Score           int      `json:"score"`
	Observations    []string `json:"observations,omitempty"`
	Recommendations []string `json:"recommendations,omitempty"`
}

// HealthReport contains non-blocking repository sustainability observations.
type HealthReport struct {
	Level           string           `json:"level"`
	Score           int              `json:"score"`
	Categories      []HealthCategory `json:"categories"`
	Metrics         map[string]int   `json:"metrics"`
	Recommendations []string         `json:"recommendations,omitempty"`
}

// MigrationReport describes a deterministic migration attempt.
type MigrationReport struct {
	Before     string   `json:"before"`
	After      string   `json:"after"`
	Target     string   `json:"target"`
	Success    bool     `json:"success"`
	Operations []string `json:"operations"`
	Warnings   []string `json:"warnings,omitempty"`
}

// validOutcome reports whether a caller supplied a normative terminal outcome.
func validOutcome(outcome string) bool {
	switch outcome {
	case OutcomeSucceeded, OutcomeCancelled, OutcomeSuperseded, OutcomeFailed:
		return true
	default:
		return false
	}
}
