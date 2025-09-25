package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookHandler(t *testing.T) {
	app := &App{
		config: &Config{
			Token:    "test_token",
			Port:     "8080",
			ChatID:   123456789,
			ThreadID: 1,
			Links: Links{
				Distillate: "https://example.com/distillate",
				Prices:     "https://example.com/prices",
				Soap:       "https://example.com/soap",
				Ubtan:      "https://example.com/ubtan",
			},
		},
	}

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
	}{
		{
			name: "valid telegram update",
			requestBody: `{
				"message": {
					"text": "Hello world",
					"chat": {"id": 123456789},
					"from": {
						"id": 987654321,
						"first_name": "John",
						"last_name": "Doe",
						"username": "johndoe"
					},
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
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid JSON",
			requestBody: `{
				"message": {
					"text": "Hello world",
					"chat": {"id": 123456789},
					"from": {
						"id": 987654321,
						"first_name": "John",
						"last_name": "Doe",
						"username": "johndoe"
					},
					"new_chat_members": [
						{
							"id": 111222333,
							"first_name": "Jane",
							"last_name": "Smith",
							"username": "janesmith"
						}
					]
				}
				invalid json here
			}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest, // Empty body is invalid JSON
		},
		{
			name: "update without message",
			requestBody: `{
				"update_id": 123456789
			}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/bot", bytes.NewBufferString(tt.requestBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			app.webhookHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestWebhookHandlerWithNewMember(t *testing.T) {
	app := &App{
		config: &Config{
			Token:    "test_token",
			Port:     "8080",
			ChatID:   123456789,
			ThreadID: 1,
			Links: Links{
				Distillate: "https://example.com/distillate",
				Prices:     "https://example.com/prices",
				Soap:       "https://example.com/soap",
				Ubtan:      "https://example.com/ubtan",
			},
		},
	}

	// Test with new member joining
	requestBody := `{
		"message": {
			"text": "New member joined",
			"chat": {"id": 123456789},
			"from": {
				"id": 987654321,
				"first_name": "John",
				"last_name": "Doe",
				"username": "johndoe"
			},
			"new_chat_members": [
				{
					"id": 111222333,
					"first_name": "Jane",
					"last_name": "Smith",
					"username": "janesmith"
				}
			]
		}
	}`

	req, err := http.NewRequest("POST", "/bot", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	app.webhookHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestWebhookHandlerWithWrongChat(t *testing.T) {
	app := &App{
		config: &Config{
			Token:    "test_token",
			Port:     "8080",
			ChatID:   123456789,
			ThreadID: 1,
			Links: Links{
				Distillate: "https://example.com/distillate",
				Prices:     "https://example.com/prices",
				Soap:       "https://example.com/soap",
				Ubtan:      "https://example.com/ubtan",
			},
		},
	}

	// Test with new member joining wrong chat
	requestBody := `{
		"message": {
			"text": "New member joined",
			"chat": {"id": 999999999},
			"from": {
				"id": 987654321,
				"first_name": "John",
				"last_name": "Doe",
				"username": "johndoe"
			},
			"new_chat_members": [
				{
					"id": 111222333,
					"first_name": "Jane",
					"last_name": "Smith",
					"username": "janesmith"
				}
			]
		}
	}`

	req, err := http.NewRequest("POST", "/bot", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	app.webhookHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestRegisterRoutes(t *testing.T) {
	app := &App{
		config: &Config{
			Token:    "test_token",
			Port:     "8080",
			ChatID:   123456789,
			ThreadID: 1,
		},
	}

	// This test verifies that the route is registered
	// We can't easily test the route registration without starting a server,
	// but we can verify the handler function works
	app.registerRoutes()

	// Test that the handler is accessible
	req, err := http.NewRequest("POST", "/bot", bytes.NewBufferString(`{}`))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	app.webhookHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestWebhookHandlerMalformedJSON(t *testing.T) {
	app := &App{
		config: &Config{
			Token:    "test_token",
			Port:     "8080",
			ChatID:   123456789,
			ThreadID: 1,
		},
	}

	// Test with malformed JSON
	requestBody := `{
		"message": {
			"text": "Hello world",
			"chat": {"id": 123456789},
			"from": {
				"id": 987654321,
				"first_name": "John",
				"last_name": "Doe",
				"username": "johndoe"
			},
			"new_chat_members": [
				{
					"id": 111222333,
					"first_name": "Jane",
					"last_name": "Smith",
					"username": "janesmith"
				}
			]
		}
		invalid json here
	}`

	req, err := http.NewRequest("POST", "/bot", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	app.webhookHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestWebhookHandlerValidJSON(t *testing.T) {
	app := &App{
		config: &Config{
			Token:    "test_token",
			Port:     "8080",
			ChatID:   123456789,
			ThreadID: 1,
		},
	}

	// Test with valid JSON
	requestBody := `{
		"message": {
			"text": "Hello world",
			"chat": {"id": 123456789},
			"from": {
				"id": 987654321,
				"first_name": "John",
				"last_name": "Doe",
				"username": "johndoe"
			}
		}
	}`

	req, err := http.NewRequest("POST", "/bot", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	app.webhookHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

// Test helper function to create a test update
func createTestUpdate(chatID int64, hasNewMembers bool) string {
	update := map[string]interface{}{
		"message": map[string]interface{}{
			"text": "Test message",
			"chat": map[string]interface{}{
				"id": chatID,
			},
			"from": map[string]interface{}{
				"id":         987654321,
				"first_name": "John",
				"last_name":  "Doe",
				"username":   "johndoe",
			},
		},
	}

	if hasNewMembers {
		update["message"].(map[string]interface{})["new_chat_members"] = []map[string]interface{}{
			{
				"id":         111222333,
				"first_name": "Jane",
				"last_name":  "Smith",
				"username":   "janesmith",
			},
		}
	}

	jsonData, _ := json.Marshal(update)
	return string(jsonData)
}

func TestWebhookHandlerWithTestHelper(t *testing.T) {
	app := &App{
		config: &Config{
			Token:    "test_token",
			Port:     "8080",
			ChatID:   123456789,
			ThreadID: 1,
		},
	}

	// Test with new members
	requestBody := createTestUpdate(123456789, true)

	req, err := http.NewRequest("POST", "/bot", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	app.webhookHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Test without new members
	requestBody = createTestUpdate(123456789, false)

	req, err = http.NewRequest("POST", "/bot", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr = httptest.NewRecorder()
	app.webhookHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}
