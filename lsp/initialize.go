package lsp

type InitializeRequest struct {
	Params InitializeRequestParams `json:"params"`
	Request
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Result InitializeResult `json:"result"`
	Response
}

type InitializeResult struct {
	ServerInfo   *ServerInfo        `json:"serverInfo"`
	Capabilities ServerCapabilities `json:"capabilities"`
}

type ServerCapabilities struct {
	TextDocumentSync int `json:"textDocumentSync"`
}

type ServerInfo struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			ID: &id,
			Message: Message{
				RPC: "2.0",
			},
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: 1,
			},
			ServerInfo: &ServerInfo{
				Version: "0.0.0",
				Name:    "Test LSP",
			},
		},
	}
}
