package config

import (
	"github.com/dharmasastra/learn-echolabstack-twitter/app/controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func NewRouter() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} | ${status} | ${latency_human} | ${remote_ip} | ${method} | ${uri} | ${error}\n",
	}))
	e.Use(middleware.Recover())

	// Initialize handler
	h := &controllers.Handler{DB: conn()}

	// Routes
	e.POST("/signup", h.SignUp)
	e.POST("/login", h.Login)
	e.POST("/follow/:id", h.Follow)
	e.POST("/posts", h.CreatePost)
	e.GET("/feed", h.FetchPost)

	return e
}
