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

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

func NewRange(startLine, startCharacter, endLine, endCharacter int) Range {
	return Range{
		Start: Position{Line: startLine, Character: startCharacter},
		End:   Position{Line: endLine, Character: endCharacter},
	}
}

type TextEdit struct {
	// The string to be inserted. For delete operations use an empty string.
	NewText string `json:"newText"`

	// The range of the text document to be manipulated. To insert text into a
	// document create a range where start === end.
	Range Range `json:"range"`
}

type WorkspaceEdit struct {
	Changes map[string][]TextEdit `json:"changes"`
}
