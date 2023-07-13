package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// omitempty is used to omit empty fields in json
type User struct {
	ID        primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string               `json:"name,omitempty" bson:"name"`
	Username  string               `json:"username" bson:"username"`
	Email     string               `json:"email" bson:"email"`
	Password  string               `json:"password" bson:"password"`
	Online    bool                 `json:"online,omitempty" bson:"online"`
	Role      string               `json:"role" bson:"role" default:"user"`
	Bio       string               `json:"bio,omitempty" bson:"bio"`
	Picture   string               `json:"picture,omitempty" bson:"picture"`
	Likes     []primitive.ObjectID `json:"likes,omitempty" bson:"likes"`
	Followers []primitive.ObjectID `json:"followers,omitempty" bson:"followers"`
	Posts     []primitive.ObjectID `json:"posts,omitempty" bson:"posts"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time            `json:"updated_at" bson:"updated_at"`
}
