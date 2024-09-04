package context

import (
	"inventory-service-go/auth"
	"inventory-service-go/person"
)

type ApplicationContext struct {
	personService person.PersonService
	authProvider  auth.AuthProvider
}

func NewApplicationContext() ApplicationContext {
	return ApplicationContext{
		personService: person.InitializePersonService(),
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
