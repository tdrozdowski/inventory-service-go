package context

import (
	"go.uber.org/mock/gomock"
	"inventory-service-go/auth"
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
}

func TestMockApplicationContext(t *testing.T) {
	controller := gomock.NewController(t)
	mockPersonService := person.NewMockPersonService(controller)
	mockItemService := item.NewMockItemService(controller)
	appCtx := MockApplicationContext(mockPersonService, *mockItemService)
	if _, ok := appCtx.AuthProvider().(auth.AuthProvider); !ok {
		t.Error("AuthProvider in mocked context should be of type *auth.AuthProvider")
	}
	if _, ok := appCtx.PersonService().(person.PersonService); !ok {
		t.Error("Mocked PersonService should be of type *person.MockPersonService")
	}
}
