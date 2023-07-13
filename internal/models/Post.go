package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Author    primitive.ObjectID   `json:"author" bson:"author"`
	Content   PostContent          `json:"content" bson:"content"`
	Repost    []primitive.ObjectID `json:"repost" bson:"repost"`
	Starred   []primitive.ObjectID `json:"starred" bson:"starred"`
	Comments  []primitive.ObjectID `json:"comments" bson:"comments"`
	Likes     []primitive.ObjectID `json:"likes" bson:"likes"`
	CreatedAt time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt" bson:"updatedAt"`
}

type PostContent struct {
	Text  string   `json:"text" bson:"text"`
	Media []string `json:"media" bson:"media"`
}
