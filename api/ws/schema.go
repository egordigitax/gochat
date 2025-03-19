package ws_api

type ActionType string

type GetMessageFromClientRequest struct {
	ActionType ActionType `json:"action_type"`
	Text       string     `json:"text"`
}

type SendMessageToClientResponse struct {
	ActionType ActionType `json:"action_type"`
	Text       string     `json:"text"`
	AuthorId   string     `json:"author_id"`
	Nickname   string     `json:"nickname"`
	CreatedAt  string     `json:"created_at"`
}

type Chat struct {
	Title       string `json:"title"`
	UnreadCount int    `json:"unread_count"`
	LastMessage string `json:"last_message"`
	LastAuthor  string `json:"last_author"`
	MediaUrl    string `json:"media_url"`
}

type GetChatsResponse struct {
	ActionType ActionType `json:"action_type"`
	Items      []Chat     `json:"items"`
}
