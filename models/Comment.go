package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID 		string	`json:"_id" bson:"_id"`			//Id of comment
	Author 	string	`json:"author" bson:"Author"`	//Id of who wrote the comment
	Report 	string	`json:"report" bson:"Report"`	//Id of the report the comment is under
	Content	string	`json:"content" bson:"Content"`	//Content of the comment
}