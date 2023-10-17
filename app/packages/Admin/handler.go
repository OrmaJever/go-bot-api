package Admin

import (
	"log"
	"main/models"
	"main/services"
	"main/telegram"
	"regexp"
	"strconv"
)

type callback func(data *telegram.Data, tgApi *services.Telegram, bot *models.Bot)

var Commands map[string]callback

func init() {
	Commands = make(map[string]callback)
	Commands["send_message"] = sendMessage
}

func sendMessage(data *telegram.Data, tgApi *services.Telegram, bot *models.Bot) {
	match := regexp.
		MustCompile("^/send_message (-?\\d+) (?m)(.+?)$").
		FindStringSubmatch(data.Message.Text)

	if len(match) < 3 {
		log.Println("len(match) < 3")
		return
	}

	chatId, err := strconv.Atoi(match[1])
	if err != nil {
		log.Println(err)
		return
	}

	tgApi.SendMessage(int64(chatId), match[2], true, false)
}

func Message(_ *telegram.Data, _ *services.Telegram, _ *models.Bot) {

}
