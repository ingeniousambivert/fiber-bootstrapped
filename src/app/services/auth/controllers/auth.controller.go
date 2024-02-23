package auth

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ingeniousambivert/fiber-bootstrapped/src/app/helpers"
	auth_schema "github.com/ingeniousambivert/fiber-bootstrapped/src/app/schemas/auth"
	auth_manage_schema "github.com/ingeniousambivert/fiber-bootstrapped/src/app/schemas/auth/manage"
	users_schema "github.com/ingeniousambivert/fiber-bootstrapped/src/app/schemas/users"
	auth_utils "github.com/ingeniousambivert/fiber-bootstrapped/src/app/services/auth/utils"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/app/utils"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/core"
)

func Create(params map[string]interface{}) error {
	config := core.Configuration()
	c, ok := params["ctx"].(*fiber.Ctx)
	if !ok {
		return helpers.Unexpected("missing ctx")
	}
	h, ok := params["handler"].(core.Handler)
	if !ok {
		return helpers.Unexpected("missing handler")
	}

	payload := new(auth_schema.Request)
	err := c.BodyParser(payload)
	if err != nil {
		return helpers.Unexpected(err.Error())
	}
	if !utils.IsString(payload.Email) {
		return helpers.BadRequest("missing payload: email")
	}
	if !utils.IsString(payload.Password) {
		return helpers.BadRequest("missing payload: password")
	}

	filter := map[string]interface{}{"email": strings.ToLower(payload.Email)}
	findOptions := options.FindOneOptions{}
	findResponse := h.Get(filter, &findOptions)
	if findResponse.Exception != nil {
		if findResponse.Exception == mongo.ErrNoDocuments {
			return helpers.NotFound("document not found")
		} else {
			return helpers.Unexpected(findResponse.Exception.Error())
		}
	}
	var user users_schema.Raw
	findResponse.Result.Decode(&user)
	err = utils.VerifyPassword(user.Password, payload.Password)
	if err != nil {
		return helpers.Unauthorized("invalid password")
	}
	claims := jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * time.Duration(config.JWT_EXPIRY)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		return helpers.Unexpected("could not generate token")
	}
	response := auth_schema.Response{Token: jwt, ID: user.ID}
	c.Locals("response", response)
	return c.
		Status(utils.HttpStatusOK).
		JSON(response)
}

func Patch(params map[string]interface{}) error {
	c, ok := params["ctx"].(*fiber.Ctx)
	if !ok {
		return helpers.Unexpected("missing ctx")
	}
	h, ok := params["handler"].(core.Handler)
	if !ok {
		return helpers.Unexpected("missing handler")
	}
	var user users_schema.Raw
	payload := new(auth_manage_schema.Request)
	err := c.BodyParser(payload)
	if err != nil {
		return helpers.Unexpected(err.Error())
	}
	if !utils.IsString(payload.Action) {
		return helpers.BadRequest("missing/invalid payload: action")
	}
	if !utils.IsMap(payload.Data) {
		return helpers.BadRequest("missing/invalid payload: data")
	}

	switch payload.Action {
	case auth_manage_schema.EmailVerificationComplete:
		{
			filter := map[string]interface{}{}
			if !utils.IsString(payload.Data["token"].(string)) {
				return helpers.BadRequest("missing/invalid payload: token")
			}
			filter["verify_token"] = payload.Data["token"]
			findOptions := options.FindOneOptions{}
			findResponse := h.Get(filter, &findOptions)
			if findResponse.Exception != nil {
				if findResponse.Exception == mongo.ErrNoDocuments {
					return helpers.NotFound("user not found")
				}
				return helpers.Unexpected(findResponse.Exception.Error())
			}
			err = findResponse.Result.Decode(&user)
			if err != nil {
				return helpers.Unexpected(err.Error())
			}
			if utils.IsPast(user.VerifyExpires) {
				return helpers.Unauthorized("expired token")
			}
			_, err = utils.VerifyUUIDs(user.VerifyToken, payload.Data["token"].(string))
			if err != nil {
				return helpers.Unauthorized("invalid token")
			}
			patchOptions := options.FindOneAndUpdateOptions{}
			update := map[string]interface{}{
				"verified":       true,
				"verify_token":   nil,
				"verify_expires": nil,
			}
			patchResponse := h.Patch(filter, update, &patchOptions)
			if patchResponse.Exception != nil {
				if patchResponse.Exception == mongo.ErrNoDocuments {
					return helpers.NotFound("document not found")
				} else {
					return helpers.Unexpected(patchResponse.Exception.Error())
				}
			}
			patchResponse.Result.Decode(&user)
		}

	case auth_manage_schema.SendPasswordReset:
		{
			filter := map[string]interface{}{}
			if !utils.IsString(payload.Data["email"].(string)) {
				return helpers.BadRequest("missing/invalid payload: data['email']")
			}
			filter["email"] = payload.Data["email"].(string)
			findOptions := options.FindOneOptions{}
			findResponse := h.Get(filter, &findOptions)
			if findResponse.Exception != nil {
				if findResponse.Exception == mongo.ErrNoDocuments {
					return helpers.NotFound("user not found")
				}
				return helpers.Unexpected(findResponse.Exception.Error())
			}
			findResponse.Result.Decode(&user)

			patchOptions := options.FindOneAndUpdateOptions{}
			update := map[string]interface{}{
				"reset_token":   uuid.New().String(),
				"reset_expires": time.Now().Add(time.Hour * 24),
			}
			patchResponse := h.Patch(filter, update, &patchOptions)
			if patchResponse.Exception != nil {
				if patchResponse.Exception == mongo.ErrNoDocuments {
					return helpers.NotFound("document not found")
				} else {
					return helpers.Unexpected(patchResponse.Exception.Error())
				}
			}
			patchResponse.Result.Decode(&user)
		}

	case auth_manage_schema.PasswordResetComplete:
		{
			filter := map[string]interface{}{}
			if !utils.IsString(payload.Data["token"].(string)) {
				return helpers.BadRequest("missing/invalid payload: token")
			}
			if !utils.IsString(payload.Data["newPassword"]) {
				return helpers.BadRequest("missing/invalid payload: newPassword")
			}

			filter["reset_token"] = payload.Data["token"].(string)
			findOptions := options.FindOneOptions{}
			findResponse := h.Get(filter, &findOptions)
			if findResponse.Exception != nil {
				if findResponse.Exception == mongo.ErrNoDocuments {
					return helpers.NotFound("user not found")
				}
				return helpers.Unexpected(findResponse.Exception.Error())
			}
			findResponse.Result.Decode(&user)

			if utils.IsPast(user.ResetExpires) {
				return helpers.Unauthorized("expired token")
			}
			_, err = utils.VerifyUUIDs(user.ResetToken, payload.Data["token"].(string))
			if err != nil {
				return helpers.Unauthorized("invalid token")
			}
			hashedPassword, err := utils.HashPassword(payload.Data["newPassword"].(string))
			if err != nil {
				return helpers.Unexpected(err.Error())
			}
			patchOptions := options.FindOneAndUpdateOptions{}
			update := map[string]interface{}{
				"password":      hashedPassword,
				"reset_token":   nil,
				"reset_expires": nil,
			}
			patchResponse := h.Patch(filter, update, &patchOptions)
			if patchResponse.Exception != nil {
				if patchResponse.Exception == mongo.ErrNoDocuments {
					return helpers.NotFound("document not found")
				} else {
					return helpers.Unexpected(patchResponse.Exception.Error())
				}
			}
			patchResponse.Result.Decode(&user)
		}

	case auth_manage_schema.EmailUpdate:
		{
			filter := map[string]interface{}{}
			if !utils.IsString(payload.Data["email"]) {
				return helpers.BadRequest("missing/invalid payload: email")
			}
			if !utils.IsString(payload.Data["newEmail"]) {
				return helpers.BadRequest("missing/invalid payload: newEmail")
			}
			if !utils.IsString(payload.Data["password"]) {
				return helpers.BadRequest("missing/invalid payload: password")
			}
			filter["email"] = payload.Data["email"].(string)
			findOptions := options.FindOneOptions{}
			findResponse := h.Get(filter, &findOptions)
			if findResponse.Exception != nil {
				if findResponse.Exception == mongo.ErrNoDocuments {
					return helpers.NotFound("user not found")
				}
				return helpers.Unexpected(findResponse.Exception.Error())
			}
			findResponse.Result.Decode(&user)
			err = utils.VerifyPassword(user.Password, payload.Data["password"].(string))
			if err != nil {
				return helpers.Unauthorized("invalid password")
			}
			patchOptions := options.FindOneAndUpdateOptions{}
			update := map[string]interface{}{
				"verified":       false,
				"verify_token":   uuid.New().String(),
				"verify_expires": time.Now().Add(time.Hour * 168),
				"email":          strings.ToLower(payload.Data["newEmail"].(string)),
			}
			patchResponse := h.Patch(filter, update, &patchOptions)

			if patchResponse.Exception != nil {
				if patchResponse.Exception == mongo.ErrNoDocuments {
					return helpers.NotFound("document not found")
				}
				if er, ok := patchResponse.Exception.(mongo.CommandError); ok && er.Code == 11000 {
					return helpers.Conflict("email already exists")
				}

				return helpers.Unexpected(patchResponse.Exception.Error())
			}
			patchResponse.Result.Decode(&user)
		}

	case auth_manage_schema.PasswordUpdate:
		{
			filter := map[string]interface{}{}
			if !utils.IsString(payload.Data["email"]) {
				return helpers.BadRequest("missing/invalid payload: email")
			}
			if !utils.IsString(payload.Data["password"]) {
				return helpers.BadRequest("missing/invalid payload: password")
			}
			if !utils.IsString(payload.Data["newPassword"]) {
				return helpers.BadRequest("missing/invalid payload: newPassword")
			}
			filter["email"] = payload.Data["email"].(string)
			findOptions := options.FindOneOptions{}
			findResponse := h.Get(filter, &findOptions)
			if findResponse.Exception != nil {
				if findResponse.Exception == mongo.ErrNoDocuments {
					return helpers.NotFound("user not found")
				}
				return helpers.Unexpected(findResponse.Exception.Error())
			}
			findResponse.Result.Decode(&user)
			err = utils.VerifyPassword(user.Password, payload.Data["password"].(string))
			if err != nil {
				return helpers.Unauthorized("invalid password")
			}
			patchOptions := options.FindOneAndUpdateOptions{}
			hashedPassword, err := utils.HashPassword(payload.Data["newPassword"].(string))
			if err != nil {
				return helpers.Unexpected(err.Error())
			}
			update := map[string]interface{}{
				"password": hashedPassword,
			}
			patchResponse := h.Patch(filter, update, &patchOptions)
			if patchResponse.Exception != nil {
				if patchResponse.Exception == mongo.ErrNoDocuments {
					return helpers.NotFound("document not found")
				} else {
					return helpers.Unexpected(patchResponse.Exception.Error())
				}
			}
			patchResponse.Result.Decode(&user)
		}
	}

	if utils.IsZeroOrNil(user) {
		filter := map[string]interface{}{}
		if !utils.IsString(payload.Data["email"].(string)) {
			return helpers.BadRequest("missing/invalid payload: data['email']")
		}
		filter["email"] = payload.Data["email"].(string)
		findOptions := options.FindOneOptions{}
		findResponse := h.Get(filter, &findOptions)
		if findResponse.Exception != nil {
			if findResponse.Exception == mongo.ErrNoDocuments {
				return helpers.NotFound("user not found")
			}
			return helpers.Unexpected(findResponse.Exception.Error())
		}
		findResponse.Result.Decode(&user)

	}

	payload.Data["user"] = users_schema.GenerateResponse(&user)
	result, err := auth_utils.Notifier(*payload)
	if err != nil {
		return err
	}
	response := auth_manage_schema.Response{Link: result}
	c.Locals("response", response)
	return c.
		Status(utils.HttpStatusOK).
		JSON(response)

}
