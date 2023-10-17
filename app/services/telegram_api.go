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

type BadResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int16  `json:"result"`
	Description string `json:"description"`
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

	res, code := t.sendRequest("sendMessage", data)

	if code == 400 {
		var badResponse BadResponse
		json.Unmarshal(res, &badResponse)
		log.Println(badResponse.Description)
	}
}

func (t *Telegram) SendPhoto(chatId int64, fileId string, text string, protect bool, silence bool) {
	var data = make(map[string]string)
	data["chat_id"] = strconv.FormatInt(chatId, 10)
	data["caption"] = text
	data["parse_mode"] = "Markdown"
	data["disable_web_page_preview"] = "true"
	data["photo"] = fileId

	if protect {
		data["protect_content"] = "true"
	}

	if silence {
		data["disable_notification"] = "true"
	}

	res, code := t.sendRequest("sendPhoto", data)

	if code == 400 {
		var badResponse BadResponse
		json.Unmarshal(res, &badResponse)
		log.Println(badResponse.Description)
	}
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

	res, code := t.sendRequest("sendVoice", data)

	if code == 400 {
		var badResponse BadResponse
		json.Unmarshal(res, &badResponse)
		log.Println(badResponse.Description)
	}
}

func (t *Telegram) GetFile(fileId string) telegram.File {
	var data = make(map[string]string)
	data["file_id"] = fileId

	res, code := t.sendRequest("getFile", data)

	if code == 400 {
		var badResponse BadResponse
		json.Unmarshal(res, &badResponse)
		log.Println(badResponse.Description)
		return telegram.File{}
	}

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

	text = strings.ReplaceAll(text, "_", "\\_")
	text = strings.ReplaceAll(text, "*", "\\*")

	return text
}

func (t *Telegram) sendRequest(method string, params any) ([]byte, int) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", t.token, method)

	body, err := json.Marshal(params)

	if err != nil {
		log.Println(err)
		return []byte{}, 0
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Println(err)
		return []byte{}, 0
	}

	resBody, err := io.ReadAll(response.Body)

	if err != nil {
		log.Println(err)
		return []byte{}, 0
	}

	return resBody, response.StatusCode
}
