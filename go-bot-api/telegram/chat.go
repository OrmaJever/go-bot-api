package telegram

type Chat struct {
	Id                                 int64         `json:"id"`
	Type                               string        `json:"type"`
	Title                              string        `json:"title"`
	Username                           string        `json:"username"`
	FirstName                          string        `json:"first_name"`
	LastName                           string        `json:"last_name"`
	isForum                            string        `json:"is_forum"`
	Photo                              *ChatPhoto    `json:"photo"`
	ActiveUsernames                    []string      `json:"active_usernames"`
	EmojiStatusCustomEmojiId           string        `json:"emoji_status_custom_emoji_id"`
	bio                                string        `json:"bio"`
	HasPrivateForwards                 bool          `json:"has_private_forwards"`
	HasRestrictedVoiceAndVideoMessages bool          `json:"has_restricted_voice_and_video_messages"`
	JoinToSendMessages                 bool          `json:"join_to_send_messages"`
	JoinByRequest                      bool          `json:"join_by_request"`
	Description                        string        `json:"description"`
	InviteLink                         string        `json:"invite_link"`
	PinnedMessage                      *Message      `json:"pinned_message"`
	Permissions                        interface{}   `json:"permissions"`
	SlowModeDelay                      int64         `json:"slow_mode_delay"`
	MessageAutoDeleteTime              int64         `json:"message_auto_delete_time"`
	HasAggressiveAntiSpamEnabled       bool          `json:"has_aggressive_anti_spam_enabled"`
	HasHiddenMembers                   bool          `json:"has_hidden_members"`
	HasProtectedContent                bool          `json:"has_protected_content"`
	StickerSetName                     string        `json:"sticker_set_name"`
	CanSetStickerSet                   bool          `json:"can_set_sticker_set"`
	LinkedChatId                       int64         `json:"linked_chat_id"`
	Location                           *ChatLocation `json:"location"`
}
