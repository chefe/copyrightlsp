package diagnostics

import (
	"github.com/chefe/copyrightlsp/analysis"
	"github.com/chefe/copyrightlsp/lsp"
	"github.com/chefe/copyrightlsp/state"
)

// CalculateDiagnostics returns all the diagnostics for the given document.
func CalculateDiagnostics(lspState *state.State, document string) []lsp.Diagnostic {
	doc, found := lspState.Documents[document]
	if !found {
		return []lsp.Diagnostic{}
	}

	templateLines, found := lspState.Templates[doc.Language]
	if !found {
		return []lsp.Diagnostic{}
	}

	searchRange := lspState.GetSearchRange(doc.Language)
	if analysis.ContainsCopyrightString(doc.Content, templateLines, searchRange) {
		return []lsp.Diagnostic{}
	}

	return []lsp.Diagnostic{
		lsp.NewErrorDiagnostic("No copyright header found!"),
	}
}
