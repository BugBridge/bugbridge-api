package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID 			string 		`json:"_id" bson:"_id"`				//Id of user
	Projects	[]string	`json:"projects" bson:"Projects"` 	//Array of IDs of projects
	Reports 	[]string	`json:"reports" bson:"Reports"`		//Array of reports done by user
	Solutions 	[]string	`json:"solutions" bson:"Solutions"`	//Array of solutions posted by user
	Username	string		`json:"username" bson:"Username"`	//Username of user
	Email		string		`json:"email" bson:"Email"`			//Email of user
	Password	string 		`json:"-" bson:"Password"`			//Password of user, it will not be sent over API?
}