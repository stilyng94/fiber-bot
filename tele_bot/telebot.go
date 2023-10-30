package tele_bot

import (
	"fmt"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/fiber-bot/config"
	"github.com/fiber-bot/repos"
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

func ConfigBot(envConfig config.EnvConfig, offline bool) (*gotgbot.Bot, error) {
	bot, err := gotgbot.NewBot(envConfig.TelegramBotToken, &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{UseTestEnvironment: offline},
	})
	if err != nil {
		return nil, err
	}
	return bot, nil
}

func ConnectBot(envConfig config.EnvConfig, logger *slog.Logger, db *bun.DB, offline bool) (*gotgbot.Bot, *ext.Updater, error) {
	botLogger := logger.With(slog.Attr{Key: "context", Value: slog.StringValue("TeleBot")})
	bot, err := ConfigBot(envConfig, offline)
	if err != nil {
		botLogger.Error(errorMessage, slog.Any("error", err))
		return nil, nil, err
	}

	_, err = bot.SetMyCommands([]gotgbot.BotCommand{
		{Command: "/help", Description: "Show all commands"},
		{Command: "/start", Description: "Start bot"},
	},
		nil)
	if err != nil {
		botLogger.Error(errorMessage, slog.Any("error", err))
		return nil, nil, err
	}

	productRepo := repos.NewProductRepo(db)

	updater := ext.NewUpdater(&ext.UpdaterOpts{
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			// If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				botLogger.Error(errorMessage, slog.String("error", err.Error()))
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		}),
		UnhandledErrFunc: func(err error) {
			botLogger.Error(errorMessage, slog.Any("error", err))
		},
	})

	logger.Info("BOT up and running")
	dispatcher := updater.Dispatcher
	dispatcher.AddHandler(handlers.NewCommand(gotgbot.UpdateTypePreCheckoutQuery, onCheckout))
	dispatcher.AddHandler(handlers.NewMessage(message.Equal(gotgbot.UpdateTypePreCheckoutQuery), onCheckout))
	dispatcher.AddHandler(handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		return onStart(b, ctx, envConfig.AppUrl)
	}))
	dispatcher.AddHandler(handlers.NewCommand("help", onHelp))
	dispatcher.AddHandler(handlers.NewCommand("successful_payment", onPayment))
	dispatcher.AddHandler(handlers.NewCommand("payment", onPayment))
	dispatcher.AddHandler(handlers.NewMessage(message.Equal("successful_payment"), onPayment))

	dispatcher.AddHandler(handlers.NewConversation(
		[]ext.Handler{handlers.NewCommand("uploadProduct", func(b *gotgbot.Bot, ctx *ext.Context) error {
			return initAddProduct(b, ctx, envConfig)
		})},
		map[string][]ext.Handler{
			string(StateName):        {handlers.NewMessage(noCommands, productName)},
			string(StateDescription): {handlers.NewMessage(noCommands, productDescription)},
			string(StateImage): {handlers.NewMessage(func(msg *gotgbot.Message) bool {
				return message.Photo(msg) && !message.Text(msg) && !message.Command(msg)
			}, productImage)},
			string(StatePrice): {handlers.NewMessage(noCommands, productPrice)},
			string(StateUpload): {handlers.NewMessage(noCommands, func(b *gotgbot.Bot, ctx *ext.Context) error {
				return productUpload(b, ctx, productRepo, envConfig)
			})},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewCommand("cancel", cancel), handlers.NewCommand("restart", restart)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
			AllowReEntry: true,
		},
	))

	dispatcher.AddHandler(handlers.NewMessage(message.All, func(b *gotgbot.Bot, ctx *ext.Context) error {
		_, err := ctx.EffectiveMessage.Reply(b, "Unknown command. Type /help to see all commands.", nil)
		if err != nil {
			return fmt.Errorf("failed handle unknown commands: %w", err)
		}
		return nil
	}))

	return bot, updater, nil
}
