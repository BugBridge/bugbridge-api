package models

type User struct {
	ID         string   `json:"_id" 	     bson:"_id"`         //Id of user
	ProjectIDs []string `json:"projects"  bson:"Projects"`    //Array of IDs of projects
	ReportIDs  []string `json:"reportIds"   bson:"ReportIDs"` //Array of reports submitted by user
	Username   string   `json:"username"  bson:"Username"`    //Username of user
	Email      string   `json:"email" 	 bson:"Email"`         //Email of user
	Password   string   `json:"-" 		 bson:"Password"`         //Password of user, it will not be sent over API?
}

// User details to save to user in DB, received in API call
type UserDetails struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
