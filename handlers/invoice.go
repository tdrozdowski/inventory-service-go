package handlers

import (
	"github.com/labstack/echo/v4"
	"inventory-service-go/context"
	"net/http"
)

func InvoiceRoutes(g *echo.Group, a context.ApplicationContext) {
	g.GET("/invoices", GetAllInvoices(a))
}

func GetAllInvoices(a context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		pagination := paginationFromRequest(c)
		results, err := a.InvoiceService().GetAllInvoices(pagination)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, results)
	}
}
