package ws_api

type GetMessageFromClientPayload struct {
	Text string `json:"text"`
}

type SendMessageToClientPayload struct {
	Text      string `json:"text"`
	AuthorId  string `json:"author_id"`
	Nickname  string `json:"nickname"`
	CreatedAt string `json:"created_at"`
}
