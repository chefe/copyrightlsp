package lsp

type TextDocumentItem struct {
	// The text document's URI.
	URI string `json:"uri"`
	//  The text document's language identifier.
	LanguageID string `json:"languageId"`
	// The version number of this document (it will increase after each change,
	// including undo/redo).
	Text string `json:"text"`
	// The content of the opened text document.
	Version int `json:"version"`
}

type TextDocumentIdentifier struct {
	// The text document's URI.
	URI string `json:"uri"`
}

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	// The version number of this document.
	Version int `json:"version"`
}

type Position struct {
	// Line position in a document (zero-based).
	Line uint `json:"line"`
	// Character offset on a line in a document (zero-based).
	Character uint `json:"character"`
}

type Range struct {
	// The range's start position.
	Start Position `json:"start"`
	// The range's end position.
	End Position `json:"end"`
}

type TextEdit struct {
	// The string to be inserted. For delete operations use an empty string.
	NewText string `json:"newText"`
	// The range of the text document to be manipulated. To insert text into a
	// document create a range where start === end.
	Range Range `json:"range"`
}

type WorkspaceEdit struct {
	// Holds changes to existing resources.
	Changes map[string][]TextEdit `json:"changes"`
}

func NewRange(startLine, startCharacter, endLine, endCharacter uint) Range {
	return Range{
		Start: Position{Line: startLine, Character: startCharacter},
		End:   Position{Line: endLine, Character: endCharacter},
	}
}
