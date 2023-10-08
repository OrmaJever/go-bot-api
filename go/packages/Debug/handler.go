package Debug

import (
	"main/models"
	"main/services"
	"main/telegram"
)

type callback func(data *telegram.Data, tgApi *services.Telegram, bot *models.Bot)

var Commands map[string]callback

func init() {
	Commands = make(map[string]callback)
	Commands["img"] = getImage
}

func getImage(data *telegram.Data, tgApi *services.Telegram, _ *models.Bot) {

	if data.Message == nil || data.Message.Chat.Id < 0 {
		return
	}

	tgApi.SendPhoto(data.Message.Chat.Id, data.Message.Text[5:], "", false, false)
}

func Message(data *telegram.Data, tgApi *services.Telegram, _ *models.Bot) {
	if data.Message == nil || data.Message.Chat.Id < 0 {
		return
	}

	if len(data.Message.Photo) > 0 {
		tgApi.SendMessage(data.Message.Chat.Id, data.Message.Photo[0].FileId, false, false)
	}

	if data.Message.Voice != nil {
		tgApi.SendMessage(data.Message.Chat.Id, data.Message.Voice.FileId, false, false)
	}
}
