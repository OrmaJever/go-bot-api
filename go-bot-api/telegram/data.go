package telegram

import "time"

type Data struct {
	UpdateId           int64               `json:"update_id" bson:"update_id"`
	Message            *Message            `json:"message" bson:"message"`
	EditedMessage      *Message            `json:"edited_message" bson:"edited_message"`
	ChannelPost        *Message            `json:"channel_post" bson:"channel_post"`
	EditedChannelPost  *Message            `json:"edited_channel_post" bson:"edited_channel_post"`
	InlineQuery        *InlineQuery        `json:"inline_query" bson:"inline_query"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result" bson:"chosen_inline_result"`
	CallbackQuery      *CallbackQuery      `json:"callback_query" bson:"callback_query"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query" bson:"shipping_query"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query" bson:"pre_checkout_query"`
	Poll               *Poll               `json:"poll" bson:"poll"`
	PollAnswer         *PollAnswer         `json:"poll_answer" bson:"poll_answer"`
	MyChatMember       *ChatMemberUpdated  `json:"my_chat_member" bson:"my_chat_member"`
	ChatMember         *ChatMemberUpdated  `json:"chat_member" bson:"chat_member"`
	ChatJoinRequest    *ChatJoinRequest    `json:"chat_join_request" bson:"chat_join_request"`
	CreatedAt          time.Time           `json:"created_at" bson:"created_at"`
}
