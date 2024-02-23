package services

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ingeniousambivert/fiber-bootstrapped/src/app/helpers"
	auth "github.com/ingeniousambivert/fiber-bootstrapped/src/app/services/auth/build"
	users "github.com/ingeniousambivert/fiber-bootstrapped/src/app/services/users/build"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/core"
)

func BindRouter(server *core.Server) map[string]*core.Service {
	AuthService := auth.Build(server)
	UsersService := users.Build(server)

	var services = map[string]*core.Service{}
	services[AuthService.Name] = AuthService
	services[UsersService.Name] = UsersService

	app := server.Engine
	router := app.Group("/api/v1")

	for _, service := range services {
		for method, route := range service.Router {
			controller := service.Bind(route.Controller, server)
			switch method {
			case "FIND":
				router.Get(route.Path, helpers.Validate(route.Extras.Authenticate, route.Extras.Authorize), controller)
			case "GET":
				router.Get(route.Path, helpers.Validate(route.Extras.Authenticate, route.Extras.Authorize), controller)
			case "CREATE":
				router.Post(route.Path, helpers.Validate(route.Extras.Authenticate, route.Extras.Authorize), controller)
			case "PATCH":
				router.Patch(route.Path, helpers.Validate(route.Extras.Authenticate, route.Extras.Authorize), controller)
			case "PUT":
				router.Put(route.Path, helpers.Validate(route.Extras.Authenticate, route.Extras.Authorize), controller)
			case "DELETE":
				router.Delete(route.Path, helpers.Validate(route.Extras.Authenticate, route.Extras.Authorize), controller)
			}
		}
	}

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendString("route not found")
	})

	return services
}
