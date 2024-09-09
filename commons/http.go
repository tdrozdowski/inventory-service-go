package commons

import (
	"database/sql"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HandleServiceError(c echo.Context, err error) error {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, err)
		} else {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}
	return err
}
