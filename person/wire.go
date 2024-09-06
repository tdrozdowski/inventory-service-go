//go:build wireinject
// +build wireinject

package person

import (
	"github.com/google/wire"
	"inventory-service-go/commons"
)

func InitializePersonService() (PersonService, error) {
	wire.Build(
		NewPersonService,
		NewPersonRepository,
		commons.GetDB,
		wire.Bind(new(PersonRepository), new(*PersonRepositoryImpl)),
		wire.Bind(new(PersonService), new(*PersonServiceImpl)),
	)
	return nil, nil
}
