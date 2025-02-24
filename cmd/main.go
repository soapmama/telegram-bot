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
}

type Chat struct {
	ID int64 `json:"id"`
}

func (app *App) handleChannelPost(update *Update) {
	if update.Message != nil && strings.Contains(update.Message.Text, "ботик") {
		chatID := update.Message.Chat.ID

		// Send message using HTTP POST request
		url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", app.config.Token)
		payload := map[string]interface{}{
			"chat_id": chatID,
			"text":    "Привет!\n\nВы пришли в мастерскую крафтового мыла \"Мыльная Мама\", которая специализируется на натуральной и безопасной продукции. Делаем своими руками, из своих трав и по своим рецептам.",
			"reply_markup": map[string]interface{}{
				"inline_keyboard": [][]map[string]string{
					{
						{
							"text": "Что такое крафтовое мыло",
							"url":  app.config.Links.Soap,
						},
					},
					{
						{
							"text": "Как сделать заказ",
							"url":  app.config.Links.Prices,
						},
					},
					{
						{
							"text": "Что такое гидролат",
							"url":  app.config.Links.Distillate,
						},
					},
				},
			},
		}

		jsonData, _ := json.Marshal(payload)
		resp, _ := http.Post(url, "application/json", strings.NewReader(string(jsonData)))
		defer resp.Body.Close()

		fmt.Printf("Sent Message, Status: %v\n", resp.Status)
	}
}

func main() {
	config := newConfig()

	app := &App{
		config: config,
	}

	if config.GoEnv == "development" {
		fmt.Println("Development mode not implemented for standard library version")
		os.Exit(1)
	} else if config.GoEnv == "production" {
		http.HandleFunc("/bot", func(w http.ResponseWriter, r *http.Request) {
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

			app.handleChannelPost(&update)
			w.WriteHeader(http.StatusOK)
		})

		fmt.Printf("Starting webhook server on port %s\n", config.Port)
		err := http.ListenAndServe(":"+config.Port, nil)
		if err != nil {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
	}
}
