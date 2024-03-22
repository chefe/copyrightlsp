package lsp

type TextDocumentItem struct {
	// The text document's URI.
	URI string `json:"uri"`

	//  The text document's language identifier.
	LanguageID string `json:"languageId"`

	// The version number of this document
	// (it will increase after each change, including undo/redo).
	Text string `json:"text"`

	// The content of the opened text document.
	Version int `json:"version"`
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}
