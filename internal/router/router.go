package router

import (
	"calendar/internal/handlers"

	"github.com/gin-gonic/gin"
)

func GetRouter(handlers *handlers.EventHandler) *gin.Engine {
	router := gin.Default()
	router.POST("/events", handlers.CreateEvent)
	router.PUT("/events/:id", handlers.UpdateEvent)
	router.DELETE("/events/:id", handlers.DeleteEvent)
	router.GET("/events", handlers.EventsGet)
	// router.GET("/events/week", handlers.EventsForWeek)
	// router.GET("/events/month", handlers.EventsForMonth)
	return router
}
