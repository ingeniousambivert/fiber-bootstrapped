package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type Request struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required,min=8"`
}
type Response struct {
	Token string             `json:"token" bson:"token"`
	ID    primitive.ObjectID `json:"id" bson:"_id"`
}
