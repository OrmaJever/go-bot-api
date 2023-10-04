package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/telegram"
	"net/http"
	"strconv"
	"strings"
)

type Telegram struct {
	token string
}

func CreateTelegram(token string) Telegram {
	return Telegram{token: token}
}

func (t *Telegram) SendMessage(chatId int64, text string, protect bool, silence bool) {
	var data = make(map[string]string)
	data["chat_id"] = strconv.FormatInt(chatId, 10)
	data["text"] = text
	data["parse_mode"] = "Markdown"
	data["disable_web_page_preview"] = "true"

	if protect {
		data["protect_content"] = "true"
	}

	if silence {
		data["disable_notification"] = "true"
	}

	t.sendRequest("sendMessage", data)
}

func (t *Telegram) SendPhoto(chatId int64, fileId string, text string, protect bool, silence bool) {
	var data = make(map[string]string)
	data["chat_id"] = strconv.FormatInt(chatId, 10)
	data["text"] = text
	data["parse_mode"] = "Markdown"
	data["disable_web_page_preview"] = "true"
	data["photo"] = fileId

	if protect {
		data["protect_content"] = "true"
	}

	if silence {
		data["disable_notification"] = "true"
	}

	fmt.Printf("%s\n", t.sendRequest("sendPhoto", data))
}

func (t *Telegram) SendVoice(chatId int64, fileId string, protect bool, silence bool) {
	var data = make(map[string]string)
	data["chat_id"] = strconv.FormatInt(chatId, 10)
	data["parse_mode"] = "Markdown"
	data["disable_web_page_preview"] = "true"
	data["voice"] = fileId

	if silence {
		data["disable_notification"] = "true"
	}

	if protect {
		data["protect_content"] = "true"
	}

	t.sendRequest("sendVoice", data)
}

func (t *Telegram) GetFile(fileId string) telegram.File {
	var data = make(map[string]string)
	data["file_id"] = fileId

	res := t.sendRequest("getFile", data)

	result := struct {
		Ok     bool          `json:"ok"`
		Result telegram.File `json:"result"`
	}{}
	json.Unmarshal(res, &result)

	if result.Ok {
		return result.Result
	} else {
		return telegram.File{}
	}
}

func (t *Telegram) Format(text string) string {
	if len(text) == 0 {
		return "-"
	}

	text = strings.Replace(text, "_", "\\_", -1)
	text = strings.Replace(text, "*", "\\*", -1)

	return text
}

func (t *Telegram) sendRequest(method string, params any) []byte {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", t.token, method)

	body, err := json.Marshal(params)

	if err != nil {
		log.Println(err)
		return []byte{}
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Println(err)
		return []byte{}
	}

	resBody, err := io.ReadAll(response.Body)

	if err != nil {
		log.Println(err)
		return []byte{}
	}

	return resBody
}
