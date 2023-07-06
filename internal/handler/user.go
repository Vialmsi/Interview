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

func (h *Handler) RegisterUser(c *gin.Context) {
	var user entity.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		h.logger.Errorf("[RegisterUser] Error while parse user struct: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data, use correct",
		})
		return
	}

	userID, err := h.service.RegisterUser(c, user)
	if err != nil {
		h.logger.Errorf("[RegisterUser] Error while register user: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, userID)
}

func (h *Handler) LoginUser(c *gin.Context) {
	var user entity.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		h.logger.Errorf("[LoginUser] Error while parse user struct: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data, use correct",
		})
		return
	}

	loggedUser, err := h.service.LoginUser(c, user)
	if err != nil {
		h.logger.Errorf("[LoginUser] Error while login user: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.jwtService.GenerateToken(loggedUser.ID)
	if err != nil {
		h.logger.Errorf("[LoginUser] Error while generate token: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't generate token for you, try later",
		})
		return
	}

	c.JSON(http.StatusOK, token)
}
