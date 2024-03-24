package state

import (
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
}

func NewState() State {
	return State{
		Documents: map[string]documentInfo{},
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

func buildCopyrightLine(template string) string {
	year := strconv.Itoa(time.Now().Year())
	return strings.ReplaceAll(template, "{year}", year)
}

func (s *State) CalculateCodeActions(document string, start lsp.Position, end lsp.Position) []lsp.CodeAction {
	doc, ok := s.Documents[document]
	if !ok {
		return []lsp.CodeAction{}
	}

	lines := strings.Split(doc.Content, "\n")
	actions := []lsp.CodeAction{}

	limit := 10
	if len(lines) < limit {
		limit = len(lines)
	}

	if !strings.Contains(strings.Join(lines[:limit], "\n"), "Copyright") {
		changes := map[string][]lsp.TextEdit{}
		changes[document] = []lsp.TextEdit{
			{
				NewText: buildCopyrightLine("// Copyright (C) {year} AUTHOR\n"),
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
