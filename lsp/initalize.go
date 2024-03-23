package lsp

type IntitializeRequest struct {
	Params IntitializeParams `json:"params"`
	Request
}

type IntitializeParams struct {
	// Information about the client.
	ClientInfo *ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	// The client's version as defined by the client.
	Version *string `json:"version,omitempty"`
	// The name of the client as defined by the client.
	Name string `json:"name"`
}

type IntitializeResponse struct {
	Response
	Result IntitializeResult `json:"result"`
}

type IntitializeResult struct {
	// Information about the server.
	ServerInfo ServerInfo `json:"serverInfo"`
	// The capabilities the language server provides.
	Capabilities ServerCapabilities `json:"capabilities"`
}

type ServerCapabilities struct {
	// Defines how text documents are synced.
	TextDocumentSync TextDocumentSyncKind `json:"textDocumentSync"`
	// The server provides code actions.
	CodeActionProvider bool `json:"codeActionProvider"`
}

type ServerInfo struct {
	// The name of the server as defined by the server.
	Name string `json:"name"`
	// The server's version as defined by the server.
	Version string `json:"version"`
}

// Defines how the host (editor) should sync document changes to the language
// server.
type TextDocumentSyncKind int

const (
	// Documents should not be synced at all.
	TextDocumentSyncKindNone TextDocumentSyncKind = 0
	// Documents are synced by always sending the full content of the document.
	TextDocumentSyncKindFull TextDocumentSyncKind = 1
	// Documents are synced by sending the full content on open. After that only
	// incremental updates to the document are sent.
	TextDocumentSyncKindIncremental TextDocumentSyncKind = 2
)

func NewInitializeResponse(id int) IntitializeResponse {
	return IntitializeResponse{
		Response: NewResponse(id),
		Result: IntitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:   TextDocumentSyncKindFull,
				CodeActionProvider: true,
			},
			ServerInfo: ServerInfo{
				Name:    "copyrightlsp",
				Version: "0.0.0",
			},
		},
	}
}
