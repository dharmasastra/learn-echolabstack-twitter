package main

import (
	"github.com/dharmasastra/learn-echolabstack-twitter/app/controllers"
	"github.com/dharmasastra/learn-echolabstack-twitter/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := config.NewRouter()
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(controllers.Key),
		Skipper: func(c echo.Context) bool {
			// Skip authentication for and signup login requests
			if c.Path() == "/login" || c.Path() == "/signup" {
				return true
			}
			return false
		},
	}))

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
