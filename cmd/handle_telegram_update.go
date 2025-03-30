package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func (app *App) isNewMemberJoined(message *Message) bool {
	return message != nil &&
		len(message.NewChatMembers) > 0 &&
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

func createWelcomeMessageForNewMembers(newMembers []User) string {
	// Format all user mentions
	var mentions []string
	for _, member := range newMembers {
		mentions = append(mentions, formatUserMention(&member))
	}

	// Join all mentions with commas and "and" for the last one
	var userMentions string
	if len(mentions) == 1 {
		userMentions = mentions[0]
	} else if len(mentions) == 2 {
		userMentions = mentions[0] + " и " + mentions[1]
	} else {
		userMentions = strings.Join(mentions[:len(mentions)-1], ", ") + " и " + mentions[len(mentions)-1]
	}

	return fmt.Sprintf("Привет, %s!\n\nВы пришли в мастерскую крафтового мыла «Мыльная Мама», которая специализируется на натуральной и безопасной продукции. Делаем своими руками, из своих трав и по своим рецептам.", userMentions)
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

func (app *App) buildNewMembersMessagePayload(newMembers []User) *strings.Reader {
	payload := map[string]any{
		"chat_id":      app.config.ChatID,
		"text":         createWelcomeMessageForNewMembers(newMembers),
		"reply_markup": createButtonsMarkup(&app.config.Links),
	}
	if app.config.ThreadID > 1 {
		payload["message_thread_id"] = app.config.ThreadID
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
	if app.isNewMemberJoined(update.Message) {
		url := buildSendMessageUrl(app.config.Token)
		payload := app.buildNewMembersMessagePayload(update.Message.NewChatMembers)
		sendMessage(url, payload)
	}
}
