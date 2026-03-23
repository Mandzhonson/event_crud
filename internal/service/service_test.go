package service_test

import (
	"calendar/internal/apperr"
	"calendar/internal/dto"
	"calendar/internal/models"
	mock_repo "calendar/internal/repository/mocks"
	"calendar/internal/service"
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateEvent(t *testing.T) {
	type mockBehaviour func(r *mock_repo.MockRepo, event models.Events)

	testTable := []struct {
		name          string
		inputDTO      dto.RequestDTO
		inputModel    models.Events
		mockBehaviour mockBehaviour
		expectedID    int
		expectedErr   error
	}{
		{
			name: "OK",
			inputDTO: dto.RequestDTO{
				UserID: 1,
				Event:  "meeting",
				Date:   "2024-01-01",
			},
			inputModel: models.Events{
				UserID: 1,
				Event:  "meeting",
				Date:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			mockBehaviour: func(r *mock_repo.MockRepo, event models.Events) {
				r.EXPECT().
					CreateEvent(gomock.Any(), event).
					Return(1, nil)
			},
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name: "empty event",
			inputDTO: dto.RequestDTO{
				UserID: 1,
				Event:  "",
				Date:   "2024-01-01",
			},
			mockBehaviour: func(r *mock_repo.MockRepo, event models.Events) {},
			expectedID:    0,
			expectedErr:   apperr.InvalidReqParams,
		},
		{
			name: "empty date",
			inputDTO: dto.RequestDTO{
				UserID: 1,
				Event:  "meeting",
				Date:   "",
			},
			mockBehaviour: func(r *mock_repo.MockRepo, event models.Events) {},
			expectedID:    0,
			expectedErr:   apperr.InvalidReqParams,
		},
		{
			name: "invalid userID",
			inputDTO: dto.RequestDTO{
				UserID: 0,
				Event:  "meeting",
				Date:   "2024-01-01",
			},
			mockBehaviour: func(r *mock_repo.MockRepo, event models.Events) {},
			expectedID:    0,
			expectedErr:   apperr.InvalidReqParams,
		},
		{
			name: "invalid date format",
			inputDTO: dto.RequestDTO{
				UserID: 1,
				Event:  "meeting",
				Date:   "01-01-2024",
			},
			mockBehaviour: func(r *mock_repo.MockRepo, event models.Events) {},
			expectedID:    0,
			expectedErr:   apperr.InvalidReqParams,
		},
		{
			name: "repo internal error",
			inputDTO: dto.RequestDTO{
				UserID: 1,
				Event:  "meeting",
				Date:   "2024-01-01",
			},
			inputModel: models.Events{
				UserID: 1,
				Event:  "meeting",
				Date:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			mockBehaviour: func(r *mock_repo.MockRepo, event models.Events) {
				r.EXPECT().
					CreateEvent(gomock.Any(), event).
					Return(0, apperr.InternalServErr)
			},
			expectedID:  0,
			expectedErr: apperr.InternalServErr,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repo.NewMockRepo(ctrl)
			testCase.mockBehaviour(repo, testCase.inputModel)

			svc := service.NewEventService(repo)
			id, err := svc.CreateEvent(context.Background(), testCase.inputDTO)

			assert.Equal(t, testCase.expectedID, id)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestService_UpdateEvent(t *testing.T) {
	type mockBehaviour func(r *mock_repo.MockRepo, event models.Events)

	testTable := []struct {
		name          string
		inputDTO      dto.RequestDTO
		inputModel    models.Events
		mockBehaviour mockBehaviour
		expectedErr   error
	}{
		{
			name: "OK",
			inputDTO: dto.RequestDTO{
				EventID: 1,
				Event:   "meeting",
				Date:    "2024-01-01",
			},
			inputModel: models.Events{
				EventID: 1,
				Event:   "meeting",
				Date:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			mockBehaviour: func(r *mock_repo.MockRepo, event models.Events) {
				r.EXPECT().
					FindEvents(gomock.Any(), event.EventID).
					Return(nil)
				r.EXPECT().
					UpdateEvent(gomock.Any(), event).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "empty event and date",
			inputDTO: dto.RequestDTO{
				EventID: 1,
				Event:   "",
				Date:    "",
			},
			mockBehaviour: func(r *mock_repo.MockRepo, event models.Events) {},
			expectedErr:   apperr.InternalServErr,
		},
		{
			name: "event not found",
			inputDTO: dto.RequestDTO{
				EventID: 1,
				Event:   "meeting",
				Date:    "2024-01-01",
			},
			mockBehaviour: func(r *mock_repo.MockRepo, event models.Events) {
				r.EXPECT().
					FindEvents(gomock.Any(), 1).
					Return(pgx.ErrNoRows)
			},
			expectedErr: apperr.EventNotFound,
		},
		{
			name: "invalid date format",
			inputDTO: dto.RequestDTO{
				EventID: 1,
				Event:   "meeting",
				Date:    "01-01-2024",
			},
			mockBehaviour: func(r *mock_repo.MockRepo, event models.Events) {
				r.EXPECT().
					FindEvents(gomock.Any(), 1).
					Return(nil)
			},
			expectedErr: apperr.InvalidReqParams,
		},
		{
			name: "repo internal error on update",
			inputDTO: dto.RequestDTO{
				EventID: 1,
				Event:   "meeting",
				Date:    "2024-01-01",
			},
			inputModel: models.Events{
				EventID: 1,
				Event:   "meeting",
				Date:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			mockBehaviour: func(r *mock_repo.MockRepo, event models.Events) {
				r.EXPECT().
					FindEvents(gomock.Any(), event.EventID).
					Return(nil)
				r.EXPECT().
					UpdateEvent(gomock.Any(), event).
					Return(apperr.InternalServErr)
			},
			expectedErr: apperr.InternalServErr,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repo.NewMockRepo(ctrl)
			testCase.mockBehaviour(repo, testCase.inputModel)

			svc := service.NewEventService(repo)
			err := svc.UpdateEvent(context.Background(), testCase.inputDTO)

			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestService_DeleteEvent(t *testing.T) {
	type mockBehaviour func(r *mock_repo.MockRepo, id int)

	testTable := []struct {
		name          string
		inputID       int
		mockBehaviour mockBehaviour
		expectedErr   error
	}{
		{
			name:    "OK",
			inputID: 1,
			mockBehaviour: func(r *mock_repo.MockRepo, id int) {
				r.EXPECT().
					DeleteEvent(gomock.Any(), id).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "repo internal error",
			inputID: 1,
			mockBehaviour: func(r *mock_repo.MockRepo, id int) {
				r.EXPECT().
					DeleteEvent(gomock.Any(), id).
					Return(apperr.InternalServErr)
			},
			expectedErr: apperr.InternalServErr,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repo.NewMockRepo(ctrl)
			testCase.mockBehaviour(repo, testCase.inputID)

			svc := service.NewEventService(repo)
			err := svc.DeleteEvent(context.Background(), testCase.inputID)

			if testCase.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, testCase.expectedErr.Error())
			}
		})
	}
}
func TestService_EventsForDay(t *testing.T) {
	type mockBehaviour func(r *mock_repo.MockRepo, userID int, date time.Time)

	testTable := []struct {
		name            string
		inputUserID     int
		inputDate       string
		inputDateParsed time.Time
		mockBehaviour   mockBehaviour
		expectedEvents  []models.Events
		expectedErr     error
	}{
		{
			name:            "OK",
			inputUserID:     1,
			inputDate:       "2024-01-01",
			inputDateParsed: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			mockBehaviour: func(r *mock_repo.MockRepo, userID int, date time.Time) {
				r.EXPECT().
					EventsForDay(gomock.Any(), userID, date).
					Return([]models.Events{{EventID: 1, UserID: 1, Event: "meeting", Date: date}}, nil)
			},
			expectedEvents: []models.Events{{EventID: 1, UserID: 1, Event: "meeting", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}},
			expectedErr:    nil,
		},
		{
			name:           "invalid date format",
			inputUserID:    1,
			inputDate:      "01-01-2024",
			mockBehaviour:  func(r *mock_repo.MockRepo, userID int, date time.Time) {},
			expectedEvents: nil,
			expectedErr:    apperr.InvalidReqParams,
		},
		{
			name:            "repo internal error",
			inputUserID:     1,
			inputDate:       "2024-01-01",
			inputDateParsed: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			mockBehaviour: func(r *mock_repo.MockRepo, userID int, date time.Time) {
				r.EXPECT().
					EventsForDay(gomock.Any(), userID, date).
					Return(nil, apperr.InternalServErr)
			},
			expectedEvents: nil,
			expectedErr:    apperr.InternalServErr,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repo.NewMockRepo(ctrl)
			testCase.mockBehaviour(repo, testCase.inputUserID, testCase.inputDateParsed)

			svc := service.NewEventService(repo)
			events, err := svc.EventsForDay(context.Background(), testCase.inputUserID, testCase.inputDate)

			assert.Equal(t, testCase.expectedEvents, events)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestService_EventsForWeek(t *testing.T) {
	type mockBehaviour func(r *mock_repo.MockRepo, userID int, date time.Time)

	testTable := []struct {
		name            string
		inputUserID     int
		inputDate       string
		inputDateParsed time.Time
		mockBehaviour   mockBehaviour
		expectedEvents  []models.Events
		expectedErr     error
	}{
		{
			name:            "OK",
			inputUserID:     1,
			inputDate:       "2024-01-01",
			inputDateParsed: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			mockBehaviour: func(r *mock_repo.MockRepo, userID int, date time.Time) {
				r.EXPECT().
					EventsForWeek(gomock.Any(), userID, date).
					Return([]models.Events{{EventID: 1, UserID: 1, Event: "meeting", Date: date}}, nil)
			},
			expectedEvents: []models.Events{{EventID: 1, UserID: 1, Event: "meeting", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}},
			expectedErr:    nil,
		},
		{
			name:           "invalid date format",
			inputUserID:    1,
			inputDate:      "01-01-2024",
			mockBehaviour:  func(r *mock_repo.MockRepo, userID int, date time.Time) {},
			expectedEvents: nil,
			expectedErr:    apperr.InvalidReqParams,
		},
		{
			name:            "repo internal error",
			inputUserID:     1,
			inputDate:       "2024-01-01",
			inputDateParsed: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			mockBehaviour: func(r *mock_repo.MockRepo, userID int, date time.Time) {
				r.EXPECT().
					EventsForWeek(gomock.Any(), userID, date).
					Return(nil, apperr.InternalServErr)
			},
			expectedEvents: nil,
			expectedErr:    apperr.InternalServErr,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repo.NewMockRepo(ctrl)
			testCase.mockBehaviour(repo, testCase.inputUserID, testCase.inputDateParsed)

			svc := service.NewEventService(repo)
			events, err := svc.EventsForWeek(context.Background(), testCase.inputUserID, testCase.inputDate)

			assert.Equal(t, testCase.expectedEvents, events)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestService_EventsForMonth(t *testing.T) {
	type mockBehaviour func(r *mock_repo.MockRepo, userID int, date time.Time)

	testTable := []struct {
		name            string
		inputUserID     int
		inputDate       string
		inputDateParsed time.Time
		mockBehaviour   mockBehaviour
		expectedEvents  []models.Events
		expectedErr     error
	}{
		{
			name:            "OK",
			inputUserID:     1,
			inputDate:       "2024-01-01",
			inputDateParsed: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			mockBehaviour: func(r *mock_repo.MockRepo, userID int, date time.Time) {
				r.EXPECT().
					EventsForMonth(gomock.Any(), userID, date).
					Return([]models.Events{{EventID: 1, UserID: 1, Event: "meeting", Date: date}}, nil)
			},
			expectedEvents: []models.Events{{EventID: 1, UserID: 1, Event: "meeting", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}},
			expectedErr:    nil,
		},
		{
			name:           "invalid date format",
			inputUserID:    1,
			inputDate:      "01-01-2024",
			mockBehaviour:  func(r *mock_repo.MockRepo, userID int, date time.Time) {},
			expectedEvents: nil,
			expectedErr:    apperr.InvalidReqParams,
		},
		{
			name:            "repo internal error",
			inputUserID:     1,
			inputDate:       "2024-01-01",
			inputDateParsed: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			mockBehaviour: func(r *mock_repo.MockRepo, userID int, date time.Time) {
				r.EXPECT().
					EventsForMonth(gomock.Any(), userID, date).
					Return(nil, apperr.InternalServErr)
			},
			expectedEvents: nil,
			expectedErr:    apperr.InternalServErr,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repo.NewMockRepo(ctrl)
			testCase.mockBehaviour(repo, testCase.inputUserID, testCase.inputDateParsed)

			svc := service.NewEventService(repo)
			events, err := svc.EventsForMonth(context.Background(), testCase.inputUserID, testCase.inputDate)

			assert.Equal(t, testCase.expectedEvents, events)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}
