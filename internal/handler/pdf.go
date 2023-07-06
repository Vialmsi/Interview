package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/Vialmsi/Interview/internal/entity"

	"github.com/gin-gonic/gin"
)

var (
	PDFNotExistError = errors.New("couldn't find file with this name")
)

type PDFService interface {
	LoadPDFFromBarcode(userID int, barcode string) (string, error)
	LoadPDFFromName(fileName string) (string, error)
	GeneratePDF(ctx context.Context, userID int, product entity.Product) (string, error)
}

func (h *Handler) GetPdfFromBarcode(ctx *gin.Context) {
	productID := ctx.Param("barcode")
	if productID == "" {
		h.logger.Errorf("[GetPdf] Couldn't get product id from query param")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't get product id from query param",
		})
		return
	}

	ctxID, ok := ctx.Get(userIDCtx)
	if !ok {
		h.logger.Errorf("[GetPdf] Couldn't get user id from context")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't get your credentials, try later",
		})
		return
	}

	userID, ok := ctxID.(int)
	if !ok {
		h.logger.Errorf("[GetPdf] Couldn't parse user id to int")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't get your credentials, try later",
		})
		return
	}

	loadedPDF, err := h.pdfService.LoadPDFFromBarcode(userID, productID)
	if err != nil {
		h.logger.Warnf("[GetPdf] Error while load pdf attempt: %s", err.Error())
		if !errors.As(err, &PDFNotExistError) {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		product, err := h.service.RetrieveProduct(ctx, productID, userID)
		if err != nil {
			h.logger.Errorf("[GetPdf] Error while retrieving product: %s", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		generatedPDF, err := h.pdfService.GeneratePDF(ctx, userID, *product)
		if err != nil {
			h.logger.Errorf("[GetPdf] Error while generating pdf: %s", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.File(generatedPDF)
		return
	}

	ctx.File(loadedPDF)
}

func (h *Handler) GetPdfFromName(ctx *gin.Context) {
	path := ctx.Query("path")

	pdf, err := h.pdfService.LoadPDFFromName(path)
	if err != nil {
		h.logger.Errorf("[GetPdfFromName] Error while getting pdf: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.File(pdf)
}
