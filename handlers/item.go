package handlers

import (
	"github.com/labstack/echo/v4"
	"inventory-service-go/context"
	"net/http"
)

func ItemRoutes(p *echo.Group, appContext context.ApplicationContext) {
	p.GET("", AllItems(appContext))
}

func AllItems(appContext context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		pagination := paginationFromRequest(c)
		itemService := appContext.ItemService()
		items, err := itemService.GetItems(pagination)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, items)
	}
}
