package handlers

// test the handlers in person.go

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"inventory-service-go/commons"
	"inventory-service-go/context"
	"inventory-service-go/invoice"
	"inventory-service-go/item"
	"inventory-service-go/person"
	"io"
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

func TestPersonRoutes(t *testing.T) {
	mockApp := context.MockApplicationContext(nil, nil, nil)
	e := echo.New()
	t.Run("successful route registration", func(t *testing.T) {
		PersonRoutes(e.Group("/test"), mockApp)
		routes := e.Routes()
		assert.Equal(t, 5, len(routes))
	})
}

func TestGetAll(t *testing.T) {
	pagination := &commons.Pagination{LastId: 0, PageSize: 10}
	controller := gomock.NewController(t)
	mockPersonService := person.NewMockPersonService(controller)
	mockItemService := item.NewMockItemService(controller)
	applicationContext := context.MockApplicationContext(mockPersonService, mockItemService, nil)
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
			handler := GetAllPersons(tt.appContext)

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

func TestGetById(t *testing.T) {
	tests := []struct {
		name         string
		uuid         string
		expectedCode int
	}{
		{
			name:         "Valid Request",
			uuid:         uuid.NewString(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid Request",
			uuid:         "1",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Internal Error Scenario",
			uuid:         uuid.NewString(),
			expectedCode: http.StatusInternalServerError,
		},
	}
	controller := gomock.NewController(t)
	for _, tt := range tests {
		mockPersonService := person.NewMockPersonService(controller)
		mockItemService := item.NewMockItemService(controller)
		applicationContext := context.MockApplicationContext(mockPersonService, mockItemService, nil)
		expectedPerson := personFixture()
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedCode == http.StatusInternalServerError {
				mockPersonService.EXPECT().GetById(gomock.Any()).Return(nil, errors.New("Internal Error"))
			} else if tt.expectedCode == http.StatusOK {
				mockPersonService.EXPECT().GetById(gomock.Any()).Return(&expectedPerson, nil)
			}
			uri := fmt.Sprintf("/%s", tt.uuid)
			req := httptest.NewRequest(http.MethodGet, uri, nil)
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.uuid)
			handler := GetPersonById(applicationContext)
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

func TestCreate(t *testing.T) {
	tests := []struct {
		name          string
		createRequest person.CreatePersonRequest
		expectedCode  int
	}{
		{
			name: "Valid Request",
			createRequest: person.CreatePersonRequest{
				Name:      "John Doe",
				Email:     "john.doe@test.com",
				CreatedBy: "unit_test",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:          "Internal Error Scenario",
			createRequest: person.CreatePersonRequest{},
			expectedCode:  http.StatusInternalServerError,
		},
		{
			name:          "Invalid Request Body",
			createRequest: person.CreatePersonRequest{},
			expectedCode:  http.StatusBadRequest,
		},
	}
	controller := gomock.NewController(t)
	for _, tt := range tests {
		mockPersonService := person.NewMockPersonService(controller)
		mockItemService := item.NewMockItemService(controller)
		applicationContext := context.MockApplicationContext(mockPersonService, mockItemService, nil)
		expectedPerson := personFixture()
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedCode == http.StatusInternalServerError {
				mockPersonService.EXPECT().Create(gomock.Any()).Return(nil, errors.New("Internal Error"))
			} else if tt.expectedCode == http.StatusCreated {
				mockPersonService.EXPECT().Create(gomock.Any()).Return(&expectedPerson, nil)
			}
			var requestBody []byte
			if tt.expectedCode == http.StatusBadRequest {
				requestBody = []byte(`bad request`)
			} else {
				requestBody, _ = json.Marshal(tt.createRequest)
			}
			req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)
			handler := CreatePerson(applicationContext)
			err := handler(c)
			if err != nil {
				return
			}
			if err != nil {
				t.Errorf("Handler returned error: %v", err)
			}
			if tt.expectedCode != rec.Result().StatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, rec.Result().StatusCode)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name          string
		updateRequest person.UpdatePersonRequest
		expectedCode  int
	}{
		{
			name: "Valid Request",
			updateRequest: person.UpdatePersonRequest{
				Id:            uuid.UUID{},
				Name:          "Test User",
				Email:         "test.user@test.com",
				LastChangedBy: "unit_test",
			},
			expectedCode: http.StatusOK,
		},
		{
			name:          "Internal Error Scenario",
			updateRequest: person.UpdatePersonRequest{},
			expectedCode:  http.StatusInternalServerError,
		},
		{
			name:          "Invalid Request Body",
			updateRequest: person.UpdatePersonRequest{},
			expectedCode:  http.StatusBadRequest,
		},
	}
	controller := gomock.NewController(t)
	for _, tt := range tests {
		mockPersonService := person.NewMockPersonService(controller)
		mockItemService := item.NewMockItemService(controller)
		applicationContext := context.MockApplicationContext(mockPersonService, mockItemService, nil)
		expectedPerson := personFixture()
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedCode == http.StatusInternalServerError {
				mockPersonService.EXPECT().Update(gomock.Any()).Return(nil, errors.New("Internal Error"))
			} else if tt.expectedCode == http.StatusOK {
				mockPersonService.EXPECT().Update(gomock.Any()).Return(&expectedPerson, nil)
			}
			var requestBody []byte
			if tt.expectedCode == http.StatusBadRequest {
				requestBody = []byte(`bad request`)
			} else {
				requestBody, _ = json.Marshal(tt.updateRequest)
			}
			uri := fmt.Sprintf("/%s", tt.updateRequest.Id)
			req := httptest.NewRequest(http.MethodPut, uri, io.NopCloser(bytes.NewReader(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.updateRequest.Id.String())
			handler := UpdatePerson(applicationContext)
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

func TestDelete(t *testing.T) {
	tests := []struct {
		name         string
		uuid         string
		expectedCode int
	}{
		{
			name:         "Valid Request",
			uuid:         uuid.NewString(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid Request",
			uuid:         "1",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Internal Error Scenario",
			uuid:         uuid.NewString(),
			expectedCode: http.StatusInternalServerError,
		},
	}
	controller := gomock.NewController(t)
	for _, tt := range tests {
		mockPersonService := person.NewMockPersonService(controller)
		mockItemService := item.NewMockItemService(controller)
		mockInvoiceService := invoice.NewMockInvoiceService(controller)
		applicationContext := context.MockApplicationContext(mockPersonService, mockItemService, mockInvoiceService)
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedCode == http.StatusInternalServerError {
				mockPersonService.EXPECT().DeleteByUuid(gomock.Any()).Return(nil, errors.New("Internal Error"))
			} else if tt.expectedCode == http.StatusOK {
				mockPersonService.EXPECT().DeleteByUuid(gomock.Any()).Return(&commons.DeleteResult{}, nil)
			}
			uri := fmt.Sprintf("/%s", tt.uuid)
			req := httptest.NewRequest(http.MethodDelete, uri, nil)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.uuid)
			handler := DeletePerson(applicationContext)
			err := handler(c)
			if err != nil {
				msg := fmt.Sprintf("Handler returned error: %v", err)
				assert.Fail(t, msg)
			}
			if tt.expectedCode != rec.Result().StatusCode {
				msg := fmt.Sprintf("Expected status code %d, got %d", tt.expectedCode, rec.Result().StatusCode)
				assert.FailNow(t, msg)
			}
		})
	}
}

func TestPaginationFromRequest(t *testing.T) {
	tests := []struct {
		name         string
		request      echo.Context
		expectedPage *commons.Pagination
	}{
		{
			name: "No pagination parameters",
			request: echo.New().NewContext(
				httptest.NewRequest(http.MethodGet, "/", nil),
				httptest.NewRecorder(),
			),
			expectedPage: nil,
		},
		{
			name: "Valid pagination parameters",
			request: echo.New().NewContext(
				httptest.NewRequest(http.MethodGet, "/?last_id=5&page_size=15", nil),
				httptest.NewRecorder(),
			),
			expectedPage: &commons.Pagination{
				LastId:   5,
				PageSize: 15,
			},
		},
		{
			name: "Invalid LastId, valid PageSize",
			request: echo.New().NewContext(
				httptest.NewRequest(http.MethodGet, "/?last_id=invalid&page_size=15", nil),
				httptest.NewRecorder(),
			),
			expectedPage: &commons.Pagination{
				LastId:   0,
				PageSize: 15,
			},
		},
		{
			name: "Valid LastId, invalid PageSize",
			request: echo.New().NewContext(
				httptest.NewRequest(http.MethodGet, "/?last_id=5&page_size=invalid", nil),
				httptest.NewRecorder(),
			),
			expectedPage: &commons.Pagination{
				LastId:   5,
				PageSize: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualPage := paginationFromRequest(tt.request)
			assert.Equal(t, tt.expectedPage, actualPage)
		})
	}
}
