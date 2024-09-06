package handlers

import (
	"github.com/labstack/echo/v4"
	"inventory-service-go/context"
	"inventory-service-go/item"
	"net/http"
)

func ItemRoutes(p *echo.Group, appContext context.ApplicationContext) {
	p.GET("/items", AllItems(appContext))
	p.POST("/items", CreateItem(appContext))
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

func CreateItem(appContext context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		var createItemRequest item.CreateItemRequest
		err := c.Bind(&createItemRequest)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		itemService := appContext.ItemService()
		results, err := itemService.CreateItem(createItemRequest)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, results)
	}
}
