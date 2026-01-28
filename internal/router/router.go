package router

import (
	"calendar/internal/handlers"

	"github.com/gin-gonic/gin"
)

func GetRouter(handlers *handlers.EventHandler) *gin.Engine {
	router := gin.Default()
	router.POST("/events", handlers.CreateEvent)
	router.PUT("/events", handlers.UpdateEvent)
	router.DELETE("/events", handlers.DeleteEvent)
	router.GET("/events/day", handlers.EventsForDay)
	router.GET("/events/week", handlers.EventsForWeek)
	router.GET("/events/month", handlers.EventsForMonth)
	return router
}
