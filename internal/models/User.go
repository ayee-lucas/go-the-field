package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// omitempty is used to omit empty fields in json
type User struct {
	ID            primitive.ObjectID   `json:"id,omitempty"        bson:"_id,omitempty"`
	Name          string               `json:"name,omitempty"      bson:"name"`
	Username      string               `json:"username"            bson:"username"`
	Email         string               `json:"email"               bson:"email"`
	Password      string               `json:"password"            bson:"password"`
	Online        bool                 `json:"online,omitempty"    bson:"online"`
	Finished      bool                 `json:"finished"            bson:"finished"`
	Role          string               `json:"role"                bson:"role"          default:"user"`
	Bio           string               `json:"bio,omitempty"       bson:"bio"`
	Verified      bool                 `json:"verified"            bson:"verified"`
	Picture       Picture              `json:"picture,omitempty"   bson:"picture"`
	Org           primitive.ObjectID   `json:"org,omitempty"       bson:"org"`
	Athlete       primitive.ObjectID   `json:"athlete,omitempty"   bson:"athlete"`
	Sport         []primitive.ObjectID `json:"sport,omitempty"     bson:"sport"`
	Conversations []primitive.ObjectID `json:"conversations"       bson:"conversations"`
	Likes         []primitive.ObjectID `json:"likes,omitempty"     bson:"likes"`
	Followers     []primitive.ObjectID `json:"followers,omitempty" bson:"followers"`
	Following     []primitive.ObjectID `json:"following,omitempty" bson:"following"`
	Posts         []primitive.ObjectID `json:"posts,omitempty"     bson:"posts"`
	CreatedAt     time.Time            `json:"created_at"          bson:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"          bson:"updated_at"`
}

type Picture struct {
	PictureKey string `json:"pictureKey" bson:"pictureKey"`
	PictureURL string `json:"pictureURL" bson:"pictureURL"`
}

type Org struct {
	ID       primitive.ObjectID `json:"id,omitempty"       bson:"_id,omitempty"`
	Official bool               `json:"official"           bson:"official"`
	Country  string             `json:"country"            bson:"country"`
	Email    string             `json:"email"              bson:"email"`
	City     string             `json:"city"               bson:"city"`
	Website  string             `json:"website,omitempty"  bson:"website"`
	Sport    []string           `json:"sport"              bson:"sport"`
	Sponsors []string           `json:"sponsors,omitempty" bson:"sponsor"`
}

type Athlete struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Nationality  string             `json:"nationality"  bson:"nationality"`
	Gender       string             `json:"gender"       bson:"gender"`
	Sport        string             `json:"sport"        bson:"sport"`
	Sponsors     []string           `json:"sponsors"     bson:"sponsors"`
	CurrentTeam  string             `json:"current_team" bson:"current_team"`
	Height       int                `json:"height"       bson:"height"`
	Weight       int                `json:"weight"       bson:"weight"`
	Achievements string             `json:"achievements" bson:"achievements"`
	Contact      string             `json:"contact"      bson:"contact"`
}
