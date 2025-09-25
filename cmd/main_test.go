package main

import (
	"testing"
)

func TestAppStruct(t *testing.T) {
	config := &Config{
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
	}

	app := &App{
		config: config,
	}

	// Test that app has config
	if app.config == nil {
		t.Error("Expected app.config to be non-nil")
	}

	// Test that config values are accessible
	if app.config.Token == "" {
		t.Error("Expected app.config.Token to be non-empty")
	}

	if app.config.Port == "" {
		t.Error("Expected app.config.Port to be non-empty")
	}

	if app.config.ChatID == 0 {
		t.Error("Expected app.config.ChatID to be non-zero")
	}

	// Test that links are accessible
	if app.config.Links.Distillate == "" {
		t.Error("Expected app.config.Links.Distillate to be non-empty")
	}

	if app.config.Links.Prices == "" {
		t.Error("Expected app.config.Links.Prices to be non-empty")
	}

	if app.config.Links.Soap == "" {
		t.Error("Expected app.config.Links.Soap to be non-empty")
	}

	if app.config.Links.Ubtan == "" {
		t.Error("Expected app.config.Links.Ubtan to be non-empty")
	}
}
