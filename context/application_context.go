package context

import (
	"inventory-service-go/auth"
	"inventory-service-go/invoice"
	"inventory-service-go/item"
	"inventory-service-go/person"
)

type ApplicationContext struct {
	personService  person.PersonService
	itemService    item.ItemService
	invoiceService invoice.InvoiceService
	authProvider   auth.AuthProvider
}

func NewApplicationContext() ApplicationContext {
	p, err := person.InitializePersonService()
	if err != nil {
		panic(err)
	}
	i, err := item.InitializeItemService()
	if err != nil {
		panic(err)
	}
	inv, err := invoice.InitializeInvoiceService()
	if err != nil {
		panic(err)
	}
	return ApplicationContext{
		personService:  p,
		itemService:    i,
		invoiceService: inv,
		authProvider:   auth.NewAuthProvider(),
	}
}

func MockApplicationContext(mockPersonService person.PersonService, mockItemService item.ItemService, mockInvoiceService invoice.InvoiceService) ApplicationContext {
	return ApplicationContext{
		personService:  mockPersonService,
		itemService:    mockItemService,
		invoiceService: mockInvoiceService,
		authProvider:   auth.NewJwtAuthProvider("dummy_secret"),
	}
}

func (a ApplicationContext) PersonService() person.PersonService {
	return a.personService
}

func (a ApplicationContext) AuthProvider() auth.AuthProvider {
	return a.authProvider
}

func (a ApplicationContext) ItemService() item.ItemService {
	return a.itemService
}

func (a ApplicationContext) InvoiceService() invoice.InvoiceService {
	return a.invoiceService
}
