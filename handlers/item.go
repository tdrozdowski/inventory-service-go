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

// AllItems
//
//		@Summary		List Items
//		@Description	List all Items
//		@Id				all_items
//		@Tags			item
//		@Produce		json
//		@Param			last_id		query		int	false	"last seq id"
//	 	@Param			page_size 	query		int false 	"number of items per page"
//		@Success		200	{array}		item.Item			"OK"
//		@Failure		400	{string}	string 				"Bad Request"
//		@Failure		500	{string}	string 				"Internal Server Error"
//		@Router			/items [get]
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

// CreateItem
//
//		@Summary		Create Item
//		@Description	Create an Item
//		@ID				create_item
//		@Tags			item
//		@Accept			json
//		@Produce		json
//	    @Param 			request body 		item.CreateItemRequest	true 	"Create Item Request"
//		@Success		201		{object}	item.Item						"Created"
//		@Failure		400		{string}	string					"Bad Request"
//		@Failure		401		{string}	string					"Unauthorized (invalid credentials)"
//		@Failure		500		{object}	error					"Internal Server Error"
//		@Router			/items [post]
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
		return c.JSON(http.StatusCreated, results)
	}
}

// UpdateItem
//
//		@Summary		Update Item
//		@Description	Update an Item
//		@ID				update_item
//		@Tags			item
//		@Accept			json
//		@Produce		json
//	    @Param 			request body 		item.UpdateItemRequest		true 	"Update Item Request"
//		@Param			id	path			uuid.Uuid					true	"Invoice Id"
//		@Success		200		{object}	item.Item				"OK"
//		@Failure		400		{string}	string					"Bad Request"
//		@Failure		401		{string}	string					"Unauthorized (invalid credentials)"
//		@Failure		500		{object}	error					"Internal Server Error"
//		@Router			/items/{id} [put]
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

// DeleteItem
//
//		@Summary		Delete Item
//		@Description	Remove a specific Item
//		@Id				delete_item
//		@Tags			item
//		@Produce		json
//	 	@Param			id				query		uuid.Uuid 	true 	"id of the item to be deleted"
//		@Success		200	{array}		commons.DeleteResult	"OK"
//		@Failure		400	{string}	string 					"Bad Request"
//		@Failure		404 {string} 	string					"Not Found"
//		@Failure		500	{string}	string 					"Internal Server Error"
//		@Router			/items/{id} [delete]
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

// GetItem
//
//		@Summary		Get Item
//		@Description	Get a specific Item
//		@Id				get_item
//		@Tags			item
//		@Produce		json
//	 	@Param			id				query		uuid.Uuid 	true 	"id of the item requested"
//		@Success		200	{array}		item.Item		 	"OK"
//		@Failure		400	{string}	string 				"Bad Request"
//		@Failure		404 {string} 	string				"Not Found"
//		@Failure		500	{string}	string 				"Internal Server Error"
//		@Router			/items/{id} [get]
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
