package handlers_test

import (
	"bytes"
	"calendar/internal/apperr"
	"calendar/internal/dto"
	"calendar/internal/handlers"
	"calendar/internal/models"
	mock_service "calendar/internal/service/mocks"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandlers_CreateEvent(t *testing.T) {
	type mockBehaivour func(s *mock_service.MockEventService, event dto.RequestDTO)
	testTable := []struct {
		name                string
		inputBody           string
		inputEvent          dto.RequestDTO
		mockBehaivour       mockBehaivour
		expectedStatusCode  int
		expectedBodyRequest string
	}{
		{
			name: "OK",
			inputBody: `{
			"user_id": 2,
			"event": "meeting",
			"event_date": "2026-03-10"
		}`,
			inputEvent: dto.RequestDTO{
				UserID: 2,
				Event:  "meeting",
				Date:   "2026-03-10",
			},
			mockBehaivour: func(s *mock_service.MockEventService, event dto.RequestDTO) {
				s.EXPECT().
					CreateEvent(gomock.Any(), event).
					Return(10, nil)
			},
			expectedStatusCode:  http.StatusCreated,
			expectedBodyRequest: `{"result":{"event_id":10,"user_id":2,"event_date":"2026-03-10","event":"meeting"}}`,
		},
		{
			name: "invalid json",
			inputBody: `{
			"user_id": 2,
			"event": "meeting",
		}`,
			mockBehaivour:       func(s *mock_service.MockEventService, event dto.RequestDTO) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedBodyRequest: `{"error":"invalid request parameters"}`,
		},
		{
			name: "service invalid params",
			inputBody: `{
			"user_id": 0,
			"event": "",
			"event_date": ""
		}`,
			inputEvent: dto.RequestDTO{},
			mockBehaivour: func(s *mock_service.MockEventService, event dto.RequestDTO) {
				s.EXPECT().
					CreateEvent(gomock.Any(), gomock.Any()).
					Return(0, apperr.InvalidReqParams)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedBodyRequest: `{"error":"invalid request parameters"}`,
		},
		{
			name: "service internal error",
			inputBody: `{
			"user_id": 2,
			"event": "meeting",
			"event_date": "2026-03-10"
		}`,
			inputEvent: dto.RequestDTO{
				UserID: 2,
				Event:  "meeting",
				Date:   "2026-03-10",
			},
			mockBehaivour: func(s *mock_service.MockEventService, event dto.RequestDTO) {
				s.EXPECT().
					CreateEvent(gomock.Any(), event).
					Return(0, apperr.InternalServErr)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedBodyRequest: `{"error":"internal server error"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := mock_service.NewMockEventService(ctrl)
			testCase.mockBehaivour(serv, testCase.inputEvent)

			handler := &handlers.EventHandler{
				Service: serv,
			}

			r := gin.New()
			gin.SetMode(gin.TestMode)
			r.POST("/events", handler.CreateEvent)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedBodyRequest)
		})
	}
}

func TestHandlers_UpdateEvent(t *testing.T) {
	type mockBehaivour func(s *mock_service.MockEventService, event dto.RequestDTO)

	testTable := []struct {
		name                string
		inputID             string
		inputBody           string
		inputEvent          dto.RequestDTO
		mockBehaivour       mockBehaivour
		expectedStatusCode  int
		expectedBodyRequest string
	}{
		{
			name:    "OK",
			inputID: "1",
			inputBody: `{
				"user_id": 2,
				"event": "meeting",
				"event_date": "2026-03-10"
			}`,
			inputEvent: dto.RequestDTO{
				EventID: 1,
				UserID:  2,
				Event:   "meeting",
				Date:    "2026-03-10",
			},
			mockBehaivour: func(s *mock_service.MockEventService, event dto.RequestDTO) {
				s.EXPECT().UpdateEvent(gomock.Any(), event).Return(nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedBodyRequest: `{"result":{"event_id":1,"user_id":2,"event_date":"2026-03-10","event":"meeting"}}`,
		},
		{
			name:                "invalid id",
			inputID:             "abc",
			inputBody:           `{}`,
			mockBehaivour:       func(s *mock_service.MockEventService, event dto.RequestDTO) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedBodyRequest: `{"error":"bad request"}`,
		},
		{
			name:    "invalid json",
			inputID: "1",
			inputBody: `{
				"user_id": 2,
				"event": "meeting",
			}`,
			mockBehaivour:       func(s *mock_service.MockEventService, event dto.RequestDTO) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedBodyRequest: `{"error":"invalid request parameters"}`,
		},
		{
			name:    "service invalid params",
			inputID: "1",
			inputBody: `{
				"user_id": 0,
				"event": "",
				"event_date": ""
			}`,
			inputEvent: dto.RequestDTO{
				EventID: 1,
			},
			mockBehaivour: func(s *mock_service.MockEventService, event dto.RequestDTO) {
				s.EXPECT().UpdateEvent(gomock.Any(), event).Return(apperr.InvalidReqParams)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedBodyRequest: `{"error":"invalid request parameters"}`,
		},
		{
			name:    "event not found",
			inputID: "1",
			inputBody: `{
				"user_id": 2,
				"event": "meeting",
				"event_date": "2026-03-10"
			}`,
			inputEvent: dto.RequestDTO{
				EventID: 1,
				UserID:  2,
				Event:   "meeting",
				Date:    "2026-03-10",
			},
			mockBehaivour: func(s *mock_service.MockEventService, event dto.RequestDTO) {
				s.EXPECT().UpdateEvent(gomock.Any(), event).Return(apperr.EventNotFound)
			},
			expectedStatusCode:  http.StatusNotFound,
			expectedBodyRequest: `{"error":"event not found"}`,
		},
		{
			name:    "service internal error",
			inputID: "1",
			inputBody: `{
				"user_id": 2,
				"event": "meeting",
				"event_date": "2026-03-10"
			}`,
			inputEvent: dto.RequestDTO{
				EventID: 1,
				UserID:  2,
				Event:   "meeting",
				Date:    "2026-03-10",
			},
			mockBehaivour: func(s *mock_service.MockEventService, event dto.RequestDTO) {
				s.EXPECT().UpdateEvent(gomock.Any(), event).Return(apperr.InternalServErr)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedBodyRequest: `{"error":"internal server error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			serv := mock_service.NewMockEventService(ctrl)
			testCase.mockBehaivour(serv, testCase.inputEvent)

			handlers := &handlers.EventHandler{
				Service: serv,
			}

			r := gin.New()
			gin.SetMode(gin.TestMode)

			r.PUT("/events/:id", handlers.UpdateEvent)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/events/%s", testCase.inputID), bytes.NewBufferString(testCase.inputBody))
			r.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedBodyRequest, w.Body.String())
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestHandlers_DeleteEvent(t *testing.T) {
	type mockBehaivour func(s *mock_service.MockEventService, id int)

	testTable := []struct {
		name                string
		inputID             string
		mockBehaivour       mockBehaivour
		expectedStatusCode  int
		expectedBodyRequest string
	}{
		{
			name:    "OK",
			inputID: "1",
			mockBehaivour: func(s *mock_service.MockEventService, eventID int) {
				s.EXPECT().DeleteEvent(gomock.Any(), eventID).Return(nil)
			},
			expectedStatusCode:  http.StatusNoContent,
			expectedBodyRequest: "",
		},
		{
			name:                "invalid id",
			inputID:             "abc",
			mockBehaivour:       func(s *mock_service.MockEventService, eventID int) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedBodyRequest: `{"error":"bad request"}`,
		},
		{
			name:    "event not found",
			inputID: "1",
			mockBehaivour: func(s *mock_service.MockEventService, eventID int) {
				s.EXPECT().DeleteEvent(gomock.Any(), eventID).Return(apperr.EventNotFound)
			},
			expectedStatusCode:  http.StatusNotFound,
			expectedBodyRequest: `{"error":"event not found"}`,
		},
		{
			name:    "service internal error",
			inputID: "1",
			mockBehaivour: func(s *mock_service.MockEventService, eventID int) {
				s.EXPECT().DeleteEvent(gomock.Any(), eventID).Return(apperr.InternalServErr)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedBodyRequest: `{"error":"internal server error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv := mock_service.NewMockEventService(ctrl)
			id, _ := strconv.Atoi(testCase.inputID)

			testCase.mockBehaivour(srv, id)
			handlers := &handlers.EventHandler{
				Service: srv,
			}

			r := gin.New()
			gin.SetMode(gin.TestMode)

			r.DELETE("/events/:id", handlers.DeleteEvent)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/events/%s", testCase.inputID), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedBodyRequest)
		})
	}
}

func TestHandlers_EventsGet(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockEventService, userID int, date string)

	testTable := []struct {
		name                string
		queryParams         string
		inputUserID         int
		inputDate           string
		mockBehaviour       mockBehaviour
		expectedStatusCode  int
		expectedBodyRequest string
	}{
		{
			name:        "OK day",
			queryParams: "?period=day&user_id=1&event_date=2024-01-01",
			inputUserID: 1,
			inputDate:   "2024-01-01",
			mockBehaviour: func(s *mock_service.MockEventService, userID int, date string) {
				s.EXPECT().
					EventsForDay(gomock.Any(), userID, date).
					Return([]models.Events{
						{
							EventID: 1,
							UserID:  1,
							Date:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
							Event:   "meeting",
						},
					}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedBodyRequest: `[{"event_id":1,"user_id":1,"event_date":"2024-01-01T00:00:00Z","event":"meeting"}]`,
		},
		{
			name:        "OK week",
			queryParams: "?period=week&user_id=1&event_date=2024-01-01",
			inputUserID: 1,
			inputDate:   "2024-01-01",
			mockBehaviour: func(s *mock_service.MockEventService, userID int, date string) {
				s.EXPECT().
					EventsForWeek(gomock.Any(), userID, date).
					Return([]models.Events{
						{
							EventID: 1,
							UserID:  1,
							Date:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
							Event:   "meeting",
						},
					}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedBodyRequest: `[{"event_id":1,"user_id":1,"event_date":"2024-01-01T00:00:00Z","event":"meeting"}]`,
		},
		{
			name:        "OK month",
			queryParams: "?period=month&user_id=1&event_date=2024-01-01",
			inputUserID: 1,
			inputDate:   "2024-01-01",
			mockBehaviour: func(s *mock_service.MockEventService, userID int, date string) {
				s.EXPECT().
					EventsForMonth(gomock.Any(), userID, date).
					Return([]models.Events{
						{
							EventID: 1,
							UserID:  1,
							Date:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
							Event:   "meeting",
						},
					}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedBodyRequest: `[{"event_id":1,"user_id":1,"event_date":"2024-01-01T00:00:00Z","event":"meeting"}]`,
		},
		{
			name:                "invalid query params",
			queryParams:         "?period=day&user_id=abc&event_date=2024-01-01",
			inputUserID:         0,
			inputDate:           "",
			mockBehaviour:       func(s *mock_service.MockEventService, userID int, date string) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedBodyRequest: `{"error":"bad request"}`,
		},
		{
			name:                "unknown period",
			queryParams:         "?period=year&user_id=1&event_date=2024-01-01",
			inputUserID:         1,
			inputDate:           "2024-01-01",
			mockBehaviour:       func(s *mock_service.MockEventService, userID int, date string) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedBodyRequest: `{"error":"bad request"}`,
		},
		{
			name:        "event not found",
			queryParams: "?period=day&user_id=1&event_date=2024-01-01",
			inputUserID: 1,
			inputDate:   "2024-01-01",
			mockBehaviour: func(s *mock_service.MockEventService, userID int, date string) {
				s.EXPECT().
					EventsForDay(gomock.Any(), userID, date).
					Return(nil, apperr.EventNotFound)
			},
			expectedStatusCode:  http.StatusNotFound,
			expectedBodyRequest: `{"error":"event not found"}`,
		},
		{
			name:        "invalid request params from service",
			queryParams: "?period=day&user_id=1&event_date=2024-01-01",
			inputUserID: 1,
			inputDate:   "2024-01-01",
			mockBehaviour: func(s *mock_service.MockEventService, userID int, date string) {
				s.EXPECT().
					EventsForDay(gomock.Any(), userID, date).
					Return(nil, apperr.InvalidReqParams)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedBodyRequest: `{"error":"bad request"}`,
		},
		{
			name:        "service internal error",
			queryParams: "?period=day&user_id=1&event_date=2024-01-01",
			inputUserID: 1,
			inputDate:   "2024-01-01",
			mockBehaviour: func(s *mock_service.MockEventService, userID int, date string) {
				s.EXPECT().
					EventsForDay(gomock.Any(), userID, date).
					Return(nil, apperr.InternalServErr)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedBodyRequest: `{"error":"internal server error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv := mock_service.NewMockEventService(ctrl)
			testCase.mockBehaviour(srv, testCase.inputUserID, testCase.inputDate)

			handler := &handlers.EventHandler{
				Service: srv,
			}

			r := gin.New()
			gin.SetMode(gin.TestMode)
			r.GET("/events", handler.EventsGet)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/events"+testCase.queryParams, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedBodyRequest, w.Body.String())
		})
	}
}
