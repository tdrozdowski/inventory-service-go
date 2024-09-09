package handlers

import (
	uuid2 "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"inventory-service-go/commons"
	"inventory-service-go/context"
	"inventory-service-go/person"
	"net/http"
	"strconv"
)

func paginationFromRequest(c echo.Context) *commons.Pagination {
	lastIdStr := c.QueryParam("last_id")
	pageSizeStr := c.QueryParam("page_size")
	if lastIdStr == "" && pageSizeStr == "" {
		return nil
	} else {
		lastId, err := strconv.Atoi(lastIdStr)
		if err != nil {
			lastId = 0
		}
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			pageSize = 10
		}
		return &commons.Pagination{
			LastId:   lastId,
			PageSize: pageSize,
		}
	}
}

func PersonRoutes(p *echo.Group, appContext context.ApplicationContext) {
	p.GET("/persons", GetAllPersons(appContext))
	p.GET("/persons/:id", GetPersonById(appContext))
	p.POST("/persons", CreatePerson(appContext))
	p.PUT("/persons/:id", UpdatePerson(appContext))
	p.DELETE("/persons/:id", DeletePerson(appContext))
}

// GetAllPersons
//
//		@Summary		List Persons
//		@Description	List all Persons
//		@Id				all_persons
//		@Tags			person
//		@Produce		json
//		@Param			last_id		query		int	false	"last seq id"
//	 	@Param			page_size 	query		int false 	"number of persons per page"
//		@Success		200	{array}		person.Person		"OK"
//		@Failure		400	{string}	string 				"Bad Request"
//		@Failure		500	{string}	string 				"Internal Server Error"
//		@Router			/persons [get]
func GetAllPersons(appContext context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		pagination := paginationFromRequest(c)
		personService := appContext.PersonService()
		persons, err := personService.GetAll(pagination)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, persons)
	}
}

// GetPersonById
//
//		@Summary		Get Person
//		@Description	Get a specific Person
//		@Id				get_person
//		@Tags			person
//		@Produce		json
//	 	@Param			id				query		uuid.Uuid 	true 	"id of the Person requested"
//		@Success		200	{array}		person.Person	 	"OK"
//		@Failure		400	{string}	string 				"Bad Request"
//		@Failure		404 {string} 	string				"Not Found"
//		@Failure		500	{string}	string 				"Internal Server Error"
//		@Router			/persons/{id} [get]
func GetPersonById(appContext context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")
		uuid, err := uuid2.Parse(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		personService := appContext.PersonService()
		p, err := personService.GetById(uuid)
		err2 := commons.HandleServiceError(c, err)
		if err2 != nil {
			return err2
		}
		return c.JSON(http.StatusOK, p)
	}
}

// CreatePerson
//
//		@Summary		Create Person
//		@Description	Create an Person
//		@ID				create_person
//		@Tags			person
//		@Accept			json
//		@Produce		json
//	    @Param 			request body 		person.CreatePersonRequest	true 	"Create Person Request"
//		@Success		201		{object}	person.Person						"Created"
//		@Failure		400		{string}	string					"Bad Request"
//		@Failure		401		{string}	string					"Unauthorized (invalid credentials)"
//		@Failure		500		{object}	error					"Internal Server Error"
//		@Router			/persons [post]
func CreatePerson(appContext context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		var createPersonRequest person.CreatePersonRequest
		if err := c.Bind(&createPersonRequest); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		personService := appContext.PersonService()
		results, err := personService.Create(createPersonRequest)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusCreated, results)
	}
}

// UpdatePerson
//
//		@Summary		Update Person
//		@Description	Update an Person
//		@ID				update_person
//		@Tags			person
//		@Accept			json
//		@Produce		json
//	    @Param 			request body 		person.UpdatePersonRequest		true 	"Update Person Request"
//		@Param			id	path			uuid.Uuid					true	"Person Id"
//		@Success		200		{object}	person.Person				"OK"
//		@Failure		400		{string}	string					"Bad Request"
//		@Failure		401		{string}	string					"Unauthorized (invalid credentials)"
//		@Failure		500		{object}	error					"Internal Server Error"
//		@Router			/persons/{id} [put]
func UpdatePerson(appContext context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		var updatePersonRequest person.UpdatePersonRequest
		if err := c.Bind(&updatePersonRequest); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		personService := appContext.PersonService()
		results, err := personService.Update(updatePersonRequest)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, results)
	}
}

// DeletePerson
//
//		@Summary		Delete Person
//		@Description	Remove a specific Person
//		@Id				delete_person
//		@Tags			person
//		@Produce		json
//	 	@Param			id				query		uuid.Uuid 	true 	"id of the Person to be deleted"
//		@Success		200	{array}		commons.DeleteResult	"OK"
//		@Failure		400	{string}	string 					"Bad Request"
//		@Failure		404 {string} 	string					"Not Found"
//		@Failure		500	{string}	string 					"Internal Server Error"
//		@Router			/persons/{id} [delete]
func DeletePerson(appContext context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")
		uuid, err := uuid2.Parse(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		personService := appContext.PersonService()
		results, err := personService.DeleteByUuid(uuid)
		err2 := commons.HandleServiceError(c, err)
		if err2 != nil {
			return err2
		}
		return c.JSON(http.StatusOK, results)
	}
}
