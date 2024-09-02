package person

import (
	"github.com/labstack/echo/v4"
	"inventory-service-go/commons"
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

func GetAll(c echo.Context) error {
	pagination := paginationFromRequest(c)
	personService, ok := c.Get("personService").(PersonService)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "peronService not found")
	}
	persons, err := personService.GetAll(pagination)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, persons)
}
