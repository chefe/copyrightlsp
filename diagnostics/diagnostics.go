package diagnostics

import (
	"github.com/chefe/copyrightlsp/analysis"
	"github.com/chefe/copyrightlsp/lsp"
	"github.com/chefe/copyrightlsp/state"
)

func CalculateDiagnostics(state *state.State, document string) []lsp.Diagnostic {
	doc, ok := state.Documents[document]
	if !ok {
		return []lsp.Diagnostic{}
	}

	templateLines, ok := state.Templates[doc.Language]
	if !ok {
		return []lsp.Diagnostic{}
	}

	if analysis.ContainsCopyrightString(doc.Content, templateLines) {
		return []lsp.Diagnostic{}
	}

	return []lsp.Diagnostic{
		lsp.NewErrorDiagnostic("No copyright header found!"),
	}
}
