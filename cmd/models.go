package main

type App struct {
	config *Config
}

type Update struct {
	Message *Message `json:"message"`
}

type Message struct {
	Text            string `json:"text"`
	Chat            Chat   `json:"chat"`
	From            User   `json:"from"`
	MessageThreadID int64  `json:"message_thread_id,omitempty"`
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
