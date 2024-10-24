package router

import (
	// "main/middleware"

	"main/controller"
	"main/middleware"
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

		auth := api.Group("/auth")
		auth.POST("", c.AuthCtrl.Login)

		io := api.Group("/io")
		io.GET("/", c.SocketCtrl.TestSocket)

		static := api.Group("/static", m.RoleMid.IsLoggedIn)
		static.StaticFS("/img", http.Dir("./img"))

		user := api.Group("/user")
		user.GET("", m.RoleMid.IsLoggedIn, c.UserCtrl.GetUserList)
		user.GET(":id", m.RoleMid.IsLoggedIn, c.UserCtrl.GetUserById)
		user.POST("/me", m.RoleMid.IsLoggedIn, c.UserCtrl.Me)
		user.POST("", m.RoleMid.IsLoggedIn, c.UserCtrl.RegisterUser)
		user.POST(":id", m.RoleMid.IsLoggedIn, c.UserCtrl.RegisterUser)
		user.POST("/pwd", m.RoleMid.IsLoggedIn, c.UserCtrl.UpdatePassword)
		user.POST("/pwdreset", m.RoleMid.IsLoggedIn, m.RoleMid.CmpSuperAdminOnly, c.UserCtrl.ResetPassword)
		user.POST("/approve/:id", m.RoleMid.IsLoggedIn, m.RoleMid.CmpSuperAdminOnly, c.UserCtrl.ApproveUser)
		user.POST("/UpdateDriverPic", m.RoleMid.IsLoggedIn, c.UserCtrl.UpdateDriverPic)
		user.PUT(":id", m.RoleMid.IsLoggedIn, c.UserCtrl.UpdateUser)
		user.DELETE("", m.RoleMid.IsLoggedIn, c.UserCtrl.DeleteUser)

		cmp := api.Group("/cmp")
		cmp.GET("", m.RoleMid.IsLoggedIn, c.CmpCtrl.GetCmp)
		cmp.GET("/all", m.RoleMid.IsLoggedIn, c.CmpCtrl.GetAllCmp)
		cmp.POST("", m.RoleMid.IsLoggedIn, c.CmpCtrl.RegisterCmp)
		cmp.PUT("", m.RoleMid.IsLoggedIn, c.CmpCtrl.UpdateCmp)
		cmp.DELETE("", m.RoleMid.IsLoggedIn, c.CmpCtrl.DeleteCmp)

		jobs := api.Group("/jobs")
		jobs.POST("/all", m.RoleMid.IsLoggedIn, c.JobsCtrl.GetAllJob)
		jobs.POST("", m.RoleMid.IsLoggedIn, m.RoleMid.CmpSuperAdminOnly, c.JobsCtrl.CreateJob)
		jobs.PUT("", m.RoleMid.IsLoggedIn, m.RoleMid.CmpSuperAdminOnly, c.JobsCtrl.UpdateJob)
		jobs.DELETE(":id", m.RoleMid.IsLoggedIn, m.RoleMid.CmpSuperAdminOnly, c.JobsCtrl.DeleteJob)

		claimed := api.Group("/claimed")
		// sec?
		claimed.GET("", m.RoleMid.IsLoggedIn, c.JobsCtrl.GetAllClaimedJobs)
		claimed.GET("/cj", m.RoleMid.IsLoggedIn, c.JobsCtrl.GetCJDate)
		claimed.GET("/list", m.RoleMid.IsLoggedIn, c.JobsCtrl.GetClaimedJobByDriverID)
		claimed.GET("/pending", m.RoleMid.IsLoggedIn, c.JobsCtrl.GetClaimedJobByDriverID)
		claimed.GET("/userwitpendingjob", m.RoleMid.IsLoggedIn, m.RoleMid.CmpAdminOnly, c.JobsCtrl.GetUserWithPendingJob)
		claimed.POST("/current", m.RoleMid.IsLoggedIn, c.JobsCtrl.GetCurrentClaimedJob)
		claimed.POST(":id", m.RoleMid.IsLoggedIn, c.JobsCtrl.ClaimJob)
		claimed.POST("/finish/:id", m.RoleMid.IsLoggedIn, c.JobsCtrl.FinishClaimJob)
		claimed.POST("/approve/:id", m.RoleMid.IsLoggedIn, m.RoleMid.CmpSuperAdminOnly, c.JobsCtrl.ApproveClaimedJob)
		claimed.DELETE(":id", m.RoleMid.IsLoggedIn, c.JobsCtrl.CancelClaimJob)

		repair := api.Group("/repair")
		repair.GET("", m.RoleMid.IsLoggedIn, c.RepairCtrl.GetRepair)
		repair.GET("/cj", m.RoleMid.IsLoggedIn, c.RepairCtrl.GetRepairDate)
		repair.POST("", m.RoleMid.IsLoggedIn, c.RepairCtrl.CreateNewRepair)
		repair.POST("/approve/:id", m.RoleMid.IsLoggedIn, c.RepairCtrl.ApproveRepair)
		repair.DELETE(":id", m.RoleMid.IsLoggedIn, c.RepairCtrl.DeleteRepair)

		gas := api.Group("/gas")
		gas.GET("", m.RoleMid.IsLoggedIn, c.GasCtrl.GetGas)
		gas.GET("/cj", m.RoleMid.IsLoggedIn, c.GasCtrl.GetGasDate)
		gas.POST("", m.RoleMid.IsLoggedIn, c.GasCtrl.CreateNewGas)
		gas.POST("/approve/:id", m.RoleMid.IsLoggedIn, c.GasCtrl.ApproveGas)
		gas.DELETE(":id", m.RoleMid.IsLoggedIn, c.GasCtrl.DeleteGas)

		alert := api.Group("/alert")
		alert.GET("", m.RoleMid.IsLoggedIn, c.AlertCtrl.GetAlert)
		alert.POST("", m.RoleMid.IsLoggedIn, m.RoleMid.CmpSuperAdminOnly, c.AlertCtrl.CreateAlert)
		alert.PUT("", m.RoleMid.IsLoggedIn, c.AlertCtrl.CheckNewAlert)
		alert.DELETE(":id", m.RoleMid.IsLoggedIn, c.AlertCtrl.DeleteAlert)

		revenue := api.Group("/revenue")
		revenue.GET("", m.RoleMid.IsLoggedIn, c.RevenueCtrl.RevenueDriver)
		revenue.GET("/excel", m.RoleMid.IsLoggedIn, m.RoleMid.CmpSuperAdminOnly, c.RevenueCtrl.RevenueExcel)
		api.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"version": os.Getenv("version")})
		})
	}
	router.Run(":8080")
}
