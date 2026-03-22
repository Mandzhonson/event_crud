package handlers_test

import (
	"bytes"
	"calendar/internal/apperr"
	"calendar/internal/dto"
	"calendar/internal/handlers"
	mock_service "calendar/internal/mocks"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

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
