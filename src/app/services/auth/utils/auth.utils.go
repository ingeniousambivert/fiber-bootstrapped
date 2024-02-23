package auth

import (
	"github.com/ingeniousambivert/fiber-bootstrapped/src/app/helpers"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/app/modules"
	auth_manage_schema "github.com/ingeniousambivert/fiber-bootstrapped/src/app/schemas/auth/manage"
	users_schema "github.com/ingeniousambivert/fiber-bootstrapped/src/app/schemas/users"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/app/utils"
	"github.com/ingeniousambivert/fiber-bootstrapped/src/core"
)

func GenerateLink(baseURL string, action string, hash ...string) string {
	token := ""
	if len(hash) > 0 {
		token = "?token=" + hash[0]
	}

	if action == "signin" {
		return baseURL + "/signin"
	} else {
		return baseURL + "/" + action + token
	}
}

func Notifier(payload auth_manage_schema.Request) (string, error) {
	if !utils.IsString(payload.Action) {
		return "", helpers.BadRequest("missing/invalid param: action")
	}
	if utils.IsNil(payload.Data) {
		return "", helpers.BadRequest("missing param: data")
	}
	if utils.IsNil(payload.Data["user"]) {
		return "", helpers.BadRequest("missing param: data['user']")
	}

	config := core.Configuration()
	baseURL := config.AUDIENCE
	user := payload.Data["user"].(users_schema.Response)
	from := config.MAILER.FROM
	if !utils.IsString(from) || from == "" {
		return "", helpers.Unexpected("missing/invalid MAILER['FROM']")
	}
	mailer := modules.Mailer{
		From: from,
	}

	switch payload.Action {
	case auth_manage_schema.SendEmailVerification:
		{
			result := GenerateLink(baseURL, "verify-email", user.VerifyToken)
			params := map[string]interface{}{
				"to":      user.Email,
				"subject": "SendEmailVerification",
				"body":    result,
			}
			mailer.Send(params)
			return result, nil
		}
	case auth_manage_schema.EmailVerificationComplete:
		{
			result := GenerateLink(baseURL, "signin")
			params := map[string]interface{}{
				"to":      user.Email,
				"subject": "EmailVerificationComplete",
				"body":    result,
			}
			mailer.Send(params)
			return result, nil
		}
	case auth_manage_schema.SendPasswordReset:
		{
			result := GenerateLink(baseURL, "reset-password", user.ResetToken)
			params := map[string]interface{}{
				"to":      user.Email,
				"subject": "SendPasswordReset",
				"body":    result,
			}
			mailer.Send(params)
			return result, nil
		}
	case auth_manage_schema.PasswordResetComplete:
		{
			result := GenerateLink(baseURL, "signin")
			params := map[string]interface{}{
				"to":      user.Email,
				"subject": "PasswordResetComplete",
				"body":    result,
			}
			mailer.Send(params)
			return result, nil
		}
	case auth_manage_schema.EmailUpdate:
		{
			result := GenerateLink(baseURL, "verify", user.VerifyToken)
			params := map[string]interface{}{
				"to":      user.Email,
				"subject": "EmailUpdate",
				"body":    result,
			}
			mailer.Send(params)
			return result, nil
		}
	case auth_manage_schema.PasswordUpdate:
		{
			result := GenerateLink(baseURL, "signin")
			params := map[string]interface{}{
				"to":      user.Email,
				"subject": "PasswordUpdate",
				"body":    result,
			}
			mailer.Send(params)
			return result, nil
		}

	default:
		return "", helpers.BadRequest("invalid action")
	}
}
