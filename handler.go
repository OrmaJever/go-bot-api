package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"main/models"
	"main/packages/ChatStatistic"
	"main/packages/Debug"
	"main/packages/GetId"
	"main/packages/SelectUser"
	"main/packages/SpeechToText"
	"main/services"
	"main/telegram"
	"time"
)

func Handler(c *gin.Context) {
	botSecret := c.GetHeader("X-Telegram-Bot-Api-Secret-Token")

	bot := models.Bot{}
	Postgres.Model(&bot).Where("secret = ?", botSecret).Select()

	if bot.Id == 0 {
		c.JSON(200, JSON{
			"ok": true,
		})
		return
	}

	var data telegram.Data

	c.ShouldBindJSON(&data)

	go CreateWebhook(data)
	go CallPackages(&bot, &data)

	c.JSON(200, JSON{
		"ok": true,
	})
}

func CallPackages(bot *models.Bot, data *telegram.Data) {
	tgApi := services.CreateTelegram(bot.Token)

	command := parseCommand(data)

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
			if Debug.Commands[command] != nil {
				Debug.Commands[command](data, &tgApi, bot)
				result = true
			} else {
				Debug.Message(data, &tgApi, bot)
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
		}
	}

	if !result && data.Message.Chat.Id > 0 {
		tgApi.SendMessage(data.Message.Chat.Id, "Undefined command", false, false)
	}
}

func parseCommand(data *telegram.Data) string {
	var mess *telegram.Message

	if data.Message != nil {
		mess = data.Message
	} else if data.EditedMessage != nil {
		mess = data.EditedMessage
	}

	if len(mess.Entities) == 0 || mess.Entities[0].Type != "bot_command" {
		return ""
	}

	return mess.Text[mess.Entities[0].Offset+1 : mess.Entities[0].Length]
}

func CreateWebhook(data telegram.Data) {
	data.CreatedAt = time.Now()
	_, err := MongoCollection.InsertOne(context.TODO(), data)

	if err != nil {
		log.Println(err)
	}
}
