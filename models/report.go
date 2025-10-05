package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Report struct {
	ID         primitive.ObjectID `json:"_id"        bson:"_id"`        //report id
	AuthorID   string             `json:"author"     bson:"author"`     //ID of author
	Title      string             `json:"title"      bson:"title"`      // title of id
	Desc       string             `json:"desc"       bson:"desc"`       //description of report
	Severity   int                `json:"severity"   bson:"severity"`   //severity of report
	Resolved   bool               `json:"resolved"   bson:"resolved"`   //array of IDs of solutions
	CommentIDs []string           `json:"commentIds" bson:"commentIds"` //Array of comments under report
}

// Data structure of the json object received in POST to create report
type ReportDetails struct {
	AuthorID string `json:"author"` //ID of author
	Title    string `json:"title"`  // title of the report
	Desc     string `json:"desc"`   //description of report
}
