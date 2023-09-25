package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"log/slog"

	"github.com/fiber-bot/config"
	"github.com/fiber-bot/internal"
	"github.com/fiber-bot/user"
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

//go:embed www/fiber-bot/dist
var embedDirStatic embed.FS

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
	user.CreateTable(db, logger)
	app := initApp(envConfig, db)

	signalCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		<-signalCtx.Done()
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
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
	pathPrefix := "www/fiber-bot/dist"

	sessionStore := session.New(session.Config{
		CookieSecure: false, CookieHTTPOnly: true, CookiePath: "/",
		CookieDomain: "localhost", CookieSameSite: "Lax", KeyGenerator: utils.UUIDv4,
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
	user.InitializeUser(app, customValidator, sessionStore, db)
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(embedDirStatic),
		PathPrefix:   pathPrefix,
		MaxAge:       345600,
		NotFoundFile: "index.html",
	}))

	return app
}
