// This file evaluates advisory repository sustainability without blocking operations.
package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Health evaluates the six normative Health categories without modifying the repository.
func (e *Engine) Health() (HealthReport, error) {
	report, err := e.Validate()
	if err != nil {
		return HealthReport{}, err
	}
	works, err := e.listWorks()
	if err != nil {
		return HealthReport{}, err
	}
	knowledge, err := e.countKnowledgeDocuments()
	if err != nil {
		return HealthReport{}, err
	}
	active := countState(works, "active")
	completed := countState(works, "completed")
	validationErrors := countIssues(report, "error")
	relationshipErrors := countCategoryIssues(report, "Relationships")
	structureErrors := countCategoryIssues(report, "Repository Structure")
	categoryCounts := e.knowledgeCategoryCounts()
	coveredCategories := 0
	for _, count := range categoryCounts {
		if count > 0 {
			coveredCategories++
		}
	}

	categories := []HealthCategory{
		healthKnowledge(knowledge, categoryCounts),
		healthRuntime(active, completed),
		healthLifecycle(validationErrors),
		healthRelationships(relationshipErrors),
		healthCoverage(coveredCategories, len(categoryCounts)),
		healthRepository(structureErrors),
	}
	sort.Slice(categories, func(i, j int) bool { return categories[i].Name < categories[j].Name })
	total := 0
	for _, category := range categories {
		total += category.Score
	}
	score := 0
	if len(categories) > 0 {
		score = total / len(categories)
	}
	result := HealthReport{
		Level:      healthLevel(score),
		Score:      score,
		Categories: categories,
		Metrics: map[string]int{
			"knowledge_documents": knowledge,
			"active_works":        active,
			"completed_works":     completed,
			"validation_errors":   validationErrors,
		},
	}
	result.Recommendations = healthRecommendations(categories)
	return result, nil
}

// knowledgeCategoryCounts counts documents in each normative Knowledge category.
func (e *Engine) knowledgeCategoryCounts() map[string]int {
	counts := map[string]int{}
	for _, category := range e.Profile.KnowledgeEntries {
		count := 0
		entries, err := e.listMarkdown(category.Directory)
		if err == nil {
			count = len(entries)
		}
		counts[category.Name] = count
	}
	return counts
}

// listMarkdown lists category Markdown files without modifying the repository.
func (e *Engine) listMarkdown(relative string) ([]string, error) {
	entries, err := e.walkFiles(relative)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(entries))
	for _, path := range entries {
		base := filepath.Base(path)
		if strings.HasSuffix(path, ".md") && base != "AGENTS.md" && base != "CLAUDE.md" {
			result = append(result, path)
		}
	}
	return result, nil
}

// walkFiles returns regular files beneath a profile-relative directory.
func (e *Engine) walkFiles(relative string) ([]string, error) {
	var paths []string
	err := filepath.Walk(e.path(relative), func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	sort.Strings(paths)
	return paths, err
}

// healthKnowledge builds the Knowledge advisory category.
func healthKnowledge(total int, counts map[string]int) HealthCategory {
	score := 100
	observations := []string{}
	for name, count := range counts {
		if count == 0 {
			score -= 15
			observations = append(observations, fmt.Sprintf("Knowledge category %s has no documents", name))
		}
	}
	if total == 0 {
		score -= 25
		observations = append(observations, "Knowledge domain is empty")
	}
	sort.Strings(observations)
	return HealthCategory{Name: "Knowledge", Level: healthLevel(score), Score: clampScore(score), Observations: observations}
}

// healthRuntime builds the Runtime advisory category.
func healthRuntime(active, completed int) HealthCategory {
	score := 100 - active*10
	observations := []string{}
	if active > 0 {
		observations = append(observations, fmt.Sprintf("%d Work(s) remain active", active))
	}
	if active > 5 {
		observations = append(observations, "active Runtime accumulation is high")
		score -= 15
	}
	return HealthCategory{Name: "Runtime", Level: healthLevel(score), Score: clampScore(score), Observations: observations}
}

// healthLifecycle builds the lifecycle advisory category.
func healthLifecycle(errors int) HealthCategory {
	score := 100 - errors*15
	observations := []string{}
	if errors > 0 {
		observations = append(observations, fmt.Sprintf("Validation reports %d structural error(s)", errors))
	}
	return HealthCategory{Name: "Lifecycle", Level: healthLevel(score), Score: clampScore(score), Observations: observations}
}

// healthRelationships builds the relationship advisory category.
func healthRelationships(errors int) HealthCategory {
	score := 100 - errors*20
	observations := []string{}
	if errors > 0 {
		observations = append(observations, fmt.Sprintf("%d relationship error(s) need attention", errors))
	}
	return HealthCategory{Name: "Relationships", Level: healthLevel(score), Score: clampScore(score), Observations: observations}
}

// healthCoverage builds the Knowledge coverage advisory category.
func healthCoverage(covered, total int) HealthCategory {
	score := 100
	if total > 0 {
		score = covered * 100 / total
	}
	observations := []string{}
	if covered < total {
		observations = append(observations, fmt.Sprintf("%d of %d Knowledge categories contain documents", covered, total))
	}
	return HealthCategory{Name: "Coverage", Level: healthLevel(score), Score: clampScore(score), Observations: observations}
}

// healthRepository builds the repository organization advisory category.
func healthRepository(errors int) HealthCategory {
	score := 100 - errors*20
	observations := []string{}
	if errors > 0 {
		observations = append(observations, fmt.Sprintf("%d repository structure issue(s) need attention", errors))
	}
	return HealthCategory{Name: "Repository", Level: healthLevel(score), Score: clampScore(score), Observations: observations}
}

// countIssues counts validation issues at one severity.
func countIssues(report ValidationReport, severity string) int {
	count := 0
	for _, issue := range report.Issues {
		if issue.Severity == severity {
			count++
		}
	}
	return count
}

// countCategoryIssues counts validation issues in one category.
func countCategoryIssues(report ValidationReport, category string) int {
	count := 0
	for _, issue := range report.Issues {
		if issue.Category == category {
			count++
		}
	}
	return count
}

// clampScore bounds advisory scores to the documented 0–100 range.
func clampScore(score int) int {
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

// healthLevel maps a numeric score to the normative Health level vocabulary.
func healthLevel(score int) string {
	switch {
	case score >= 90:
		return "Excellent"
	case score >= 75:
		return "Good"
	case score >= 60:
		return "Fair"
	case score >= 40:
		return "Poor"
	default:
		return "Critical"
	}
}

// healthRecommendations creates stable advisory next steps.
func healthRecommendations(categories []HealthCategory) []string {
	var recommendations []string
	for _, category := range categories {
		if category.Score < 75 {
			recommendations = append(recommendations, fmt.Sprintf("Review %s Health observations", category.Name))
		}
	}
	sort.Strings(recommendations)
	return recommendations
}
