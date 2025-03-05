package domain

type Chat struct {
	Id          int      `json:"id" db:"id"`
	Title       string   `json:"title" db:"title"`
	MediaURL    string   `json:"media_url" db:"media_url"`
	UsersIds    []string `json:"users_ids" db:"users_ids"`
	LastMessage Message  `json:"last_message" db:"last_message"`
	UnreadCount int      `json:"unread_count" db:"unread_count"`
}
