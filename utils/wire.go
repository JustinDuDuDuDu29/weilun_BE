// go:build wireinject
//go:build wireinject
// +build wireinject

package utils

import (
	"main/controller"
	"main/service"
	db "main/sql"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5"
)

var cmpServSet = wire.NewSet(
	service.CmpServInit,
	wire.Bind(new(service.CmpServ), new(*service.CmpServImpl)),
)

var userServSet = wire.NewSet(
	service.UserServInit,
	wire.Bind(new(service.UserServ), new(*service.UserServImpl)),
)

var userCtrlSet = wire.NewSet(
	controller.UserCtrlInit,
	wire.Bind(new(controller.UserCtrl), new(*controller.UserCtrlImpl)),
)

var authCtrlSet = wire.NewSet(
	controller.AuthCtrlInit,
	wire.Bind(new(controller.AuthCtrl), new(*controller.AuthCtrlImpl)),
)

var cmpCtrlSet = wire.NewSet(
	controller.CmpCtrlInit,
	wire.Bind(new(controller.CmpCtrl), new(*controller.CmpCtrlImpl)),
)

func Init(q *db.Queries, conn *pgx.Conn) *controller.AppControllerImpl {
	wire.Build(
		cmpCtrlSet,
		controller.AppControllerInit,
		userCtrlSet,
		userServSet,
		authCtrlSet,
		cmpServSet,
		service.AppServiceInit,
	)

	return nil
}
