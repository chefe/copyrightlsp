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
}
