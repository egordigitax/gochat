package ws_api

import (
	"encoding/json"
)
type ActionType string

type RootMessage struct {
	ActionType ActionType      `json:"action_type"`
	RawPayload json.RawMessage `json:"payload"`
}

type GetMessageFromClientRequest struct {
	Text string `json:"text"`
}

type SendMessageToClientResponse struct {
	Text      string `json:"text"`
	AuthorId  string `json:"author_id"`
	Nickname  string `json:"nickname"`
	CreatedAt string `json:"created_at"`
}

type Chat struct {
	Title       string `json:"title"`
	UnreadCount int    `json:"unread_count"`
	LastMessage string `json:"last_message"`
	LastAuthor  string `json:"last_author"`
	MediaUrl    string `json:"media_url"`
}

type GetChatsResponse struct {
	Items []Chat `json:"items"`
}
