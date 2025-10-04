package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Solution struct {
	ID 		string	`json:"_id" bson:"_id"`			//Id of solution
	Report 	string	`json:"report" bson:"Report"`	//Id of report solution is for
	Author 	string	`json:"author" bson:"Author"`	//Id of user that wrote it
	Code 	string	`json:"code" bson:"Code"`		//Solution code
}