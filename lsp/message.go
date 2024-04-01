package lsp

// The language server protocol always uses "2.0" as the jsonrpc version.
const jsonRPCVersion = "2.0"

type Request struct {
	// The version of "jsonrpc" to use.
	RPC string `json:"jsonrpc"`
	//  The method to be invoked.
	Method string `json:"method"`
	// The request id.
	ID int `json:"id"`
}

type Response struct {
	// The request id.
	ID *int `json:"id,omitempty"`
	// The version of "jsonrpc" to use.
	RPC string `json:"jsonrpc"`
}

type Notification struct {
	// The version of "jsonrpc" to use.
	RPC string `json:"jsonrpc"`
	//  The method to be invoked.
	Method string `json:"method"`
}

func NewResponse(id int) Response {
	return Response{ID: &id, RPC: jsonRPCVersion}
}

func NewNotification(method string) Notification {
	return Notification{RPC: jsonRPCVersion, Method: method}
}
