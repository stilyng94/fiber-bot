package tgbot

import (
	"context"
	"log/slog"

	"github.com/fiber-bot/config"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/uptrace/bun"
)

const errorMessage = "TeleBot error"

type BotState string

const (
	StateName        BotState = "name"
	StateDescription BotState = "description"
	StatePrice       BotState = "price"
	StateUpload      BotState = "upload"
	StateImage       BotState = "image"
)

func configBot(logger *slog.Logger) []bot.Option {
	opts := []bot.Option{
		bot.WithDebug(),
		// bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {
		// 	b.SendMessage(ctx, &bot.SendMessageParams{
		// 		ChatID: update.Message.Chat.ID,
		// 		Text:   "Unknown command. Type /help to see all commands.",
		// 	})
		// }),
		bot.WithErrorsHandler(func(err error) {
			logger.Error(errorMessage, slog.String("error", err.Error()))
		}),
	}
	return opts
}

func ConfigBot(logger *slog.Logger, envConfig config.EnvConfig, db *bun.DB) (*bot.Bot, error) {
	botLogger := logger.With(slog.Attr{Key: "context", Value: slog.StringValue("TeleBot")})
	opts := configBot(botLogger)
	b, err := bot.New(envConfig.TelegramBotToken, opts...)
	if err != nil {
		botLogger.Error(errorMessage, slog.Any("error", err))
		return nil, err
	}
	_, err = b.SetMyCommands(context.Background(), &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "/help", Description: "Show all commands"},
			{Command: "/start", Description: "Start bot"},
		},
	})
	if err != nil {
		botLogger.Error(errorMessage, slog.Any("error", err))
		return nil, err
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, onHelp)
	b.RegisterHandler(bot.HandlerTypeMessageText, "pre_checkout_query", bot.MatchTypeExact, onCheckout)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		onStart(ctx, b, update, envConfig.AppUrl)
	})
	b.RegisterHandler(bot.HandlerTypeMessageText, "/successful_payment", bot.MatchTypeExact, onPayment)

	return b, nil
}
