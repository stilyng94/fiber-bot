package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"
)

func InitLogger(envConfig EnvConfig) *slog.Logger {
	var loggerHandler slog.Handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	if envConfig.IsProd() {
		logPath := fmt.Sprintf("logs/%s.log", time.Now().Format(time.DateOnly))
		if _, err := os.Stat(logPath); err != nil {
			if _, err := os.Create(logPath); err != nil {
				log.Fatalln(err)
			}
		}
		out, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		loggerHandler = slog.NewJSONHandler(out, &slog.HandlerOptions{AddSource: true})
	}
	logger := slog.New(loggerHandler).With(slog.String("version", envConfig.Version), slog.String("env", envConfig.Environment))
	return logger
}

//TODO: create log dir on fly
