package users

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ingeniousambivert/fiber-bootstrapped/src/app/helpers"
	schema "github.com/ingeniousambivert/fiber-bootstrapped/src/app/schemas/users"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/app/utils"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/core"
)

func Find(params map[string]interface{}) error {
	c, ok := params["ctx"].(*fiber.Ctx)
	if !ok {
		return helpers.Unexpected("missing ctx")
	}
	e, ok := params["entity"].(core.Entity)
	if !ok {
		return helpers.Unexpected("missing entity")
	}
	h, ok := params["handler"].(core.Handler)
	if !ok {
		return helpers.Unexpected("missing handler")
	}

	filter := c.Queries()

	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		limit = utils.Limit
	}
	delete(filter, "limit")

	skip, err := strconv.ParseInt(c.Query("skip"), 10, 64)
	if err != nil {
		skip = utils.Skip
	}
	delete(filter, "skip")

	opts := options.Find().SetLimit(limit).SetSkip(skip)
	findResponse := h.Find(filter, opts)

	if findResponse.Exception != nil {
		return helpers.Unexpected(findResponse.Exception.Error())
	}
	defer func() {
		if err := findResponse.Result.Close(e.Ctx); err != nil {
			log.Errorf("Error closing cursor:", err)
		}
	}()
	var results []schema.Response
	for findResponse.Result.Next(e.Ctx) {
		var user schema.Raw
		findResponse.Result.Decode(&user)
		results = append(results, schema.GenerateResponse(&user))
	}
	if err := findResponse.Result.Err(); err != nil {
		return helpers.Unexpected(err.Error())
	}
	total, err := e.Collection.CountDocuments(e.Ctx, filter)
	if err != nil {
		return helpers.Unexpected(err.Error())
	}

	response := map[string]interface{}{
		"data":  results,
		"total": total,
		"limit": limit,
		"skip":  skip,
	}
	c.Locals("response", response)
	return c.
		Status(utils.HttpStatusOK).
		JSON(response)
}

func Get(params map[string]interface{}) error {
	c, ok := params["ctx"].(*fiber.Ctx)
	if !ok {
		return helpers.Unexpected("missing ctx")
	}
	h, ok := params["handler"].(core.Handler)
	if !ok {
		return helpers.Unexpected("missing handler")
	}
	id := c.Params("id")
	if id != "" {
		current := c.Locals("user").(string)
		if id != current {
			id = current
		}
		oid, _ := primitive.ObjectIDFromHex(id)
		filter := map[string]interface{}{"_id": oid}
		findOptions := options.FindOneOptions{}
		findResponse := h.Get(filter, &findOptions)
		if findResponse.Exception != nil {
			if findResponse.Exception == mongo.ErrNoDocuments {
				return helpers.NotFound("document not found")
			} else {
				return helpers.Unexpected(findResponse.Exception.Error())
			}
		}
		var user schema.Raw
		findResponse.Result.Decode(&user)

		response := schema.GenerateResponse(&user)
		c.Locals("response", response)
		return c.
			Status(utils.HttpStatusOK).
			JSON(response)
	} else {
		return helpers.BadRequest("missing params: id")
	}
}

func Create(params map[string]interface{}) error {
	c, ok := params["ctx"].(*fiber.Ctx)
	if !ok {
		return helpers.Unexpected("missing ctx")
	}
	ue, ok := params["entity"].(core.Entity)
	if !ok {
		return helpers.Unexpected("missing entity")
	}
	h, ok := params["handler"].(core.Handler)
	if !ok {
		return helpers.Unexpected("missing handler")
	}
	payload := new(schema.Request)
	err := c.BodyParser(payload)
	if err != nil {
		return helpers.Unexpected(err.Error())
	}
	if !utils.IsString(payload.Firstname) {
		return helpers.BadRequest("missing payload: firstname")
	}
	if !utils.IsString(payload.Lastname) {
		return helpers.BadRequest("missing payload: lastname")
	}
	if !utils.IsString(payload.Email) {
		return helpers.BadRequest("missing payload: email")
	}
	if !utils.IsString(payload.Password) {
		return helpers.BadRequest("missing payload: password")
	}
	payload.Email = utils.SanitizeString(payload.Email)
	payload.Role = schema.UserRole
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = payload.CreatedAt
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return helpers.Unexpected(err.Error())
	}
	payload.Password = hashedPassword
	createOptions := options.InsertOneOptions{}
	createResponse := h.Create(payload, &createOptions)
	if createResponse.Exception != nil {
		if er, ok := createResponse.Exception.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return helpers.Conflict("email already exists")
		}
		return helpers.Unexpected(createResponse.Exception.Error())
	}
	utils.CreateIndex(ue, "email")

	filter := map[string]interface{}{"_id": createResponse.Result.InsertedID}
	findOptions := options.FindOneOptions{}
	findResponse := h.Get(filter, &findOptions)
	if findResponse.Exception != nil {
		return helpers.Unexpected(findResponse.Exception.Error())
	}
	var newUser schema.Raw
	findResponse.Result.Decode(&newUser)
	response := schema.GenerateResponse(&newUser)
	c.Locals("response", response)
	return c.
		Status(utils.HttpStatusCreated).
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
	id := c.Params("id")
	if id != "" {
		current := c.Locals("user").(string)
		if id != current {
			id = current
		}
		oid, _ := primitive.ObjectIDFromHex(id)
		body := c.Body()
		var payload bson.M
		err := bson.UnmarshalExtJSON(body, true, &payload)
		if err != nil {
			return helpers.Unexpected(err.Error())
		}
		delete(payload, "id")
		delete(payload, "email")
		delete(payload, "password")
		delete(payload, "verified")
		delete(payload, "role")
		delete(payload, "verify_token")
		delete(payload, "verify_expires")
		delete(payload, "reset_token")
		delete(payload, "reset_expires")
		delete(payload, "created_at")
		delete(payload, "updated_at")

		filter := map[string]interface{}{"_id": oid}
		patchOptions := options.FindOneAndUpdateOptions{}
		patchResponse := h.Patch(filter, payload, &patchOptions)
		if patchResponse.Exception != nil {
			if patchResponse.Exception == mongo.ErrNoDocuments {
				return helpers.NotFound("document not found")
			} else {
				return helpers.Unexpected(patchResponse.Exception.Error())
			}
		}
		var updatedUser schema.Response
		patchResponse.Result.Decode(&updatedUser)
		response := updatedUser
		c.Locals("response", response)
		return c.
			Status(utils.HttpStatusOK).
			JSON(response)
	} else {
		return helpers.Unexpected("missing params: id")
	}
}

func Delete(params map[string]interface{}) error {
	c, ok := params["ctx"].(*fiber.Ctx)
	if !ok {
		return helpers.Unexpected("missing ctx")
	}
	h, ok := params["handler"].(core.Handler)
	if !ok {
		return helpers.Unexpected("missing handler")
	}
	id := c.Params("id")
	if id != "" {
		current := c.Locals("user").(string)
		if id != current {
			id = current
		}
		oid, _ := primitive.ObjectIDFromHex(id)
		filter := map[string]interface{}{"_id": oid}
		deleteOptions := options.FindOneAndDeleteOptions{}
		deleteResponse := h.Delete(filter, &deleteOptions)
		if deleteResponse.Exception != nil {
			if deleteResponse.Exception == mongo.ErrNoDocuments {
				return helpers.NotFound("document not found")
			} else {
				return helpers.Unexpected(deleteResponse.Exception.Error())
			}
		}
		var deletedUser schema.Response
		deleteResponse.Result.Decode(&deletedUser)
		response := deletedUser
		c.Locals("response", response)
		return c.
			Status(utils.HttpStatusOK).
			JSON(response)

	} else {
		return helpers.BadRequest("missing params: id")
	}
}
