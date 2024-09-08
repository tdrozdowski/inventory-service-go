package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"inventory-service-go/context"
	"inventory-service-go/invoice"
	"net/http"
)

func InvoiceRoutes(g *echo.Group, a context.ApplicationContext) {
	g.POST("/invoices/items", AddItemsToInvoice(a))
	g.POST("/invoices", CreateInvoice(a))
	g.DELETE("/invoices/:id", DeleteInvoice(a))
	g.GET("/invoices", GetAllInvoices(a))
	g.GET("/invoices/:id", GetInvoice(a))
	g.GET("/invoices/user/:userId", GetAllInvoicesForUser(a))
	g.PUT("/invoices/:id", UpdateInvoice(a))
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

func UpdateInvoice(a context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		var request invoice.UpdateInvoiceRequest
		err := c.Bind(&request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		if id != request.Id {
			return c.JSON(http.StatusBadRequest, "id in path does not match id in body")
		}
		result, err := a.InvoiceService().UpdateInvoice(request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, result)
	}
}

func GetInvoice(a context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		withItems := c.QueryParam("withItems") == "true"
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		result, err := a.InvoiceService().GetInvoice(id, withItems)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, result)
	}
}

func DeleteInvoice(a context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		results, err := a.InvoiceService().DeleteInvoice(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, results)
	}
}

func GetAllInvoicesForUser(a context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		userIdParam := c.Param("userId")
		userId, err := uuid.Parse(userIdParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		results, err := a.InvoiceService().GetInvoicesForUser(userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, results)
	}
}

func AddItemsToInvoice(a context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		var request invoice.ItemsToInvoiceRequest
		err := c.Bind(&request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		response, err := a.InvoiceService().AddItemsToInvoice(request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, response)
	}
}
