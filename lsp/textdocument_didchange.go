package lsp

type DidChangeTextDocumentNotification struct {
	Notification
	Params DidChangeTextDocumentParams `json:"params"`
}

type DidChangeTextDocumentParams struct {
	// The document that did change. The version number points to the version
	// after all provided content changes have been applied.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The actual content changes. The content changes describe single state
	// changes to the document. So if there are two content changes c1 (at array
	// index 0) and c2 (at array index 1) for a document in state S then c1 moves
	// the document from S to S' and c2 from S' to S''. So c1 is computed on the
	// state S and c2 is computed on the state S'.
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type TextDocumentContentChangeEvent struct {
	// The new text of the whole document.
	Text string `json:"text"`
}
