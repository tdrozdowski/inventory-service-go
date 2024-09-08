package invoice

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"inventory-service-go/commons"
	"testing"
	"time"
)

func TestInvoiceService_GetInvoice(t *testing.T) {
	invoiceUuid := uuid.New()
	now := time.Now()
	invoiceRowFixture := InvoiceRow{
		Id:            1,
		AltId:         invoiceUuid,
		UserId:        uuid.New(),
		Total:         10.0,
		Paid:          false,
		CreatedBy:     "Unit Test",
		CreatedAt:     now,
		LastChangedBy: "Unit Test",
		LastUpdate:    now,
	}
	invoiceFixture := fromRow(invoiceRowFixture)
	invoiceItemRowFixture := []InvoiceItemRow{{
		Id:                1,
		AltId:             invoiceUuid,
		UserId:            uuid.New(),
		Total:             10.0,
		Paid:              false,
		CreatedBy:         "Unit Test",
		CreatedAt:         now,
		LastChangedBy:     "Unit Test",
		LastUpdate:        now,
		ItemSeqId:         sql.NullInt64{Int64: 1, Valid: true},
		ItemAltId:         uuid.New(),
		ItemName:          sql.NullString{String: "Item 1", Valid: true},
		ItemDescription:   sql.NullString{String: "Item 1 Description", Valid: true},
		ItemUnitPrice:     sql.NullFloat64{Float64: 10.0, Valid: true},
		ItemCreatedBy:     sql.NullString{String: "Unit Test", Valid: true},
		ItemCreatedAt:     sql.NullTime{Time: now, Valid: true},
		ItemLastChangedBy: sql.NullString{String: "Unit Test", Valid: true},
		ItemLastUpdate:    sql.NullTime{Time: now, Valid: true},
	}}
	invoiceItemRowWithNoItemsFixture := []InvoiceItemRow{
		{
			Id:                1,
			AltId:             invoiceUuid,
			UserId:            uuid.New(),
			Total:             10.0,
			Paid:              false,
			CreatedBy:         "Unit Test",
			CreatedAt:         now,
			LastChangedBy:     "Unit Test",
			LastUpdate:        now,
			ItemSeqId:         sql.NullInt64{Valid: false},
			ItemAltId:         uuid.Nil,
			ItemName:          sql.NullString{Valid: false},
			ItemDescription:   sql.NullString{Valid: false},
			ItemUnitPrice:     sql.NullFloat64{Float64: 10.0, Valid: false},
			ItemCreatedBy:     sql.NullString{Valid: false},
			ItemCreatedAt:     sql.NullTime{Valid: false},
			ItemLastChangedBy: sql.NullString{Valid: false},
			ItemLastUpdate:    sql.NullTime{Valid: false},
		},
	}
	invoiceFixtureWithItems := fromRowWithItems(invoiceItemRowFixture)
	invoiceFixtureWithNoItems := fromRowWithItems(invoiceItemRowWithNoItemsFixture)
	emptyInvoiceRowFixture := InvoiceRow{}
	emptyInvoiceFixture := Invoice{}
	testCases := []struct {
		name      string
		want      Invoice
		wantErr   bool
		withItems bool
		noItems   bool
	}{
		{
			name:      "Get Invoice No Items",
			want:      invoiceFixture,
			wantErr:   false,
			withItems: false,
			noItems:   false,
		},
		{
			name:      "Get Invoice With Items",
			want:      invoiceFixtureWithItems,
			wantErr:   false,
			withItems: true,
			noItems:   false,
		},
		{
			name:      "Get Invoice With Items - No Items",
			want:      invoiceFixtureWithNoItems,
			wantErr:   false,
			withItems: true,
			noItems:   true,
		},
		{
			name:      "Get Invoice Error",
			want:      emptyInvoiceFixture,
			wantErr:   true,
			withItems: false,
			noItems:   false,
		},
	}
	controller := gomock.NewController(t)
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockInvoiceRepository(controller)
			if tt.wantErr {
				mockRepo.EXPECT().GetInvoice(invoiceUuid).Return(emptyInvoiceRowFixture, errors.New("Boom"))
			} else if tt.withItems && !tt.noItems {
				mockRepo.EXPECT().GetInvoiceWithItems(invoiceUuid).Return(invoiceItemRowFixture, nil)
			} else if tt.withItems && tt.noItems {
				mockRepo.EXPECT().GetInvoiceWithItems(invoiceUuid).Return(invoiceItemRowWithNoItemsFixture, nil)
			} else {
				mockRepo.EXPECT().GetInvoice(invoiceUuid).Return(invoiceRowFixture, nil)
			}
			service := NewInvoiceService(mockRepo)
			results, err := service.GetInvoice(invoiceUuid, tt.withItems)
			if (err != nil) != tt.wantErr {
				t.Errorf("InvoiceService.GetInvoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NotNil(t, results)
				assert.Equal(t, tt.want, results)
			}
		})
	}
}

func TestInvoiceService_GetInvoicesForUser(t *testing.T) {
	userId := uuid.New()
	invoiceRowFixture1 := InvoiceRow{
		Id:            1,
		AltId:         uuid.New(),
		UserId:        userId,
		Total:         10.0,
		Paid:          false,
		CreatedBy:     "Unit Test",
		CreatedAt:     time.Now(),
		LastChangedBy: "Unit Test",
		LastUpdate:    time.Now(),
	}
	invoiceRowFixture2 := InvoiceRow{
		Id:            2,
		AltId:         uuid.New(),
		UserId:        userId,
		Total:         15.0,
		Paid:          false,
		CreatedBy:     "Unit Test",
		CreatedAt:     time.Now(),
		LastChangedBy: "Unit Test",
		LastUpdate:    time.Now(),
	}
	emptyUuid := uuid.UUID{}
	emptyInvoiceRowFixture := []InvoiceRow{}
	invoicesFixture := []Invoice{fromRow(invoiceRowFixture1), fromRow(invoiceRowFixture2)}
	emptyInvoicesFixture := []Invoice{}
	testCases := []struct {
		name     string
		userId   uuid.UUID
		mockFunc func(mockRepo *MockInvoiceRepository, userId uuid.UUID)
		want     []Invoice
		wantErr  bool
	}{
		{
			name:   "Get All Invoices For User",
			userId: userId,
			mockFunc: func(mockRepo *MockInvoiceRepository, userId uuid.UUID) {
				mockRepo.EXPECT().GetAllForUser(userId).Return([]InvoiceRow{invoiceRowFixture1, invoiceRowFixture2}, nil)
			},
			want:    invoicesFixture,
			wantErr: false,
		},
		{
			name:   "Get All Invoices For User - Error",
			userId: userId,
			mockFunc: func(mockRepo *MockInvoiceRepository, userId uuid.UUID) {
				mockRepo.EXPECT().GetAllForUser(userId).Return(nil, errors.New("Boom"))
			},
			want:    emptyInvoicesFixture,
			wantErr: true,
		},
		{
			name:   "Get All Invoices For User No Invoices",
			userId: emptyUuid,
			mockFunc: func(mockRepo *MockInvoiceRepository, userId uuid.UUID) {
				mockRepo.EXPECT().GetAllForUser(userId).Return(emptyInvoiceRowFixture, nil)
			},
			want:    emptyInvoicesFixture,
			wantErr: false,
		},
	}
	controller := gomock.NewController(t)
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockInvoiceRepository(controller)
			tt.mockFunc(mockRepo, tt.userId)
			service := NewInvoiceService(mockRepo)
			results, err := service.GetInvoicesForUser(tt.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("InvoiceService.GetInvoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NotNil(t, results)
				assert.Equal(t, tt.want, results)
			}
		})
	}
}

func TestInvoiceService_CreateInvoice(t *testing.T) {
	userId := uuid.New()
	createdBy := "Unit Test"
	createInvoiceRequest := CreateInvoiceRequest{
		UserId:    userId,
		Paid:      false,
		Total:     0.0,
		CreatedBy: createdBy,
	}

	invoiceRow := InvoiceRow{
		Id:            1,
		AltId:         uuid.New(),
		UserId:        userId,
		Total:         0.0,
		Paid:          false,
		CreatedBy:     createdBy,
		CreatedAt:     time.Now(),
		LastChangedBy: createdBy,
		LastUpdate:    time.Now(),
	}
	invoice := fromRow(invoiceRow)

	testCases := []struct {
		name     string
		request  CreateInvoiceRequest
		want     Invoice
		wantErr  bool
		mockFunc func(mockRepo *MockInvoiceRepository, request CreateInvoiceRequest)
	}{
		{
			name:    "Create Invoice Successfully",
			request: createInvoiceRequest,
			want:    invoice,
			wantErr: false,
			mockFunc: func(mockRepo *MockInvoiceRepository, request CreateInvoiceRequest) {
				mockRepo.EXPECT().CreateInvoice(request).Return(invoiceRow, nil)
			},
		},
		{
			name:    "Create Invoice - Repo Error",
			request: createInvoiceRequest,
			want:    Invoice{},
			wantErr: true,
			mockFunc: func(mockRepo *MockInvoiceRepository, request CreateInvoiceRequest) {
				mockRepo.EXPECT().CreateInvoice(request).Return(InvoiceRow{}, errors.New("Repo Error"))
			},
		},
	}
	controller := gomock.NewController(t)
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockInvoiceRepository(controller)
			tt.mockFunc(mockRepo, tt.request)
			service := NewInvoiceService(mockRepo)
			result, err := service.CreateInvoice(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("InvoiceService.CreateInvoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestInvoiceService_UpdateInvoice(t *testing.T) {
	invoiceRow := InvoiceRow{
		Id:            1,
		AltId:         uuid.New(),
		UserId:        uuid.New(),
		Total:         10.0,
		Paid:          true,
		CreatedBy:     "Unit Test",
		CreatedAt:     time.Now(),
		LastChangedBy: "Unit Test",
		LastUpdate:    time.Now(),
	}
	emptyInvoiceRow := InvoiceRow{}
	invoice := fromRow(invoiceRow)
	updateInvoiceRequest := UpdateInvoiceRequest{
		Id:            invoiceRow.AltId,
		Paid:          true,
		Total:         50.0,
		LastChangedBy: "Unit Test Update",
	}
	emptyInvoice := Invoice{}
	testCases := []struct {
		name     string
		request  UpdateInvoiceRequest
		want     Invoice
		wantErr  bool
		mockFunc func(mockRepo *MockInvoiceRepository, request UpdateInvoiceRequest)
	}{
		{
			name:    "Update Invoice Successfully",
			request: updateInvoiceRequest,
			want:    invoice,
			wantErr: false,
			mockFunc: func(mockRepo *MockInvoiceRepository, request UpdateInvoiceRequest) {
				mockRepo.EXPECT().UpdateInvoice(request).Return(invoiceRow, nil)
			},
		},
		{
			name:    "Update Invoice - Repo Error",
			request: updateInvoiceRequest,
			want:    emptyInvoice,
			wantErr: true,
			mockFunc: func(mockRepo *MockInvoiceRepository, request UpdateInvoiceRequest) {
				mockRepo.EXPECT().UpdateInvoice(request).Return(emptyInvoiceRow, errors.New("Repo Error"))
			},
		},
	}
	controller := gomock.NewController(t)
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockInvoiceRepository(controller)
			tt.mockFunc(mockRepo, tt.request)
			service := NewInvoiceService(mockRepo)
			result, err := service.UpdateInvoice(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("InvoiceService.UpdateInvoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestInvoiceService_DeleteInvoice(t *testing.T) {
	testCases := []struct {
		name      string
		prepare   func(m *MockInvoiceRepository)
		want      commons.DeleteResult
		wantError bool
	}{
		{
			name: "Delete Invoice Successfully",
			prepare: func(m *MockInvoiceRepository) {
				m.EXPECT().DeleteInvoice(gomock.Any()).Return(commons.DeleteResult{Deleted: true}, nil).AnyTimes()
			},
			want:      commons.DeleteResult{Deleted: true},
			wantError: false,
		},
		{
			name: "Delete Invoice - Repo Error",
			prepare: func(m *MockInvoiceRepository) {
				m.EXPECT().DeleteInvoice(gomock.Any()).Return(commons.DeleteResult{}, errors.New("Repo Error")).AnyTimes()
			},
			want:      commons.DeleteResult{},
			wantError: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			mockRepo := NewMockInvoiceRepository(controller)
			tt.prepare(mockRepo)
			service := NewInvoiceService(mockRepo)
			result, err := service.DeleteInvoice(uuid.New())
			if (err != nil) != tt.wantError {
				t.Errorf("InvoiceService.DeleteInvoice() error = %v, wantErr %v", err, tt.wantError)
			}
			if !tt.wantError {
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestInvoiceService_GetAllInvoices(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepo := NewMockInvoiceRepository(controller)
	pag := &commons.Pagination{LastId: 1, PageSize: 5}

	testCases := []struct {
		name     string
		want     []Invoice
		wantErr  bool
		mockFunc func(mockRepo *MockInvoiceRepository)
	}{
		{
			name:    "Get All Invoices Successfully",
			want:    make([]Invoice, 0),
			wantErr: false,
			mockFunc: func(mockRepo *MockInvoiceRepository) {
				mockRepo.EXPECT().GetAll(pag).Return([]InvoiceRow{}, nil)
			},
		},
		{
			name:    "Get All Invoices - Repo Error",
			want:    make([]Invoice, 0),
			wantErr: true,
			mockFunc: func(mockRepo *MockInvoiceRepository) {
				mockRepo.EXPECT().GetAll(pag).Return(nil, errors.New("Repo Error"))
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockRepo)
			service := NewInvoiceService(mockRepo)
			result, err := service.GetAllInvoices(pag)
			if (err != nil) != tt.wantErr {
				t.Errorf("InvoiceService.GetAllInvoices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestInvoiceService_AddItemsToInvoice(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepo := NewMockInvoiceRepository(controller)
	invoiceUuid := uuid.New()
	itemUuid1 := uuid.New()
	itemUuid2 := uuid.New()

	testCases := []struct {
		name     string
		request  ItemsToInvoiceRequest
		want     ItemsToInvoiceResponse
		wantErr  bool
		mockFunc func(mockRepo *MockInvoiceRepository, request ItemsToInvoiceRequest)
	}{
		{
			name:    "Add Items To Invoice Successfully",
			request: ItemsToInvoiceRequest{InvoiceId: invoiceUuid, Items: []uuid.UUID{itemUuid1, itemUuid2}},
			want:    ItemsToInvoiceResponse{InvoiceId: invoiceUuid, Items: []uuid.UUID{itemUuid1, itemUuid2}, Success: true},
			wantErr: false,
			mockFunc: func(mockRepo *MockInvoiceRepository, request ItemsToInvoiceRequest) {
				mockRepo.EXPECT().AddItemsToInvoice(request).Return(ItemsToInvoiceResponse{InvoiceId: invoiceUuid, Items: []uuid.UUID{itemUuid1, itemUuid2}, Success: true}, nil)
			},
		},
		{
			name:    "Add Items To Invoice - Repo Error",
			request: ItemsToInvoiceRequest{InvoiceId: invoiceUuid, Items: []uuid.UUID{itemUuid1, itemUuid2}},
			want:    ItemsToInvoiceResponse{},
			wantErr: true,
			mockFunc: func(mockRepo *MockInvoiceRepository, request ItemsToInvoiceRequest) {
				mockRepo.EXPECT().AddItemsToInvoice(request).Return(ItemsToInvoiceResponse{}, errors.New("Repo Error"))
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockRepo, tt.request)
			service := NewInvoiceService(mockRepo)
			result, err := service.AddItemsToInvoice(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("InvoiceService.AddItemsToInvoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestInvoiceService_RemoveItemFromInvoice(t *testing.T) {
	controller := gomock.NewController(t)
	mockRepo := NewMockInvoiceRepository(controller)
	invoiceUuid := uuid.New()
	itemUuid := uuid.New()

	testCases := []struct {
		name     string
		request  SimpleInvoiceItem
		mockFunc func(mockRepo *MockInvoiceRepository)
		want     ItemsToInvoiceResponse
		wantErr  bool
	}{
		{
			name:    "Remove Item From Invoice Successfully",
			request: SimpleInvoiceItem{InvoiceId: invoiceUuid, ItemId: itemUuid},
			mockFunc: func(mockRepo *MockInvoiceRepository) {
				mockRepo.EXPECT().RemoveItemFromInvoice(gomock.Eq(SimpleInvoiceItem{InvoiceId: invoiceUuid, ItemId: itemUuid})).Return(ItemsToInvoiceResponse{InvoiceId: invoiceUuid, Items: []uuid.UUID{}, Success: true}, nil)
			},
			want:    ItemsToInvoiceResponse{InvoiceId: invoiceUuid, Items: []uuid.UUID{}, Success: true},
			wantErr: false,
		},
		{
			name:    "Remove Item From Invoice - Repo Error",
			request: SimpleInvoiceItem{InvoiceId: invoiceUuid, ItemId: itemUuid},
			mockFunc: func(mockRepo *MockInvoiceRepository) {
				mockRepo.EXPECT().RemoveItemFromInvoice(gomock.Eq(SimpleInvoiceItem{InvoiceId: invoiceUuid, ItemId: itemUuid})).Return(ItemsToInvoiceResponse{}, errors.New("Repo Error"))
			},
			want:    ItemsToInvoiceResponse{},
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockRepo)
			service := NewInvoiceService(mockRepo)
			result, err := service.RemoveItemFromInvoice(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("InvoiceService.RemoveItemFromInvoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, result)
			}
		})
	}
}
