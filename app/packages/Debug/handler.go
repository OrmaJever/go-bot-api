package Debug

import (
	"fmt"
	"main/models"
	"main/services"
	"main/telegram"
)

var Commands map[string]services.CallbackT

func init() {
	Commands = make(map[string]services.CallbackT)
	Commands["img"] = getImage
}

func getImage(data *telegram.Data, tgApi *services.Telegram, _ *models.Bot) {

	if data.Message == nil || data.Message.IsChat() {
		return
	}

	tgApi.SendPhoto(data.Message.Chat.Id, data.Message.Text[5:], "", false, false)
}

func Message(data *telegram.Data, tgApi *services.Telegram, bot *models.Bot) {
	if data.Message == nil || data.Message.IsChat() || data.Message.From.Id != bot.AdminId {
		return
	}

	var text string

	if data.Message.ForwardFromMessageId != 0 {
		text += fmt.Sprintf("Message id: `%d`\n", data.Message.ForwardFromMessageId)
	}

	if data.Message.ForwardFromChat != nil {
		text += fmt.Sprintf("Chat \\[%s] id: `%d`\n", data.Message.ForwardFromChat.Title, data.Message.ForwardFromChat.Id)
	}

	if data.Message.ForwardFrom != nil {
		text += fmt.Sprintf("User \\[%s] id: `%d`\n", data.Message.ForwardFrom.FirstName, data.Message.ForwardFrom.Id)
	}

	if len(data.Message.Photo) > 0 {
		text += fmt.Sprintf("Photo file id: `%s`\n", data.Message.Photo[0].FileId)
	}

	if data.Message.Voice != nil {
		text += fmt.Sprintf("Voice file id: `%s`\n", data.Message.Voice.FileId)
	}

	if len(text) > 0 {
		tgApi.SendMessage(data.Message.Chat.Id, text, false, false)
	}
}
