// go:build wireinject
//go:build wireinject
// +build wireinject

package utils

import (
	"database/sql"
	"main/controller"

	"main/middleware"
	"main/service"
	db "main/sql"

	"github.com/google/wire"
)

var alertServSet = wire.NewSet(
	service.AlertServInit,
	wire.Bind(new(service.AlertServ), new(*service.AlertServImpl)),
)

var cmpServSet = wire.NewSet(
	service.CmpServInit,
	wire.Bind(new(service.CmpServ), new(*service.CmpServImpl)),
)
var jobsServSet = wire.NewSet(
	service.JobsServInit,
	wire.Bind(new(service.JobsServ), new(*service.JobsServImpl)),
)

var userServSet = wire.NewSet(
	service.UserServInit,
	wire.Bind(new(service.UserServ), new(*service.UserServImpl)),
)

var repairServSet = wire.NewSet(
	service.RepairServInit,
	wire.Bind(new(service.RepairServ), new(*service.RepairServImpl)),
)

var userCtrlSet = wire.NewSet(
	controller.UserCtrlInit,
	wire.Bind(new(controller.UserCtrl), new(*controller.UserCtrlImpl)),
)

var alertCtrlSet = wire.NewSet(
	controller.AlertCtrlInit,
	wire.Bind(new(controller.AlertCtrl), new(*controller.AlertCtrlImpl)),
)

var jobsCtrlSet = wire.NewSet(
	controller.JobsCtrlInit,
	wire.Bind(new(controller.JobsCtrl), new(*controller.JobsCtrlImpl)),
)

var authCtrlSet = wire.NewSet(
	controller.AuthCtrlInit,
	wire.Bind(new(controller.AuthCtrl), new(*controller.AuthCtrlImpl)),
)

var cmpCtrlSet = wire.NewSet(
	controller.CmpCtrlInit,
	wire.Bind(new(controller.CmpCtrl), new(*controller.CmpCtrlImpl)),
)

var repairCtrlSet = wire.NewSet(
	controller.RepairCtrlInit,
	wire.Bind(new(controller.RepairCtrl), new(*controller.RepairCtrlImpl)),
)

var socketCtrlSet = wire.NewSet(
	controller.SocketCtrlInit,
	wire.Bind(new(controller.SocketCtrl), new(*controller.SocketCtrlImpl)),
)

var roleMidSet = wire.NewSet(
	middleware.RoleMidInit,
	wire.Bind(new(middleware.RoleMid), new(*middleware.RoleMidImpl)),
)

func Init(q *db.Queries, conn *sql.DB) *controller.AppControllerImpl {
	wire.Build(
		cmpCtrlSet,
		controller.AppControllerInit,
		userCtrlSet,
		userServSet,
		authCtrlSet,
		cmpServSet,
		service.AppServiceInit,
		jobsCtrlSet,
		jobsServSet,
		repairCtrlSet,
		repairServSet,
		alertServSet,
		alertCtrlSet,
		socketCtrlSet,
	)

	return nil
}

func MInit(q *db.Queries, conn *sql.DB) *middleware.AppMiddlewareImpl {
	wire.Build(
		userServSet,
		cmpServSet,
		service.AppServiceInit,
		jobsServSet,
		repairServSet,
		middleware.AppMiddlewareInit,
		roleMidSet,
		alertServSet,
	)

	return nil
}
