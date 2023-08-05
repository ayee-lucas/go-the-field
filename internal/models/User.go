package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// omitempty is used to omit empty fields in json
type User struct {
	ID            primitive.ObjectID `json:"id,omitempty"      bson:"_id,omitempty"`
	Username      string             `json:"username"          bson:"username"`
	Email         string             `json:"email"             bson:"email"`
	EmailVerified bool               `json:"email_verified"    bson:"email_verified"`
	Password      string             `json:"password"          bson:"password"`
	Role          string             `json:"role"              bson:"role"`
	Verified      bool               `json:"verified"          bson:"verified"`
	Picture       Picture            `json:"picture,omitempty" bson:"picture"`
	ProfileID     primitive.ObjectID `json:"profile_id"        bson:"profile_id"`
	CreatedAt     time.Time          `json:"created_at"        bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"        bson:"updated_at"`
}

type Profile struct {
	ID             primitive.ObjectID `json:"id,omitempty"        bson:"_id,omitempty"`
	Name           string             `json:"name,omitempty"      bson:"name"`
	Bio            string             `json:"bio,omitempty"       bson:"bio"`
	Type           primitive.ObjectID `json:"type,omitempty"      bson:"type"`
	PreferedSports []string           `json:"prefered_sports"     bson:"prefered_sports"`
	Online         bool               `json:"online,omitempty"    bson:"online"`
	Finished       bool               `json:"finished"            bson:"finished"`
	Followers      primitive.ObjectID `json:"followers,omitempty" bson:"followers"`
	Following      primitive.ObjectID `json:"following,omitempty" bson:"following"`
	Posts          primitive.ObjectID `json:"posts,omitempty"     bson:"posts"`
	Likes          primitive.ObjectID `json:"likes,omitempty"     bson:"likes"`
	Chats          primitive.ObjectID `json:"chats,omitempty"     bson:"chats"`
	CreatedAt      time.Time          `json:"created_at"          bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"          bson:"updated_at"`
}

type Picture struct {
	PictureKey string `json:"pictureKey" bson:"pictureKey"`
	PictureURL string `json:"pictureURL" bson:"pictureURL"`
}

type Team struct {
	ID       primitive.ObjectID `json:"id,omitempty"       bson:"_id,omitempty"`
	Official bool               `json:"official"           bson:"official"`
	Country  string             `json:"country"            bson:"country"`
	Email    string             `json:"email"              bson:"email"`
	City     string             `json:"city"               bson:"city"`
	Links    []string           `json:"links,omitempty"    bson:"links"`
	Sport    string             `json:"sport"              bson:"sport"`
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
