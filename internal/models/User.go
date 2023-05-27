package models

import "time"

type User struct {
	ID        string    `json: "id,omitempty" bson: "_id,omitempty"`
	USERNAME  string    `json: "username" bson:  "username"`
	EMAIL     string    `json: "email" bson: "email"`
	PASSWORD  string    `json: "password,omitempty" bson: "password,omitempty"`
	CreatedAt time.Time `json: "created_at,omitempty" bson: "created_at,omitempty"`
	UpdatedAt time.Time `json: "updated_at,omitempty" bson: "updated_at,omitempty"`
}
