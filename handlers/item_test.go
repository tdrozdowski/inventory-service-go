package handlers

import (
	"errors"
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
	for _, tt := range tests {
		mockItemService := item.NewMockItemService(controller)
		mockApplicationContext := context.MockApplicationContext(nil, *mockItemService)
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedStatusCode == http.StatusInternalServerError {
				mockItemService.EXPECT().GetItems(tt.pagination).Return([]item.Item{}, errors.New("error"))
			} else {
				mockItemService.EXPECT().GetItems(tt.pagination).Return([]item.Item{}, nil)
			}
			req := httptest.NewRequest(http.MethodGet, "/?last_id=0&page_size=10", nil)
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)
			handler := AllItems(mockApplicationContext)
			if err := handler(c); (err != nil) != (tt.expectedStatusCode == http.StatusInternalServerError) {
				t.Errorf("AllItems() error = %v, expectedStatusCode %v", err, tt.expectedStatusCode)
			}
			if rec.Code != tt.expectedStatusCode {
				t.Errorf("AllItems() = %v, expectedStatusCode %v", rec.Code, tt.expectedStatusCode)
			}
		})
	}

}
