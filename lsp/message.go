package lsp

type Message struct {
	RPC string `json:"jsonrpc"`
}

type Request struct {
	Message
	Method string `json:"method"`
	ID     int    `json:"id"`
}

type Response struct {
	ID *int `json:"id,omitempty"`
	Message
}

type Notification struct {
	Message
	Method string `json:"method"`
}
