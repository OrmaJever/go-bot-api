package SpeechToText

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/models"
	"main/services"
	"main/telegram"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
)

type callback func(data *telegram.Data, tgApi *services.Telegram, bot *models.Bot)

type WitResult struct {
	Text string `json:"_text"`
}

var Commands map[string]callback

func init() {
	Commands = make(map[string]callback)
}

func Message(data *telegram.Data, tgApi *services.Telegram, bot *models.Bot) {
	if data.Message == nil ||
		data.Message.Voice == nil ||
		data.Message.Voice.Duration >= 20 ||
		!data.Message.IsChat() {
		return
	}

	// 1) Получаем путь файла от телеграмма
	file := tgApi.GetFile(data.Message.Voice.FileId)

	// 2) Скачиваем его и пишет в файл в папку tmp
	response, err := http.Get(fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath))

	if err != nil {
		log.Println(err)
		return
	}

	resBody, err := io.ReadAll(response.Body)

	tmpFilePath := fmt.Sprintf("%svoice_%d", os.TempDir(), rand.Int())

	os.WriteFile(tmpFilePath, resBody, 0644)

	// 3) Конвертируем файл с помощью утилиты ffmpeg в ogg формат
	err = exec.Command("bash", "-c", fmt.Sprintf("ffmpeg -i %s -f ogg %s.ogg", tmpFilePath, tmpFilePath)).Run()
	go os.Remove(tmpFilePath)

	if err != nil {
		log.Println(err)
		return
	}

	fileData, err := os.ReadFile(fmt.Sprintf("%s.ogg", tmpFilePath))
	go os.Remove(fmt.Sprintf("%s.ogg", tmpFilePath))

	if err != nil {
		log.Println(err)
		return
	}

	// 4) Отдаем ogg файл в wit.ai
	req, err := http.NewRequest("POST", "https://api.wit.ai/speech?v=20200422", bytes.NewBuffer(fileData))

	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Add("Content-Type", "audio/ogg")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("WIT_TOKEN")))

	client := http.Client{}
	response, err = client.Do(req)

	if err != nil {
		log.Println(err)
		return
	}

	// 5) Получаем результат и отправляем в чат телеграмма
	resBody, err = io.ReadAll(response.Body)

	if err != nil {
		log.Println(err)
		return
	}

	var jsonResult WitResult

	json.Unmarshal(resBody, &jsonResult)

	if len(jsonResult.Text) == 0 {
		return
	}

	text := fmt.Sprintf("*%s* сказав:\n`%s`", data.Message.From.FirstName, jsonResult.Text)
	tgApi.SendMessage(data.Message.Chat.Id, text, true, true)
}
