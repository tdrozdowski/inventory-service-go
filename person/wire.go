//go:build wireinject
// +build wireinject

package person

import (
	"github.com/google/wire"
	"inventory-service-go/commons"
)

func InitializePersonService() PersonService {
	wire.Build(NewPersonService, NewPersonRepository, commons.GetDB)
	return PersonService{}
}
