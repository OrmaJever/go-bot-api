package telegram

type Message struct {
	MessageId                     int64                          `json:"message_id"`
	MessageThreadId               int64                          `json:"message_thread_id"`
	Date                          int64                          `json:"date"`
	Text                          string                         `json:"text"`
	Chat                          *Chat                          `json:"chat"`
	From                          *User                          `json:"from"`
	SenderChat                    *Chat                          `json:"sender_chat"`
	ForwardFrom                   *User                          `json:"forward_from"`
	ForwardFromChat               *Chat                          `json:"forward_from_chat"`
	ForwardFromMessageId          int64                          `json:"forward_from_message_id"`
	ForwardSignature              string                         `json:"forward_signature"`
	ForwardSenderName             string                         `json:"forward_sender_name"`
	ForwardDate                   int64                          `json:"forward_date"`
	IsTopicMessage                bool                           `json:"is_topic_message"`
	IsAutomaticForward            bool                           `json:"is_automatic_forward"`
	ReplyToMessage                *Message                       `json:"reply_to_message"`
	ViaBot                        *User                          `json:"via_bot"`
	EditDate                      int64                          `json:"edit_date"`
	HasProtectedContent           bool                           `json:"has_protected_content"`
	MediaGroupId                  string                         `json:"media_group_id"`
	AuthorSignature               string                         `json:"author_signature"`
	Entities                      []*MessageEntity               `json:"entities"`
	Animation                     *Animation                     `json:"animation"`
	Audio                         *Audio                         `json:"audio"`
	Document                      *Document                      `json:"document"`
	Photo                         []*PhotoSize                   `json:"photo"`
	Sticker                       *Sticker                       `json:"sticker"`
	Video                         *Video                         `json:"video"`
	VideoNote                     *VideoNote                     `json:"video_note"`
	Voice                         *Voice                         `json:"voice"`
	Caption                       string                         `json:"caption"`
	CaptionEntities               []*MessageEntity               `json:"caption_entities"`
	HasMediaSpoiler               bool                           `json:"has_media_spoiler"`
	Contact                       *Contact                       `json:"contact"`
	Dice                          *Dice                          `json:"dice"`
	Game                          *Game                          `json:"game"`
	Poll                          *Poll                          `json:"poll"`
	Venue                         *Venue                         `json:"venue"`
	Location                      *Location                      `json:"location"`
	NewChatMembers                []*User                        `json:"new_chat_members"`
	LeftChatMember                *User                          `json:"left_chat_member"`
	NewChatTitle                  string                         `json:"new_chat_title"`
	NewChatPhoto                  []*PhotoSize                   `json:"new_chat_photo"`
	DeleteChatPhoto               bool                           `json:"delete_chat_photo"`
	GroupChatCreated              bool                           `json:"group_chat_created"`
	SupergroupChatCreated         bool                           `json:"supergroup_chat_created"`
	ChannelChatCreated            bool                           `json:"channel_chat_created"`
	MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed"`
	MigrateToChatId               int64                          `json:"migrate_to_chat_id"`
	MigrateFromChatId             int64                          `json:"migrate_from_chat_id"`
	PinnedMessage                 *Message                       `json:"pinned_message"`
	Invoice                       *Invoice                       `json:"invoice"`
	SuccessfulPayment             interface{}                    `json:"successful_payment"` // SuccessfulPayment
	UserShared                    interface{}                    `json:"user_shared"`        // UserShared
	ChatShared                    interface{}                    `json:"chat_shared"`        // ChatShared
	ConnectedWebsite              string                         `json:"connected_website"`
	WriteAccessAllowed            interface{}                    `json:"write_access_allowed"`            // WriteAccessAllowed
	PassportData                  interface{}                    `json:"passport_data"`                   // PassportData
	ProximityAlertTriggered       interface{}                    `json:"proximity_alert_triggered"`       // ProximityAlertTriggered
	ForumTopicCreated             interface{}                    `json:"forum_topic_created"`             // ForumTopicCreated
	ForumTopicEdited              interface{}                    `json:"forum_topic_edited"`              // ForumTopicEdited
	ForumTopicClosed              interface{}                    `json:"forum_topic_closed"`              // ForumTopicClosed
	ForumTopicReopened            interface{}                    `json:"forum_topic_reopened"`            // ForumTopicReopened
	GeneralForumTopicHidden       interface{}                    `json:"general_forum_topic_hidden"`      // GeneralForumTopicHidden
	GeneralForumTopicUnhidden     interface{}                    `json:"general_forum_topic_unhidden"`    // GeneralForumTopicUnhidden
	VideoChatScheduled            interface{}                    `json:"video_chat_scheduled"`            // VideoChatScheduled
	VideoChatStarted              interface{}                    `json:"video_chat_started"`              // VideoChatStarted
	VideoChatEnded                interface{}                    `json:"video_chat_ended"`                // VideoChatEnded
	VideoChatParticipantsInvited  interface{}                    `json:"video_chat_participants_invited"` // VideoChatParticipantsInvited
	WebAppData                    *WebAppData                    `json:"web_app_data"`
	ReplyMarkup                   *InlineKeyboardMarkup          `json:"reply_markup"`
}

func (m *Message) IsChat() bool {
	return m.Chat != nil && m.Chat.Id < 0
}
