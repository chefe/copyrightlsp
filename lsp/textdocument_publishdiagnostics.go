package lsp

type PublishDiagnosticsNotification struct {
	Notification
	Params PublishDiagnosticsParams `json:"params"`
}

type PublishDiagnosticsParams struct {
	// The URI for which diagnostic information is reported.
	URI string `json:"uri"`
	// An array of diagnostic information items.
	Diagnostics []Diagnostic `json:"diagnostics"`
}

func NewPublishDiagnosticsNotification(uri string, diagnostics []Diagnostic) PublishDiagnosticsNotification {
	return PublishDiagnosticsNotification{
		Notification: NewNotification("textDocument/publishDiagnostics"),
		Params: PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: diagnostics,
		},
	}
}
