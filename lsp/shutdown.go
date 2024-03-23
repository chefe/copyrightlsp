package lsp

type ShutdownRequest struct {
	Request
}

type ShutdownResponse struct {
	Response
}

func NewShudownResponse(id int) ShutdownResponse {
	return ShutdownResponse{Response: NewResponse(id)}
}
