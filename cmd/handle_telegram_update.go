package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func (app *App) containsKeyword(message *Message) bool {
	return message != nil &&
		strings.Contains(strings.ToLower(message.Text), "ботик") &&
		message.Chat.ID == app.config.ChatID
}

func formatUserMention(user *User) string {
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

func buildSendMessageUrl(token string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
}

func createWelcomeMessage(message *Message) string {
	userMention := formatUserMention(&message.From)
	return fmt.Sprintf("Привет, %s!\n\nВы пришли в мастерскую крафтового мыла «Мыльная Мама», которая специализируется на натуральной и безопасной продукции. Делаем своими руками, из своих трав и по своим рецептам.", userMention)
}

func createButtonsMarkup(links *Links) map[string]any {
	return map[string]any{
		"inline_keyboard": [][]map[string]string{
			{
				{
					"text": "Как сделать заказ",
					"url":  links.Prices,
				},
			},
			{
				{
					"text": "Что такое крафтовое мыло",
					"url":  links.Soap,
				},
			},
			{
				{
					"text": "Что такое гидролат",
					"url":  links.Distillate,
				},
			},
			{
				{
					"text": "Что такое убтан",
					"url":  links.Ubtan,
				},
			},
		},
	}
}

func (app *App) buildMessagePayload(message *Message) *strings.Reader {
	payload := map[string]any{
		"chat_id":           app.config.ChatID,
		"text":              createWelcomeMessage(message),
		"reply_markup":      createButtonsMarkup(&app.config.Links),
		"message_thread_id": app.config.ThreadID,
	}
	jsonData, _ := json.Marshal(payload)
	return strings.NewReader(string(jsonData))
}

func sendMessage(url string, payload *strings.Reader) {
	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		slog.Error("Error sending message", "error", err)
		return
	}
	defer resp.Body.Close()
	slog.Info("Sent message", "status", resp.Status)
}

func (app *App) handleTelegramUpdate(update *Update) {
	if app.containsKeyword(update.Message) {
		url := buildSendMessageUrl(app.config.Token)
		payload := app.buildMessagePayload(update.Message)
		sendMessage(url, payload)
	}
}
