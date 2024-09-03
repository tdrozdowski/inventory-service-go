package handlers

import (
	"github.com/labstack/echo/v4"
	"inventory-service-go/auth"
	"inventory-service-go/context"
	"net/http"
	"time"
)

type TokenCredentials struct {
	Token     string `json:"token"`
	CreatedAt int64  `json:"createdAt"`
}

func Authorize(appContext context.ApplicationContext) func(context2 echo.Context) error {
	return func(c echo.Context) error {
		// read the auth.Credentials from the request body
		credentials := new(auth.Credentials)
		if err := c.Bind(credentials); err != nil {
			return c.String(http.StatusBadRequest, "Invalid request body")
		}
		token, err := appContext.AuthProvider().Authenticate(credentials.ClientId, credentials.ClientSecret)
		if err != nil {
			return c.String(http.StatusUnauthorized, "Invalid client credentials")
		} else {
			return c.JSON(http.StatusOK, TokenCredentials{
				Token:     token,
				CreatedAt: time.Now().Unix(),
			})
		}
	}
}
