//go:build wireinject
// +build wireinject

package invoice

import (
	"github.com/google/wire"
	"inventory-service-go/commons"
)

func InitializeInvoiceService() (InvoiceService, error) {
	wire.Build(
		NewInvoiceService,
		NewInvoiceRepository,
		commons.GetDB,
		wire.Bind(new(InvoiceService), new(*InvoiceServiceImpl)),
		wire.Bind(new(InvoiceRepository), new(*InvoiceRepositoryImpl)),
	)
	return nil, nil
}
