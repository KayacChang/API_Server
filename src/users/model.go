package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct {
	ID       primitive.ObjectID `json:"_id, omitempty" bson:"_id, omitempty"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
}