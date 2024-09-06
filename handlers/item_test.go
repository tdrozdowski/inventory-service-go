package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
	"inventory-service-go/commons"
	"inventory-service-go/context"
	"inventory-service-go/item"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandlers_AllItems(t *testing.T) {
	tests := []struct {
		name               string
		pagination         *commons.Pagination
		expectedStatusCode int
	}{
		{
			name:               "OK with pagination",
			pagination:         &commons.Pagination{LastId: 1, PageSize: 1},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "OK without pagination",
			pagination:         nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "should return 500",
			pagination:         nil,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	controller := gomock.NewController(t)
	mockItemService := item.NewMockItemService(controller)
	mockApplicationContext := context.MockApplicationContext(nil, mockItemService)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedStatusCode == http.StatusInternalServerError {
				mockItemService.EXPECT().GetItems(nil).Return([]item.Item{}, errors.New("error"))
			} else {
				if tt.pagination == nil {
					mockItemService.EXPECT().GetItems(nil).Return([]item.Item{}, nil)
				} else {
					mockItemService.EXPECT().GetItems(tt.pagination).Return([]item.Item{}, nil)
				}
			}
			var target string
			if tt.pagination == nil {
				target = fmt.Sprintf("/")
			} else {
				target = fmt.Sprintf("/?last_id=%d&page_size=%d", tt.pagination.LastId, tt.pagination.PageSize)
			}
			req := httptest.NewRequest(http.MethodGet, target, nil)
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)
			handler := AllItems(mockApplicationContext)
			err := handler(c)
			if err != nil {
				t.Errorf("AllItems() error = %v, expectedStatusCode %v", err, tt.expectedStatusCode)
			}
			if rec.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("AllItems() = %v, expectedStatusCode %v", rec.Code, tt.expectedStatusCode)
			}
		})
	}

}

func TestHandlers_CreateItem(t *testing.T) {
	now := time.RFC3339
	tests := []struct {
		name               string
		createItemRequest  item.CreateItemRequest
		expectedResults    item.Item
		expectedStatusCode int
	}{
		{
			name:              "OK with a valid request",
			createItemRequest: item.CreateItemRequest{Name: "TV", Description: "TV Description", UnitPrice: 10.99, CreatedBy: "Unit Test"},
			expectedResults: item.Item{Name: "TV", Description: "TV Description", UnitPrice: 10.99, AuditInfo: commons.AuditInfo{
				CreatedBy:     "Unit Test",
				CreatedAt:     now,
				LastUpdate:    "Unit Test",
				LastChangedBy: now,
			}},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Fail with an invalid request",
			createItemRequest:  item.CreateItemRequest{Name: "", Description: "", UnitPrice: -1},
			expectedResults:    item.Item{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Fail with Internal Server Error",
			createItemRequest:  item.CreateItemRequest{Name: "", Description: "", UnitPrice: -1},
			expectedResults:    item.Item{},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	controller := gomock.NewController(t)
	for _, tt := range tests {
		mockItemService := item.NewMockItemService(controller)
		mockApplicationContext := context.MockApplicationContext(nil, mockItemService)
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedStatusCode == http.StatusInternalServerError {
				mockItemService.EXPECT().CreateItem(gomock.Any()).Return(nil, errors.New("error"))
			} else if tt.expectedStatusCode == http.StatusOK {
				mockItemService.EXPECT().CreateItem(tt.createItemRequest).Return(&tt.expectedResults, nil)
			}
			e := echo.New()
			var req *http.Request
			if tt.expectedStatusCode == http.StatusBadRequest {
				invalidPayload := []byte(`not json`)
				req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(invalidPayload))
			} else {
				requestJson, err := json.Marshal(tt.createItemRequest)
				if err != nil {
					t.Errorf("CreateItem() error = %v", err)
				}
				req = httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewReader(requestJson)))
			}
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := CreateItem(mockApplicationContext)(c)
			if err != nil {
				t.Errorf("CreateItem() error = %v, expectedStatusCode %v", err, tt.expectedStatusCode)
			}
			if rec.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("CreateItem() = %v, expectedStatusCode %v", rec.Code, tt.expectedStatusCode)
			}
		})
	}
}
