package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"main/models"
	"main/packages/Admin"
	"main/packages/Alerts"
	"main/packages/Birthday"
	"main/packages/ChatStatistic"
	DebugPackage "main/packages/Debug"
	"main/packages/GetId"
	"main/packages/SelectUser"
	"main/packages/SpeechToText"
	"main/services"
	"main/telegram"
	"strings"
	"time"
)

func Handler(c *gin.Context) {
	go c.JSON(200, JSON{
		"ok": true,
	})

	botSecret := c.GetHeader("X-Telegram-Bot-Api-Secret-Token")

	bot := models.Bot{}
	Postgres.Model(&bot).
		Where("secret = ?", botSecret).
		First()

	if bot.Id == 0 {
		log.Printf("Can not find bot with secret [%s]\n", botSecret)
		return
	}

	var data telegram.Data
	err := c.ShouldBindJSON(&data)

	if err != nil {
		log.Println(err)
		return
	}

	go CreateWebhook(data)

	// сообщение старое, приходило когда бот упал, игнорируем
	if data.Message != nil && time.Now().Sub(time.Unix(data.Message.Date, 0)) > 30*time.Second {
		return
	}

	go CallPackages(&bot, &data)
}

func CallPackages(bot *models.Bot, data *telegram.Data) {
	tgApi := services.CreateTelegram(bot.Token)

	command := parseCommand(data, bot.Name)

	result := len(command) == 0 // если это не комманда то сразу true

	for _, name := range bot.Packages {
		switch name {
		case "GetId":
			if GetId.Commands[command] != nil {
				GetId.Commands[command](data, &tgApi, bot)
				result = true
			} else {
				GetId.Message(data, &tgApi, bot)
			}
			break

		case "Debug":
			if DebugPackage.Commands[command] != nil {
				DebugPackage.Commands[command](data, &tgApi, bot)
				result = true
			} else {
				DebugPackage.Message(data, &tgApi, bot)
			}
			break

		case "ChatStatistic":
			if ChatStatistic.Commands[command] != nil {
				ChatStatistic.Commands[command](data, &tgApi, bot)
				result = true
			} else {
				ChatStatistic.Message(data, &tgApi, bot)
			}
			break

		case "SelectUser":
			if SelectUser.Commands[command] != nil {
				SelectUser.Commands[command](data, &tgApi, bot)
				result = true
			} else {
				SelectUser.Message(data, &tgApi, bot)
			}
			break

		case "SpeechToText":
			if SpeechToText.Commands[command] != nil {
				SpeechToText.Commands[command](data, &tgApi, bot)
				result = true
			} else {
				SpeechToText.Message(data, &tgApi, bot)
			}
			break

		case "Admin":
			if Admin.Commands[command] != nil {
				Admin.Commands[command](data, &tgApi, bot)
				result = true
			} else {
				Admin.Message(data, &tgApi, bot)
			}

		case "Alerts":
			if Alerts.Commands[command] != nil {
				Alerts.Commands[command](data, &tgApi, bot)
				result = true
			} else {
				Alerts.Message(data, &tgApi, bot)
			}

		case "Birthday":
			if Birthday.Commands[command] != nil {
				Birthday.Commands[command](data, &tgApi, bot)
				result = true
			} else {
				Birthday.Message(data, &tgApi, bot)
			}
		default:
			log.Printf("Can not find package [%s]\n", name)
		}
	}

	if !result && !data.Message.IsChat() {
		tgApi.SendMessage(data.Message.Chat.Id, "Undefined command", false, false)
	}
}

func parseCommand(data *telegram.Data, name string) string {
	var mess *telegram.Message

	if data.Message != nil {
		mess = data.Message
	} else if data.EditedMessage != nil {
		mess = data.EditedMessage
	}

	if len(mess.Entities) == 0 || mess.Entities[0].Type != "bot_command" {
		return ""
	}

	command := mess.Text[mess.Entities[0].Offset+1 : mess.Entities[0].Length]

	return strings.ReplaceAll(command, "@"+name, "")
}

func CreateWebhook(data telegram.Data) {
	data.CreatedAt = time.Now()
	_, err := MongoCollection.InsertOne(context.TODO(), data)

	if err != nil {
		log.Println(err)
	}
}
