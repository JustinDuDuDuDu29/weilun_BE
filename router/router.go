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
		auth.POST("", c.AuthCtrl.Login)

		user := api.Group("/user")
		user.GET("", c.UserCtrl.RegisterUser)
		user.POST("", c.UserCtrl.RegisterUser)
		user.PUT("", c.UserCtrl.RegisterUser)
		user.DELETE("", c.UserCtrl.DeleteUser)

		cmp := api.Group("/cmp")
		cmp.GET("", c.UserCtrl.RegisterUser)
		cmp.POST("", c.UserCtrl.RegisterUser)
		cmp.PUT("")
		cmp.DELETE("", c.UserCtrl.DeleteUser)

		jobs := api.Group("/jobs")
		jobs.GET("", c.UserCtrl.RegisterUser)
		jobs.POST("", c.UserCtrl.RegisterUser)
		jobs.PUT("")
		jobs.DELETE("", c.UserCtrl.DeleteUser)
	}
	router.Run(":8080")
}
