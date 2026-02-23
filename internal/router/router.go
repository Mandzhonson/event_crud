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
	router.POST("/events", handlers.CreateEvent)
	router.PUT("/events/:id", handlers.UpdateEvent)
	router.DELETE("/events/:id", handlers.DeleteEvent)
	router.GET("/events", handlers.EventsGet)
	return router
}
