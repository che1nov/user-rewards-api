package main

import (
	"log/slog"
	"os"

	"user-rewards-api/internal/app"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	application, err := app.NewApp()
	if err != nil {
		slog.Error("Ошибка инициализации приложения", "error", err)
		os.Exit(1)
	}
	defer application.Close()

	if err := application.Run(); err != nil {
		slog.Error("Ошибка запуска приложения", "error", err)
		os.Exit(1)
	}
}
