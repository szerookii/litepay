package router

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/phuslu/log"
	"github.com/szerookii/litepay/frontend"
	"github.com/szerookii/litepay/router/middleware"
	"github.com/szerookii/litepay/router/routes/payment"
	"github.com/szerookii/litepay/utils"
	"mime"
	"regexp"
	"strings"
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

	app.All("/*", func(ctx fiber.Ctx) error {
		file := ctx.Params("*")

		re := regexp.MustCompile(`\.\w+$`)
		if file != "" && re.MatchString(file) {
			embeddedFile, err := frontend.Assets.ReadFile("build/" + file)
			if err != nil {
				return ctx.SendStatus(404)
			}

			ext := strings.Split(file, ".")[len(strings.Split(file, "."))-1]
			ctx.Set("Content-Type", mime.TypeByExtension("."+ext))

			return ctx.Send(embeddedFile)
		}

		embeddedFile, err := frontend.Assets.ReadFile("build/index.html")
		if err != nil {
			return ctx.SendStatus(404)
		}

		ctx.Set("Content-Type", "text/html")
		return ctx.Send(embeddedFile)
	})

	app.Hooks().OnListen(func(data fiber.ListenData) error {
		log.Info().Msgf("Started web server on %s:%s", data.Host, data.Port)
		return nil
	})

	return app
}
