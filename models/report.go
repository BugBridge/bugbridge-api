package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Report struct {
	ID        primitive.ObjectID `json:"_id"       bson:"_id"`       //report id
	Author    string             `json:"author"    bson:"Author"`    //ID of author
	Desc      string             `json:"desc"      bson:"Desc"`      //description of report
	Severity  int                `json:"severity"  bson:"Severity"`  //severity of report
	Solutions []string           `json:"solutions" bson:"Solutions"` //array of IDs of solutions
	Comments  []string           `json:"comments"  bson:"Comments"`  //Array of comments under report
	Anonymous bool               `json:"anonymous" bson:"Anonymous"` //Dictates if solutions can be done anonymously or not
}
