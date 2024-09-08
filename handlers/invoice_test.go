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
	"testing"
)

func TestInvoiceRoutes(t *testing.T) {
	mockApp := context.MockApplicationContext(nil, nil, nil)
	e := echo.New()
	t.Run("successful route registration", func(t *testing.T) {
		InvoiceRoutes(e.Group("/test"), mockApp)
		routes := e.Routes()
		assert.Equal(t, 8, len(routes))
	})
}

func TestGetAllInvoices(t *testing.T) {
	controller := gomock.NewController(t)
	mockInvoiceService := invoice.NewMockInvoiceService(controller)
	paginationFixture := commons.Pagination{
		LastId:   0,
		PageSize: 10,
	}
	expectedInvoices := []invoice.Invoice{
		{
			Seq:       1,
			Id:        uuid.UUID{},
			UserId:    uuid.UUID{},
			Total:     0,
			Paid:      false,
			Items:     nil,
			AuditInfo: commons.AuditInfo{},
		},
		{
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
		nilBody       bool
		expectBody    invoice.Invoice
		expectErrCode int
	}{

		{
			name: "successful update",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().UpdateInvoice(updateInvoiceRequest).Return(expectedInvoice, nil)
			},
			nilBody:       false,
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
			nilBody:       false,
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
			nilBody:       true,
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
			nilBody:       false,
			inputBody:     updateInvoiceRequest,
			paramId:       "bad-id",
			expectBody:    invoice.Invoice{},
			expectErrCode: http.StatusBadRequest,
		},
		{
			name: "bad request: mismatched ids",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().UpdateInvoice(updateInvoiceRequest).Times(0)
			},
			nilBody:       false,
			inputBody:     updateInvoiceRequest,
			paramId:       uuid.New().String(),
			expectBody:    invoice.Invoice{},
			expectErrCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApp := context.MockApplicationContext(nil, nil, mockInvoiceService)
			tt.mockFunc(mockInvoiceService)
			e := echo.New()
			var body []byte
			if tt.nilBody {
				body = []byte("invalid body")
			} else {
				body, _ = json.Marshal(tt.inputBody)
			}
			req := httptest.NewRequest(http.MethodPost, "/"+tt.paramId, bytes.NewReader(body))
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

func TestDeleteInvoice(t *testing.T) {
	gomockController := gomock.NewController(t)
	mockInvoiceService := invoice.NewMockInvoiceService(gomockController)
	id := uuid.New()
	mockApp := context.MockApplicationContext(nil, nil, mockInvoiceService)
	expectedResults := commons.DeleteResult{
		id,
		true,
	}
	tests := []struct {
		name          string
		funcSetup     func()
		paramId       string
		expectErrCode int
	}{
		{
			name: "successful deletion",
			funcSetup: func() {
				mockInvoiceService.EXPECT().DeleteInvoice(id).Return(expectedResults, nil)
			},
			paramId:       id.String(),
			expectErrCode: http.StatusOK,
		},
		{
			name: "internal server error",
			funcSetup: func() {
				mockInvoiceService.EXPECT().DeleteInvoice(id).Return(commons.DeleteResult{}, errors.New("BOOM"))
			},
			paramId:       id.String(),
			expectErrCode: http.StatusInternalServerError,
		},
		{
			name: "bad request: id",
			funcSetup: func() {
				mockInvoiceService.EXPECT().DeleteInvoice(id).Times(0)
			},
			paramId:       "bad-id",
			expectErrCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.funcSetup()
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/"+tt.paramId, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramId)
			if assert.NoError(t, DeleteInvoice(mockApp)(c)) {
				assert.Equal(t, tt.expectErrCode, rec.Code)
			}
		})
	}
	gomockController.Finish()
}

func TestGetAllInvoicesForUser(t *testing.T) {
	controller := gomock.NewController(t)
	mockInvoiceService := invoice.NewMockInvoiceService(controller)
	userId := uuid.New()
	expectedInvoices := []invoice.Invoice{
		{
			Seq:       1,
			Id:        uuid.UUID{},
			UserId:    userId,
			Total:     0,
			Paid:      false,
			Items:     nil,
			AuditInfo: commons.AuditInfo{},
		},
		{
			Seq:       2,
			Id:        uuid.UUID{},
			UserId:    userId,
			Total:     0,
			Paid:      false,
			Items:     nil,
			AuditInfo: commons.AuditInfo{},
		},
	}
	tests := []struct {
		name          string
		mockFunc      func(mockService *invoice.MockInvoiceService)
		paramUserId   string
		expectBody    []invoice.Invoice
		expectErrCode int
	}{
		{
			name: "successful retrieval",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().GetInvoicesForUser(userId).Return(expectedInvoices, nil)
			},
			paramUserId:   userId.String(),
			expectBody:    expectedInvoices,
			expectErrCode: http.StatusOK,
		},
		{
			name: "service error",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().GetInvoicesForUser(userId).Return([]invoice.Invoice{}, errors.New("BOOM"))
			},
			paramUserId:   userId.String(),
			expectErrCode: http.StatusInternalServerError,
		},
		{
			name: "bad request: userId",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().GetInvoicesForUser(userId).Times(0)
			},
			paramUserId:   "bad-id",
			expectErrCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApp := context.MockApplicationContext(nil, nil, mockInvoiceService)
			tt.mockFunc(mockInvoiceService)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/"+tt.paramUserId, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("userId")
			c.SetParamValues(tt.paramUserId)
			if assert.NoError(t, GetAllInvoicesForUser(mockApp)(c)) {
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

func TestAddItemsToInvoice(t *testing.T) {
	controller := gomock.NewController(t)
	mockInvoiceService := invoice.NewMockInvoiceService(controller)
	invoiceId := uuid.New()
	items := []uuid.UUID{uuid.New(), uuid.New()}
	addItemsRequest := invoice.ItemsToInvoiceRequest{
		InvoiceId: invoiceId,
		Items:     items,
	}
	expectedResult := invoice.ItemsToInvoiceResponse{
		InvoiceId: invoiceId,
		Items:     items,
		Success:   true,
	}
	tests := []struct {
		name          string
		mockFunc      func(mockService *invoice.MockInvoiceService)
		inputBody     invoice.ItemsToInvoiceRequest
		expectBody    invoice.ItemsToInvoiceResponse
		expectErrCode int
	}{
		{
			name: "successful addition",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().AddItemsToInvoice(addItemsRequest).Return(expectedResult, nil)
			},
			inputBody:     addItemsRequest,
			expectBody:    expectedResult,
			expectErrCode: http.StatusOK,
		},
		{
			name: "internal server error",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().AddItemsToInvoice(addItemsRequest).Return(invoice.ItemsToInvoiceResponse{}, errors.New("BOOM"))
			},
			inputBody:     addItemsRequest,
			expectErrCode: http.StatusInternalServerError,
		},
		{
			name: "bad request: body missing",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().AddItemsToInvoice(addItemsRequest).Times(0)
			},
			inputBody:     invoice.ItemsToInvoiceRequest{},
			expectErrCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApp := context.MockApplicationContext(nil, nil, mockInvoiceService)
			tt.mockFunc(mockInvoiceService)
			e := echo.New()
			var body []byte
			if tt.expectErrCode == http.StatusBadRequest {
				body = []byte("not json")
			} else {
				body, _ = json.Marshal(tt.inputBody)
			}
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if assert.NoError(t, AddItemsToInvoice(mockApp)(c)) {
				assert.Equal(t, tt.expectErrCode, rec.Code)
				if tt.expectErrCode == http.StatusOK {
					var body invoice.ItemsToInvoiceResponse
					err := json.NewDecoder(rec.Body).Decode(&body)
					assert.NoError(t, err)
					assert.Equal(t, tt.expectBody, body)
				}
			}
		})
	}
}

func TestRemoveItemFromInvoice(t *testing.T) {
	controller := gomock.NewController(t)
	mockInvoiceService := invoice.NewMockInvoiceService(controller)
	invoiceId := uuid.New()
	itemId := uuid.New()
	expectedResult := invoice.ItemsToInvoiceResponse{
		InvoiceId: invoiceId,
		Success:   true,
	}
	tests := []struct {
		name           string
		mockFunc       func(mockService *invoice.MockInvoiceService)
		paramInvoiceId string
		paramItemId    string
		expectBody     invoice.ItemsToInvoiceResponse
		expectErrCode  int
	}{
		{
			name: "successful removal",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().RemoveItemFromInvoice(invoice.SimpleInvoiceItem{invoiceId, itemId}).Return(expectedResult, nil)
			},
			paramInvoiceId: invoiceId.String(),
			paramItemId:    itemId.String(),
			expectBody:     expectedResult,
			expectErrCode:  http.StatusOK,
		},
		{
			name: "internal server error",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().RemoveItemFromInvoice(invoice.SimpleInvoiceItem{invoiceId, itemId}).Return(invoice.ItemsToInvoiceResponse{}, errors.New("BOOM"))
			},
			paramInvoiceId: invoiceId.String(),
			paramItemId:    itemId.String(),
			expectErrCode:  http.StatusInternalServerError,
		},
		{
			name: "bad request: invoiceId",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().RemoveItemFromInvoice(invoice.SimpleInvoiceItem{invoiceId, itemId}).Times(0)
			},
			paramInvoiceId: "bad-id",
			paramItemId:    itemId.String(),
			expectErrCode:  http.StatusBadRequest,
		},
		{
			name: "bad request: itemId",
			mockFunc: func(mockService *invoice.MockInvoiceService) {
				mockService.EXPECT().RemoveItemFromInvoice(invoice.SimpleInvoiceItem{invoiceId, itemId}).Times(0)
			},
			paramInvoiceId: invoiceId.String(),
			paramItemId:    "bad-id",
			expectErrCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApp := context.MockApplicationContext(nil, nil, mockInvoiceService)
			tt.mockFunc(mockInvoiceService)
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/"+tt.paramInvoiceId+"/"+tt.paramItemId, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id", "itemId")
			c.SetParamValues(tt.paramInvoiceId, tt.paramItemId)
			if assert.NoError(t, RemoveItemFromInvoice(mockApp)(c)) {
				assert.Equal(t, tt.expectErrCode, rec.Code)
				if tt.expectErrCode == http.StatusOK {
					var body invoice.ItemsToInvoiceResponse
					err := json.NewDecoder(rec.Body).Decode(&body)
					assert.NoError(t, err)
					assert.Equal(t, tt.expectBody, body)
				}
			}
		})
	}
}
