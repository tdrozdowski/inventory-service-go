package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"inventory-service-go/context"
	"inventory-service-go/handlers"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	e := echo.New()
	appContext := context.NewApplicationContext()
	e.GET("/api/v1/persons", handlers.GetAll(appContext))
	//middelwares
	e.Use(middleware.CORS())
	// Start the server
	err := e.Start(":8080")
	if err != nil {
		return
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
