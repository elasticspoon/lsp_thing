package lsp

import "context"

type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type Position struct {
	/**
	 * Line position in a document (zero-based).
	 */
	Line int `json:"line"`

	/**
	 * Character offset on a line in a document (zero-based). The meaning of this
	 * offset is determined by the negotiated `PositionEncodingKind`.
	 *
	 * If the character value is greater than the line length it defaults back
	 * to the line length.
	 */
	Character int `json:"character"`
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

type HoverRequest struct {
	Request
	Params HoverParams `json:"params"`
}

type HoverParams struct {
	TextDocumentPositionParams
}

type HoverResponse struct {
	Response
	Result HoverResult `json:"result"`
}

type HoverResult struct {
	Contents MarkupContent `json:"contents"`
}

type MarkupContent struct {
	// options are "markdown" or "plaintext"
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type HoverResponseFunc func(context.Context, *HoverParams) (*HoverResponse, error)

func NewHoverResponse(id int) HoverResponse {
	return HoverResponse{
		Response: Response{
			ID: &id,
			Message: Message{
				RPC: "2.0",
			},
		},
		Result: HoverResult{
			Contents: MarkupContent{
				Kind:  "plaintext",
				Value: "some random response",
			},
		},
	}
}
