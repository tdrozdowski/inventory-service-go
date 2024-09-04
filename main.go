package main

import (
	"fmt"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"inventory-service-go/context"
	"inventory-service-go/handlers"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	e := echo.New()
	appContext := context.NewApplicationContext()
	apiV1 := e.Group("/api/v1")
	apiV1.POST("/authorize", handlers.Authorize(appContext))
	handlers.PersonRoutes(apiV1, appContext)

	//middlewares
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echojwt.WithConfig(echojwt.Config{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/api/v1/authorize"
		},
		SigningMethod: "HS256",
		SigningKey:    appContext.AuthProvider().GetSecret(),
		ErrorHandler: func(c echo.Context, err error) error {
			fmt.Printf("Authorization header: %v\n", c.Request().Header.Get("Authorization"))
			fmt.Printf("JWT error: %v\n", err)
			return c.String(401, "Unauthorized")
		},
	}))
	// Start the server
	err = e.Start(":8080")
	if err != nil {
		return
	}
}
