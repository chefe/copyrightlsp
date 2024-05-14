package analysis

import (
	"fmt"
	"regexp"
	"strings"
)

func matchesTemplateLine(line, template string) bool {
	pattern := fmt.Sprintf("^%s$", regexp.QuoteMeta(template))
	pattern = strings.ReplaceAll(pattern, "\\{year\\}", "[0-9]{4}")

	return regexp.MustCompile(pattern).MatchString(line)
}

func containsTemplateLines(lines, template []string) bool {
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

func ContainsCopyrightString(content string, templateLines []string, searchRange uint8) bool {
	limit := int(searchRange) + len(templateLines)

	// split content into n+1 lines, because the last item contains the
	// remaining part of the string, which should be ignored
	lines := strings.SplitN(content, "\n", limit+1)

	// handle case if the number of lines is smaller then the limit
	if len(lines) < limit {
		limit = len(lines)
	}

	return containsTemplateLines(lines[:limit], templateLines)
}
