package app

import (
	"github.com/ingeniousambivert/fiber-bootstrapped/src/app/services"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/core"
)

func Init(server *core.Server) {
	server.App.InitServices(services.BindRouter(server))
}
