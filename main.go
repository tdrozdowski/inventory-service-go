package main

import (
	"fmt"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mvrilo/go-redoc"
	echoredoc "github.com/mvrilo/go-redoc/echo"
	"inventory-service-go/context"
	"inventory-service-go/handlers"
	"log"
	"slices"
)

// @title Inventory Service API
// @version 1.0
// @description This is an implementation of the imventory-service in Go.

// @contact.name API Support
// @contact.url http://localhost/support

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
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
	handlers.ItemRoutes(apiV1, appContext)
	handlers.InvoiceRoutes(apiV1, appContext)

	//middlewares
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echoredoc.New(doc()))
	e.Use(echojwt.WithConfig(echojwt.Config{
		Skipper: func(c echo.Context) bool {
			pathsToSkip := []string{"/api/v1/authorize", "", "/docs", "/docs/swagger.json", "/redoc.standalone.js.map"}
			log.Printf("Path: '%s'", c.Path())
			log.Printf("Will Skip: %v", slices.Contains(pathsToSkip, c.Path()))
			return slices.Contains(pathsToSkip, c.Path())
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

func doc() redoc.Redoc {
	return redoc.Redoc{
		Title:       "Example API",
		Description: "Example API Description",
		SpecFile:    "./docs/swagger.json", // "./openapi.yaml"
		SpecPath:    "/swagger.json",       // "/openapi.yaml"
		DocsPath:    "/docs",
	}
}
