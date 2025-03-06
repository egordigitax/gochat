package domain

type Message struct {
	Uid       string `json:"uid" db:"uid"`
	ChatUid   string `json:"chat_uid" db:"chat_uid"`
	UserUid   string `json:"user_uid" db:"user_uid"`
	Text      string `json:"text" db:"message"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UserInfo  User   `json:"user_info" db:"user_info"`
}
