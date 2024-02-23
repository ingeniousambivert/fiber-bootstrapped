package schemas

type Action string

const (
	SendEmailVerification     Action = "SendEmailVerification"
	EmailVerificationComplete Action = "EmailVerificationComplete"
	SendPasswordReset         Action = "SendPasswordReset"
	PasswordResetComplete     Action = "PasswordResetComplete"
	EmailUpdate               Action = "EmailUpdate"
	PasswordUpdate            Action = "PasswordUpdate"
)

type Request struct {
	Action Action                 `json:"action" bson:"action" binding:"required"`
	Data   map[string]interface{} `json:"data" bson:"data" binding:"required"`
}

type Response struct {
	Link string `json:"link" bson:"link"`
}
