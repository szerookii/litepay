package router

import (
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"
	"github.com/szerookii/litepay/backend/router/middleware"
	routeauth "github.com/szerookii/litepay/backend/router/routes/auth"
	"github.com/szerookii/litepay/backend/router/routes/config"
	"github.com/szerookii/litepay/backend/router/routes/payment"
	routeuser "github.com/szerookii/litepay/backend/router/routes/user"
	"github.com/szerookii/litepay/backend/utils"
	"github.com/szerookii/litepay/frontend"
)

func Init() *gin.Engine {
	log.Info().Msg("Initializing router...")

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173"
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(allowedOrigins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           86400,
	}))

	api := r.Group("/api")
	api.Use(middleware.GlobalLimiter.Middleware())
	api.GET("/status", func(c *gin.Context) {
		utils.SendJSON(c, http.StatusOK, gin.H{"status": "ok"})
	})
	api.GET("/config", config.Get)
	api.GET("/payment/:id", payment.Get)

	authGroup := api.Group("/auth")
	authGroup.Use(middleware.AuthLimiter.Middleware())
	authGroup.POST("/register", routeauth.Register)
	authGroup.POST("/login", routeauth.Login)
	authGroup.POST("/logout", routeauth.Logout)

	protected := api.Group("")
	protected.Use(middleware.APIKey)
	protected.POST("/payment", payment.Post)

	userGroup := api.Group("/user")
	userGroup.Use(middleware.JWT)
	userGroup.GET("/me", routeuser.Me)
	userGroup.GET("/payments", routeuser.Payments)
	userGroup.GET("/balance", routeuser.Balance)
	userGroup.PUT("/wallets", routeuser.UpdateWallets)
	userGroup.POST("/api-key", routeuser.RegenerateKey)
	userGroup.POST("/webhook-secret", routeuser.RotateWebhookSecret)
	userGroup.POST("/cashout", routeuser.Cashout)

	paymentGroup := api.Group("/payment")
	paymentGroup.Use(middleware.JWT)
	paymentGroup.POST("/:id/refund", payment.Refund)

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
			return
		}

		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		data, err := frontend.Assets.ReadFile("build/" + path)
		if err == nil {
			contentType := mime.TypeByExtension(path[strings.LastIndex(path, "."):])
			if contentType == "" {
				contentType = "application/octet-stream"
			}
			c.Data(http.StatusOK, contentType, data)
			return
		}

		data, err = frontend.Assets.ReadFile("build/index.html")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "text/html", data)
	})

	return r
}
