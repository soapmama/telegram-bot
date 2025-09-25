package main

import (
	"testing"
)

func TestConfigStruct(t *testing.T) {
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

	if config.Token == "" {
		t.Error("Token should not be empty")
	}

	if config.Port == "" {
		t.Error("Port should not be empty")
	}

	if config.ChatID == 0 {
		t.Error("ChatID should not be zero")
	}

	if config.Links.Distillate == "" {
		t.Error("Distillate link should not be empty")
	}
}

func TestLinksStruct(t *testing.T) {
	links := Links{
		Distillate: "https://example.com/distillate",
		Prices:     "https://example.com/prices",
		Soap:       "https://example.com/soap",
		Ubtan:      "https://example.com/ubtan",
	}

	if links.Distillate == "" {
		t.Error("Distillate link should not be empty")
	}

	if links.Prices == "" {
		t.Error("Prices link should not be empty")
	}

	if links.Soap == "" {
		t.Error("Soap link should not be empty")
	}

	if links.Ubtan == "" {
		t.Error("Ubtan link should not be empty")
	}
}
