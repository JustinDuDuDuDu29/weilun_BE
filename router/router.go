package router

import (
	"main/controller"
	"main/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	app := gin.Default()
	app.GET("/login", controller.Login)
	app.POST("/", middleware.IsLoggedIn, controller.Login)
	app.POST("/register", controller.RegisterNewUser)
	app.Run()
	return app
}
