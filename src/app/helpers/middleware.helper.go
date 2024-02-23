package helpers

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/ingeniousambivert/fiber-bootstrapped/src/core"
)

func Validate(authenticate bool, authorize bool) fiber.Handler {
	config := core.Configuration()
	if authenticate && authorize {
		return jwtware.New(jwtware.Config{
			ContextKey: "auth",
			SigningKey: jwtware.SigningKey{Key: []byte(config.JWT_SECRET)},
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return Unauthorized(err.Error())
			},
			SuccessHandler: func(c *fiber.Ctx) error {
				auth := c.Locals("auth").(*jwt.Token)
				claims := auth.Claims.(jwt.MapClaims)
				c.Locals("user", claims["id"].(string))
				role := claims["role"].(string)
				if role == "admin" {
					return c.Next()
				} else {
					return Forbidden("user not authorized")
				}
			},
		})
	} else if authenticate {
		return jwtware.New(jwtware.Config{
			ContextKey: "auth",
			SigningKey: jwtware.SigningKey{Key: []byte(config.JWT_SECRET)},
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return Unauthorized(err.Error())
			},
			SuccessHandler: func(c *fiber.Ctx) error {
				auth := c.Locals("auth").(*jwt.Token)
				claims := auth.Claims.(jwt.MapClaims)
				c.Locals("user", claims["id"].(string))
				return c.Next()
			},
		})
	} else {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}
}
