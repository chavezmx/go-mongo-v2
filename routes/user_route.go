package routes

import (
	"go-mongo-v2/controllers"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	e.GET("/user/:userId", controllers.GetUser)
	e.POST("/user", controllers.CreateUser)
	e.GET("/users", controllers.GetAllUsers)
}