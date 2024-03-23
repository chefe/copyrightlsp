package lsp

type DidCloseTextDocumentNotification struct {
	Notification
	Params DidCloseTextDocumentParams `json:"params"`
}

type DidCloseTextDocumentParams struct {
	// The document that was closed.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}
