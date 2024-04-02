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
		user := api.Group("/user")
		user.GET("", c.UserCtrl.RegisterNewUser)
	}
	router.Run(":8080")
}
