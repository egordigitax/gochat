package resources

type Message struct {
	Username  string `json:"username"`
	AuthorUid string `json:"author_uid"`
	ChatUid   string `json:"chat_uid"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

type BaseMessage struct {
	Action ActionType  `json:"type"`
	Data   interface{} `json:"data"`
}

type ActionType string
