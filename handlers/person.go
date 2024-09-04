package handlers

import (
	uuid2 "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"inventory-service-go/commons"
	"inventory-service-go/context"
	"net/http"
	"strconv"
)

func paginationFromRequest(c echo.Context) *commons.Pagination {
	lastIdStr := c.QueryParam("last_id")
	lastId, err := strconv.Atoi(lastIdStr)
	if err != nil {
		lastId = 0
	}
	pageSizeStr := c.QueryParam("page_size")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}
	return &commons.Pagination{
		LastId:   lastId,
		PageSize: pageSize,
	}
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
		person, err := personService.GetById(uuid)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, person)
	}
}
