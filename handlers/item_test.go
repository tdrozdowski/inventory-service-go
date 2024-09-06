package handlers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
	"inventory-service-go/commons"
	"inventory-service-go/context"
	"inventory-service-go/item"
	"net/http"
	"net/http/httptest"
	"testing"
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
