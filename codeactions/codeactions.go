package codeactions

import (
	"strconv"
	"strings"
	"time"

	"github.com/chefe/copyrightlsp/analysis"
	"github.com/chefe/copyrightlsp/lsp"
	"github.com/chefe/copyrightlsp/state"
)

func buildCopyrightString(template []string) string {
	year := strconv.Itoa(time.Now().Year())
	content := strings.Join(template, "\n")

	return strings.ReplaceAll(content, "{year}", year)
}

func CalculateCodeActions(lspState *state.State, document string, start lsp.Position, end lsp.Position) []lsp.CodeAction {
	// show code action only on first line of a document
	if start.Line > 0 || end.Line > 0 {
		return []lsp.CodeAction{}
	}

	doc, found := lspState.Documents[document]
	if !found {
		return []lsp.CodeAction{}
	}

	templateLines, found := lspState.Templates[doc.Language]
	if !found {
		return []lsp.CodeAction{}
	}

	searchRange := lspState.GetSearchRange(doc.Language)
	if analysis.ContainsCopyrightString(doc.Content, templateLines, searchRange) {
		return []lsp.CodeAction{}
	}

	changes := map[string][]lsp.TextEdit{}
	changes[document] = []lsp.TextEdit{
		{
			NewText: buildCopyrightString(templateLines) + "\n",
			Range:   lsp.NewRange(0, 0, 0, 0),
		},
	}

	return []lsp.CodeAction{
		{
			Title: "Add copyright header",
			Edit:  lsp.WorkspaceEdit{Changes: changes},
		},
	}
}
