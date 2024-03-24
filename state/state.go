package state

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chefe/copyrightlsp/lsp"
)

type documentInfo struct {
	Language string
	Content  string
}

type State struct {
	// Map file names to document infos
	Documents map[string]documentInfo
	// Map languages to template lines
	Templates map[string][]string
}

func NewState() State {
	templates := map[string][]string{}
	templates["markdown"] = []string{"<!-- Copyright (C) {year} MARKDOWN -->"}
	templates["sh"] = []string{"# Copyright (C) {year} SHELLSCRIPT"}

	return State{
		Documents: map[string]documentInfo{},
		Templates: templates,
	}
}

func (s *State) OpenDocument(document, text, language string) {
	s.Documents[document] = documentInfo{
		Language: language,
		Content:  text,
	}
}

func (s *State) UpdateDocument(document, text string) {
	doc, ok := s.Documents[document]
	if !ok {
		return
	}

	s.Documents[document] = documentInfo{
		Language: doc.Language,
		Content:  text,
	}
}

func (s *State) CloseDocument(document string) {
	delete(s.Documents, document)
}

func matchesTemplateLine(line string, template string) bool {
	pattern := fmt.Sprintf("^%s$", regexp.QuoteMeta(template))
	pattern = strings.ReplaceAll(pattern, "\\{year\\}", "[0-9]{4}")
	return regexp.MustCompile(pattern).Match([]byte(line))
}

func containsCopyrightString(lines []string, template []string) bool {
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

func buildCopyrightString(template []string) string {
	year := strconv.Itoa(time.Now().Year())
	content := strings.Join(template, "\n")
	return strings.ReplaceAll(content, "{year}", year)
}

func (s *State) CalculateCodeActions(document string, start lsp.Position, end lsp.Position) []lsp.CodeAction {
	doc, ok := s.Documents[document]
	if !ok {
		return []lsp.CodeAction{}
	}

	templateLines, ok := s.Templates[doc.Language]
	if !ok {
		return []lsp.CodeAction{}
	}

	lines := strings.Split(doc.Content, "\n")
	actions := []lsp.CodeAction{}

	limit := 10 + len(templateLines)
	if len(lines) < limit {
		limit = len(lines)
	}

	if !containsCopyrightString(lines[:limit], templateLines) {
		changes := map[string][]lsp.TextEdit{}
		changes[document] = []lsp.TextEdit{
			{
				NewText: fmt.Sprintf("%s\n", buildCopyrightString(templateLines)),
				Range:   lsp.NewRange(0, 0, 0, 0),
			},
		}

		actions = append(actions, lsp.CodeAction{
			Title: "Add copyright header",
			Edit:  lsp.WorkspaceEdit{Changes: changes},
		})
	}

	return actions
}
