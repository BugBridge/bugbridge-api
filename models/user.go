package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `json:"_id"        bson:"_id"`        // Id of user
	ProjectIDs []string           `json:"projectIds" bson:"projectIds"` // Project IDs that user is a member of
	Username   string             `json:"username"   bson:"username"`   // Username of user
	Email      string             `json:"email"      bson:"email"`      // Email of user
	Password   string             `json:"-"          bson:"password"`   // Password of user, it will not be sent over API?
}

// Data structure of the json object received in POST to create user
type UserDetails struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
