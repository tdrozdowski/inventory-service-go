package invoice

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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
		ItemSeqId:         1,
		ItemAltId:         uuid.New(),
		ItemName:          "Item 1",
		ItemDescription:   "Item 1 Description",
		ItemUnitPrice:     10.0,
		ItemCreatedBy:     "Unit Test",
		ItemCreatedAt:     now,
		ItemLastChangedBy: "Unit Test",
		ItemLastUpdate:    now,
	}}
	invoiceFixtureWithItems := fromRowWithItems(invoiceItemRowFixture)
	emptyInvoiceRowFixture := InvoiceRow{}
	emptyInvoiceFixture := Invoice{}
	testCases := []struct {
		name      string
		want      Invoice
		wantErr   bool
		withItems bool
	}{
		{
			name:      "Get Invoice No Items",
			want:      invoiceFixture,
			wantErr:   false,
			withItems: false,
		},
		{
			name:      "Get Invoice With Items",
			want:      invoiceFixtureWithItems,
			wantErr:   false,
			withItems: true,
		},
		{
			name:      "Get Invoice Error",
			want:      emptyInvoiceFixture,
			wantErr:   true,
			withItems: false,
		},
	}
	controller := gomock.NewController(t)
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockInvoiceRepository(controller)
			if tt.wantErr {
				mockRepo.EXPECT().GetInvoice(invoiceUuid).Return(emptyInvoiceRowFixture, errors.New("Boom"))
			} else if tt.withItems {
				mockRepo.EXPECT().GetInvoiceWithItems(invoiceUuid).Return(invoiceItemRowFixture, nil)
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
