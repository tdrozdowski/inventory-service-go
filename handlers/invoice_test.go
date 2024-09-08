package handlers

import (
	"bytes"
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
	"sort"
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
		sort.SliceStable(routes, func(i, j int) bool {
			return routes[i].Name < routes[j].Name
		})
		assert.Equal(t, 4, len(routes))
		assert.Equal(t, "/test/invoices", routes[0].Path)
		assert.Equal(t, http.MethodPost, routes[0].Method)
		assert.Equal(t, "/test/invoices", routes[1].Path)
		assert.Equal(t, http.MethodGet, routes[1].Method)
	})
}
func TestCreateInvoice(t *testing.T) {
	controller := gomock.NewController(t)
	mockInvoiceService := invoice.NewMockInvoiceService(controller)
	userId := uuid.New()
	createInvoiceRequest := invoice.CreateInvoiceRequest{
		UserId:    userId,
		Paid:      false,
		Total:     10.0,
		CreatedBy: "unit test",
	}
	expectedInvoice := invoice.Invoice{
		Seq:       1,
		Id:        uuid.UUID{},
		UserId:    userId,
		Total:     10.0,
		Paid:      false,
		Items:     nil,
		AuditInfo: commons.AuditInfo{},
	}
	tests := []struct {
		name          string
		mockFunc      func(mockService *invoice.MockInvoiceService)
		inputBody     invoice.CreateInvoiceRequest
		expectBody    invoice.Invoice
		expectErrCode int
	}{

		{
			name: "successful creation",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().CreateInvoice(createInvoiceRequest).Return(expectedInvoice, nil)
			},
			inputBody:     createInvoiceRequest,
			expectBody:    expectedInvoice,
			expectErrCode: http.StatusOK,
		},
		{
			name: "internal server error",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().CreateInvoice(createInvoiceRequest).Return(invoice.Invoice{}, errors.New("BOOM"))
			},
			inputBody:     createInvoiceRequest,
			expectBody:    invoice.Invoice{},
			expectErrCode: http.StatusInternalServerError,
		},
		{
			name: "bad request",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().CreateInvoice(createInvoiceRequest).Times(0)
			},
			inputBody:     invoice.CreateInvoiceRequest{},
			expectBody:    invoice.Invoice{},
			expectErrCode: http.StatusBadRequest,
		},
	}
	{
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockApp := context.MockApplicationContext(nil, nil, mockInvoiceService)
				tt.mockFunc(mockInvoiceService)
				e := echo.New()
				var jsonBody []byte
				if tt.expectErrCode == http.StatusBadRequest {
					jsonBody = []byte("bad request")
				} else {
					jsonBody, _ = json.Marshal(tt.inputBody)
				}
				req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				if assert.NoError(t, CreateInvoice(mockApp)(c)) {
					assert.Equal(t, tt.expectErrCode, rec.Code)
					if tt.expectErrCode == http.StatusOK {
						var body invoice.Invoice
						err := json.NewDecoder(rec.Body).Decode(&body)
						assert.NoError(t, err)
						assert.Equal(t, tt.expectBody, body)
					}
				}
			})
		}
	}
}

func TestUpdateInvoice(t *testing.T) {
	controller := gomock.NewController(t)
	mockInvoiceService := invoice.NewMockInvoiceService(controller)
	id := uuid.New()
	userId := uuid.New()
	updateInvoiceRequest := invoice.UpdateInvoiceRequest{
		Id:            id,
		Paid:          true,
		Total:         20.0,
		LastChangedBy: "unit test",
	}
	expectedInvoice := invoice.Invoice{
		Seq:       1,
		Id:        id,
		UserId:    userId,
		Total:     20.0,
		Paid:      true,
		Items:     nil,
		AuditInfo: commons.AuditInfo{},
	}
	tests := []struct {
		name          string
		mockFunc      func(mockService *invoice.MockInvoiceService)
		inputBody     invoice.UpdateInvoiceRequest
		paramId       string
		expectBody    invoice.Invoice
		expectErrCode int
	}{

		{
			name: "successful update",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().UpdateInvoice(updateInvoiceRequest).Return(expectedInvoice, nil)
			},
			inputBody:     updateInvoiceRequest,
			paramId:       id.String(),
			expectBody:    expectedInvoice,
			expectErrCode: http.StatusOK,
		},
		{
			name: "internal server error",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().UpdateInvoice(updateInvoiceRequest).Return(invoice.Invoice{}, errors.New("BOOM"))
			},
			inputBody:     updateInvoiceRequest,
			paramId:       id.String(),
			expectBody:    invoice.Invoice{},
			expectErrCode: http.StatusInternalServerError,
		},
		{
			name: "bad request: body",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().UpdateInvoice(updateInvoiceRequest).Times(0)
			},
			inputBody:     invoice.UpdateInvoiceRequest{},
			paramId:       id.String(),
			expectBody:    invoice.Invoice{},
			expectErrCode: http.StatusBadRequest,
		},
		{
			name: "bad request: id",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().UpdateInvoice(updateInvoiceRequest).Times(0)
			},
			inputBody:     updateInvoiceRequest,
			paramId:       "bad-id",
			expectBody:    invoice.Invoice{},
			expectErrCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApp := context.MockApplicationContext(nil, nil, mockInvoiceService)
			tt.mockFunc(mockInvoiceService)
			e := echo.New()
			jsonBody, _ := json.Marshal(tt.inputBody)
			req := httptest.NewRequest(http.MethodPost, "/"+tt.paramId, bytes.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramId)
			if assert.NoError(t, UpdateInvoice(mockApp)(c)) {
				assert.Equal(t, tt.expectErrCode, rec.Code)
				if tt.expectErrCode == http.StatusOK {
					var body invoice.Invoice
					err := json.NewDecoder(rec.Body).Decode(&body)
					assert.NoError(t, err)
					assert.Equal(t, tt.expectBody, body)
				}
			}
		})
	}
}

func TestGetInvoice(t *testing.T) {
	controller := gomock.NewController(t)
	mockInvoiceService := invoice.NewMockInvoiceService(controller)
	id := uuid.New()
	expectedInvoice := invoice.Invoice{
		Seq:       1,
		Id:        id,
		UserId:    uuid.New(),
		Total:     10.0,
		Paid:      false,
		Items:     nil,
		AuditInfo: commons.AuditInfo{},
	}
	tests := []struct {
		name           string
		mockFunc       func(mockService *invoice.MockInvoiceService)
		paramId        string
		queryWithItems string
		expectBody     invoice.Invoice
		expectErrCode  int
	}{
		{
			name: "successful retrieval",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().GetInvoice(id, true).Return(expectedInvoice, nil)
			},
			paramId:        id.String(),
			queryWithItems: "true",
			expectBody:     expectedInvoice,
			expectErrCode:  http.StatusOK,
		},
		{
			name: "internal server error",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().GetInvoice(id, true).Return(invoice.Invoice{}, errors.New("BOOM"))
			},
			paramId:        id.String(),
			queryWithItems: "true",
			expectErrCode:  http.StatusInternalServerError,
		},
		{
			name: "bad request: id",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().GetInvoice(id, true).Times(0)
			},
			paramId:        "bad-id",
			queryWithItems: "true",
			expectErrCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApp := context.MockApplicationContext(nil, nil, mockInvoiceService)
			tt.mockFunc(mockInvoiceService)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/"+tt.paramId+"?withItems="+tt.queryWithItems, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramId)
			if assert.NoError(t, GetInvoice(mockApp)(c)) {
				assert.Equal(t, tt.expectErrCode, rec.Code)
				if tt.expectErrCode == http.StatusOK {
					var body invoice.Invoice
					err := json.NewDecoder(rec.Body).Decode(&body)
					assert.NoError(t, err)
					assert.Equal(t, tt.expectBody, body)
				}
			}
		})
	}
}
