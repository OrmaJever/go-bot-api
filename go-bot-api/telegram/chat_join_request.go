package telegram

type ChatJoinRequest struct {
	Chat       *Chat           `json:"chat"`
	From       *User           `json:"from"`
	UserChatId int64           `json:"user_chat_id"`
	Date       int64           `json:"date"`
	Bio        string          `json:"bio"`
	InviteLink *ChatInviteLink `json:"invite_link"`
}
