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

func Init(q *db.Queries, conn *pgx.Conn) *controller.AppControllerImpl {
	wire.Build(
		controller.AppControllerInit,
		userCtrlSet,
		userServSet,
		authCtrlSet,
		service.AppServiceInit,
	)

	return nil
}
