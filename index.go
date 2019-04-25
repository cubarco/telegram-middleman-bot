package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const BASE_URL = "https://api.telegram.org/bot"
const POLL_TIMEOUT_SEC = 180

var token string
var key string
var initialized int

// model.go

type InMessage struct {
	ChatId string `json:"chat_id"`
	Text   string `json:"text"`
	Key    string `json:"key"`
}

// Only required fields are implemented
type TelegramUser struct {
	Id        int    `json:"id`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Only required fields are implemented
type TelegramChat struct {
	Id   int    `json:"id`
	Type string `json:"type"`
}

// Only required fields are implemented
type TelegramOutMessage struct {
	ChatId    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// Only required fields are implemented
type TelegramInMessage struct {
	MessageId int          `json:"message_id"`
	From      TelegramUser `json:"from"`
	Date      int          `json:"date"`
	Chat      TelegramChat `json:"chat"`
	Text      string       `json:"text"`
}

// Only required fields are implemented
type TelegramUpdate struct {
	UpdateId int               `json:"update_id"`
	Message  TelegramInMessage `json:"message"`
}

type TelegramUpdateResponse struct {
	Ok     bool             `json:"ok"`
	Result []TelegramUpdate `json:"result"`
}

type BotConfig struct {
	Token string
	Key   string
}

// main.go

func getApiUrl() string {
	return BASE_URL + token
}

func sendMessage(chatId, text string) error {
	m, err := json.Marshal(&TelegramOutMessage{ChatId: chatId, Text: text, ParseMode: "Markdown"})
	if err != nil {
		return err
	}
	reader := strings.NewReader(string(m))
	resp, err := http.Post(getApiUrl()+"/sendMessage", "application/json", reader)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func processUpdate(update TelegramUpdate) {
	var text string
	chatId := update.Message.Chat.Id
	if strings.HasPrefix(update.Message.Text, "/start") {
		text = "Here is your chatId you can use to send messages to your Telegram account:\n\n_" + strconv.Itoa(chatId) + "_"
	} else {
		text = "Please use the _/start_ command to fetch a new token."
	}
	err := sendMessage(strconv.Itoa(chatId), text)
	if err != nil {
		log.Println(err)
	}
}

func getConfig() BotConfig {
	tokenPtr := os.Getenv("TELEGRAM_BOT_TOKEN")
	key := os.Getenv("TELEGRAM_BOT_KEY")

	return BotConfig{
		Token: tokenPtr,
		Key:   key}
}

func initall() {
	if initialized == 1 {
		return
	}

	initialized = 1

	config := getConfig()
	token = config.Token
	key = config.Key
}

func Handler(w http.ResponseWriter, r *http.Request) {
	initall()
	if strings.Contains(r.URL.RequestURI(), "/api/messages") {
		messageHandler(w, r)
	} else if strings.Contains(r.URL.RequestURI(), "/api/updates") {
		webhookUpdateHandler(w, r)
	} else {
		w.WriteHeader(404)
	}
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(415)
		return
	}
	dec := json.NewDecoder(r.Body)
	var m InMessage
	err := dec.Decode(&m)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	if len(m.ChatId) == 0 || len(m.Text) == 0 || len(m.Key) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("You need to pass chat_id, key and text parameters."))
		return
	}

	if key != m.Key {
		w.WriteHeader(403)
		w.Write([]byte("Wrong key."))
		return
	}

	chatId := m.ChatId

	err = sendMessage(chatId, m.Text)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func webhookUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(415)
		return
	}
	dec := json.NewDecoder(r.Body)
	var u TelegramUpdate
	err := dec.Decode(&u)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	processUpdate(u)
	w.WriteHeader(200)
}
