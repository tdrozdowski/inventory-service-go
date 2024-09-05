package context

import (
	"inventory-service-go/auth"
	"inventory-service-go/item"
	"inventory-service-go/person"
)

type ApplicationContext struct {
	personService person.PersonService
	itemService   item.ItemService
	authProvider  auth.AuthProvider
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
	return ApplicationContext{
		personService: p,
		itemService:   i,
		authProvider:  auth.NewAuthProvider(),
	}
}

func MockApplicationContext(mockPersonService person.PersonService) ApplicationContext {
	return ApplicationContext{
		personService: mockPersonService,
		authProvider:  auth.NewJwtAuthProvider("dummy_secret"),
	}
}

func (a ApplicationContext) PersonService() person.PersonService {
	return a.personService
}

func (a ApplicationContext) AuthProvider() auth.AuthProvider {
	return a.authProvider
}
