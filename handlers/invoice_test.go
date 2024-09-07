package handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"inventory-service-go/commons"
	"inventory-service-go/context"
	"inventory-service-go/invoice"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllInvoices(t *testing.T) {
	controller := gomock.NewController(t)
	mockInvoiceService := invoice.NewMockInvoiceService(controller)
	paginationFixture := commons.Pagination{
		LastId:   0,
		PageSize: 10,
	}
	expectedInvoices := []invoice.Invoice{
		invoice.Invoice{
			Seq:       1,
			Id:        uuid.UUID{},
			UserId:    uuid.UUID{},
			Total:     0,
			Paid:      false,
			Items:     nil,
			AuditInfo: commons.AuditInfo{},
		},
		invoice.Invoice{
			Seq:       2,
			Id:        uuid.UUID{},
			UserId:    uuid.UUID{},
			Total:     0,
			Paid:      false,
			Items:     nil,
			AuditInfo: commons.AuditInfo{},
		},
	}
	tests := []struct {
		name          string
		mockFunc      func(mockService *invoice.MockInvoiceService)
		expectBody    []invoice.Invoice
		expectErrCode int
	}{
		{
			name: "successful retrieval",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().GetAllInvoices(&paginationFixture).Return(expectedInvoices, nil)
			},
			expectBody:    expectedInvoices,
			expectErrCode: http.StatusOK,
		},
		{
			name: "service error",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().GetAllInvoices(&paginationFixture).Return([]invoice.Invoice{}, errors.New("BOOM"))
			},
			expectBody:    nil,
			expectErrCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			mockApp := context.MockApplicationContext(nil, nil, mockInvoiceService)
			tt.mockFunc(mockInvoiceService)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/?last_id=0&page_size=10", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if assert.NoError(t, GetAllInvoices(mockApp)(c)) {
				assert.Equal(t, tt.expectErrCode, rec.Code)
				if tt.expectErrCode == http.StatusOK {
					var body []invoice.Invoice
					err := json.NewDecoder(rec.Body).Decode(&body)
					assert.NoError(t, err)
					assert.Equal(t, tt.expectBody, body)
				}
			}
		})
	}
}

func TestInvoiceRoutes(t *testing.T) {
	controller := gomock.NewController(t)
	mockInvoiceService := invoice.NewMockInvoiceService(controller)
	mockApp := context.MockApplicationContext(nil, nil, mockInvoiceService)
	e := echo.New()

	t.Run("successful route registration", func(t *testing.T) {
		InvoiceRoutes(e.Group("/test"), mockApp)
		routes := e.Routes()
		assert.Equal(t, 1, len(routes))
		assert.Equal(t, "/test/invoices", routes[0].Path)
		assert.Equal(t, http.MethodGet, routes[0].Method)
	})
}
