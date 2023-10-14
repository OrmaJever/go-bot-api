package GetId

import (
	"fmt"
	"main/models"
	"main/services"
	"main/telegram"
)

type callback func(data *telegram.Data, tgApi *services.Telegram, _ *models.Bot)

var Commands map[string]callback

func init() {
	Commands = make(map[string]callback)
	Commands["get_my_id"] = getMyId
}

func getMyId(data *telegram.Data, tgApi *services.Telegram, _ *models.Bot) {
	if data.Message == nil || data.Message.IsChat() {
		return
	}

	tgApi.SendMessage(data.Message.Chat.Id, fmt.Sprintf("Your telegram id: %d", data.Message.From.Id), false, false)
}

func Message(_ *telegram.Data, _ *services.Telegram, _ *models.Bot) {}
