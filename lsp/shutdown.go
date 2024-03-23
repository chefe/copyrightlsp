package lsp

type ShutdownRequest struct {
	Request
}

type ShutdownResponse struct {
	// Result *bool `json:"result"`
	Response
}

func NewShudownResponse(id int) ShutdownResponse {
	return ShutdownResponse{
		Response: Response{
			ID:  &id,
			RPC: "2.0",
		},
	}
}
