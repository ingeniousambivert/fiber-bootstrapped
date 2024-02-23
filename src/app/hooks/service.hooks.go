package hooks

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ingeniousambivert/fiber-bootstrapped/src/app/helpers"
	auth_manage_schema "github.com/ingeniousambivert/fiber-bootstrapped/src/app/schemas/auth/manage"
	users_schema "github.com/ingeniousambivert/fiber-bootstrapped/src/app/schemas/users"
	auth_utils "github.com/ingeniousambivert/fiber-bootstrapped/src/app/services/auth/utils"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/app/utils"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/core"
)

func AddVerfication(params map[string]interface{}) (interface{}, error) {
	c, ok := params["ctx"].(*fiber.Ctx)
	if !ok {
		return nil, helpers.Unexpected("missing ctx")
	}
	h, ok := params["handler"].(core.Handler)
	if !ok {
		return nil, helpers.Unexpected("missing handler")
	}
	result := c.Locals("response").(users_schema.Response)
	if utils.IsNil(result) {
		return nil, helpers.Unexpected("missing/invalid : c.Locals('response')")
	}

	payload := map[string]interface{}{
		"verified":       false,
		"verify_token":   uuid.New().String(),
		"verify_expires": time.Now().Add(time.Hour * 168),
		"reset_token":    nil,
		"reset_expires":  nil,
	}

	filter := map[string]interface{}{"_id": result.ID}
	patchOptions := options.FindOneAndUpdateOptions{}
	patchResponse := h.Patch(filter, payload, &patchOptions)
	if patchResponse.Exception != nil {
		if patchResponse.Exception == mongo.ErrNoDocuments {
			return nil, helpers.NotFound("document not found")
		} else {
			return nil, helpers.Unexpected(patchResponse.Exception.Error())
		}
	}
	var updatedUser users_schema.Response
	patchResponse.Result.Decode(&updatedUser)
	response := updatedUser
	c.Locals("response", response)
	return response, nil

}

func NotifyVerfication(params map[string]interface{}) (interface{}, error) {
	c, ok := params["ctx"].(*fiber.Ctx)
	if !ok {
		return nil, helpers.Unexpected("missing ctx")
	}
	payload := auth_manage_schema.Request{
		Action: auth_manage_schema.SendEmailVerification,
		Data:   map[string]interface{}{},
	}
	response := c.Locals("response").(users_schema.Response)
	if utils.IsNil(response) {
		return nil, helpers.Unexpected("missing/invalid : c.Locals('response')")
	}
	payload.Data["user"] = response
	result, err := auth_utils.Notifier(payload)
	return result, err
}
