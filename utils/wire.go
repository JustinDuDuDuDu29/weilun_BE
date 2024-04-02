// go:build wireinject
//go:build wireinject
// +build wireinject

package utils

import (
	"main/controller"
	"main/service"
	db "main/sql"

	"github.com/google/wire"
)

// InitializeEvent creates an Event. It will error if the Event is staffed with
// a grumpy greeter.

var userServSet = wire.NewSet(
	service.UserServInit,
	wire.Bind(new(service.UserServ), new(*service.UserServImpl)),
)

var userCtrlSet = wire.NewSet(
	controller.UserCtrlInit,
	wire.Bind(new(controller.UserCtrl), new(*controller.UserCtrlImpl)),
)

func Init(q *db.Queries) *controller.AppControllerImpl {
	wire.Build(
		controller.AppControllerInit,
		userCtrlSet,
		userServSet,
		service.AppServiceInit,
	)

	return nil
}
