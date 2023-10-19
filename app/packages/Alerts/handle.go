package Alerts

import (
	"database/sql"
	"encoding/json"
	"github.com/joho/godotenv"
	"io"
	"log"
	"main/models"
	"main/services"
	"main/telegram"
	"net/http"
	"os"
	"time"
)

const (
	region = "Чернігівська_область"
	URL    = "https://alarmmap.online/assets/json/_alarms/siren.json"
)

var (
	Commands map[string]services.CallbackT
	status   bool
	bot      models.Bot
)

type District struct {
	District  string `json:"district"`
	Start     string `json:"start"`
	SirenType string `json:"sirenType"`
}

func init() {
	Commands = make(map[string]services.CallbackT)

	// Load env
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalln(err)
	}

	bot = getBot()

	go run()
}

func Message(_ *telegram.Data, _ *services.Telegram, _ *models.Bot) {}

func run() {
	for {
		getData()

		time.Sleep(3 * time.Second)
	}
}

func getData() {
	response, err := http.Get(URL)

	if err != nil {
		return
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return
	}

	parseAlert(body)
}

func parseAlert(body []byte) {
	var res []District

	err := json.Unmarshal(body, &res)

	if err != nil {
		return
	}

	var exist *District

	for _, district := range res {
		if district.District == region {
			exist = &district
		}
	}

	if exist != nil && !status {
		date, _ := time.Parse("2006-01-02 15:04:05Z07:00", exist.Start)
		status = true
		AlertOn(date)
	}

	if exist == nil && status {
		status = false
		AlertOff()
	}
}

func AlertOn(start time.Time) {
	if time.Now().Sub(start) > 20*time.Minute {
		return // тревога началась давно, уже нет смысла постить
	}

	text := "🛑Повітряна тривога Чернігів!"
	tgApi := services.CreateTelegram(bot.Token)
	tgApi.SendMessage(bot.AdminId, text, false, false)
}

func AlertOff() {
	text := "🟢Відбій повітряної тривоги Чернігів!"
	tgApi := services.CreateTelegram(bot.Token)
	tgApi.SendMessage(bot.AdminId, text, false, false)
}

func getBot() models.Bot {
	botName := os.Getenv("BOT_NAME")

	postgres := services.ConnectToPostgres(os.Getenv("PG_DATABASE"))

	defer postgres.Close()

	var bot models.Bot
	err := postgres.Model(&bot).
		Where("name = ?", botName).
		First()

	if err == sql.ErrNoRows {
		log.Printf("Cannot get bot [%s]\n", botName)
		return bot
	}

	return bot
}
