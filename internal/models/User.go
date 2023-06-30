package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// omitempty is used to omit empty fields in json
type User struct {
	ID        string               `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string               `json:"name,omitempty" bson:"name,omitempty"`
	Username  string               `json:"username" bson:"username"`
	Email     string               `json:"email" bson:"email"`
	Password  string               `json:"password" bson:"password"`
	Online    bool                 `json:"online,omitempty" bson:"online,omitempty"`
	Role      string               `json:"role" bson:"role" default:"user"`
	Bio       string               `json:"bio,omitempty" bson:"bio,omitempty"`
	Likes     []primitive.ObjectID `json:"likes,omitempty" bson:"likes,omitempty"`
	Followers []primitive.ObjectID `json:"followers,omitempty" bson:"followers,omitempty"`
	Posts     []primitive.ObjectID `json:"posts,omitempty" bson:"posts,omitempty"`
	CreatedAt time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
