package handlers_test

import (
	"bytes"
	"calendar/internal/apperr"
	"calendar/internal/dto"
	"calendar/internal/handlers"
	mock_service "calendar/internal/mocks"
	"net/http"
	"net/http/httptest"
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
