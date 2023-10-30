package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/fiber-bot/bootstraps"
	"github.com/fiber-bot/cmd/frontend"
	"github.com/fiber-bot/config"
	"github.com/fiber-bot/internal"
	"github.com/fiber-bot/models"
	"github.com/fiber-bot/tele_bot"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/uptrace/bun"
)

func main() {
	envConfig, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalln(err)
	}

	logger := config.InitLogger(envConfig)

	db, err := config.InitDB(envConfig)
	if err != nil {
		logger.Error("Database connection error", slog.Any("error", err))
		panic("database error")
	}

	if envConfig.IsDev() {
		models := []interface{}{
			&models.UserModel{}, &models.OrderModel{}, &models.ProductModel{},
		}
		config.ResetTable(db, logger, models...)
	}

	app := initApp(envConfig, db)

	bot, botUpdater, err := tele_bot.ConnectBot(envConfig, logger, db, false)
	if err == nil {
		go func() {
			err = botUpdater.StartPolling(bot, &ext.PollingOpts{DropPendingUpdates: true})
			if err != nil {
				panic("failed to start polling: " + err.Error())
			}

			botUpdater.Idle()
		}()
	}

	signalCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		<-signalCtx.Done()
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := botUpdater.Stop(); err != nil {
			logger.Error("Bot Shutdown", slog.String("error", err.Error()))
		}
		if err := app.ShutdownWithContext(signalCtx); err != nil {
			logger.Error("Server Shutdown", slog.String("error", err.Error()))
		}
	}()

	if err := app.Listen(fmt.Sprintf("0.0.0.0:%v", envConfig.Port)); err != nil && err != http.ErrServerClosed {
		logger.Error(fmt.Sprintf("listen: %s\n", err))
	}
}

func initApp(envConfig config.EnvConfig, db *bun.DB) *fiber.App {
	customValidator := internal.NewXValidator()

	sessionStore := session.New(session.Config{
		CookieSecure: false, CookieHTTPOnly: true, CookiePath: "/",
		CookieDomain: envConfig.Domain, CookieSameSite: "Lax", KeyGenerator: utils.UUIDv4,
	})

	app := fiber.New(fiber.Config{
		EnableTrustedProxyCheck: true, EnablePrintRoutes: envConfig.IsDev(),
		PassLocalsToViews:     true,
		DisableStartupMessage: envConfig.IsProd(),
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
	})
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(logger.New(logger.Config{
		TimeZone: "UTC",
	}))
	app.Use(encryptcookie.New(encryptcookie.Config{Key: envConfig.CookieSecret}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, DELETE, PUT, PATCH",
		AllowHeaders:     "Origin,ACCEPT,AUTHORIZATION,CONTENT-TYPE,X-CSRF-TOKEN",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           int((1 * time.Hour).Seconds()),
	}))
	app.Use(compress.New(compress.Config{Level: compress.LevelDefault}))
	app.Use(requestid.New(requestid.Config{Generator: utils.UUIDv4}))
	if envConfig.IsProd() {
		app.Use(limiter.New(limiter.Config{
			LimiterMiddleware: limiter.SlidingWindow{},
		}))
	}
	app.Use(etag.New())
	app.Use(helmet.New(helmet.Config{}))
	bootstraps.BootstrapUser(app, customValidator, sessionStore, db, envConfig)
	bootstraps.BootstrapOrder(app, customValidator, sessionStore, db, envConfig)
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         frontend.BuildHTTPFS(),
		MaxAge:       345600,
		NotFoundFile: "index.html",
	}))

	return app
}
