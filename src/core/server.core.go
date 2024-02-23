package core

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Server struct {
	Port     int
	Engine   *fiber.App
	App      *App
	Database *Database
}

func (s *Server) Boot() error {
	return s.Engine.Listen(fmt.Sprintf(":%v", s.Port))
}

var server *Server

const (
	TimeoutInSeconds = 30
)

type ServerError struct {
	Status  int    `json:"status"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (e *ServerError) Error() string {
	return e.Message
}

func Build() *Server {
	if server == nil {
		stage := Configuration().STAGE

		var allowOrigins string
		if stage == "development" {
			allowOrigins = "*"
		} else {
			allowOrigins = Configuration().AUDIENCE
		}

		database := InitDatabase()
		port := Configuration().PORT
		app := InitApp()
		engine := fiber.New(fiber.Config{
			ReadTimeout:  time.Duration(TimeoutInSeconds) * time.Second,
			WriteTimeout: time.Duration(TimeoutInSeconds) * time.Second,
			ErrorHandler: func(ctx *fiber.Ctx, err error) error {
				if e, ok := err.(*ServerError); ok {
					return ctx.Status(e.Status).JSON(e)
				} else if e, ok := err.(*fiber.Error); ok {
					return ctx.Status(e.Code).JSON(ServerError{Status: e.Code, Title: "internal-server", Message: e.Message})
				} else {
					return ctx.Status(500).JSON(ServerError{Status: 500, Title: "internal-server", Message: err.Error()})
				}
			},
			DisableStartupMessage: true,
		})
		engine.Use(idempotency.New())
		engine.Use(recover.New())
		engine.Use(requestid.New())
		engine.Use(logger.New(logger.Config{
			Format: "[${time}] [${locals:requestid}] ${method} ${path} ${status} [${ip}] ${error} \n",
		}))
		engine.Use(cors.New(cors.Config{
			AllowOrigins: allowOrigins,
		}))
		engine.Get("/ping", func(c *fiber.Ctx) error {
			return c.SendString("Pong!")
		})

		server = &Server{
			Port:     port,
			Engine:   engine,
			App:      app,
			Database: database,
		}
		log.Infof("%s server listening on :%d\n", stage, port)
	}
	return server
}
