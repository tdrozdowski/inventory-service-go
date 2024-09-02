package person

import "github.com/jmoiron/sqlx"
import "github.com/google/wire"

func InitializePersonService(db sqlx.DB) PersonServiceImpl {
	wire.Build(NewPersonService, NewPersonRepository)
	return PersonServiceImpl{}
}
