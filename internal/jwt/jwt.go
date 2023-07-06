package jwt

import (
	"fmt"
	"strconv"

	"github.com/Vialmsi/Interview/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type TokenClient interface {
	Generate(userID int) (string, error)
	Validate(token string) error
}

type JWTService struct {
	logger       *logrus.Logger
	tokenService TokenClient
	signingKey   string
}

func NewJWTService(logger *logrus.Logger, tokenService TokenClient, tokenCredentials config.TokenCredentials) *JWTService {
	return &JWTService{
		logger:       logger,
		tokenService: tokenService,
		signingKey:   tokenCredentials.SigningKey,
	}
}

func (j *JWTService) GenerateToken(userID int) (string, error) {
	j.logger.Info("[GenerateToken] started")

	token, err := j.tokenService.Generate(userID)
	if err != nil {
		j.logger.Errorf("[GenerateToken] Error while generating token: %s", err)
		return "", err
	}

	j.logger.Info(token)
	j.logger.Info("[GenerateToken] ended")

	return token, nil
}

func (j *JWTService) ValidateToken(token string) error {
	j.logger.Info("[ValidateToken] started")

	err := j.tokenService.Validate(token)
	if err != nil {
		j.logger.Errorf("[ValidateToken] Error while validating token")
		return err
	}

	j.logger.Info("[ValidateToken] ended")

	return nil
}

func (j *JWTService) ParseToken(tokenString string) (int, error) {
	j.logger.Info("[ParseToken] started")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("[ParseToken] Unexpected signing method: %v", token.Header["alg"])
		}

		return j.signingKey, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		j.logger.Errorf("[ParseToken] Error while getting userID from claims: %s", err.Error())
		return 0, fmt.Errorf("auth error, try later")
	}

	userID, err := strconv.Atoi(claims["login"].(string))
	if err != nil {
		j.logger.Errorf("[ParseToken] Error while parse userID to int: %s", err.Error())
		return 0, fmt.Errorf("auth error, try later")
	}

	j.logger.Info(userID)
	j.logger.Info("[ParseToken] ended")

	return userID, nil
}
