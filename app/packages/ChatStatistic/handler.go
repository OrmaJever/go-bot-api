package ChatStatistic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"main/models"
	"main/services"
	"main/telegram"
	"os"
	"strings"
	"time"
)

type callback func(data *telegram.Data, tgApi *services.Telegram, bot *models.Bot)

var Commands map[string]callback
var mongoCollection *mongo.Collection
var postgres *pg.DB
var lang map[string]string

const chatId int64 = -1001524992976
const botName = "pidor_bp_bot"

type user struct {
	TgId      int64  `bson:"_id"`
	FirstName string `bson:"first_name"`
	Username  string `bson:"username"`
	Count     int32  `bson:"count"`
}

func init() {
	Commands = make(map[string]callback)

	// Load env
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalln(err)
	}

	// Load lang package
	lang = loadLang("ua")

	// run cron goroutine
	go services.Schedule("23:59", calculateStatistic)
}

func Message(_ *telegram.Data, _ *services.Telegram, _ *models.Bot) {}

func calculateStatistic() {
	// Connect to Postgres
	pgTg := services.ConnectToPostgres(os.Getenv("PG_DATABASE"))
	postgres = services.ConnectToPostgres("select_user")

	// Connect to Mongo
	var mongoClient *mongo.Client
	mongoCollection, mongoClient = services.ConnectToMongo()

	defer postgres.Close()
	defer pgTg.Close()
	defer mongoClient.Disconnect(context.Background())

	var bot models.Bot
	err := pgTg.Model(&bot).
		Where("name = ?", botName).
		Select()

	if err == sql.ErrNoRows {
		log.Printf("Cannot get bot [%s]\n", botName)
		return
	}

	text := getFormattedText(chatId)

	if len(text) > 0 {
		tgApi := services.CreateTelegram(bot.Token)
		tgApi.SendMessage(chatId, text, true, true)
	}
}

func getFormattedText(chatId int64) string {
	allMessageCount := getAllMessage(chatId)
	voiceMessageCount := getVoiceMessage(chatId)
	videoNotes := getVideoNotes(chatId)
	forwardedCount := getForwardedCount(chatId)
	topUsers := getTopUsers(chatId)
	inactiveUsers := getInactiveUsers(chatId, topUsers)
	forwarded := getForwarded(chatId)

	if allMessageCount == 0 {
		return ""
	}

	text := trans("header")
	text += fmt.Sprintf(trans("all_messages"), allMessageCount)

	if voiceMessageCount > 0 {
		text += fmt.Sprintf(trans("voice_messages"), voiceMessageCount)
	}
	if videoNotes > 0 {
		text += fmt.Sprintf(trans("video_notes"), videoNotes)
	}

	if forwardedCount > 0 {
		text += fmt.Sprintf(trans("forwarded_count"), forwardedCount)
	}

	text += trans("active_users")
	for _, user := range topUsers {
		text += fmt.Sprintf(trans("active_user_line"), user.FirstName, user.Username, user.Count)
	}

	text += trans("inactive_users")
	for _, user := range inactiveUsers {
		text += fmt.Sprintf("*%s* (`%s`), ", user.FirstName, user.Username)
	}

	text = strings.TrimSuffix(text, ", ")

	if forwarded.TgId != 0 {
		text += fmt.Sprintf(trans("forwarded_messages"), forwarded.FirstName, forwarded.Username, forwarded.Count)
	}

	text += trans("footer")

	return text
}

func getAllMessage(chatId int64) int64 {
	filter := getFilters(chatId)
	count, err := mongoCollection.CountDocuments(context.Background(), filter)

	if err != nil {
		log.Println(err)
		return 0
	}

	return count
}

func getVoiceMessage(chatId int64) int64 {
	filter := getFilters(chatId)
	filter["message.voice"] = bson.M{"$ne": nil}

	count, err := mongoCollection.CountDocuments(context.Background(), filter)

	if err != nil {
		log.Println(err)
		return 0
	}

	return count
}

func getVideoNotes(chatId int64) int64 {
	filter := getFilters(chatId)
	filter["message.videonote"] = bson.M{"$ne": nil}

	count, err := mongoCollection.CountDocuments(context.Background(), filter)

	if err != nil {
		log.Println(err)
		return 0
	}

	return count
}

func getForwardedCount(chatId int64) int64 {
	filter := getFilters(chatId)
	filter["message.forwardfrom"] = bson.M{"$ne": nil}

	count, err := mongoCollection.CountDocuments(context.Background(), filter)

	if err != nil {
		log.Println(err)
		return 0
	}

	return count
}

func getTopUsers(chatId int64) []user {
	match := bson.D{{"$match", getFilters(chatId)}}
	group := bson.D{{"$group", bson.D{
		{"_id", "$message.from.id"},
		{"count", bson.M{"$count": bson.D{}}},
		{"first_name", bson.M{"$max": "$message.from.firstname"}},
		{"username", bson.M{"$max": "$message.from.username"}},
	}}}
	sort := bson.D{{"$sort", bson.M{"count": -1}}}
	cursor, err := mongoCollection.Aggregate(context.Background(), mongo.Pipeline{match, group, sort})

	var result []user

	if err != nil {
		log.Println(err)
		return result
	}

	cursor.All(context.Background(), &result)

	return result
}

func getInactiveUsers(chatId int64, activeUsers []user) []user {

	var activeIds []int64
	if len(activeUsers) > 0 {
		activeIds = make([]int64, len(activeUsers))

		for _, us := range activeUsers {
			activeIds = append(activeIds, us.TgId)
		}
	}

	var users []user

	query := postgres.Model(&users).
		ColumnExpr("tg_id as _id, first_name, username").
		Where("chat_id = ?", chatId)

	if len(activeIds) > 0 {
		query.WhereIn("tg_id in (?)", activeIds)
	}

	err := query.Select()

	if err != nil {
		log.Println(err)
	}

	return users
}

func getForwarded(chatId int64) user {
	filters := getFilters(chatId)
	filters["message.forward_from"] = bson.M{"$exists": true}

	match := bson.D{{"$match", filters}}
	group := bson.D{{"$group", bson.D{
		{"_id", "$message.from.id"},
		{"count", bson.M{"$count": bson.D{}}},
		{"first_name", bson.M{"$max": "$message.from.firstname"}},
		{"username", bson.M{"$max": "$message.from.username"}},
	}}}
	sort := bson.D{{"$sort", bson.M{"count": -1}}}
	limit := bson.D{{"$limit", 1}}

	cursor, err := mongoCollection.Aggregate(context.Background(), mongo.Pipeline{match, group, sort, limit})

	var result user

	if err != nil {
		log.Println(err)
		return result
	}

	cursor.Decode(&result)

	return result
}

func getFilters(chatId int64) bson.M {
	y, m, d := time.Now().Date()
	start := time.Date(y, m, d, 0, 0, 0, 0, time.Now().Location())
	end := time.Date(y, m, d, 23, 59, 59, 0, time.Now().Location())

	return bson.M{
		"created_at": bson.M{
			"$gte": start,
			"$lte": end,
		},
		"message.chat.id":         chatId,
		"message":                 bson.M{"$ne": nil},
		"message.entities.0.type": bson.M{"$ne": "bot_command"},
	}
}

func loadLang(lang string) map[string]string {
	data, err := os.ReadFile(fmt.Sprintf("packages/ChatStatistic/lang/%s.json", lang))

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
