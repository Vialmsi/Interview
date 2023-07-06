package handler

import (
	"context"
	"net/http"

	"github.com/Vialmsi/Interview/internal/entity"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	RegisterUser(ctx context.Context, user entity.User) (int, error)
	LoginUser(ctx context.Context, user entity.User) (*entity.User, error)
}

func (h *Handler) RegisterUser(ctx *gin.Context) {
	var user entity.User

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		h.logger.Errorf("[RegisterUser] Error while parse user struct: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data, use correct",
		})
		return
	}

	userID, err := h.service.RegisterUser(ctx, user)
	if err != nil {
		h.logger.Errorf("[RegisterUser] Error while register user: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, userID)
}

func (h *Handler) LoginUser(ctx *gin.Context) {
	var user entity.User

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		h.logger.Errorf("[LoginUser] Error while parse user struct: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data, use correct",
		})
		return
	}

	loggedUser, err := h.service.LoginUser(ctx, user)
	if err != nil {
		h.logger.Errorf("[LoginUser] Error while login user: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.jwtService.GenerateToken(loggedUser.ID)
	if err != nil {
		h.logger.Errorf("[LoginUser] Error while generate token: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't generate token for you, try later",
		})
		return
	}

	ctx.JSON(http.StatusOK, token)
}
