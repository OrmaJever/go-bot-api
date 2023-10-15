package Admin

import (
	"main/models"
	"main/services"
	"main/telegram"
)

type callback func(data *telegram.Data, tgApi *services.Telegram, bot *models.Bot)

var Commands map[string]callback

func init() {
	Commands = make(map[string]callback)
	Commands["send_message"] = sendMessage
	Commands["remove_message"] = removeMessage
}

func sendMessage(data *telegram.Data, tgApi *services.Telegram, bot *models.Bot) {

}

func removeMessage(data *telegram.Data, tgApi *services.Telegram, bot *models.Bot) {

}
