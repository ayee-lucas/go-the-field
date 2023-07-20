package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSession struct {
	ID        string             `json:"id" bson:"_id"`
	Sub       primitive.ObjectID `json:"sub" bson:"sub"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Role      string             `json:"role" bson:"role"`
	Image     string             `json:"image" bson:"image"`
	Picture   Picture            `json:"picture" bson:"picture"`
	ExpireOn  int64              `json:"expire_on" bson:"expire_on"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
