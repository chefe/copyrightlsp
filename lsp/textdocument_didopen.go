package lsp

type DidOpenTextDocumentNotification struct {
	Notification
	Params DidOpenTextDocumentParams `json:"params"`
}

type DidOpenTextDocumentParams struct {
	// The document that was opened.
	TextDocument TextDocumentItem `json:"textDocument"`
}
