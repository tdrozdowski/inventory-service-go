package context

import (
	"go.uber.org/mock/gomock"
	"inventory-service-go/auth"
	"inventory-service-go/invoice"
	"inventory-service-go/item"
	"inventory-service-go/person"
	"testing"
)

func TestNewApplicationContext(t *testing.T) {
	appCtx := NewApplicationContext()
	if _, ok := appCtx.AuthProvider().(auth.AuthProvider); !ok {
		t.Error("AuthProvider should be of type auth.AuthProvider")
	}
	if _, ok := appCtx.PersonService().(person.PersonService); !ok {
		t.Error("PersonService should be of type person.PersonService")
	}
	if _, ok := appCtx.ItemService().(item.ItemService); !ok {
		t.Error("ItemService should be of type item.ItemService")
	}
	if _, ok := appCtx.InvoiceService().(invoice.InvoiceService); !ok {
		t.Error("InvoiceService should be of type invoice.InvoiceService")
	}
}

func TestMockApplicationContext(t *testing.T) {
	controller := gomock.NewController(t)
	mockPersonService := person.NewMockPersonService(controller)
	mockItemService := item.NewMockItemService(controller)
	mockInvoiceService := invoice.NewMockInvoiceService(controller)
	appCtx := MockApplicationContext(mockPersonService, mockItemService, mockInvoiceService)
	if _, ok := appCtx.AuthProvider().(auth.AuthProvider); !ok {
		t.Error("AuthProvider in mocked context should be of type *auth.AuthProvider")
	}
	if _, ok := appCtx.PersonService().(person.PersonService); !ok {
		t.Error("Mocked PersonService should be of type *person.MockPersonService")
	}
	if _, ok := appCtx.ItemService().(item.ItemService); !ok {
		t.Error("ItemService should be of type item.ItemService")
	}
	if _, ok := appCtx.InvoiceService().(invoice.InvoiceService); !ok {
		t.Error("InvoiceService should be of type invoice.InvoiceService")
	}
}
