package telegram

type InlineQuery struct {
	Id       int64     `json:"id"`
	From     *User     `json:"from"`
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
	ChatType string    `json:"chat_type"`
	Location *Location `json:"location"`
}
