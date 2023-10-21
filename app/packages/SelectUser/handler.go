package SelectUser

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
	"log"
	"main/models"
	"main/services"
	"main/telegram"
	"math/rand"
	"os"
	"time"
)

var (
	Commands map[string]services.CallbackT
	postgres *pg.DB
	lang     map[string]string
	sleep    time.Time
)

func init() {
	Commands = make(map[string]services.CallbackT)
	Commands["reg"] = reg
	Commands["check"] = check
	Commands["run"] = run
	Commands["userlist"] = userlist

	err := godotenv.Load(".env")

	if err != nil {
		log.Println(err)
	}

	lang = loadLang("ua")

	postgres = pg.Connect(&pg.Options{
		Addr:     os.Getenv("PG_ADDR"),
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		Database: "select_user",
	})

	if os.Getenv("GIN_MODE") == "debug" {
		postgres.AddQueryHook(services.PostgresLogger{})
	}
}

func reg(data *telegram.Data, tgApi *services.Telegram, _ *models.Bot) {
	if data.Message == nil {
		return
	}
	user := User{}

	postgres.Model(&user).
		Where("tg_id = ? and chat_id = ?", data.Message.From.Id, data.Message.Chat.Id).
		Select()

	if user.Id > 0 {
		text := fmt.Sprintf("Already reg in %s", user.CreatedAt)
		tgApi.SendMessage(data.Message.Chat.Id, text, true, true)
	} else {
		user.TgId = data.Message.From.Id
		user.ChatId = data.Message.Chat.Id
		user.FirstName = fmt.Sprintf("%s %s", data.Message.From.FirstName, data.Message.From.LastName)
		user.Username = data.Message.From.Username

		_, err := postgres.Model(&user).Insert()
		trans("reg_success")

		if err == nil {
			tgApi.SendMessage(data.Message.Chat.Id, trans("reg_success"), true, true)
		} else {
			log.Println(err)
			tgApi.SendMessage(data.Message.Chat.Id, trans("error"), true, true)
		}
	}
}

func check(data *telegram.Data, tgApi *services.Telegram, _ *models.Bot) {
	if data.Message == nil {
		return
	}

	var selectedUser []struct {
		TgId        int64
		Id          int32
		ChatId      int64
		CustomizeId int32
		Count       int16
		FirstName   string
		Username    string
	}

	query := `
		SELECT 
		    su.tg_id, 
		    max(su.id) id, 
		    max(su.chat_id) chat_id, 
		    max(su.customize_id) customize_id, 
		    count(su.tg_id) count,
			max(u.first_name) first_name,
			max(u.username) username
		FROM selected_users su
		JOIN users u ON u.tg_id = su.tg_id
		WHERE su.chat_id = ? and u.chat_id = su.chat_id
		GROUP BY su.tg_id
		ORDER BY count(su.*) desc
	`

	_, err := postgres.Query(&selectedUser, query, data.Message.Chat.Id)

	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		tgApi.SendMessage(data.Message.Chat.Id, trans("error"), true, true)
		return
	}

	text := trans("pidor_result_header")

	for i, user := range selectedUser {
		text += fmt.Sprintf(trans("pidor_result_line"), i, user.FirstName, user.Username, user.Count)
	}

	tgApi.SendMessage(data.Message.Chat.Id, text, true, true)
}

func run(data *telegram.Data, tgApi *services.Telegram, _ *models.Bot) {
	if data.Message == nil {
		return
	}

	if time.Now().Sub(sleep).Minutes() < 3 && os.Getenv("GIN_MODE") != "debug" {
		tgApi.SendMessage(data.Message.Chat.Id, trans("too_fast"), false, true)
		return
	} else {
		sleep = time.Now()
	}

	var todayUser SelectedUser

	err := postgres.Model(&todayUser).
		Where("su.chat_id = ?", data.Message.Chat.Id).
		Where("date(su.created_at) = ?", time.Now().Format("2006-01-02")).
		Relation("User").
		Relation("Customize").
		First()

	if todayUser.Id > 0 {
		text := fmt.Sprintf(
			trans("already_run"),
			todayUser.User.FirstName,
			tgApi.Format(todayUser.User.Username),
		)

		if len(todayUser.Customize.Image) > 0 {
			tgApi.SendPhoto(data.Message.Chat.Id, todayUser.Customize.Image, text, true, true)
		} else {
			tgApi.SendMessage(data.Message.Chat.Id, text, true, true)
		}
		return
	}

	var user User

	err = postgres.Model(&user).
		Where("chat_id = ?", data.Message.Chat.Id).
		OrderExpr("random()").
		Select()

	if err == sql.ErrNoRows {
		tgApi.SendMessage(data.Message.Chat.Id, trans("empty_users"), true, true)
		return
	}

	var customize Customize

	err = postgres.Model(&customize).
		Where("user_id = ?", user.Id).
		OrderExpr("random()").
		Select()

	selectedUser := SelectedUser{
		TgId:        user.TgId,
		ChatId:      data.Message.Chat.Id,
		Type:        1,
		CustomizeId: customize.Id,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	_, err = postgres.Model(&selectedUser).Insert()

	if err != nil {
		log.Println(err)
		tgApi.SendMessage(data.Message.Chat.Id, trans("error"), true, true)
		return
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	text := trans(fmt.Sprintf("step1_%d", r1.Intn(4)))
	tgApi.SendMessage(data.Message.Chat.Id, text, true, true)

	time.Sleep(1 * time.Second)

	text = fmt.Sprintf(trans("pidor_text"), user.FirstName, tgApi.Format(user.Username))

	if customize.Id != 0 {
		tgApi.SendPhoto(data.Message.Chat.Id, customize.Image, text, true, true)
	} else {
		tgApi.SendMessage(data.Message.Chat.Id, text, true, true)
	}

	time.Sleep(1 * time.Second)

	voices := [2]string{
		"AwACAgIAAxkBAAIDOmR69n4ZyrEheO7XqWVRkBgAAdzzgAACrjIAAuu8QUvpPsJCKMeEaS8E",
		"AwACAgIAAxkBAAIBt2MjnLL_Bry2-vVz5giGxIp-68OMAALIHgACPVwhSVIQXq0eHGbVKQQ",
	}

	tgApi.SendVoice(data.Message.Chat.Id, voices[r1.Intn(1)], true, true)
}

func userlist(data *telegram.Data, tgApi *services.Telegram, _ *models.Bot) {
	if data.Message == nil {
		return
	}

	var users []User
	err := postgres.Model(&users).
		Where("chat_id = ?", data.Message.Chat.Id).
		Order("created_at desc").
		Select()

	if err != nil {
		log.Println(err)
		tgApi.SendMessage(data.Message.Chat.Id, trans("error"), true, true)
		return
	}

	text := trans("userlist_result_header")

	for i, user := range users {
		date, _ := time.Parse("2006-01-02 15:04:05", user.CreatedAt)
		text += fmt.Sprintf(trans("userlist_result_line"), i, user.FirstName, user.Username, date.Format("01 Jan 15:04"))
	}

	tgApi.SendMessage(data.Message.Chat.Id, text, true, true)
}

func Message(_ *telegram.Data, _ *services.Telegram, _ *models.Bot) {}

func loadLang(lang string) map[string]string {
	data, err := os.ReadFile(fmt.Sprintf("packages/SelectUser/lang/%s.json", lang))

	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]string

	err = json.Unmarshal(data, &result)

	if err != nil {
		log.Fatalln(err)
	}

	return result
}

func trans(key string) string {
	return fmt.Sprintf("%s", lang[key])
}
