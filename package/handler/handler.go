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

	auth := router.Group("/user")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/reset-password", h.resetPassword)
		auth.GET("/verify-confirmcode", h.verifyCode)
		auth.PUT("/set-password", h.setPassword)
		auth.PUT("/change-password", h.changePassword)
	}

	api := router.Group("/api", h.userIdentity)
	{
		trips := api.Group("/trips")
		{
			trips.POST("/", h.createTrip)
			trips.GET("/", h.getAllTrips)
			trips.GET("/:id", h.getTripById)
			trips.DELETE("/:id", h.deleteTrip)

			points := trips.Group(":id/points")
			{
				points.POST("/", h.createPoint)
				points.GET("/", h.getAllPoints)
				points.GET("/:point_id", h.getPointById)
			}
		}
	}

	return router
}