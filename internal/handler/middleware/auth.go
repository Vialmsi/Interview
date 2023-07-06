package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	userIDCtx = "userID"
)

var (
	TokenExpiredErr = errors.New("token is expired, pls login again")
)

type TokenService interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) error
	ParseToken(tokenString string) (int, error)
}

type AuthMiddleware struct {
	logger       *logrus.Logger
	tokenService TokenService
}

func NewAuthMiddleware(logger *logrus.Logger, tokenService TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		logger:       logger,
		tokenService: tokenService,
	}
}

func (a *AuthMiddleware) UserIdentity(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		a.logger.Error("[UserIdentity]: auth header is empty")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "auth header is empty"})
		ctx.Abort()
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		a.logger.Error("[UserIdentity] invalid auth header")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid auth header"})
		ctx.Abort()
		return
	}

	stringToken := headerParts[1]

	err := a.tokenService.ValidateToken(stringToken)
	if err != nil {
		fmt.Println(err)
		fmt.Println(TokenExpiredErr)
		a.logger.Errorf("[UserIdentity] error while validating token %s", err)
		if errors.As(err, &TokenExpiredErr) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "token is expired, pls login again"})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while process request"})
		ctx.Abort()
		return
	}

	userID, err := a.tokenService.ParseToken(stringToken)
	if err != nil {
		a.logger.Error("[UserIdentity] couldn't parse token")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't parse token"})
		ctx.Abort()
		return
	}

	ctx.Set(userIDCtx, userID)
}
