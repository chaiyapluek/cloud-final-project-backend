package main

import (
	"fmt"
	"time"

	"dev.chaiyapluek.cloud.final.backend/src/config"
	"dev.chaiyapluek.cloud.final.backend/src/handler"
	myMiddleware "dev.chaiyapluek.cloud.final.backend/src/middleware"
	"dev.chaiyapluek.cloud.final.backend/src/pkg/database"
	"dev.chaiyapluek.cloud.final.backend/src/repository"
	"dev.chaiyapluek.cloud.final.backend/src/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	mail "github.com/xhit/go-simple-mail/v2"
)

func smtpServer(cfg *config.SMTPConfig) *mail.SMTPServer {
	mailServer := mail.NewSMTPClient()
	mailServer.Host = cfg.Host
	mailServer.Port = cfg.Port
	mailServer.ConnectTimeout = time.Duration(cfg.ConnectionTimeout) * time.Second
	mailServer.SendTimeout = time.Duration(cfg.SendTimeout) * time.Second
	switch cfg.Encryption {
	case "SSLTLS":
		mailServer.Encryption = mail.EncryptionSSLTLS
	case "SSL":
		mailServer.Encryption = mail.EncryptionSSL
	case "STARTTLS":
		mailServer.Encryption = mail.EncryptionSTARTTLS
	default:
		mailServer.Encryption = mail.EncryptionNone
	}
	if cfg.Auth {
		mailServer.Username = cfg.Username
		mailServer.Password = cfg.Password
	}
	return mailServer
}

func main() {
	cfg := config.LoadEnv()
	e := echo.New()

	conn := database.GetMongoConnection(cfg.DB.Host)

	emailRepo := repository.NewEmailRepository(conn, cfg.DB.DBName, cfg.DB.Collection)

	authRepo := repository.NewAuthRepository(conn, cfg.DB.DBName, cfg.DB.Collection)
	userRepo := repository.NewUserRepository(conn, cfg.DB.DBName, cfg.DB.Collection)
	locationRepo := repository.NewLocationRepository(conn, cfg.DB.DBName, cfg.DB.Collection)
	cartRepo := repository.NewCartRepository(conn, cfg.DB.DBName, cfg.DB.Collection)

	mailServer := smtpServer(cfg.SMTP)
	mailService := service.NewEmailService(emailRepo, cfg.SMTP.MaxSendPerDay, mailServer, cfg.SMTP.Sender)
	authService := service.NewAuthService(mailService, authRepo, userRepo)
	locationService := service.NewLocationService(locationRepo)
	cartService := service.NewCartService(cartRepo, userRepo, mailService)

	authHandler := handler.NewAuthHandler(authService)
	locationHandler := handler.NewLocationHandler(locationService)
	cartHandler := handler.NewCartHandler(cartService)

	e.Use(middleware.Logger())
	e.Use(myMiddleware.ErrorHandlerMiddleware)

	authRoute := e.Group("auth")
	authRoute.POST("/login", authHandler.Login)
	authRoute.POST("/login-attempt", authHandler.LoginAttempt)
	authRoute.POST("/register", authHandler.Register)
	authRoute.POST("/register-attempt", authHandler.RegisterAttempt)

	locationRoute := e.Group("locations")
	locationRoute.GET("", locationHandler.GetAllLocation)
	locationRoute.GET("/:id/menus", locationHandler.GetLocationById)
	locationRoute.GET("/:locationId/menus/:menuId", locationHandler.GetLocationMenu)

	userReoute := e.Group("users")
	userReoute.GET("/:userId", authHandler.GetMe)
	userReoute.GET("/:userId/carts", cartHandler.GetCartByUserId)

	cartRoute := e.Group("carts")
	cartRoute.POST("", cartHandler.CreateCart)
	cartRoute.POST("/:cartId/items", cartHandler.AddCartItem)
	cartRoute.DELETE("/:cartId/items/:itemId", cartHandler.DeleteCartItem)

	e.POST("/checkout", cartHandler.Checkout)

	e.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
}
