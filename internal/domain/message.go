package domain

type Message struct {
	ChatID    string `json:"chat_id" db:"chat_uid"`
	UserID    string `json:"user_id" db:"user_uid"`
	Text      string `json:"text" db:"message"`
	CreatedAt string `json:"created_at" db:"created_at"`
}
