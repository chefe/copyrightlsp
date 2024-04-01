package diagnostics

import (
	"github.com/chefe/copyrightlsp/analysis"
	"github.com/chefe/copyrightlsp/lsp"
	"github.com/chefe/copyrightlsp/state"
)

func CalculateDiagnostics(state *state.State, document string) []lsp.Diagnostic {
	doc, found := state.Documents[document]
	if !found {
		return []lsp.Diagnostic{}
	}

	templateLines, found := state.Templates[doc.Language]
	if !found {
		return []lsp.Diagnostic{}
	}

	searchRange := state.GetSearchRange(doc.Language)
	if analysis.ContainsCopyrightString(doc.Content, templateLines, searchRange) {
		return []lsp.Diagnostic{}
	}

	return []lsp.Diagnostic{
		lsp.NewErrorDiagnostic("No copyright header found!"),
	}
}
