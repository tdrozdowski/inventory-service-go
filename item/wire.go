//go:build wireinject
// +build wireinject

package item

import (
	"github.com/google/wire"
	"inventory-service-go/commons"
)

func InitializeItemService() (ItemService, error) {
	wire.Build(
		NewItemService,
		NewItemRepository,
		commons.GetDB,
		wire.Bind(new(ItemService), new(*ItemServiceImpl)),
		wire.Bind(new(ItemRepository), new(*ItemRepositoryImpl)),
	)
	return nil, nil
}
