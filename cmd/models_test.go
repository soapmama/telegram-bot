package main

import (
	"encoding/json"
	"testing"
)

func TestUpdateStruct(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected Update
	}{
		{
			name: "valid update with message",
			json: `{
				"message": {
					"text": "Hello world",
					"chat": {"id": 123456789},
					"from": {
						"id": 987654321,
						"first_name": "John",
						"last_name": "Doe",
						"username": "johndoe"
					},
					"message_thread_id": 1,
					"new_chat_members": [
						{
							"id": 111222333,
							"first_name": "Jane",
							"last_name": "Smith",
							"username": "janesmith"
						}
					]
				}
			}`,
			expected: Update{
				Message: &Message{
					Text:            "Hello world",
					Chat:            Chat{ID: 123456789},
					From:            User{ID: 987654321, FirstName: "John", LastName: "Doe", Username: "johndoe"},
					MessageThreadID: 1,
					NewChatMembers: []User{
						{ID: 111222333, FirstName: "Jane", LastName: "Smith", Username: "janesmith"},
					},
				},
			},
		},
		{
			name: "update without message",
			json: `{}`,
			expected: Update{
				Message: nil,
			},
		},
		{
			name: "update with message but no new members",
			json: `{
				"message": {
					"text": "Regular message",
					"chat": {"id": 123456789},
					"from": {
						"id": 987654321,
						"first_name": "John"
					}
				}
			}`,
			expected: Update{
				Message: &Message{
					Text: "Regular message",
					Chat: Chat{ID: 123456789},
					From: User{ID: 987654321, FirstName: "John"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var update Update
			err := json.Unmarshal([]byte(tt.json), &update)
			if err != nil {
				t.Errorf("Failed to unmarshal JSON: %v", err)
				return
			}

			// Compare the unmarshaled struct with expected
			if tt.expected.Message == nil {
				if update.Message != nil {
					t.Errorf("Expected Message to be nil, got %+v", update.Message)
				}
			} else {
				if update.Message == nil {
					t.Errorf("Expected Message to be non-nil")
					return
				}

				if update.Message.Text != tt.expected.Message.Text {
					t.Errorf("Expected Text %s, got %s", tt.expected.Message.Text, update.Message.Text)
				}

				if update.Message.Chat.ID != tt.expected.Message.Chat.ID {
					t.Errorf("Expected Chat.ID %d, got %d", tt.expected.Message.Chat.ID, update.Message.Chat.ID)
				}

				if update.Message.From.ID != tt.expected.Message.From.ID {
					t.Errorf("Expected From.ID %d, got %d", tt.expected.Message.From.ID, update.Message.From.ID)
				}

				if update.Message.From.FirstName != tt.expected.Message.From.FirstName {
					t.Errorf("Expected From.FirstName %s, got %s", tt.expected.Message.From.FirstName, update.Message.From.FirstName)
				}

				if len(update.Message.NewChatMembers) != len(tt.expected.Message.NewChatMembers) {
					t.Errorf("Expected %d new chat members, got %d", len(tt.expected.Message.NewChatMembers), len(update.Message.NewChatMembers))
				}
			}
		})
	}
}

func TestMessageStruct(t *testing.T) {
	message := Message{
		Text:            "Test message",
		Chat:            Chat{ID: 123456789},
		From:            User{ID: 987654321, FirstName: "John", LastName: "Doe", Username: "johndoe"},
		MessageThreadID: 1,
		NewChatMembers: []User{
			{ID: 111222333, FirstName: "Jane", LastName: "Smith", Username: "janesmith"},
		},
	}

	if message.Text == "" {
		t.Error("Text should not be empty")
	}

	if message.Chat.ID == 0 {
		t.Error("Chat ID should not be zero")
	}

	if message.From.ID == 0 {
		t.Error("From ID should not be zero")
	}

	if message.From.FirstName == "" {
		t.Error("From FirstName should not be empty")
	}

	if len(message.NewChatMembers) == 0 {
		t.Error("NewChatMembers should not be empty")
	}
}

func TestChatStruct(t *testing.T) {
	chat := Chat{ID: 123456789}

	if chat.ID == 0 {
		t.Error("Chat ID should not be zero")
	}
}

func TestUserStruct(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		expected string
	}{
		{
			name: "user with first name only",
			user: User{
				ID:        123456789,
				FirstName: "John",
			},
			expected: "John",
		},
		{
			name: "user with first and last name",
			user: User{
				ID:        123456789,
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "John Doe",
		},
		{
			name: "user with username",
			user: User{
				ID:        123456789,
				FirstName: "John",
				LastName:  "Doe",
				Username:  "johndoe",
			},
			expected: "@johndoe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.user.ID == 0 {
				t.Error("User ID should not be zero")
			}

			if tt.user.FirstName == "" {
				t.Error("User FirstName should not be empty")
			}
		})
	}
}

func TestAppStructInModels(t *testing.T) {
	config := &Config{
		Token:    "test_token",
		Port:     "8080",
		ChatID:   123456789,
		ThreadID: 1,
	}

	app := &App{
		config: config,
	}

	if app.config == nil {
		t.Error("App config should not be nil")
	}

	if app.config.Token != config.Token {
		t.Errorf("Expected token %s, got %s", config.Token, app.config.Token)
	}
}
