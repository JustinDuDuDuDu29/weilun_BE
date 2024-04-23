package router

import (
	// "main/middleware"
	"main/controller"
	"main/middleware"

	"github.com/gin-gonic/gin"
)

func RouterInit(c *controller.AppControllerImpl, m *middleware.AppMiddlewareImpl) {

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
		jobs.GET("", middleware.IsLoggedIn, m.RoleMid.SuperAdminOnly, c.JobsCtrl.GetAllJob)
		jobs.POST("", middleware.IsLoggedIn, m.RoleMid.SuperAdminOnly, c.JobsCtrl.CreateJob)
		jobs.PUT("", middleware.IsLoggedIn, m.RoleMid.SuperAdminOnly, c.JobsCtrl.UpdateJob)
		jobs.DELETE(":id", middleware.IsLoggedIn, m.RoleMid.SuperAdminOnly, c.JobsCtrl.DeleteJob)

		claimed := api.Group("/claimed")
		claimed.GET("", middleware.IsLoggedIn, m.RoleMid.SuperAdminOnly, c.JobsCtrl.GetAllClaimedJobs)
		claimed.GET("/current", middleware.IsLoggedIn, c.JobsCtrl.GetCurrentClaimedJob)
		claimed.POST(":id", middleware.IsLoggedIn, c.JobsCtrl.ClaimJob)
		claimed.POST("/finish/:id", middleware.IsLoggedIn, c.JobsCtrl.FinishClaimJob)
		claimed.DELETE(":id", middleware.IsLoggedIn, c.JobsCtrl.CancelClaimJob)

		repair := api.Group("/repair")
		repair.GET("", middleware.IsLoggedIn)
		repair.POST(":id", middleware.IsLoggedIn, c.RepairCtrl.CreateNewRepair)
		repair.POST("/finish/:id", middleware.IsLoggedIn)
		repair.DELETE(":id", middleware.IsLoggedIn)

	}
	router.Run(":8080")
}
