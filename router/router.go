package router

import (
	// "main/middleware"
	"main/controller"
	"main/middleware"

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
		user.GET("", middleware.IsLoggedIn, c.UserCtrl.GetUserList)
		user.GET(":id", middleware.IsLoggedIn, c.UserCtrl.GetUserById)
		user.POST("", c.UserCtrl.RegisterUser)
		user.PUT("", c.UserCtrl.RegisterUser)
		user.DELETE("", c.UserCtrl.DeleteUser)

		cmp := api.Group("/cmp")
		cmp.GET("", c.CmpCtrl.GetCmp)
		cmp.GET("/all", c.CmpCtrl.GetAllCmp)
		cmp.POST("", c.CmpCtrl.RegisterCmp)
		cmp.PUT("", c.CmpCtrl.UpdateCmp)
		cmp.DELETE("", c.CmpCtrl.DeleteCmp)

		jobs := api.Group("/jobs")
		jobs.GET("", c.UserCtrl.RegisterUser)
		jobs.POST("", c.UserCtrl.RegisterUser)
		jobs.PUT("")
		jobs.DELETE("", c.UserCtrl.DeleteUser)
	}
	router.Run(":8080")
}
