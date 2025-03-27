package types

type ChatType int

func (c ChatType) String() string {
	return [...]string{"Personal", "Group"}[c]
}

type Chat struct {
	Id          int      `db:"id"`
	Uid         string   `db:"uid"`
	Title       string   `db:"title"`
	MediaURL    string   `db:"media_url"`
	UpdatedAt   string   `db:"updated_at"`
	ChatType    ChatType `db:"chat_type"`
	LastMessage Message  `db:"message"`
	UsersUids   []string `db:"users_uids"`
}

type ChatRole int

func (c ChatRole) String() string {
	return [...]string{"Personal", "Group"}[c]
}
