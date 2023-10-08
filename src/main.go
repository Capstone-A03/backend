package main

import (
	"runtime"
	"time"

	"capstonea03/be/src/constants"
	"capstonea03/be/src/contracts"
	"capstonea03/be/src/libs/env"
	"capstonea03/be/src/libs/gracefulshutdown"
	applogger "capstonea03/be/src/libs/logger"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	logger := applogger.New("App")

	appName := env.Get(env.APP_NAME)
	appMode := env.Get(env.APP_MODE)
	appAddress := env.Get(env.APP_ADDRESS)
	webAddress := env.Get(env.WEB_ADDRESS)

	logger.Log("starting " + appName + " in " + appMode + " on " + runtime.Version())

	app := fiber.New(fiber.Config{
		AppName:     appName,
		Network:     fiber.NetworkTCP,
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		ReadTimeout: 30 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			switch err := err.(type) {
			case contracts.Error:
				return c.Status(err.Code).JSON(&contracts.Response{
					Error: &err,
				})
			case *fiber.Error:
				return c.Status(err.Code).JSON(&contracts.Response{
					Error: &contracts.Error{
						Code:    err.Code,
						Status:  err.Error(),
						Message: err.Message,
					},
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(&contracts.Response{
					Error: &contracts.Error{
						Code:    fiber.ErrInternalServerError.Code,
						Status:  fiber.ErrInternalServerError.Error(),
						Message: err.Error(),
					},
				})
			}
		},
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: appMode != constants.APP_MODE_RELEASE,
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: func() string {
			if appMode == constants.APP_MODE_RELEASE && len(webAddress) > 0 {
				return webAddress
			}
			return "*"
		}(),
	}))

	if env.Get(env.APP_MODE) != constants.APP_MODE_RELEASE {
		app.Use(fiberlogger.New())
	}

	module := module{app: app, logger: logger}
	module.load()

	gracefulshutdown.Add(gracefulshutdown.FnRunInShutdown{
		FnDescription: "shutting down app",
		Fn: func() {
			if err := app.Shutdown(); err != nil {
				logger.Error(err)
			}
		},
	})
	gracefulshutdown.Run()

	if err := app.Listen(appAddress); err != nil {
		logger.Panic(err)
	}
}
