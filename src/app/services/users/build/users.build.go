package users

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/ingeniousambivert/fiber-bootstrapped/src/app/hooks"
	schema "github.com/ingeniousambivert/fiber-bootstrapped/src/app/schemas/users"
	controllers "github.com/ingeniousambivert/fiber-bootstrapped/src/app/services/users/controllers"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/app/utils"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/core"
)

var Name = "users"
var Path = "/users"
var Service *core.Service
var Hooks core.Hooks

func Build(server *core.Server) *core.Service {
	ue := core.Entity{
		Ctx:        context.Background(),
		Collection: server.Database.Collection("users"),
	}

	Service = core.Create().
		SetName(Name).
		SetPath(Path).
		SetEntity(ue).
		AddProtectedRoute("FIND", controllers.Find).
		AddPrivateRoute("GET", controllers.Get, "/:id").
		AddPublicRoute("CREATE", controllers.Create).
		AddPrivateRoute("PATCH", controllers.Patch, "/:id").
		AddPrivateRoute("DELETE", controllers.Delete, "/:id").
		SetHooks(core.Hooks{
			Before: func(c *fiber.Ctx) error {
				switch c.Method() {
				case "GET":
					{

					}
				case "POST":
					{

					}
				case "PATCH":
					{

					}
				case "DELETE":
					{

					}
				}
				return nil
			},
			After: func(c *fiber.Ctx) error {
				switch c.Method() {
				case "GET":
					{

					}
				case "POST":
					{
						if utils.IsNil(c.Locals("response")) {
							log.Error("missing params: c.Locals('response')")
						} else {
							params := map[string]interface{}{
								"ctx":     c,
								"entity":  Service.Entity,
								"handler": Service.Handler,
							}
							updatedUser, err := hooks.AddVerfication(params)
							if err != nil {
								log.Errorf("failed to add verification data to user : %s", err.Error())
							}
							if utils.IsUUID(updatedUser.(schema.Response).VerifyToken) {
								_, err = hooks.NotifyVerfication(params)
								if err != nil {
									log.Errorf("failed to send verification notification to user : %s", err.Error())
								}
							}

						}
					}
				case "PATCH":
					{

					}
				case "DELETE":
					{

					}
				}
				return nil
			},
			OnError: func(c *fiber.Ctx) error {
				switch c.Method() {
				case "GET":
					{

					}
				case "POST":
					{

					}
				case "PATCH":
					{

					}
				case "DELETE":
					{

					}
				}
				return nil
			},
		})

	return Service
}
