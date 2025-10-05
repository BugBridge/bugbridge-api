package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	ID       primitive.ObjectID `json:"_id"      bson:"_id"`      //Id of comment
	AuthorID string             `json:"authorId" bson:"authorId"` //Id of who wrote the comment
	ReportID string             `json:"reportId" bson:"reportId"` //Id of the report the comment is under
	Content  string             `json:"content"  bson:"Content"`  //Content of the comment
}

// Data structure of the json object received in POST to create comment
type CommentDetails struct {
	AuthorID string `json:"authorId" validate:"required"`          //Id of who wrote the comment
	ReportID string `json:"reportId" validate:"required"`          //Id of the report the comment is under
	Content  string `json:"content"  validate:"required,max=1000"` //Content of the comment
}
