package auth

import (
	"context"

	controllers "github.com/ingeniousambivert/fiber-bootstrapped/src/app/services/auth/controllers"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/core"
)

var Name = "auth"
var Path = "/authentication"
var Service *core.Service
var Hooks core.Hooks

func Build(server *core.Server) *core.Service {
	ae := core.Entity{
		Ctx:        context.Background(),
		Collection: server.Database.Collection("users"),
	}

	Service = core.Create().
		SetName(Name).
		SetPath(Path).
		SetEntity(ae).
		AddPublicRoute("CREATE", controllers.Create).
		AddPublicRoute("PATCH", controllers.Patch)

	return Service
}
