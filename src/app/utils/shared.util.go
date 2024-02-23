package utils

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/ingeniousambivert/fiber-bootstrapped/src/core"
)

const (
	Limit int64 = 25
	Skip  int64 = 0
)
const (
	HttpStatusContinue                      int = 100
	HttpStatusSwitchingProtocols            int = 101
	HttpStatusProcessing                    int = 102
	HttpStatusEarlyHints                    int = 103
	HttpStatusOK                            int = 200
	HttpStatusCreated                       int = 201
	HttpStatusAccepted                      int = 202
	HttpStatusNonAuthoritativeInfo          int = 203
	HttpStatusNoContent                     int = 204
	HttpStatusResetContent                  int = 205
	HttpStatusPartialContent                int = 206
	HttpStatusMultiStatus                   int = 207
	HttpStatusAlreadyReported               int = 208
	HttpStatusIMUsed                        int = 226
	HttpStatusMultipleChoices               int = 300
	HttpStatusMovedPermanently              int = 301
	HttpStatusFound                         int = 302
	HttpStatusSeeOther                      int = 303
	HttpStatusNotModified                   int = 304
	HttpStatusUseProxy                      int = 305
	HttpStatusTemporaryRedirect             int = 307
	HttpStatusPermanentRedirect             int = 308
	HttpStatusBadRequest                    int = 400
	HttpStatusUnauthorized                  int = 401
	HttpStatusPaymentRequired               int = 402
	HttpStatusForbidden                     int = 403
	HttpStatusNotFound                      int = 404
	HttpStatusMethodNotAllowed              int = 405
	HttpStatusNotAcceptable                 int = 406
	HttpStatusProxyAuthRequired             int = 407
	HttpStatusRequestTimeout                int = 408
	HttpStatusConflict                      int = 409
	HttpStatusGone                          int = 410
	HttpStatusLengthRequired                int = 411
	HttpStatusPreconditionFailed            int = 412
	HttpStatusRequestEntityTooLarge         int = 413
	HttpStatusRequestURITooLong             int = 414
	HttpStatusUnsupportedMediaType          int = 415
	HttpStatusRequestedRangeNotSatisfiable  int = 416
	HttpStatusExpectationFailed             int = 417
	HttpStatusTeapot                        int = 418
	HttpStatusMisdirectedRequest            int = 421
	HttpStatusUnprocessableEntity           int = 422
	HttpStatusLocked                        int = 423
	HttpStatusFailedDependency              int = 424
	HttpStatusTooEarly                      int = 425
	HttpStatusUpgradeRequired               int = 426
	HttpStatusPreconditionRequired          int = 428
	HttpStatusTooManyRequests               int = 429
	HttpStatusRequestHeaderFieldsTooLarge   int = 431
	HttpStatusUnavailableForLegalReasons    int = 451
	HttpStatusInternalServerError           int = 500
	HttpStatusNotImplemented                int = 501
	HttpStatusBadGateway                    int = 502
	HttpStatusServiceUnavailable            int = 503
	HttpStatusGatewayTimeout                int = 504
	HttpStatusHTTPVersionNotSupported       int = 505
	HttpStatusVariantAlsoNegotiates         int = 506
	HttpStatusInsufficientStorage           int = 507
	HttpStatusLoopDetected                  int = 508
	HttpStatusNotExtended                   int = 510
	HttpStatusNetworkAuthenticationRequired int = 511
)

func CreateIndex(e core.Entity, indexKey string) error {
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{indexKey: 1}, Options: opt}
	_, err := e.Collection.Indexes().CreateOne(e.Ctx, index)
	if err != nil {
		return errors.New("could not create index in mongodb")
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

func SanitizeString(str string) string {
	return strings.TrimSpace(strings.ToLower(str))
}

func IsNil(v interface{}) bool {
	return v == nil
}

func IsZero(v interface{}) bool {
	return v == reflect.Zero(reflect.TypeOf(v)).Interface()
}

func IsZeroOrNil(v interface{}) bool {
	return IsZero(v) || IsNil(v)
}

func VerifyUUIDs(uuid1Str, uuid2Str interface{}) (bool, error) {
	if IsNil(uuid1Str) || IsNil(uuid2Str) {
		return false, nil
	}
	var uuidOne, uuidTwo uuid.UUID
	if reflect.TypeOf(uuid1Str).Kind() == reflect.String {
		uuid1, err := uuid.Parse(uuid1Str.(string))
		if err != nil {
			return false, err
		}
		uuidOne = uuid1
	}
	if reflect.TypeOf(uuid2Str).Kind() == reflect.String {
		uuid2, err := uuid.Parse(uuid2Str.(string))
		if err != nil {
			return false, err
		}
		uuidTwo = uuid2
	}
	return uuidOne == uuidTwo, nil
}

func IsPast(t time.Time) bool {
	if IsNil(t) {
		return false
	}
	tUTC := t.UTC()
	nowUTC := time.Now().UTC()
	return nowUTC.After(tUTC)
}

func IsInt(v interface{}) bool {
	if IsNil(v) {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Int
}

func IsString(v interface{}) bool {
	if IsNil(v) {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.String
}

func IsBool(v interface{}) bool {
	if IsNil(v) {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Bool
}

func IsFloat(v interface{}) bool {
	if IsNil(v) {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Float64
}

func IsArray(v interface{}) bool {
	if IsNil(v) {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Array
}

func IsSlice(v interface{}) bool {
	if IsNil(v) {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

func IsMap(v interface{}) bool {
	if IsNil(v) {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Map
}

func IsStruct(v interface{}) bool {
	if IsNil(v) {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Struct
}

func IsUUID(v interface{}) bool {
	if IsNil(v) {
		return false
	}
	str, ok := v.(string)
	if !ok {
		return false
	}
	uuidPattern := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	regex := regexp.MustCompile(uuidPattern)
	return regex.MatchString(str)
}
