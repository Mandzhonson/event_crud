package handlers

import (
	"calendar/internal/dto"
	"calendar/internal/models"
	"calendar/internal/service"
	"log/slog"
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
		slog.Error("failed to parse json", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := Hand.Service.CreateEvent(c.Request.Context(), event); err != nil {
		slog.Error("failed to create event", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("CreateEvent is work sucessfully", slog.Any("value", event))
	c.JSON(http.StatusCreated, gin.H{"result": event})
}

func (Hand *EventHandler) UpdateEvent(c *gin.Context) {
	var event models.Events
	id := c.Param("id")
	eventID, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("failed to cast string to int", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	event.EventID = eventID
	if err := c.ShouldBindJSON(&event); err != nil {
		slog.Error("failed to parse json", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := Hand.Service.UpdateEvent(c.Request.Context(), event); err != nil {
		slog.Error("failed to update event", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("UpdateEvent is work sucessfully", slog.Any("value", event))
	c.JSON(http.StatusOK, gin.H{"result": event})
}

func (Hand *EventHandler) DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	eventID, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("failed to cast string to int", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := Hand.Service.DeleteEvent(c.Request.Context(), eventID); err != nil {
		slog.Error("failed to delete user", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("DeleteEvent is work sucessfully", slog.Int("value", eventID))

	c.JSON(http.StatusNoContent, nil)

}

func (Hand *EventHandler) EventsGet(c *gin.Context) {
	var req dto.RequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to parse json", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.Compare(req.Period, "day") == 0 {
		eventsArr, err := Hand.Service.EventsForDay(c.Request.Context(), req.UserID, req.Date)
		if err != nil {
			slog.Error("failed to find events for day", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		slog.Debug("Events for day works successfully")
		c.JSON(http.StatusOK, eventsArr)
	} else if strings.Compare(req.Period, "week") == 0 {
		eventsArr, err := Hand.Service.EventsForWeek(c.Request.Context(), req.UserID, req.Date)
		if err != nil {
			slog.Error("failed to parse events for week", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		slog.Debug("Events for week works successfully")
		c.JSON(http.StatusOK, eventsArr)
	} else if strings.Compare(req.Period, "month") == 0 {
		eventsArr, err := Hand.Service.EventsForMonth(c.Request.Context(), req.UserID, req.Date)
		if err != nil {
			slog.Error("failed to parse events for month", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		slog.Debug("Events for month works successfully")
		c.JSON(http.StatusOK, eventsArr)
	}
}
