package lsp

type IntitializeRequest struct {
	Params IntitializeRequestParams `json:"params"`
	Request
}

type IntitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	Version *string `json:"version,omitempty"`
	Name    string  `json:"name"`
}

type IntitializeResponse struct {
	Response
	Result IntitializeResult `json:"result"`
}

type IntitializeResult struct {
	ServerInfo   ServerInfo         `json:"serverInfo"`
	Capabilities ServerCapabilities `json:"capabilities"`
}

type ServerCapabilities struct {
	TextDocumentSync   TextDocumentSyncKind `json:"textDocumentSync"`
	CodeActionProvider bool                 `json:"codeActionProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type TextDocumentSyncKind int

const (
	TextDocumentSyncKindNone        TextDocumentSyncKind = 0
	TextDocumentSyncKindFull        TextDocumentSyncKind = 1
	TextDocumentSyncKindIncremental TextDocumentSyncKind = 2
)

func NewInitializeResponse(id int) IntitializeResponse {
	return IntitializeResponse{
		Response: Response{
			ID:  &id,
			RPC: "2.0",
		},
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
