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
		user.POST("", middleware.IsLoggedIn, c.UserCtrl.RegisterUser)
		user.PUT("", middleware.IsLoggedIn, c.UserCtrl.RegisterUser)
		user.DELETE("", middleware.IsLoggedIn, c.UserCtrl.DeleteUser)

		cmp := api.Group("/cmp")
		cmp.GET("", middleware.IsLoggedIn, c.CmpCtrl.GetCmp)
		cmp.GET("/all", middleware.IsLoggedIn, c.CmpCtrl.GetAllCmp)
		cmp.POST("", middleware.IsLoggedIn, c.CmpCtrl.RegisterCmp)
		cmp.PUT("", middleware.IsLoggedIn, c.CmpCtrl.UpdateCmp)
		cmp.DELETE("", middleware.IsLoggedIn, c.CmpCtrl.DeleteCmp)

		jobs := api.Group("/jobs")
		jobs.GET("", middleware.IsLoggedIn, c.JobsCtrl.GetAllJob)
		jobs.POST("", middleware.IsLoggedIn, c.JobsCtrl.CreateJob)
		jobs.PUT("")
		jobs.DELETE("", middleware.IsLoggedIn, c.JobsCtrl.DeleteJob)
	}
	router.Run(":8080")
}
