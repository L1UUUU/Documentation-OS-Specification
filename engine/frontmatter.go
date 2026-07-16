// This file parses and updates YAML front matter while preserving Markdown bodies.
package engine

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var topLevelFieldPattern = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_-]*)\s*:`)

// Relationship is the front matter representation of a Documentation OS relationship.
type Relationship struct {
	Type   string `yaml:"type"`
	Target string `yaml:"target"`
}

// FrontMatter contains the fields used by the Single Repository Profile.
type FrontMatter struct {
	Values        map[string]any
	ID            string
	Status        string
	Title         string
	Outcome       string
	Relationships []Relationship
	Present       bool
}

// parseFrontMatter parses a Markdown document's leading YAML block.
func parseFrontMatter(data []byte) (FrontMatter, error) {
	result := FrontMatter{Values: map[string]any{}}
	if bytes.HasPrefix(data, []byte{0xEF, 0xBB, 0xBF}) {
		return result, fmt.Errorf("UTF-8 BOM is not allowed")
	}
	if !bytes.HasPrefix(data, []byte("---")) {
		return result, nil
	}
	text := string(data)
	firstEnd := strings.IndexByte(text, '\n')
	if firstEnd < 0 || strings.TrimSpace(text[:firstEnd]) != "---" {
		return result, nil
	}
	closingStart := -1
	searchAt := firstEnd + 1
	for searchAt <= len(text) {
		next := strings.IndexByte(text[searchAt:], '\n')
		lineEnd := len(text)
		if next >= 0 {
			lineEnd = searchAt + next
		}
		if strings.TrimSpace(text[searchAt:lineEnd]) == "---" {
			closingStart = searchAt
			break
		}
		if next < 0 {
			break
		}
		searchAt = lineEnd + 1
	}
	if closingStart < 0 {
		return result, fmt.Errorf("front matter starts with --- but has no closing ---")
	}
	raw := text[firstEnd+1 : closingStart]
	if err := detectDuplicateTopLevelFields(raw); err != nil {
		return result, err
	}
	values, err := parseSimpleYAML(raw)
	if err != nil {
		return result, err
	}
	result.Values = values
	result.Present = true
	if err := decodeFrontMatterFields(result.Values, &result); err != nil {
		return FrontMatter{}, err
	}
	return result, nil
}

// decodeFrontMatterFields converts supported YAML fields to stable Go values.
func decodeFrontMatterFields(values map[string]any, result *FrontMatter) error {
	for key, target := range map[string]*string{"id": &result.ID, "status": &result.Status, "title": &result.Title, "outcome": &result.Outcome} {
		value, ok := values[key]
		if !ok {
			continue
		}
		text, ok := value.(string)
		if !ok {
			return fmt.Errorf("front matter field %q must be a string", key)
		}
		*target = text
	}
	if value, ok := values["relationships"]; ok {
		items, ok := value.([]any)
		if !ok {
			return fmt.Errorf("front matter relationships must be a list of type/target objects")
		}
		for _, item := range items {
			fields, ok := item.(map[string]any)
			if !ok {
				return fmt.Errorf("front matter relationships must contain objects")
			}
			typeValue, typeOK := fields["type"].(string)
			targetValue, targetOK := fields["target"].(string)
			if !typeOK || !targetOK {
				return fmt.Errorf("front matter relationship objects require string type and target")
			}
			result.Relationships = append(result.Relationships, Relationship{Type: typeValue, Target: targetValue})
		}
	}
	return nil
}

// parseSimpleYAML parses the small YAML subset used by Documentation OS front matter.
func parseSimpleYAML(raw string) (map[string]any, error) {
	values := map[string]any{}
	lines := strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n")
	for index := 0; index < len(lines); index++ {
		line := strings.TrimSuffix(lines[index], "\r")
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			return nil, fmt.Errorf("unexpected indented front matter line %q", line)
		}
		separator := strings.IndexByte(line, ':')
		if separator <= 0 {
			return nil, fmt.Errorf("front matter line %q must be key: value", line)
		}
		key := strings.TrimSpace(line[:separator])
		value := strings.TrimSpace(line[separator+1:])
		if key == "relationships" && value == "" {
			items, next, err := parseRelationshipList(lines, index+1)
			if err != nil {
				return nil, err
			}
			values[key] = items
			index = next - 1
			continue
		}
		if key == "relationships" && value == "[]" {
			values[key] = []any{}
			continue
		}
		values[key] = parseYAMLScalar(value)
	}
	return values, nil
}

// parseRelationshipList parses an indented list of relationship maps.
func parseRelationshipList(lines []string, start int) ([]any, int, error) {
	items := []any{}
	index := start
	for index < len(lines) {
		line := strings.TrimSuffix(lines[index], "\r")
		if strings.TrimSpace(line) == "" {
			index++
			continue
		}
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "-") {
			if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
				return nil, 0, fmt.Errorf("relationship list item %q must begin with -", line)
			}
			break
		}
		item := map[string]any{}
		first := strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
		if first != "" {
			if err := parseRelationshipField(first, item); err != nil {
				return nil, 0, err
			}
		}
		index++
		for index < len(lines) {
			child := strings.TrimSuffix(lines[index], "\r")
			if strings.TrimSpace(child) == "" {
				index++
				continue
			}
			childTrimmed := strings.TrimSpace(child)
			if strings.HasPrefix(childTrimmed, "-") || (!strings.HasPrefix(child, " ") && !strings.HasPrefix(child, "\t")) {
				break
			}
			if err := parseRelationshipField(childTrimmed, item); err != nil {
				return nil, 0, err
			}
			index++
		}
		items = append(items, item)
	}
	return items, index, nil
}

// parseRelationshipField parses one type/target field inside a relationship item.
func parseRelationshipField(line string, item map[string]any) error {
	separator := strings.IndexByte(line, ':')
	if separator <= 0 {
		return fmt.Errorf("relationship field %q must be key: value", line)
	}
	key := strings.TrimSpace(line[:separator])
	if key != "type" && key != "target" {
		return fmt.Errorf("unsupported relationship field %q", key)
	}
	item[key] = parseYAMLScalar(strings.TrimSpace(line[separator+1:]))
	return nil
}

// parseYAMLScalar decodes quoted, boolean, numeric, and plain YAML scalar values.
func parseYAMLScalar(value string) any {
	value = strings.TrimSpace(strings.SplitN(value, " #", 2)[0])
	if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
		if value[0] == '"' {
			if decoded, err := strconv.Unquote(value); err == nil {
				return decoded
			}
		}
		return value[1 : len(value)-1]
	}
	if value == "null" || value == "~" {
		return nil
	}
	if number, err := strconv.Atoi(value); err == nil {
		return number
	}
	return value
}

// detectDuplicateTopLevelFields rejects ambiguous scalar front matter keys.
func detectDuplicateTopLevelFields(raw string) error {
	seen := map[string]bool{}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSuffix(line, "\r")
		match := topLevelFieldPattern.FindStringSubmatch(line)
		if len(match) != 2 || strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			continue
		}
		if seen[match[1]] {
			return fmt.Errorf("duplicate front matter field %q", match[1])
		}
		seen[match[1]] = true
	}
	return nil
}

// frontMatterFieldPresent reports whether a top-level field was declared.
func frontMatterFieldPresent(front FrontMatter, key string) bool {
	_, ok := front.Values[key]
	return ok
}

// setFrontMatterField updates a scalar field and preserves the document body.
func setFrontMatterField(data []byte, key, value string) ([]byte, error) {
	front, err := parseFrontMatter(data)
	if err != nil {
		return nil, err
	}
	if !front.Present {
		return nil, fmt.Errorf("cannot set front matter field %q without a YAML block", key)
	}
	text := string(data)
	firstEnd := strings.IndexByte(text, '\n')
	closingStart, _, err := frontMatterBounds(text, firstEnd)
	if err != nil {
		return nil, err
	}
	newline := "\n"
	if strings.Contains(text[:closingStart], "\r\n") {
		newline = "\r\n"
	}
	lines := strings.SplitAfter(text[firstEnd+1:closingStart], "\n")
	fieldPattern := regexp.MustCompile(`^` + regexp.QuoteMeta(key) + `\s*:`)
	found := false
	for i, line := range lines {
		trimmed := strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
		if fieldPattern.MatchString(trimmed) {
			ending := ""
			if strings.HasSuffix(line, "\r\n") {
				ending = "\r\n"
			} else if strings.HasSuffix(line, "\n") {
				ending = "\n"
			}
			lines[i] = key + ": " + value + ending
			found = true
			break
		}
	}
	updatedFront := strings.Join(lines, "")
	if !found {
		if updatedFront != "" && !strings.HasSuffix(updatedFront, newline) {
			updatedFront += newline
		}
		updatedFront += key + ": " + value + newline
	}
	return []byte(text[:firstEnd+1] + updatedFront + text[closingStart:]), nil
}

// frontMatterBounds returns the opening and closing offsets of a YAML block.
func frontMatterBounds(text string, firstEnd int) (int, int, error) {
	if firstEnd < 0 {
		return 0, 0, fmt.Errorf("front matter has no first line")
	}
	searchAt := firstEnd + 1
	for searchAt <= len(text) {
		next := strings.IndexByte(text[searchAt:], '\n')
		lineEnd := len(text)
		if next >= 0 {
			lineEnd = searchAt + next
		}
		if strings.TrimSpace(text[searchAt:lineEnd]) == "---" {
			return searchAt, lineEnd, nil
		}
		if next < 0 {
			break
		}
		searchAt = lineEnd + 1
	}
	return 0, 0, fmt.Errorf("front matter has no closing ---")
}
