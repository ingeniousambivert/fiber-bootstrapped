package core

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Name = string
type Path = string

type FindHandlerOptions = *options.FindOptions
type GetHandlerOptions = *options.FindOneOptions
type CreateHandlerOptions = *options.InsertOneOptions
type PatchHandlerOptions = *options.FindOneAndUpdateOptions
type DeleteHandlerOptions = *options.FindOneAndDeleteOptions

type FindHandlerResponse struct {
	Result    *mongo.Cursor
	Exception error
}
type GetHandlerResponse struct {
	Result    *mongo.SingleResult
	Exception error
}
type CreateHandlerResponse struct {
	Result    *mongo.InsertOneResult
	Exception error
}
type PatchHandlerResponse struct {
	Result    *mongo.SingleResult
	Exception error
}
type DeleteHandlerResponse struct {
	Result    *mongo.SingleResult
	Exception error
}

type FindHandler func(customFilter interface{}, customOptions FindHandlerOptions) FindHandlerResponse
type GetHandler func(customFilter interface{}, customOptions GetHandlerOptions) GetHandlerResponse
type CreateHandler func(customPayload interface{}, customOptions CreateHandlerOptions) CreateHandlerResponse
type PatchHandler func(customFilter interface{}, customPayload interface{}, customOptions PatchHandlerOptions) PatchHandlerResponse
type DeleteHandler func(customFilter interface{}, customOptions DeleteHandlerOptions) DeleteHandlerResponse

type Handler struct {
	Find   FindHandler
	Get    GetHandler
	Create CreateHandler
	Patch  PatchHandler
	Delete DeleteHandler
}

type Entity struct {
	Ctx        context.Context
	Collection *mongo.Collection
}

type Extras struct {
	Authenticate bool
	Authorize    bool
}

type Controller func(params map[string]interface{}) error

type Route struct {
	Path       string
	Controller Controller
	Extras     Extras
}
type Router map[string]Route
type HookFunc func(c *fiber.Ctx) error
type Hooks struct {
	Before  HookFunc
	After   HookFunc
	OnError HookFunc
}

type Service struct {
	Name    Name
	Path    Path
	Entity  Entity
	Handler Handler
	Router  Router
	Hooks   Hooks
}

func Create() *Service {
	return &Service{}
}

func (s *Service) SetName(n Name) *Service {
	s.Name = n
	return s
}
func (s *Service) SetEntity(e Entity) *Service {
	s.Entity = e

	h := Handler{}

	h.Find = func(customFilter interface{}, customOptions FindHandlerOptions) FindHandlerResponse {
		cursor, err := e.Collection.Find(e.Ctx, customFilter, customOptions)
		if err != nil {
			return FindHandlerResponse{
				Result:    nil,
				Exception: err,
			}
		}
		return FindHandlerResponse{
			Result:    cursor,
			Exception: nil,
		}

	}

	h.Get = func(customFilter interface{}, customOptions GetHandlerOptions) GetHandlerResponse {
		result := e.Collection.FindOne(e.Ctx, customFilter, customOptions)
		if result.Err() != nil {
			return GetHandlerResponse{
				Result:    nil,
				Exception: result.Err(),
			}
		}
		return GetHandlerResponse{
			Result:    result,
			Exception: nil,
		}
	}

	h.Create = func(customPayload interface{}, customOptions CreateHandlerOptions) CreateHandlerResponse {
		result, err := e.Collection.InsertOne(e.Ctx, customPayload, customOptions)
		if err != nil {
			return CreateHandlerResponse{
				Result:    nil,
				Exception: err,
			}
		}
		return CreateHandlerResponse{
			Result:    result,
			Exception: nil,
		}
	}

	h.Patch = func(customFilter interface{}, customPayload interface{}, customOptions PatchHandlerOptions) PatchHandlerResponse {
		customOptions.SetReturnDocument(options.After)
		customData := bson.M{"$set": customPayload}
		result := e.Collection.FindOneAndUpdate(e.Ctx, customFilter, customData, customOptions)
		if result.Err() != nil {
			return PatchHandlerResponse{
				Result:    nil,
				Exception: result.Err(),
			}
		}
		return PatchHandlerResponse{
			Result:    result,
			Exception: nil,
		}
	}

	h.Delete = func(customFilter interface{}, customOptions DeleteHandlerOptions) DeleteHandlerResponse {
		result := e.Collection.FindOneAndDelete(e.Ctx, customFilter, customOptions)
		if result.Err() != nil {
			return DeleteHandlerResponse{
				Result:    nil,
				Exception: result.Err(),
			}
		}
		return DeleteHandlerResponse{
			Result:    result,
			Exception: nil,
		}
	}

	s.Handler = h
	return s
}

func (s *Service) SetPath(p Path) *Service {
	s.Path = p
	return s
}

func (s *Service) AddPublicRoute(method string, controller func(params map[string]interface{}) error, path ...string) *Service {
	if s.Router == nil {
		s.Router = make(map[string]Route)
	}
	route := Route{}
	route.Path = s.Path
	if len(path) > 0 {
		route.Path += path[0]
	}
	route.Extras = Extras{
		Authenticate: false,
		Authorize:    false,
	}
	route.Controller = controller
	s.Router[method] = route
	return s
}
func (s *Service) AddPrivateRoute(method string, controller func(params map[string]interface{}) error, path ...string) *Service {
	if s.Router == nil {
		s.Router = make(map[string]Route)
	}
	route := Route{}
	route.Path = s.Path
	if len(path) > 0 {
		route.Path += path[0]
	}
	route.Controller = controller
	route.Extras = Extras{
		Authenticate: true,
		Authorize:    false,
	}
	s.Router[method] = route
	return s
}

func (s *Service) AddProtectedRoute(method string, controller func(params map[string]interface{}) error, path ...string) *Service {
	if s.Router == nil {
		s.Router = make(map[string]Route)
	}
	route := Route{}
	route.Path = s.Path
	if len(path) > 0 {
		route.Path += path[0]
	}
	route.Controller = controller
	route.Extras = Extras{
		Authenticate: true,
		Authorize:    true,
	}
	s.Router[method] = route
	return s
}

func (s *Service) SetHooks(h Hooks) *Service {
	s.Hooks = h
	return s
}

func (s *Service) Bind(controller func(params map[string]interface{}) error, server *Server) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		params := make(map[string]interface{})
		params["ctx"] = c
		params["handler"] = s.Handler
		params["entity"] = s.Entity
		params["database"] = server.Database
		params["app"] = server.App

		if s.Hooks.Before != nil {
			if err := s.Hooks.Before(c); err != nil {
				return err
			}
		}

		err := controller(params)
		if err != nil {
			if s.Hooks.OnError != nil {
				if hookErr := s.Hooks.OnError(c); hookErr != nil {
					return hookErr
				}
			}
			return err
		}

		if s.Hooks.After != nil {
			if err := s.Hooks.After(c); err != nil {
				return err
			}
		}

		return nil
	}
}
