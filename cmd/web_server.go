package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func (app *App) registerRoutes() {
	http.HandleFunc("/bot", app.webhookHandler)
}

func (app *App) startServer() {
	slog.Info("Starting webhook server", "port", app.config.Port)
	err := http.ListenAndServe(":"+app.config.Port, nil)
	if err != nil {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}
}

func (app *App) webhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Error reading request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var update Update
	if err := json.Unmarshal(body, &update); err != nil {
		slog.Error("Error unmarshaling update", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	app.handleTelegramUpdate(&update)
	w.WriteHeader(http.StatusOK)
}
