package router

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/phuslu/log"
	"github.com/szerookii/litepay/router/middleware"
	"github.com/szerookii/litepay/router/routes/payment"
	"github.com/szerookii/litepay/utils"
)

func Init() *fiber.App {
	log.Info().Msg("Initializing router...")

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			return utils.SendJSON(ctx, 500, fiber.Map{
				"message": err.Error(),
			})
		},
		JSONDecoder: sonic.Unmarshal,
		JSONEncoder: sonic.Marshal,
	})

	app.Use(cors.New())
	//app.Use(limiter.New(limiter.Config{}))

	api := app.Group("/api")
	api.Get("/status", func(ctx fiber.Ctx) error {
		return utils.SendJSON(ctx, 200, fiber.Map{
			"status": "ok",
		})
	})

	api.Get("/payment/:id", payment.Get)
	api.Use(middleware.APIKey).Post("/payment", payment.Post)

	app.Use(static.New("frontend/build"))

	app.All("/*", func(ctx fiber.Ctx) error {
		return ctx.SendFile("frontend/build/index.html")
	})

	app.Hooks().OnListen(func(data fiber.ListenData) error {
		log.Info().Msgf("Started web server on %s:%s", data.Host, data.Port)
		return nil
	})

	return app
}
