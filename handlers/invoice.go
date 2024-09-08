package handlers

import (
	"github.com/labstack/echo/v4"
	"inventory-service-go/context"
	"inventory-service-go/invoice"
	"net/http"
)

func InvoiceRoutes(g *echo.Group, a context.ApplicationContext) {
	g.POST("/invoices", CreateInvoice(a))
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

func CreateInvoice(a context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		var request invoice.CreateInvoiceRequest
		err := c.Bind(&request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		result, err := a.InvoiceService().CreateInvoice(request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, result)
	}
}
