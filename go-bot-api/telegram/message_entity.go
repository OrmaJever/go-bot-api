package telegram

type MessageEntity struct {
	Type          string `json:"type"`
	Offset        int32  `json:"offset"`
	Length        int32  `json:"length"`
	Url           string `json:"url"`
	User          *User  `json:"user"`
	Language      string `json:"language"`
	CustomEmojiId string `json:"custom_emoji_id"`
}
