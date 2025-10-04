package models

type User struct {
	ID         string   `json:"_id" 	  bson:"_id"`             //Id of user
	ProjectIDs []string `json:"projectIds"  bson:"projectIds"` //Array of IDs of projects
	ReportIDs  []string `json:"reportIds" bson:"reportIds"`    //Array of reports submitted by user
	Username   string   `json:"username"  bson:"username"`     //Username of user
	Email      string   `json:"email" 	  bson:"email"`         //Email of user
	Password   string   `json:"-" 		  bson:"password"`         //Password of user, it will not be sent over API?
}

// Data structure of the json object received in POST to create user
type UserDetails struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
