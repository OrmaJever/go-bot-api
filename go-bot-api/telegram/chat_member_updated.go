package telegram

type ChatMemberUpdated struct {
	Chat                    *Chat           `json:"chat"`
	From                    *User           `json:"from"`
	Date                    int64           `json:"date"`
	OldChatMember           interface{}     `json:"old_chat_member"`
	NewChatMember           interface{}     `json:"new_chat_member"`
	InviteLink              *ChatInviteLink `json:"invite_link"`
	ViaChatFolderInviteLink bool            `json:"via_chat_folder_invite_link"`
}
