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
	p.GET("/persons", GetAll(appContext))
	p.GET("/persons/:id", GetById(appContext))
	p.POST("/persons", CreatePerson(appContext))
	p.PUT("/persons/:id", UpdatePerson(appContext))
	p.DELETE("/persons/:id", DeletePerson(appContext))
}

func GetAll(appContext context.ApplicationContext) func(c echo.Context) error {
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

func GetById(appContext context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")
		uuid, err := uuid2.Parse(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		personService := appContext.PersonService()
		p, err := personService.GetById(uuid)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, p)
	}
}

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

func DeletePerson(appContext context.ApplicationContext) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")
		uuid, err := uuid2.Parse(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		personService := appContext.PersonService()
		results, err := personService.DeleteByUuid(uuid)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, results)
	}
}
