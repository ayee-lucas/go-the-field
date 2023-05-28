package models

import "time"

type User struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"` // omitempty is used to omit empty fields in json
	Username  string    `json:"username" bson:"username"`          // omitempty is used to omit empty fields in json
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
