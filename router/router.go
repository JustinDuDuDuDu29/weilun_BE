package router

import (
	// "main/middleware"
	"main/controller"
	"main/middleware"
	"main/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RouterInit(c *controller.AppControllerImpl, m *middleware.AppMiddlewareImpl) {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{

		api.GET("", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"version": os.Getenv("version")})
		})

		auth := api.Group("/auth")
		auth.POST("", c.AuthCtrl.Login)

		io := api.Group("/io", m.RoleMid.IsLoggedIn)
		io.GET("", utils.HandleWS)

		static := api.Group("/static", m.RoleMid.IsLoggedIn)
		static.StaticFS("/img", http.Dir("./img"))

		user := api.Group("/user")
		user.GET("", m.RoleMid.IsLoggedIn, c.UserCtrl.GetUserList)
		user.GET(":id", m.RoleMid.IsLoggedIn, c.UserCtrl.GetUserById)
		user.POST("", m.RoleMid.IsLoggedIn, c.UserCtrl.RegisterUser)
		user.POST(":id", m.RoleMid.IsLoggedIn, c.UserCtrl.RegisterUser)
		user.POST("/pwd", m.RoleMid.IsLoggedIn, c.UserCtrl.UpdatePassword)
		user.POST("/UpdateDriverPic", m.RoleMid.IsLoggedIn, c.UserCtrl.UpdateDriverPic)
		user.PUT("", m.RoleMid.IsLoggedIn, c.UserCtrl.RegisterUser)
		user.DELETE("", m.RoleMid.IsLoggedIn, c.UserCtrl.DeleteUser)

		cmp := api.Group("/cmp")
		cmp.GET("", m.RoleMid.IsLoggedIn, c.CmpCtrl.GetCmp)
		cmp.GET("/all", m.RoleMid.IsLoggedIn, c.CmpCtrl.GetAllCmp)
		cmp.POST("", m.RoleMid.IsLoggedIn, c.CmpCtrl.RegisterCmp)
		cmp.PUT("", m.RoleMid.IsLoggedIn, c.CmpCtrl.UpdateCmp)
		cmp.DELETE("", m.RoleMid.IsLoggedIn, c.CmpCtrl.DeleteCmp)

		jobs := api.Group("/jobs")
		jobs.GET("", m.RoleMid.IsLoggedIn, c.JobsCtrl.GetAllJob)
		jobs.POST("", m.RoleMid.IsLoggedIn, m.RoleMid.SuperAdminOnly, c.JobsCtrl.CreateJob)
		jobs.PUT("", m.RoleMid.IsLoggedIn, m.RoleMid.SuperAdminOnly, c.JobsCtrl.UpdateJob)
		jobs.DELETE(":id", m.RoleMid.IsLoggedIn, m.RoleMid.SuperAdminOnly, c.JobsCtrl.DeleteJob)

		claimed := api.Group("/claimed")
		claimed.GET("", m.RoleMid.IsLoggedIn, m.RoleMid.SuperAdminOnly, c.JobsCtrl.GetAllClaimedJobs)
		claimed.GET("/current", m.RoleMid.IsLoggedIn, c.JobsCtrl.GetCurrentClaimedJob)
		claimed.POST(":id", m.RoleMid.IsLoggedIn, c.JobsCtrl.ClaimJob)
		claimed.POST("/finish/:id", m.RoleMid.IsLoggedIn, c.JobsCtrl.FinishClaimJob)
		claimed.POST("/approve/:id", m.RoleMid.IsLoggedIn, c.JobsCtrl.ApproveClaimedJob)
		claimed.DELETE(":id", m.RoleMid.IsLoggedIn, c.JobsCtrl.CancelClaimJob)

		repair := api.Group("/repair")
		repair.GET("", m.RoleMid.IsLoggedIn, c.RepairCtrl.GetRepair)
		repair.POST("", m.RoleMid.IsLoggedIn, c.RepairCtrl.CreateNewRepair)
		repair.POST(":id", m.RoleMid.IsLoggedIn, c.RepairCtrl.ApproveRepair)
		repair.DELETE(":id", m.RoleMid.IsLoggedIn, c.RepairCtrl.DeleteRepair)

		alert := api.Group("/alert")
		alert.GET("", m.RoleMid.IsLoggedIn, c.AlertCtrl.GetAlert)
		alert.POST("", m.RoleMid.IsLoggedIn, m.RoleMid.SuperAdminOnly, c.AlertCtrl.CreateAlert)
		alert.PUT("", m.RoleMid.IsLoggedIn, c.AlertCtrl.CheckNewAlert)
		alert.DELETE(":id", m.RoleMid.IsLoggedIn, c.AlertCtrl.DeleteAlert)
	}
	router.Run(":8080")
}
