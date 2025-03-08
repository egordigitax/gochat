package entities

type ChatType int

func (c ChatType) String() string {
	return [...]string{"Personal", "Group"}[c]
}

type Chat struct {
	Id          int      `json:"id" db:"id"`
	Uid         string   `json:"uid" db:"uid"`
	Title       string   `json:"title" db:"title"`
	MediaURL    string   `json:"media_url" db:"media_url"`
	UsersUids   []string `json:"users_uids" db:"users_uids"`
	UpdatedAt   string   `json:"updated_at" db:"updated_at"`
	ChatType    ChatType `json:"chat_type" db:"chat_type"`
	LastMessage Message  `json:"last_message" db:"last_message"`
}

type ChatUid string
