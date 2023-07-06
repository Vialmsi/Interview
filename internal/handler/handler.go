package handler

import (
	"github.com/Vialmsi/Interview/internal/handler/middleware"
	"github.com/Vialmsi/Interview/internal/jwt"
	"github.com/Vialmsi/Interview/internal/pdfsvc"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	userIDCtx = "userID"
)

type Service interface {
	UserService
	ProductService
}

type Handler struct {
	logger     *logrus.Logger
	service    Service
	jwtService *jwt.JWTService
	pdfService PDFService
}

func NewHandler(logger *logrus.Logger, service Service, jwtService *jwt.JWTService, pdfService *pdfsvc.PDFService) *Handler {
	return &Handler{
		logger:     logger,
		service:    service,
		jwtService: jwtService,
		pdfService: pdfService,
	}
}

func (h *Handler) Mount(r *gin.Engine) {
	api := r.Group("/api")

	authRoutes := api.Group("/auth")
	authRoutes.POST("/signup", h.RegisterUser)
	authRoutes.POST("/signin", h.LoginUser)

	authMiddleware := middleware.NewAuthMiddleware(h.logger, h.jwtService)

	productRoutes := api.Group("/products", authMiddleware.UserIdentity)
	productRoutes.GET("/:barcode", h.RetrieveProduct)
	productRoutes.DELETE("/:barcode", h.DeleteProduct)
	productRoutes.POST("", h.SaveProduct)
	productRoutes.GET("", h.RetrieveProductsByUserID)

	pdfRoutes := api.Group("/prices", authMiddleware.UserIdentity)
	pdfRoutes.GET("/:barcode", h.GetPdfFromBarcode)
	pdfRoutes.GET("", h.GetPdfFromName)
}
