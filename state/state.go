package state

import (
	"strconv"
	"strings"
	"time"

	"github.com/chefe/copyrightlsp/lsp"
)

type State struct {
	// Map file names to contents
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: map[string]string{},
	}
}

func (s *State) OpenDocument(document, text string) {
	s.Documents[document] = text
}

func (s *State) UpdateDocument(document, text string) {
	s.Documents[document] = text
}

func (s *State) CloseDocument(document string) {
	delete(s.Documents, document)
}

func buildCopyrightLine(template string) string {
	year := strconv.Itoa(time.Now().Year())
	return strings.ReplaceAll(template, "{year}", year)
}

func (s *State) CalculateCodeActions(document string, start lsp.Position, end lsp.Position) []lsp.CodeAction {
	text, ok := s.Documents[document]
	if !ok {
		return []lsp.CodeAction{}
	}

	lines := strings.Split(text, "\n")
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
