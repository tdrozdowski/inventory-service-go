package handlers

// test the handlers in person.go

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
	"inventory-service-go/commons"
	"inventory-service-go/context"
	"inventory-service-go/person"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func personFixture() person.Person {
	now := time.Now().String()
	return person.Person{
		Seq:   1,
		Id:    uuid.UUID{},
		Name:  "John Doe",
		Email: "john.doe@test.com",
		AuditInfo: commons.AuditInfo{
			CreatedBy:     "unit_test",
			CreatedAt:     now,
			LastUpdate:    "unit_test",
			LastChangedBy: now,
		},
	}
}

func TestGetAll(t *testing.T) {
	pagination := &commons.Pagination{LastId: 0, PageSize: 10}
	controller := gomock.NewController(t)
	mockPersonService := person.NewMockPersonService(controller)
	applicationContext := context.MockApplicationContext(mockPersonService)
	expectedPersons := []person.Person{personFixture()}
	tests := []struct {
		name         string
		expectedCode int
		appContext   context.ApplicationContext
	}{
		{
			name:         "Valid Request",
			expectedCode: http.StatusOK,
			appContext:   applicationContext,
		},
		{
			name:         "Internal Error Scenario",
			expectedCode: http.StatusInternalServerError,
			appContext:   applicationContext,
		},
		// Add more test cases as necessary
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedCode == http.StatusInternalServerError {
				mockPersonService.EXPECT().GetAll(pagination).Return(nil, errors.New("Internal Error"))
			} else {
				mockPersonService.EXPECT().GetAll(pagination).Return(expectedPersons, nil)
			}
			req := httptest.NewRequest(http.MethodGet, "/?last_id=0&page_size=10", nil)
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)
			handler := GetAll(tt.appContext)

			err := handler(c)
			if err != nil {
				t.Errorf("Handler returned error: %v", err)
			}
			if tt.expectedCode != rec.Result().StatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, rec.Result().StatusCode)
			}
		})
	}
}
