package main

import (
	"github.com/jmoiron/sqlx"
	"inventory-service-go/person"
)

// ApplicationContext is a struct that holds all the services that are used in the application
type ApplicationContext struct {
	PersonService person.PersonService
}

func NewApplicationContext() ApplicationContext {
	// use wire to generate this
	db := sqlx.NewDb(nil, "sqlmock")
	return ApplicationContext{
		PersonService: person.NewPersonService(person.NewPersonRepository(db)),
	}
}
