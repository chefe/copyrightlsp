package lsp

type CodeActionRequest struct {
	Request
	Params CodeActionParams `json:"params"`
}

type CodeActionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
}

type CodeActionResponse struct {
	Response
	Result []CodeAction `json:"result"`
}

type CodeAction struct {
	// The workspace edit this code action performs.
	Edit WorkspaceEdit `json:"edit"`
	// A short, human-readable, title for this code action.
	Title string `json:"title"`
}

func NewCodeActionResponse(id int, actions []CodeAction) CodeActionResponse {
	return CodeActionResponse{
		Response: Response{ID: &id, RPC: "2.0"},
		Result:   actions,
	}
}
