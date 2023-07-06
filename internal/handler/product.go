package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Vialmsi/Interview/internal/entity"

	"github.com/gin-gonic/gin"
)

type ProductService interface {
	SaveProduct(ctx context.Context, product *entity.Product, userID int) (*entity.Product, error)
	RetrieveProduct(ctx context.Context, barcode string, userID int) (*entity.Product, error)
	DeleteProduct(ctx context.Context, barcode string, userID int) error
	RetrieveProductsByUserID(ctx context.Context, userID, limit, page int) ([]entity.Product, int, error)
}

func (h *Handler) SaveProduct(ctx *gin.Context) {
	var product *entity.Product

	err := ctx.ShouldBindJSON(&product)
	if err != nil {
		h.logger.Errorf("[SaveProduct] Error while parse product struct: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data, use correct",
		})
		return
	}

	ctxID, ok := ctx.Get(userIDCtx)
	if !ok {
		h.logger.Errorf("[SaveProduct] Couldn't get user id from context")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't get your credentials, try later",
		})
		return
	}

	userID, ok := ctxID.(int)
	if !ok {
		h.logger.Errorf("[SaveProduct] Couldn't parse user id to int")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't get your credentials, try later",
		})
		return
	}

	product, err = h.service.SaveProduct(ctx, product, userID)
	if err != nil {
		h.logger.Errorf("[SaveProduct] Couldnt save product: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func (h *Handler) RetrieveProduct(ctx *gin.Context) {
	productID := ctx.Param("barcode")
	if productID == "" {
		h.logger.Errorf("[RetrieveProduct] Couldn't get product id from query param")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't get product id from query param",
		})
		return
	}

	ctxID, ok := ctx.Get(userIDCtx)
	if !ok {
		h.logger.Errorf("[RetrieveProduct] Couldn't get user id from context")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't get your credentials, try later",
		})
		return
	}

	userID, ok := ctxID.(int)
	if !ok {
		h.logger.Errorf("[RetrieveProduct] Couldn't parse user id to int")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't get your credentials, try later",
		})
		return
	}

	product, err := h.service.RetrieveProduct(ctx, productID, userID)
	if err != nil {
		h.logger.Errorf("[RetrieveProduct] Error while retrieving product: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (h *Handler) DeleteProduct(ctx *gin.Context) {
	productID := ctx.Param("barcode")
	if productID == "" {
		h.logger.Errorf("[DeleteProduct] Couldn't get product id from query param")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't get product id from query param",
		})
		return
	}

	ctxID, ok := ctx.Get(userIDCtx)
	if !ok {
		h.logger.Errorf("[DeleteProduct] Couldn't get user id from context")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't get your credentials, try later",
		})
		return
	}

	userID, ok := ctxID.(int)
	if !ok {
		h.logger.Errorf("[DeleteProduct] Couldn't parse user id to int")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't get your credentials, try later",
		})
		return
	}

	err := h.service.DeleteProduct(ctx, productID, userID)
	if err != nil {
		h.logger.Errorf("[DeleteProduct] Error while deleting product: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"response": "Successfully Deleted",
	})
}

func (h *Handler) RetrieveProductsByUserID(ctx *gin.Context) {
	queryLimit := ctx.Query("limit")
	queryPage := ctx.Query("page")

	ctxID, ok := ctx.Get(userIDCtx)
	if !ok {
		h.logger.Errorf("[RetrieveProductsByUserID] Couldn't get user id from context")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't get your credentials, try later",
		})
		return
	}

	userID, ok := ctxID.(int)
	if !ok {
		h.logger.Errorf("[RetrieveProductsByUserID] Couldn't parse user id to int")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't get your credentials, try later",
		})
		return
	}

	limit, _ := strconv.Atoi(queryLimit)
	page, _ := strconv.Atoi(queryPage)

	products, count, err := h.service.RetrieveProductsByUserID(ctx, userID, limit, page)
	if err != nil {
		h.logger.Errorf("[RetrieveProductsByUserID] Couldn't get user products.sql: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products":    products,
		"total_count": count,
	})
}
