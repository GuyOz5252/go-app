package pkg

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	consoleHandler := slog.NewTextHandler(os.Stdout, nil)
	return slog.New(consoleHandler);
}
