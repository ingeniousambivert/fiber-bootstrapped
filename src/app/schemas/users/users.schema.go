package schemas

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type Request struct {
	Firstname     string      `json:"firstname" bson:"firstname" binding:"required"`
	Lastname      string      `json:"lastname" bson:"lastname" binding:"required"`
	Email         string      `json:"email" bson:"email" binding:"required"`
	Password      string      `json:"password" bson:"password" binding:"required,min=8"`
	Archived      bool        `json:"archived" bson:"archived"`
	Role          Role        `json:"role" bson:"role"`
	Verified      bool        `json:"verified" bson:"verified"`
	VerifyToken   string      `json:"verify_token" bson:"verify_token"`
	VerifyExpires time.Time   `json:"verify_expires" bson:"verify_expires"`
	ResetToken    string      `json:"reset_token" bson:"reset_token"`
	ResetExpires  time.Time   `json:"reset_expires" bson:"reset_expires"`
	CreatedAt     time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" bson:"updated_at"`
	Metadata      interface{} `json:"metadata" bson:"metadata"`
}

type Raw struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Firstname     string             `json:"firstname" bson:"firstname" `
	Lastname      string             `json:"lastname" bson:"lastname" `
	Email         string             `json:"email" bson:"email"`
	Password      string             `json:"password" bson:"password"`
	Archived      bool               `json:"archived" bson:"archived"`
	Role          Role               `json:"role" bson:"role"`
	Verified      bool               `json:"verified" bson:"verified"`
	VerifyToken   string             `json:"verify_token,omitempty" bson:"verify_token"`
	VerifyExpires time.Time          `json:"verify_expires,omitempty" bson:"verify_expires"`
	ResetToken    string             `json:"reset_token,omitempty" bson:"reset_token"`
	ResetExpires  time.Time          `json:"reset_expires,omitempty" bson:"reset_expires"`
	CreatedAt     time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	Metadata      interface{}        `json:"metadata" bson:"metadata"`
}

type Response struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Firstname     string             `json:"firstname" bson:"firstname" `
	Lastname      string             `json:"lastname" bson:"lastname" `
	Email         string             `json:"email" bson:"email"`
	Archived      bool               `json:"archived" bson:"archived"`
	Role          Role               `json:"role" bson:"role"`
	Verified      bool               `json:"verified" bson:"verified"`
	VerifyToken   string             `json:"verify_token,omitempty" bson:"verify_token"`
	VerifyExpires time.Time          `json:"verify_expires,omitempty" bson:"verify_expires"`
	ResetToken    string             `json:"reset_token,omitempty" bson:"reset_token"`
	ResetExpires  time.Time          `json:"reset_expires,omitempty" bson:"reset_expires"`
	CreatedAt     time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	Metadata      interface{}        `json:"metadata" bson:"metadata"`
}

func GenerateResponse(raw *Raw) Response {
	return Response{
		ID:            raw.ID,
		Firstname:     raw.Firstname,
		Lastname:      raw.Lastname,
		Email:         raw.Email,
		Archived:      raw.Archived,
		Role:          raw.Role,
		Verified:      raw.Verified,
		VerifyToken:   raw.VerifyToken,
		VerifyExpires: raw.VerifyExpires,
		ResetToken:    raw.ResetToken,
		ResetExpires:  raw.ResetExpires,
		CreatedAt:     raw.CreatedAt,
		UpdatedAt:     raw.UpdatedAt,
		Metadata:      raw.Metadata,
	}
}
