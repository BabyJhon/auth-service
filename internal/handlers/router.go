package handlers

import (
	"github.com/BabyJhon/auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api/v1")
	{
		api.POST("/recieve", h.recieve)
		api.POST("/refresh", h.refresh)
		api.POST("/zalupa", h.zalupa)
	}

	return router
}
