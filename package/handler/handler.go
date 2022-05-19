package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/folins/biketrackserver/package/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/get-email", h.getEmail)
		auth.POST("/check-code", h.checkConfirmCode)
	}

	return router
}