package main

import (
	"fmt"

	"github.com/Vialmsi/Interview/internal/clients/token_service"
	"github.com/Vialmsi/Interview/internal/config"
	"github.com/Vialmsi/Interview/internal/handler"
	"github.com/Vialmsi/Interview/internal/jwt"
	"github.com/Vialmsi/Interview/internal/service"
	"github.com/Vialmsi/Interview/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	cfg, err := config.Init()
	if err != nil {
		logger.Fatalf("Error while load config: %s", err.Error())
		return
	}

	userStore, err := store.NewUserStore(cfg.PSQLDatabase)
	if err != nil {
		logger.Fatalf("Error while init user store: %s", err)
	}
	productStore, err := store.NewProductStore()
	if err != nil {
		logger.Fatalf("Error while init product store: %s", err)
	}

	svc := service.NewService(logger, productStore, userStore)

	tokenService := token_service.NewTokenService(logger, cfg.TokenServiceConfig)
	err = tokenService.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}

	jwtService := jwt.NewJWTService(logger, tokenService, cfg.TokenCredentials)

	hdlr := handler.NewHandler(logger, svc, jwtService)

	server := gin.New()

	hdlr.Mount(server)

	err = server.Run(":" + cfg.Server.Port)
	if err != nil {
		logger.Fatalf("Couldn't run server: %s", err)
	}
}
