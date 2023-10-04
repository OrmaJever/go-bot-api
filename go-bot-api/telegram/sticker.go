package telegram

type Sticker struct {
	FileId           string      `json:"file_id"`
	FileUniqueId     string      `json:"file_unique_id"`
	Type             string      `json:"type"`
	Width            int16       `json:"width"`
	Height           int16       `json:"height"`
	IsAnimated       bool        `json:"is_animated"`
	IsVideo          bool        `json:"is_video"`
	Thumbnail        *PhotoSize  `json:"thumbnail"`
	Emoji            string      `json:"emoji"`
	SetName          string      `json:"set_name"`
	PremiumAnimation *File       `json:"premium_animation"`
	MaskPosition     interface{} `json:"mask_position"`
	CustomEmojiId    string      `json:"custom_emoji_id"`
	NeedsRepainting  bool        `json:"needs_repainting"`
	FileSize         int32       `json:"file_size"`
}
