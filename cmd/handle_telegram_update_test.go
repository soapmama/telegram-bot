package main

import (
	"io"
	"strings"
	"testing"
)

func TestIsNewMemberJoined(t *testing.T) {
	app := &App{
		config: &Config{
			ChatID: 123456789,
		},
	}

	tests := []struct {
		name     string
		message  *Message
		expected bool
	}{
		{
			name: "new member joined in correct chat",
			message: &Message{
				Chat: Chat{ID: 123456789},
				NewChatMembers: []User{
					{ID: 111222333, FirstName: "Jane", LastName: "Smith"},
				},
			},
			expected: true,
		},
		{
			name: "new member joined in wrong chat",
			message: &Message{
				Chat: Chat{ID: 999999999},
				NewChatMembers: []User{
					{ID: 111222333, FirstName: "Jane", LastName: "Smith"},
				},
			},
			expected: false,
		},
		{
			name: "no new members",
			message: &Message{
				Chat:           Chat{ID: 123456789},
				NewChatMembers: []User{},
			},
			expected: false,
		},
		{
			name:     "nil message",
			message:  nil,
			expected: false,
		},
		{
			name: "regular message without new members field",
			message: &Message{
				Chat: Chat{ID: 123456789},
				Text: "Hello world",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.isNewMemberJoined(tt.message)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFormatUserMention(t *testing.T) {
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
		{
			name: "user with username but no last name",
			user: User{
				ID:        123456789,
				FirstName: "John",
				Username:  "johndoe",
			},
			expected: "@johndoe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatUserMention(&tt.user)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestBuildSendMessageUrl(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected string
	}{
		{
			name:     "valid token",
			token:    "123456789:ABCdefGHIjklMNOpqrsTUVwxyz",
			expected: "https://api.telegram.org/bot123456789:ABCdefGHIjklMNOpqrsTUVwxyz/sendMessage",
		},
		{
			name:     "empty token",
			token:    "",
			expected: "https://api.telegram.org/bot/sendMessage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildSendMessageUrl(tt.token)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestCreateWelcomeMessageForNewMembers(t *testing.T) {
	tests := []struct {
		name       string
		newMembers []User
		expected   string
	}{
		{
			name: "single new member",
			newMembers: []User{
				{ID: 111222333, FirstName: "Jane", LastName: "Smith", Username: "janesmith"},
			},
			expected: "Привет, @janesmith!\n\nВы пришли в мастерскую крафтового мыла «Мыльная Мама», которая специализируется на натуральной и безопасной продукции. Делаем своими руками, из своих трав и по своим рецептам.",
		},
		{
			name: "two new members",
			newMembers: []User{
				{ID: 111222333, FirstName: "Jane", LastName: "Smith", Username: "janesmith"},
				{ID: 444555666, FirstName: "Bob", LastName: "Johnson", Username: "bobjohnson"},
			},
			expected: "Привет, @janesmith и @bobjohnson!\n\nВы пришли в мастерскую крафтового мыла «Мыльная Мама», которая специализируется на натуральной и безопасной продукции. Делаем своими руками, из своих трав и по своим рецептам.",
		},
		{
			name: "three new members",
			newMembers: []User{
				{ID: 111222333, FirstName: "Jane", LastName: "Smith", Username: "janesmith"},
				{ID: 444555666, FirstName: "Bob", LastName: "Johnson", Username: "bobjohnson"},
				{ID: 777888999, FirstName: "Alice", LastName: "Brown", Username: "alicebrown"},
			},
			expected: "Привет, @janesmith, @bobjohnson и @alicebrown!\n\nВы пришли в мастерскую крафтового мыла «Мыльная Мама», которая специализируется на натуральной и безопасной продукции. Делаем своими руками, из своих трав и по своим рецептам.",
		},
		{
			name: "member without username",
			newMembers: []User{
				{ID: 111222333, FirstName: "Jane", LastName: "Smith"},
			},
			expected: "Привет, Jane Smith!\n\nВы пришли в мастерскую крафтового мыла «Мыльная Мама», которая специализируется на натуральной и безопасной продукции. Делаем своими руками, из своих трав и по своим рецептам.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createWelcomeMessageForNewMembers(tt.newMembers)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestCreateButtonsMarkup(t *testing.T) {
	links := &Links{
		Distillate: "https://example.com/distillate",
		Prices:     "https://example.com/prices",
		Soap:       "https://example.com/soap",
		Ubtan:      "https://example.com/ubtan",
	}

	result := createButtonsMarkup(links)

	// Check that result is not nil
	if result == nil {
		t.Error("Expected result to be non-nil")
		return
	}

	// Check that inline_keyboard exists
	inlineKeyboard, exists := result["inline_keyboard"]
	if !exists {
		t.Error("Expected 'inline_keyboard' key to exist")
		return
	}

	// Type assert to check structure
	keyboard, ok := inlineKeyboard.([][]map[string]string)
	if !ok {
		t.Error("Expected inline_keyboard to be of type [][]map[string]string")
		return
	}

	// Check that we have 4 buttons
	if len(keyboard) != 4 {
		t.Errorf("Expected 4 buttons, got %d", len(keyboard))
	}

	// Check specific button texts
	expectedTexts := []string{
		"Как сделать заказ",
		"Что такое крафтовое мыло",
		"Что такое гидролат",
		"Что такое убтан",
	}

	for i, expectedText := range expectedTexts {
		if i >= len(keyboard) {
			t.Errorf("Button %d not found", i)
			continue
		}
		if len(keyboard[i]) != 1 {
			t.Errorf("Expected button %d to have 1 element, got %d", i, len(keyboard[i]))
			continue
		}
		if keyboard[i][0]["text"] != expectedText {
			t.Errorf("Expected button %d text '%s', got '%s'", i, expectedText, keyboard[i][0]["text"])
		}
	}
}

func TestBuildNewMembersMessagePayload(t *testing.T) {
	app := &App{
		config: &Config{
			ChatID:   123456789,
			ThreadID: 2, // Set to > 1 to trigger message_thread_id inclusion
			Links: Links{
				Distillate: "https://example.com/distillate",
				Prices:     "https://example.com/prices",
				Soap:       "https://example.com/soap",
				Ubtan:      "https://example.com/ubtan",
			},
		},
	}

	newMembers := []User{
		{ID: 111222333, FirstName: "Jane", LastName: "Smith", Username: "janesmith"},
	}

	result := app.buildNewMembersMessagePayload(newMembers)

	if result == nil {
		t.Error("Expected result to be non-nil")
		return
	}

	// Read the content to verify it's valid JSON
	contentBytes, err := io.ReadAll(result)
	if err != nil {
		t.Errorf("Error reading result: %v", err)
		return
	}
	content := string(contentBytes)

	// Check that content contains expected fields
	if !strings.Contains(content, "123456789") {
		t.Error("Expected content to contain chat_id")
	}

	if !strings.Contains(content, "Привет") {
		t.Error("Expected content to contain welcome message")
	}

	if !strings.Contains(content, "inline_keyboard") {
		t.Error("Expected content to contain inline_keyboard")
	}

	if !strings.Contains(content, "message_thread_id") {
		t.Error("Expected content to contain message_thread_id")
	}
}

func TestBuildNewMembersMessagePayloadWithoutThreadID(t *testing.T) {
	app := &App{
		config: &Config{
			ChatID:   123456789,
			ThreadID: 0, // No thread ID
			Links: Links{
				Distillate: "https://example.com/distillate",
				Prices:     "https://example.com/prices",
				Soap:       "https://example.com/soap",
				Ubtan:      "https://example.com/ubtan",
			},
		},
	}

	newMembers := []User{
		{ID: 111222333, FirstName: "Jane", LastName: "Smith", Username: "janesmith"},
	}

	result := app.buildNewMembersMessagePayload(newMembers)

	if result == nil {
		t.Error("Expected result to be non-nil")
		return
	}

	// Read the content to verify it's valid JSON
	contentBytes, err := io.ReadAll(result)
	if err != nil {
		t.Errorf("Error reading result: %v", err)
		return
	}
	content := string(contentBytes)

	// Check that content does NOT contain message_thread_id when ThreadID is 0
	if strings.Contains(content, "message_thread_id") {
		t.Error("Expected content to NOT contain message_thread_id when ThreadID is 0")
	}
}
