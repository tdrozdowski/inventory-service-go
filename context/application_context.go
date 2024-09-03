package context

import "inventory-service-go/person"

type ApplicationContext struct {
	personService person.PersonService
}

func NewApplicationContext() ApplicationContext {
	return ApplicationContext{
		personService: person.InitializePersonService(),
	}
}

func MockApplicationContext(mockPersonService person.PersonService) ApplicationContext {
	return ApplicationContext{
		personService: mockPersonService,
	}
}

func (a ApplicationContext) PersonService() person.PersonService {
	return a.personService
}
