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
