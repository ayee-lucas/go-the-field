package models

import "time"

type UserSession struct {
	ID        string    `json:"id" bson:"_id"`
	Username  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	ExpireOn  int64     `json:"expire_on" bson:"expire_on"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
