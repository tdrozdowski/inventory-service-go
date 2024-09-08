package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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

func TestItemRoutes(t *testing.T) {
	mockApp := context.MockApplicationContext(nil, nil, nil)
	e := echo.New()
	t.Run("successful route registration", func(t *testing.T) {
		ItemRoutes(e.Group("/test"), mockApp)
		routes := e.Routes()
		assert.Equal(t, 5, len(routes))
	})
}

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
	mockApplicationContext := context.MockApplicationContext(nil, mockItemService, nil)
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
		mockApplicationContext := context.MockApplicationContext(nil, mockItemService, nil)
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
				req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestJson))
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

func TestHandlers_UpdateItem(t *testing.T) {
	expectedUuid := uuid.New()
	now := time.RFC3339
	pathId := expectedUuid.String()
	tests := []struct {
		name               string
		updateItemRequest  item.UpdateItemRequest
		expectedResults    item.Item
		expectedStatusCode int
		pathId             string
	}{
		{
			name: "OK with a valid request",
			updateItemRequest: item.UpdateItemRequest{
				Id:            expectedUuid,
				Name:          "Updated TV",
				Description:   "Updated TV Description",
				UnitPrice:     14.99,
				LastChangedBy: "Unit Test",
			},
			expectedResults: item.Item{
				Id:          expectedUuid,
				Name:        "Updated TV",
				Description: "Updated TV Description",
				UnitPrice:   14.99,
				AuditInfo: commons.AuditInfo{
					CreatedBy:     "Unit Test",
					CreatedAt:     now,
					LastUpdate:    "Unit Test",
					LastChangedBy: now,
				},
			},
			expectedStatusCode: http.StatusOK,
			pathId:             pathId,
		},
		{
			name:               "Fail with invalid request json",
			updateItemRequest:  item.UpdateItemRequest{},
			expectedResults:    item.Item{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail with mismatched item Id",
			updateItemRequest: item.UpdateItemRequest{
				Id:            expectedUuid,
				Name:          "TV",
				Description:   "New TV Description",
				UnitPrice:     19.99,
				LastChangedBy: "Unit Test",
			},
			expectedResults:    item.Item{},
			expectedStatusCode: http.StatusBadRequest,
			pathId:             uuid.NewString(),
		},
		{
			name: "Fail with bad path id",
			updateItemRequest: item.UpdateItemRequest{
				Id:            expectedUuid,
				Name:          "TV",
				Description:   "New TV Description",
				UnitPrice:     19.99,
				LastChangedBy: "Unit Test",
			},
			expectedResults:    item.Item{},
			expectedStatusCode: http.StatusBadRequest,
			pathId:             "not a uuid",
		},
		{
			name: "Fail with Internal Server Error",
			updateItemRequest: item.UpdateItemRequest{
				Id:            expectedUuid,
				Name:          "Updated TV",
				Description:   "Updated TV Description",
				UnitPrice:     14.99,
				LastChangedBy: "Unit Test",
			},
			expectedResults:    item.Item{},
			expectedStatusCode: http.StatusInternalServerError,
			pathId:             pathId,
		},
	}
	controller := gomock.NewController(t)
	for _, tt := range tests {
		mockItemService := item.NewMockItemService(controller)
		mockApplicationContext := context.MockApplicationContext(nil, mockItemService, nil)
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedStatusCode == http.StatusInternalServerError {
				mockItemService.EXPECT().UpdateItem(gomock.Any()).Return(nil, errors.New("error"))
			} else if tt.expectedStatusCode == http.StatusOK {
				mockItemService.EXPECT().UpdateItem(tt.updateItemRequest).Return(&tt.expectedResults, nil)
			}
			e := echo.New()
			requestJson, err := json.Marshal(tt.updateItemRequest)
			if err != nil {
				t.Errorf("UpdateItem() error = %v", err)
			}
			var req *http.Request
			if tt.expectedStatusCode == http.StatusBadRequest && tt.updateItemRequest.Id == uuid.Nil {
				invalidPayload := []byte(`not json`)
				req = httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(invalidPayload))
			} else {
				req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/%v", tt.updateItemRequest.Id), io.NopCloser(bytes.NewReader(requestJson)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.pathId)
			err = UpdateItem(mockApplicationContext)(c)
			if err != nil {
				t.Errorf("UpdateItem() error = %v, expectedStatusCode %v", err, tt.expectedStatusCode)
			}
			if rec.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("UpdateItem() = %v, expectedStatusCode %v", rec.Code, tt.expectedStatusCode)
			}
		})
	}
}
func TestHandlers_GetItem(t *testing.T) {
	expectedUuid := uuid.New()
	pathId := expectedUuid.String()
	tests := []struct {
		name               string
		id                 string
		expectedItem       *item.Item
		expectedStatusCode int
	}{
		{
			name:               "OK with a valid item ID",
			id:                 pathId,
			expectedItem:       &item.Item{},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Fail with invalid item ID",
			id:                 "invalidUuid",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Fail with non-existing item",
			id:                 pathId,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Fail with Internal Server Error",
			id:                 pathId,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	controller := gomock.NewController(t)
	for _, tt := range tests {
		mockItemService := item.NewMockItemService(controller)
		mockApplicationContext := context.MockApplicationContext(nil, mockItemService, nil)
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedStatusCode == http.StatusInternalServerError {
				mockItemService.EXPECT().GetItem(expectedUuid).Return(nil, errors.New("error"))
			} else if tt.expectedStatusCode == http.StatusNotFound {
				mockItemService.EXPECT().GetItem(expectedUuid).Return(nil, sql.ErrNoRows)
			} else if tt.expectedStatusCode == http.StatusOK {
				mockItemService.EXPECT().GetItem(expectedUuid).Return(tt.expectedItem, nil)
			}
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%v", tt.id), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			err := GetItem(mockApplicationContext)(c)
			if err != nil {
				t.Errorf("GetItem() error = %v, expectedStatusCode %v", err, tt.expectedStatusCode)
			}
			if rec.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("GetItem() = %v, expectedStatusCode %v", rec.Code, tt.expectedStatusCode)
			}
		})
	}
	controller.Finish()
}

func TestHandlers_DeleteItem(t *testing.T) {
	expectedUuid := uuid.New()
	pathId := expectedUuid.String()
	expectedResults := commons.DeleteResult{
		Id:      expectedUuid,
		Deleted: true,
	}
	tests := []struct {
		name               string
		id                 string
		expectedStatusCode int
	}{
		{
			name:               "OK with a valid item ID",
			id:                 pathId,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Fail with invalid item ID",
			id:                 "invalidUuid",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Fail with non-existing item",
			id:                 pathId,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Fail with Internal Server Error",
			id:                 pathId,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	controller := gomock.NewController(t)
	for _, tt := range tests {
		mockItemService := item.NewMockItemService(controller)
		mockApplicationContext := context.MockApplicationContext(nil, mockItemService, nil)
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedStatusCode == http.StatusInternalServerError {
				mockItemService.EXPECT().DeleteItem(expectedUuid).Return(nil, errors.New("error"))
			} else if tt.expectedStatusCode == http.StatusNotFound {
				mockItemService.EXPECT().DeleteItem(expectedUuid).Return(nil, sql.ErrNoRows)
			} else if tt.expectedStatusCode == http.StatusOK {
				mockItemService.EXPECT().DeleteItem(expectedUuid).Return(&expectedResults, nil)
			}
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%v", tt.id), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			err := DeleteItem(mockApplicationContext)(c)
			if err != nil {
				t.Errorf("DeleteItem() error = %v, expectedStatusCode %v", err, tt.expectedStatusCode)
			}
			if rec.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("DeleteItem() = %v, expectedStatusCode %v", rec.Code, tt.expectedStatusCode)
			}
		})
	}
	controller.Finish()
}

// End of file
