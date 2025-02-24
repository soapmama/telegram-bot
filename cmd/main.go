package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type App struct {
	config *Config
}

type Update struct {
	Message *Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
	From User   `json:"from"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

func isBotMention(message *Message) bool {
	return message != nil && strings.Contains(strings.ToLower(message.Text), "ботик")
}

func buildUserMention(user *User) string {
	userMention := user.FirstName
	if user.LastName != "" {
		userMention += " " + user.LastName
	}

	notificationMention := userMention
	if user.Username != "" {
		notificationMention = "@" + user.Username
	}

	return notificationMention
}

func getSendMessageUrl(token string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
}

func getMessageText(message *Message) string {
	userMention := buildUserMention(&message.From)
	return fmt.Sprintf("Привет, %s!\n\nВы пришли в мастерскую крафтового мыла \"Мыльная Мама\", которая специализируется на натуральной и безопасной продукции. Делаем своими руками, из своих трав и по своим рецептам.", userMention)
}

func sendMessage(url string, jsonData []byte) {
	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Sent Message, Status: %v\n", resp.Status)
}

func getReplyMarkup(links *Links) [][]map[string]string {
	return [][]map[string]string{
		{
			{
				"text": "Что такое крафтовое мыло",
				"url":  links.Soap,
			},
		},
		{
			{
				"text": "Как сделать заказ",
				"url":  links.Prices,
			},
		},
		{
			{
				"text": "Что такое гидролат",
				"url":  links.Distillate,
			},
		},
	}
}

func (app *App) handleMessage(update *Update) {
	if isBotMention(update.Message) {
		chatID := update.Message.Chat.ID
		url := getSendMessageUrl(app.config.Token)
		payload := map[string]any{
			"chat_id":      chatID,
			"text":         getMessageText(update.Message),
			"reply_markup": getReplyMarkup(&app.config.Links),
		}

		jsonData, _ := json.Marshal(payload)
		sendMessage(url, jsonData)
	}
}

func (app *App) requestHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var update Update
	if err := json.Unmarshal(body, &update); err != nil {
		fmt.Printf("Error unmarshaling update: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	app.handleMessage(&update)
	w.WriteHeader(http.StatusOK)
}

func main() {
	config := newConfig()
	app := &App{
		config: config,
	}
	http.HandleFunc("/bot", app.requestHandler)
	fmt.Printf("Starting webhook server on port %s\n", config.Port)
	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
