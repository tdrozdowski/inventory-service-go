package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"inventory-service-go/context"
	"inventory-service-go/invoice"
	"net/http"
)

func InvoiceRoutes(g *echo.Group, a context.ApplicationContext) {
	g.POST("/invoices/:id/items", AddItemsToInvoice(a))
	g.POST("/invoices", CreateInvoice(a))
	g.DELETE("/invoices/:id", DeleteInvoice(a))
	g.GET("/invoices", GetAllInvoices(a))
	g.GET("/invoices/:id", GetInvoice(a))
	g.GET("/invoices/user/:userId", GetAllInvoicesForUser(a))
	g.DELETE("/invoices/:id/items/:itemId", RemoveItemFromInvoice(a))
	g.PUT("/invoices/:id", UpdateInvoice(a))
}

// GetAllInvoices
//
//		@Summary		List Invoices
//		@Description	List all Invoices
//		@Id				all_invoices
//		@Tags			invoice
//		@Produce		json
//		@Param			last_id		query		int	false	"last seq id"
//	 	@Param			page_size 	query		int false 	"number of invoices per page"
//		@Success		200	{array}		invoice.Invoice 	"OK"
//		@Failure		400	{string}	string 				"Bad Request"
//		@Failure		500	{string}	string 				"Internal Server Error"
//		@Router			/invoices [get]
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

// CreateInvoice
//
//		@Summary		Create Invoice
//		@Description	Create an Invoice
//		@ID				create_invoice
//		@Tags			invoice
//		@Accept			json
//		@Produce		json
//	    @Param 			request body 		invoice.CreateInvoiceRequest	true 	"Create Invoice Request"
//		@Success		201		{object}	invoice.Invoice					"Created"
//		@Failure		400		{string}	string					"Bad Request"
//		@Failure		401		{string}	string					"Unauthorized (invalid credentials)"
//		@Failure		500		{object}	error					"Internal Server Error"
//		@Router			/invoices [post]
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

// UpdateInvoice
//
//		@Summary		Update Invoice
//		@Description	Update an Invoice
//		@ID				update_invoice
//		@Tags			invoice
//		@Accept			json
//		@Produce		json
//	    @Param 			request body 		invoice.UpdateInvoiceRequest	true 	"Update Invoice Request"
//		@Param			id	path			uuid.Uuid						true	"Invoice Id"
//		@Success		200		{object}	invoice.Invoice					"OK"
//		@Failure		400		{string}	string					"Bad Request"
//		@Failure		401		{string}	string					"Unauthorized (invalid credentials)"
//		@Failure		500		{object}	error					"Internal Server Error"
//		@Router			/invoices/{id} [put]
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

// GetInvoice
//
//		@Summary		Get Invoice
//		@Description	Get a specific Invoice
//		@Id				get_invoice
//		@Tags			invoice
//		@Produce		json
//	 	@Param			id				query		uuid.Uuid 	true 	"id of the invoice requested"
//		@Param			withItems		query		bool		false	"return with Items (if there are any attached to the invoice)"
//		@Success		200	{array}		invoice.Invoice 	"OK"
//		@Failure		400	{string}	string 				"Bad Request"
//		@Failure		500	{string}	string 				"Internal Server Error"
//		@Router			/invoices/{id} [get]
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

// DeleteInvoice
//
//		@Summary		Delete Invoice
//		@Description	Remove a specific Invoice
//		@Id				delete_invoice
//		@Tags			invoice
//		@Produce		json
//	 	@Param			id				query		uuid.Uuid 	true 	"id of the invoice to be deleted"
//		@Success		200	{array}		commons.DeleteResult	"OK"
//		@Failure		400	{string}	string 					"Bad Request"
//		@Failure		500	{string}	string 					"Internal Server Error"
//		@Router			/invoices/{id} [delete]
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

// GetAllInvoicesForUser
//
//		@Summary		Get Invoices For User
//		@Description	Get all Invoices for a specific User
//		@Id				get_all_for_user_invoice
//		@Tags			invoice
//		@Produce		json
//	 	@Param			id				query		uuid.Uuid 	true 	"id of a user"
//		@Success		200	{array}		invoice.Invoice 	"OK"
//		@Failure		400	{string}	string 				"Bad Request"
//		@Failure		500	{string}	string 				"Internal Server Error"
//		@Router			/invoices/user/{id} [get]
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

// AddItemsToInvoice
//
//		@Summary		Add Items to Invoice
//		@Description	Add Items to an Invoice
//		@Id				add_items_to_invoice
//		@Tags			invoice
//		@Accept			json
//		@Produce		json
//	 	@Param			id				query		uuid.Uuid 	true 	"id of the invoice requested"
//	    @Param 			request body 	invoice.ItemsToInvoiceRequest		true 	"Add Items to Invoice Request"
//		@Success		200	{array}		invoice.ItemsToInvoiceResponse	 	"OK"
//		@Failure		400	{string}	string 								"Bad Request"
//		@Failure		500	{string}	string 								"Internal Server Error"
//		@Router			/invoices/{id}/items [post]
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

// RemoveItemFromInvoice
//
//		@Summary		Remove Item From Invoice
//		@Description	Remove a specific Item from a specific Invoice
//		@Id				remove_item_from_invoice
//		@Tags			invoice
//		@Produce		json
//	 	@Param			id				query		uuid.Uuid 	true 	"id of the invoice"
//	 	@Param			itemId			query		uuid.Uuid 	true 	"id of the item to be removed"
//		@Success		200	{array}		commons.DeleteResult	"OK"
//		@Failure		400	{string}	string 					"Bad Request"
//		@Failure		500	{string}	string 					"Internal Server Error"
//		@Router			/invoices/{id}/items/{itemId} [delete]
func RemoveItemFromInvoice(a context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		itemIdParam := c.Param("itemId")
		invoiceId, err := uuid.Parse(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		itemId, err := uuid.Parse(itemIdParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		results, err := a.InvoiceService().RemoveItemFromInvoice(invoice.SimpleInvoiceItem{invoiceId, itemId})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, results)
	}
}
