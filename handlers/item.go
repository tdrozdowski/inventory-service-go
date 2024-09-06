package handlers

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"inventory-service-go/context"
	"inventory-service-go/item"
	"net/http"
)

func ItemRoutes(p *echo.Group, appContext context.ApplicationContext) {
	p.GET("/items", AllItems(appContext))
	p.GET("/items/:id", GetItem(appContext))
	p.POST("/items", CreateItem(appContext))
	p.PUT("/items/:id", UpdateItem(appContext))
	p.DELETE("/items/:id", DeleteItem(appContext))
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

func UpdateItem(appContext context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		var updateItemRequest item.UpdateItemRequest
		err := c.Bind(&updateItemRequest)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		if id != updateItemRequest.Id {
			return c.JSON(http.StatusBadRequest, "id in path does not match id in body")
		}
		itemService := appContext.ItemService()
		results, err := itemService.UpdateItem(updateItemRequest)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, results)
	}
}

func DeleteItem(appContext context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		itemService := appContext.ItemService()
		results, err := itemService.DeleteItem(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.JSON(http.StatusNotFound, err)
			} else {
				return c.JSON(http.StatusInternalServerError, err)
			}
		}
		return c.JSON(http.StatusOK, results)
	}
}

func GetItem(appContext context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		itemService := appContext.ItemService()
		results, err := itemService.GetItem(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.JSON(http.StatusNotFound, err)
			} else {
				return c.JSON(http.StatusInternalServerError, err)
			}
		}
		return c.JSON(http.StatusOK, results)
	}
}
