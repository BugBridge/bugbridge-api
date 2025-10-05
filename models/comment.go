package models

type Comment struct {
	ID       string `json:"_id"      bson:"_id"`      //Id of comment
	AuthorID string `json:"authorId" bson:"authorId"` //Id of who wrote the comment
	ReportID string `json:"reportId" bson:"reportId"` //Id of the report the comment is under
	Content  string `json:"content"  bson:"Content"`  //Content of the comment
}

// Data structure of the json object received in POST to create comment
type CommentDetails struct {
	AuthorID string `json:"authorId"` //Id of who wrote the comment
	ReportID string `json:"reportId"` //Id of the report the comment is under
	Content  string `json:"content"`  //Content of the comment
}
