package handlers

import (
	"calendar/internal/dto"
	"calendar/internal/models"
	"calendar/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	Service service.EventService
}

func NewEventHandler(s service.EventService) *EventHandler {
	return &EventHandler{
		Service: s,
	}
}

func (Hand *EventHandler) CreateEvent(c *gin.Context) {
	var event models.Events
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := Hand.Service.CreateEvent(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"result": event})
}

func (Hand *EventHandler) UpdateEvent(c *gin.Context) {
	var event models.Events
	id := c.Param("id")
	eventID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event.EventID = eventID
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := Hand.Service.UpdateEvent(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": event})
}

func (Hand *EventHandler) DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	eventID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := Hand.Service.DeleteEvent(c.Request.Context(), eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)

}

func (Hand *EventHandler) EventsGet(c *gin.Context) {
	var req dto.RequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.Compare(req.Period, "day") == 0 {
		eventsArr, err := Hand.Service.EventsForDay(c.Request.Context(), req.UserID, req.Date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, eventsArr)
	} else if strings.Compare(req.Period, "week") == 0 {
		eventsArr, err := Hand.Service.EventsForWeek(c.Request.Context(), req.UserID, req.Date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, eventsArr)
	} else if strings.Compare(req.Period, "month") == 0 {
		eventsArr, err := Hand.Service.EventsForMonth(c.Request.Context(), req.UserID, req.Date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, eventsArr)
	}
}

// func (Hand *EventHandler) EventsForWeek(c *gin.Context) {
// }

// func (Hand *EventHandler) EventsForMonth(c *gin.Context) {
// }
