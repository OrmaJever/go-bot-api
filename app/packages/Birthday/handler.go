package Birthday

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"main/models"
	"main/packages/SelectUser"
	"main/services"
	"main/telegram"
	"os"
)

var (
	Commands map[string]services.CallbackT
	tgApi    services.Telegram
)

func init() {
	Commands = make(map[string]services.CallbackT)
	Commands["happy_birthday"] = happyBirthdayCommand

	// run cron goroutine
	go services.Schedule("08:00", happyBirthday)

	err := godotenv.Load(".env")

	if err != nil {
		log.Println(err)
	}
}

func Message(_ *telegram.Data, _ *services.Telegram, _ *models.Bot) {}

func happyBirthdayCommand(data *telegram.Data, _ *services.Telegram, bot *models.Bot) {
	if bot.AdminId != data.Message.From.Id {
		return
	}

	happyBirthday()
}

func happyBirthday() {
	// Connect to Postgres
	postgres := services.ConnectToPostgres("select_user")

	defer postgres.Close()

	botName := os.Getenv("STATISTIC_BOT_NAME")

	bot := getBot(botName)

	tgApi = services.CreateTelegram(bot.Token)

	var users []SelectUser.User
	postgres.Model(&users).
		Where("birthday is not null AND TO_CHAR(birthday, 'dd-mm') = TO_CHAR(now(), 'dd-mm')").
		Select()

	for _, user := range users {
		text := fmt.Sprintf("Happy Birthday *%s* (@%s) 🥳🥳🥳", user.FirstName, tgApi.Format(user.Username))

		tgApi.SendMessage(user.ChatId, text, false, false)
	}
}

func getBot(botName string) models.Bot {
	postgres := services.ConnectToPostgres("telegram")
	defer postgres.Close()

	bot := models.Bot{}
	postgres.Model(&bot).
		Where("name = ?", botName).
		First()

	if bot.Id == 0 {
		log.Printf("Can not find bot with name [%s]\n", botName)
	}

	return bot
}
