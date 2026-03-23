package router

import (
	"calendar/internal/handlers"
	"calendar/internal/middleware"

	"github.com/gin-gonic/gin"
)

func GetRouter(handlers *handlers.EventHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware)
	api := router.Group("api/v1")
	{
		api.POST("/events", handlers.CreateEvent)
		api.PUT("/events/:id", handlers.UpdateEvent)
		api.DELETE("/events/:id", handlers.DeleteEvent)
		api.GET("/events", handlers.EventsGet)
	}
	return router
}
