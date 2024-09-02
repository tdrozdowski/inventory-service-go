package main

import "inventory-service-go/person"

type ApplicationContext struct {
	PersonService person.PersonService
}

func NewApplicationContext() *ApplicationContext {
	return &ApplicationContext{
		PersonService: person.InitializePersonService(),
	}
}
