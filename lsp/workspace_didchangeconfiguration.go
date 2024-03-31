package lsp

type DidChangeConfigurationNotification struct {
	Params DidChangeConfigurationParams `json:"params"`
	Notification
}

type DidChangeConfigurationParams struct {
	// The actual changed settings.
	Settings CopyrightLspSettings `json:"settings"`
}

type CopyrightLspSettings struct {
	// Mappings from language to template lines
	Templates map[string][]string `json:"templates"`
	// Mappings from language to a search range. A search starts allways at the
	// top of a document and includes the amount of lines in the template. With
	// this option the search range can be increased to includes `n` additional
	// lines, which allow copyright comments on other lines then the first. If
	// no value is specified for a language then `0` is used.
	SearchRanges map[string]uint8 `json:"searchRanges"`
}
