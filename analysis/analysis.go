package analysis

import (
	"fmt"
	"regexp"
	"strings"
)

func matchesTemplateLine(line string, template string) bool {
	pattern := fmt.Sprintf("^%s$", regexp.QuoteMeta(template))
	pattern = strings.ReplaceAll(pattern, "\\{year\\}", "[0-9]{4}")
	return regexp.MustCompile(pattern).Match([]byte(line))
}

func containsTemplateLines(lines []string, template []string) bool {
	if len(template) == 0 {
		return false
	}

	for lineIndex, line := range lines {
		if !matchesTemplateLine(line, template[0]) {
			continue
		}

		if lineIndex+len(template) > len(lines) {
			// not enought lines left to match template
			return false
		}

		for templateIndex, templateLine := range template[1:] {
			if !matchesTemplateLine(lines[lineIndex+templateIndex+1], templateLine) {
				return false
			}
		}

		return true
	}

	return false
}

func ContainsCopyrightString(content string, templateLines []string) bool {
	lines := strings.Split(content, "\n")

	limit := 10 + len(templateLines)
	if len(lines) < limit {
		limit = len(lines)
	}

	return containsTemplateLines(lines[:limit], templateLines)
}
