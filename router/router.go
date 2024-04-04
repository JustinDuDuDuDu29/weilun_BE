package router

import (
	// "main/middleware"
	"main/controller"

	"github.com/gin-gonic/gin"
)

func RouterInit(c *controller.AppControllerImpl) {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		auth.GET("")
		auth.POST("", c.AuthCtrl.Login)
		auth.DELETE("")

		user := api.Group("/user")
		user.POST("", c.UserCtrl.RegisterNewUser)
		user.PUT("", c.UserCtrl.RegisterNewUser)
		user.DELETE("", c.UserCtrl.RegisterNewUser)
	}
	router.Run(":8080")
}
