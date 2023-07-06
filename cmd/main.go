package main

import (
	"github.com/Vialmsi/Interview/internal/clients/tokensvc"
	"github.com/Vialmsi/Interview/internal/config"
	"github.com/Vialmsi/Interview/internal/handler"
	"github.com/Vialmsi/Interview/internal/jwt"
	"github.com/Vialmsi/Interview/internal/pdfsvc"
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
	}

	repo, err := store.NewStore(cfg.PSQLDatabase)
	if err != nil {
		logger.Fatalf("Error while init store: %s", err.Error())
	}

	svc := service.NewService(logger, repo)

	tokenService := tokensvc.NewTokenService(logger, cfg.TokenServiceConfig)
	err = tokenService.Ping()
	if err != nil {
		logger.Fatalf("Error while init token service: %s", err)
	}

	pdfService, err := pdfsvc.NewPDFService(logger, repo)
	if err != nil {
		logger.Fatalf("Error while init pdf service: %s", err)
	}

	jwtService := jwt.NewJWTService(logger, tokenService, cfg.TokenCredentials)

	hdlr := handler.NewHandler(logger, svc, jwtService, pdfService)

	server := gin.New()

	hdlr.Mount(server)

	err = server.Run(":" + cfg.Server.Port)
	if err != nil {
		logger.Fatalf("Couldn't run server: %s", err)
	}
}
